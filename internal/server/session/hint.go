package session

import (
	"context"

	"github.com/palemoky/fight-the-landlord/internal/bot"
	"github.com/palemoky/fight-the-landlord/internal/game/card"
	"github.com/palemoky/fight-the-landlord/internal/game/rule"
	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/convert"
)

func (gs *GameSession) SuggestPlay(playerID, douzeroURL string) ([]protocol.CardInfo, string) {
	gs.mu.RLock()
	if gs.state != GameStatePlaying || gs.room == nil || gs.players[gs.currentPlayer].ID != playerID {
		gs.mu.RUnlock()
		return nil, "还没轮到你出牌"
	}

	player := gs.players[gs.currentPlayer]
	mustPlay := gs.lastPlayerIdx == gs.currentPlayer || gs.lastPlayedHand.IsEmpty()
	canBeat := mustPlay || rule.FindSmallestBeatingCards(player.Hand, gs.lastPlayedHand) != nil
	gctx := gs.buildHintContextLocked(player, mustPlay, canBeat)
	gs.mu.RUnlock()

	cards := bot.NewDouZeroEngine(douzeroURL).DecidePlay(context.Background(), player.Name, gctx)
	if !isSuggestedPlayValid(cards, gctx) {
		cards = rule.FindSmallestBeatingCards(gctx.Hand, gctx.RecentPlays[0].Played)
	}
	if len(cards) == 0 {
		return nil, "建议不出"
	}
	return convert.CardsToInfos(cards), "DouZero 推荐"
}

func (gs *GameSession) buildHintContextLocked(player *GamePlayer, mustPlay, canBeat bool) bot.GameContext {
	landlordSeat := -1
	for i, p := range gs.players {
		if p.IsLandlord {
			landlordSeat = i
			break
		}
	}
	douzeroPos := ""
	numCardsLeft := map[string]int{}
	if landlordSeat >= 0 {
		douzeroPos = seatToDouZeroPos(player.Seat, landlordSeat)
		for _, p := range gs.players {
			numCardsLeft[seatToDouZeroPos(p.Seat, landlordSeat)] = len(p.Hand)
		}
	}

	recent := [2]bot.PlayRecord{}
	lastMovePos := ""
	if !gs.lastPlayedHand.IsEmpty() && gs.lastPlayerIdx >= 0 && gs.lastPlayerIdx < len(gs.players) {
		last := gs.players[gs.lastPlayerIdx]
		recent[0] = bot.PlayRecord{
			Played:     gs.lastPlayedHand,
			PlayerName: last.Name,
			IsLandlord: last.IsLandlord,
		}
		if landlordSeat >= 0 {
			lastMovePos = seatToDouZeroPos(last.Seat, landlordSeat)
		}
	}

	return bot.GameContext{
		IsLandlord:   player.IsLandlord,
		Hand:         append([]card.Card(nil), player.Hand...),
		BottomCards:  append([]card.Card(nil), gs.bottomCards...),
		RecentPlays:  recent,
		MustPlay:     mustPlay,
		CanBeat:      canBeat,
		DouZeroPos:   douzeroPos,
		LastMovePos:  lastMovePos,
		NumCardsLeft: numCardsLeft,
	}
}

func isSuggestedPlayValid(cards []card.Card, gctx bot.GameContext) bool {
	if len(cards) == 0 {
		return !gctx.MustPlay
	}
	parsed, err := rule.ParseHand(cards)
	if err != nil || parsed.Type == rule.Invalid {
		return false
	}
	if gctx.MustPlay || gctx.RecentPlays[0].Played.IsEmpty() {
		return true
	}
	return rule.CanBeat(parsed, gctx.RecentPlays[0].Played)
}

func seatToDouZeroPos(seat, landlordSeat int) string {
	switch seat {
	case landlordSeat:
		return bot.DouZeroPosLandlord
	case (landlordSeat + 1) % 3:
		return bot.DouZeroPosLandlordDn
	default:
		return bot.DouZeroPosLandlordUp
	}
}
