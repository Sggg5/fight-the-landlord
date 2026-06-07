package handler

import (
	"context"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	"github.com/palemoky/fight-the-landlord/internal/server/storage"
	"github.com/palemoky/fight-the-landlord/internal/types"
)

// --- 排行榜处理 ---

// handleGetStats 获取个人统计
func (h *Handler) handleGetStats(client types.ClientInterface) {
	ctx := context.Background()
	playerStats, err := h.leaderboard.GetPlayerStats(ctx, client.GetID())
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "获取统计失败"))
		return
	}

	if playerStats == nil {
		// 没有统计数据，返回空数据
		client.SendMessage(codec.MustNewMessage(protocol.MsgStatsResult, protocol.StatsResultPayload{
			PlayerID:   client.GetID(),
			PlayerName: client.GetName(),
			Coins:      storage.InitialCoins,
		}))
		return
	}

	// 获取排名
	storage.NormalizeEconomy(playerStats)
	rank, _ := h.leaderboard.GetPlayerRank(ctx, client.GetID())

	winRate := 0.0
	if playerStats.TotalGames > 0 {
		winRate = float64(playerStats.Wins) / float64(playerStats.TotalGames) * 100
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgStatsResult, protocol.StatsResultPayload{
		PlayerID:      playerStats.PlayerID,
		PlayerName:    playerStats.PlayerName,
		TotalGames:    playerStats.TotalGames,
		Wins:          playerStats.Wins,
		Losses:        playerStats.Losses,
		WinRate:       winRate,
		LandlordGames: playerStats.LandlordGames,
		LandlordWins:  playerStats.LandlordWins,
		FarmerGames:   playerStats.FarmerGames,
		FarmerWins:    playerStats.FarmerWins,
		Score:         playerStats.Score,
		Coins:         playerStats.Coins,
		LastCoinChange: playerStats.LastCoinChange,
		BankruptcyGrants: playerStats.BankruptcyGrants,
		Rank:          int(rank),
		CurrentStreak: playerStats.CurrentStreak,
		MaxWinStreak:  playerStats.MaxWinStreak,
	}))
}

// handleGetLeaderboard 获取排行榜
func (h *Handler) handleGetLeaderboard(client types.ClientInterface, msg *protocol.Message) {
	payload, err := codec.ParsePayload[protocol.GetLeaderboardPayload](msg)
	if err != nil {
		// 默认获取总排行榜前 10
		payload = &protocol.GetLeaderboardPayload{
			Type:   "total",
			Offset: 0,
			Limit:  10,
		}
	}

	// 限制请求数量
	if payload.Limit <= 0 || payload.Limit > 50 {
		payload.Limit = 10
	}
	if payload.Offset < 0 {
		payload.Offset = 0
	}

	entries, err := h.leaderboard.GetLeaderboard(context.Background(), payload.Limit)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "获取排行榜失败"))
		return
	}

	// 转换为协议格式
	protocolEntries := make([]protocol.LeaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		protocolEntries = append(protocolEntries, protocol.LeaderboardEntry{
			Rank:       entry.Rank,
			PlayerID:   entry.PlayerID,
			PlayerName: entry.PlayerName,
			Score:      entry.Score,
			Wins:       entry.Wins,
			WinRate:    entry.WinRate,
		})
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgLeaderboardResult, protocol.LeaderboardResultPayload{
		Type:    payload.Type,
		Entries: protocolEntries,
	}))
}

// handleGetRoomList 获取房间列表
func (h *Handler) handleGetRoomList(client types.ClientInterface) {
	rooms := h.roomManager.GetRoomList()

	client.SendMessage(codec.MustNewMessage(protocol.MsgRoomListResult, protocol.RoomListResultPayload{
		Rooms: rooms,
	}))
}

// handleGetOnlineCount 获取在线人数（按需）
func (h *Handler) handleGetOnlineCount(client types.ClientInterface) {
	count := h.server.GetOnlineCount()

	client.SendMessage(codec.MustNewMessage(protocol.MsgOnlineCount, protocol.OnlineCountPayload{
		Count: count,
	}))
}

// handleGetMaintenanceStatus 获取维护状态
func (h *Handler) handleGetMaintenanceStatus(client types.ClientInterface) {
	maintenance := h.server.IsMaintenanceMode()

	client.SendMessage(codec.MustNewMessage(protocol.MsgMaintenancePull, protocol.MaintenanceStatusPayload{
		Maintenance: maintenance,
	}))
}

