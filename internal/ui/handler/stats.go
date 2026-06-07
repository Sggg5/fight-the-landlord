package handler

import (
	"fmt"
	tea "charm.land/bubbletea/v2"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	payloadconv "github.com/palemoky/fight-the-landlord/internal/protocol/convert/payload"
	"github.com/palemoky/fight-the-landlord/internal/ui/model"
)

func handleMsgStatsResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.StatsResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	m.Lobby().SetMyStats(&payload)
	return nil
}

func handleMsgLeaderboardResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.LeaderboardResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	m.Lobby().SetLeaderboard(payload.Entries)
	return nil
}


func handleMsgSignInResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.SignInResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)

	if payload.SignedIn {
		m.SetNotification(model.NotifyInfo, 
			fmt.Sprintf("? ???????%d???? %d ??", payload.ConsecutiveDays, payload.Reward), true)
	} else {
		// Already signed in - just show we checked
		m.Lobby().SetSignInInfo(payload.ConsecutiveDays, false, payload.Coins)
	}
	return nil
}

func handleMsgAchievementsResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.AchievementsResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	m.Lobby().SetAchievements(payload.Achievements)
	return nil
}

func handleMsgShopListResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.ShopListResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	m.Lobby().SetShopItems(payload.Items)
	return nil
}

func handleMsgPurchaseItemResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.PurchaseItemResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	if payload.Success {
		m.SetNotification(model.NotifyInfo, "购买成功！", true)
	} else {
		m.SetNotification(model.NotifyError, "购买失败: " + payload.Error, true)
	}
	return nil
}

func handleMsgDailyTasksResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.DailyTasksResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	m.Lobby().SetDailyTasks(payload.Tasks, payload.NextReset)
	return nil
}

func handleMsgClaimDailyTaskResult(m model.Model, msg *protocol.Message) tea.Cmd {
	var payload protocol.ClaimDailyTaskResultPayload
	_ = payloadconv.DecodePayload(msg.Type, msg.Payload, &payload)
	if payload.Success {
		m.SetNotification(model.NotifyInfo, "领取成功！", true)
	} else {
		m.SetNotification(model.NotifyError, "领取失败: " + payload.Error, true)
	}
	return nil
}
