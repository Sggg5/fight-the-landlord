package session

import (
	"cmp"
	"context"
	"log"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/palemoky/fight-the-landlord/internal/game/card"
	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	"github.com/palemoky/fight-the-landlord/internal/protocol/convert"
	"github.com/palemoky/fight-the-landlord/internal/server/storage"
)

// Start 开始游戏
func (gs *GameSession) Start() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if gs.room == nil {
		return
	}

	gs.startBiddingRound()
}

// maxRedeals 最大流局次数；达到后下一局随机强制指定地主，避免无限流局
const maxRedeals = 3

// startBiddingRound 发牌并进入叫地主阶段（调用方需持有 gs.mu）
func (gs *GameSession) startBiddingRound() {
	if gs.room == nil {
		return
	}
	gs.dealNewRound()
	if gs.room == nil {
		return
	}

	// 进入叫地主阶段
	gs.state = GameStateBidding
	gs.room.State = RoomStateBidding

	// 随机选择第一个叫地主的玩家
	gs.currentBidder = rand.IntN(3)

	// 通知叫地主
	gs.notifyBidTurn()
}

// dealNewRound 重置本局状态并发牌（不进入叫地主流程；调用方需持有 gs.mu）
func (gs *GameSession) dealNewRound() {
	if gs.room == nil {
		return
	}
	// 重置叫抢与倍数状态
	gs.landlordCaller = -1
	gs.landlordCandidate = -1
	gs.bidPasses = 0
	gs.grabActions = 0
	gs.bidMultiplier = 1
	gs.bombCount = 0
	gs.landlordPlays = 0
	gs.farmerPlays = 0
	for _, p := range gs.players {
		p.Hand = nil
		p.IsLandlord = false
		if rp := gs.room.Players[p.ID]; rp != nil {
			rp.IsLandlord = false
		}
	}

	// 创建并洗牌
	gs.deck = card.NewDeck()
	gs.deck.Shuffle()

	// 发牌
	gs.deal()
}

// redeal 流局（无人叫地主）重新发牌（调用方需持有 gs.mu）
// 连续流局达到 maxRedeals 次后，重新发牌并随机强制指定地主，避免无限流局。
func (gs *GameSession) redeal() {
	if gs.room == nil {
		return
	}
	gs.redealCount++

	if gs.redealCount >= maxRedeals {
		log.Printf("🔄 房间 %s 连续 %d 次流局，重新发牌并随机强制指定地主", gs.room.Code, gs.redealCount)
		gs.dealNewRound()
		gs.setLandlord(rand.IntN(3))
		return
	}

	log.Printf("🔄 房间 %s 无人叫地主，第 %d 次流局，重新发牌", gs.room.Code, gs.redealCount)
	gs.startBiddingRound()
}

// deal 发牌
func (gs *GameSession) deal() {
	if gs.room == nil {
		return
	}
	// 每人发 17 张
	for range 17 {
		for i := range 3 {
			gs.players[i].Hand = append(gs.players[i].Hand, gs.deck[0])
			gs.deck = gs.deck[1:]
		}
	}

	// 剩余 3 张为底牌
	gs.bottomCards = gs.deck

	// 排序手牌
	for _, p := range gs.players {
		slices.SortFunc(p.Hand, func(a, b card.Card) int {
			return cmp.Compare(b.Rank, a.Rank)
		})
	}

	// 发送手牌给各玩家（先不显示底牌）
	for _, p := range gs.players {
		rp := gs.room.Players[p.ID]
		if rp == nil || rp.Client == nil {
			continue
		}
		client := rp.Client
		client.SendMessage(codec.MustNewMessage(protocol.MsgDealCards, protocol.DealCardsPayload{
			Cards:       convert.CardsToInfos(p.Hand),
			BottomCards: make([]protocol.CardInfo, 3), // 暂时不显示
		}))
	}
}

// endGame 结束游戏
func (gs *GameSession) endGame(winner *GamePlayer) {
	room := gs.room
	if room == nil || winner == nil {
		return
	}
	gs.state = GameStateEnded
	room.State = RoomStateEnded

	// 计算最终倍数与各玩家得分
	multiplier := gs.finalMultiplier(winner)
	scores := gs.computeScores(winner, multiplier)

	// 收集所有玩家剩余手牌
	playerHands := make([]protocol.PlayerHand, len(gs.players))
	for i, p := range gs.players {
		playerHands[i] = protocol.PlayerHand{
			PlayerID:   p.ID,
			PlayerName: p.Name,
			Cards:      convert.CardsToInfos(p.Hand),
		}
	}

	// 广播游戏结束
	room.Broadcast(codec.MustNewMessage(protocol.MsgGameOver, protocol.GameOverPayload{
		WinnerID:    winner.ID,
		WinnerName:  winner.Name,
		IsLandlord:  winner.IsLandlord,
		PlayerHands: playerHands,
		Multiplier:  multiplier,
		Scores:      scores,
	}))

	role := "农民"
	if winner.IsLandlord {
		role = "地主"
	}
	log.Printf("🎮 游戏结束！房间 %s，获胜者: %s (%s)，倍数: %d",
		gs.room.Code, winner.Name, role, multiplier)

	// 游戏结束，解散房间
	if !room.HasBot() {
		for _, p := range gs.players {
			rp := room.Players[p.ID]
			if rp != nil && rp.Client != nil {
				rp.Client.SetRoom("")
			}
		}
	}

	// 记录游戏结果到排行榜
	gs.recordGameResults(winner, scores)
	gs.checkAchievements(winner)
	gs.updateDailyTaskProgress()
}