// handleSignIn handles daily sign-in request.
func (h *Handler) handleSignIn(client types.ClientInterface) {
	ctx := context.Background()
	playerName := client.GetName()

	reward, consecutive, err := h.leaderboard.SignIn(ctx, client.GetID(), playerName)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "????"))
		return
	}

	signedIn := reward > 0

	// Re-fetch stats for current coin total
	totalCoins := 0
	if stats, _ := h.leaderboard.GetPlayerStats(ctx, client.GetID()); stats != nil {
		totalCoins = stats.Coins
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgSignInResult, protocol.SignInResultPayload{
		SignedIn:        signedIn,
		Reward:          reward,
		ConsecutiveDays: consecutive,
		Coins:           totalCoins,
		CanSignIn:       !signedIn,
	}))
}

// handleGetAchievements handles achievement list request.
func (h *Handler) handleGetAchievements(client types.ClientInterface) {
	ctx := context.Background()
	statuses, err := h.leaderboard.GetAchievementsStatus(ctx, client.GetID())
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "??????"))
		return
	}
	if statuses == nil {
		client.SendMessage(codec.MustNewMessage(protocol.MsgAchievementsResult, protocol.AchievementsResultPayload{}))
		return
	}

	results := make([]protocol.AchievementInfo, len(statuses))
	for i, s := range statuses {
		results[i] = protocol.AchievementInfo{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Achieved:    s.Achieved,
			Progress:    s.Progress,
		}
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgAchievementsResult, protocol.AchievementsResultPayload{
		Achievements: results,
	}))
}

// handleShopList handles shop list request.
func (h *Handler) handleShopList(client types.ClientInterface) {
	items := storage.AllShopItems()
	protoItems := make([]protocol.ShopItem, len(items))
	for i, item := range items {
		protoItems[i] = protocol.ShopItem{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Category:    item.Category,
			Price:       item.Price,
		}
	}
	client.SendMessage(codec.MustNewMessage(protocol.MsgShopListResult, protocol.ShopListResultPayload{
		Items: protoItems,
	}))
}

// handlePurchaseItem handles item purchase request.
func (h *Handler) handlePurchaseItem(client types.ClientInterface, msg *protocol.Message) {
	payload, err := codec.ParsePayload[protocol.PurchaseItemPayload](msg)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeInvalidMsg, "请求格式错误"))
		return
	}

	result, err := h.leaderboard.PurchaseItem(context.Background(), client.GetID(), client.GetName(), payload.ItemID)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "购买失败"))
		return
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgPurchaseItemResult, protocol.PurchaseItemResultPayload{
		Success: result.Success,
		ItemID:  result.ItemID,
		Coins:   result.Coins,
		Error:   result.Error,
	}))
}

// handleGetDailyTasks handles daily tasks list request.
func (h *Handler) handleGetDailyTasks(client types.ClientInterface) {
	tasks, nextReset, err := h.leaderboard.GetDailyTasksStatus(context.Background(), client.GetID(), client.GetName())
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "获取任务失败"))
		return
	}

	protoTasks := make([]protocol.DailyTaskInfo, len(tasks))
	for i, t := range tasks {
		protoTasks[i] = protocol.DailyTaskInfo{
			ID:          t.Def.ID,
			Name:        t.Def.Name,
			Description: t.Def.Description,
			Requirement: t.Def.Requirement,
			Progress:    t.Progress,
			Completed:   t.Completed,
			Claimed:     t.Claimed,
			RewardCoins: t.Def.RewardCoins,
			RewardScore: t.Def.RewardScore,
		}
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgDailyTasksResult, protocol.DailyTasksResultPayload{
		Tasks:     protoTasks,
		NextReset: nextReset,
	}))
}

// handleClaimDailyTask handles daily task reward claim request.
func (h *Handler) handleClaimDailyTask(client types.ClientInterface, msg *protocol.Message) {
	payload, err := codec.ParsePayload[protocol.ClaimDailyTaskPayload](msg)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeInvalidMsg, "请求格式错误"))
		return
	}

	result, err := h.leaderboard.ClaimDailyTask(context.Background(), client.GetID(), client.GetName(), payload.TaskID)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "领取失败"))
		return
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgClaimDailyTaskResult, protocol.ClaimDailyTaskResultPayload{
		Success:     result.Success,
		TaskID:      result.TaskID,
		RewardCoins: result.RewardCoins,
		RewardScore: result.RewardScore,
		Coins:       result.Coins,
		Error:       result.Error,
	}))
}