// finalMultiplier 计算本局最终倍数：底倍 × 炸弹/王炸 × 春天/反春天
func (gs *GameSession) finalMultiplier(winner *GamePlayer) int {
	mult := max(gs.bidMultiplier, 1)

	// 炸弹与王炸：每个翻一倍
	for range gs.bombCount {
		mult *= 2
	}

	// 春天：地主获胜且农民一张牌都没出过
	// 反春天：农民获胜且地主只在首攻出过一手牌
	switch {
	case winner.IsLandlord && gs.farmerPlays == 0:
		mult *= 2
	case !winner.IsLandlord && gs.landlordPlays == 1:
		mult *= 2
	}

	return mult
}

// computeScores 按最终倍数计算各玩家得分（地主独自对抗两名农民）
func (gs *GameSession) computeScores(winner *GamePlayer, mult int) []protocol.PlayerScore {
	landlordWins := winner.IsLandlord
	scores := make([]protocol.PlayerScore, len(gs.players))
	for i, p := range gs.players {
		var score int
		switch {
		case p.IsLandlord && landlordWins:
			score = 2 * mult
		case p.IsLandlord && !landlordWins:
			score = -2 * mult
		case !p.IsLandlord && landlordWins:
			score = -mult
		default: // 农民获胜
			score = mult
		}
		scores[i] = protocol.PlayerScore{
			PlayerID:   p.ID,
			PlayerName: p.Name,
			IsLandlord: p.IsLandlord,
			Score:      score,
		}
	}
	return scores
}

// recordGameResults 记录游戏结果到排行榜
func (gs *GameSession) recordGameResults(winner *GamePlayer, scores []protocol.PlayerScore) {
	ctx := context.Background()
	leaderboard := gs.leaderboard
	room := gs.room
	if leaderboard == nil || !leaderboard.IsReady() || room == nil || winner == nil {
		return
	}
	roundScores := make(map[string]int, len(scores))
	for _, score := range scores {
		roundScores[score.PlayerID] = score.Score
	}

	// 计算获胜方
	landlordWins := winner.IsLandlord

	for _, p := range gs.players {
		rp := room.Players[p.ID]
		if (rp != nil && rp.Client != nil && rp.Client.IsBot()) || strings.HasPrefix(p.Name, "🤖") {
			continue // Bot 不计入排行榜
		}

		isWinner := false
		if landlordWins {
			isWinner = p.IsLandlord
		} else {
			isWinner = !p.IsLandlord
		}

		// 获取玩家名称
		playerName := p.Name
		if rp != nil && rp.Client != nil {
			playerName = rp.Client.GetName()
		}

		// 记录结果
		if err := leaderboard.RecordGameResultWithScore(ctx, p.ID, playerName, p.IsLandlord, isWinner, roundScores[p.ID]); err != nil {
			log.Printf("记录游戏结果失败: %v", err)
		}
	}
}

// checkAchievements checks for newly unlocked achievements after a game.
func (gs *GameSession) checkAchievements(winner *GamePlayer) {
	ctx := context.Background()
	lm := gs.leaderboard
	room := gs.room
	if lm == nil || !lm.IsReady() || room == nil {
		return
	}
	for _, p := range gs.players {
		rp := room.Players[p.ID]
		if (rp != nil && rp.Client != nil && rp.Client.IsBot()) || strings.HasPrefix(p.Name, "🤖") {
			continue
		}
		stats, err := lm.GetPlayerStats(ctx, p.ID)
		if err != nil || stats == nil {
			continue
		}
		newAchievements := storage.CheckNewAchievements(stats)
		if len(newAchievements) > 0 {
			stats.AchievedAchievements = append(stats.AchievedAchievements, newAchievements...)
			log.Printf("Player %s unlocked achievements: %v", p.Name, newAchievements)
		}
		if err := lm.SavePlayerStats(ctx, stats); err != nil {
			log.Printf("Failed to save player stats: %v", err)
		}
	}
}

// updateDailyTaskProgress updates daily task progress for all non-bot players after a game.
func (gs *GameSession) updateDailyTaskProgress() {
	ctx := context.Background()
	lm := gs.leaderboard
	room := gs.room
	if lm == nil || !lm.IsReady() || room == nil {
		return
	}
	for _, p := range gs.players {
		rp := room.Players[p.ID]
		if (rp != nil && rp.Client != nil && rp.Client.IsBot()) || strings.HasPrefix(p.Name, "🤖") {
			continue
		}
		stats, err := lm.GetPlayerStats(ctx, p.ID)
		if err != nil || stats == nil {
			continue
		}
		progress := make(map[string]int)
		progress["daily_win_3"] = stats.Wins
		progress["daily_landlord"] = stats.LandlordGames
		progress["daily_bomb"] = stats.BombsPlayed
		progress["daily_play_5"] = stats.TotalGames
		progress["daily_farmer"] = stats.FarmerWins
		for taskID, val := range progress {
			_ = lm.UpdateDailyTaskProgress(ctx, p.ID, taskID, val)
		}
	}
}
