// Package view provides UI rendering functions.
package view

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/ui/common"
	"github.com/palemoky/fight-the-landlord/internal/ui/model"
)

// getNotificationStyle returns the appropriate style for a notification type.
func getNotificationStyle(notifyType model.NotificationType) lipgloss.Style {
	switch notifyType {
	case model.NotifyError, model.NotifyRateLimit:
		return retroWarningStyle
	case model.NotifyReconnecting:
		return retroWarningStyle
	case model.NotifyReconnectSuccess:
		return retroOnlineStyle
	case model.NotifyMaintenance:
		return retroWarningStyle
	case model.NotifyOnlineCount:
		return retroOnlineStyle
	case model.NotifyInfo:
		return retroBrightStyle
	default:
		return retroTextStyle
	}
}

var (
	retroBgColor      = lipgloss.Color("#06110b")
	retroPanelColor   = lipgloss.Color("#07170d")
	retroBorderColor  = lipgloss.Color("#1f8f48")
	retroTextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6fa87d"))
	retroDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#355d42"))
	retroBrightStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#8cff8c")).Bold(true)
	retroOnlineStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#36ff62")).Bold(true)
	retroWarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd75f")).Bold(true)
	retroGlowStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#b6ff9d")).Bold(true)
	retroPanelStyle   = lipgloss.NewStyle().
				Background(retroPanelColor).
				Border(lipgloss.NormalBorder()).
				BorderForeground(retroBorderColor).
				Padding(1, 2)
)

// renderChatBox renders the chat box component for the lobby.
func renderChatBox(lobby model.LobbyAccessor, height int) string {
	var chatLines []string
	if len(lobby.ChatHistory()) > 0 {
		history := lobby.ChatHistory()
		count := len(history)
		start := 0
		if count > 5 {
			start = count - 5
		}
		for i := start; i < count; i++ {
			chatLines = append(chatLines, history[i])
		}
	} else {
		chatLines = []string{
			retroWarningStyle.Render("[系统] ") + retroTextStyle.Render("玩家 老王 进入大厅"),
			retroWarningStyle.Render("[系统] ") + retroTextStyle.Render("正在等待匹配..."),
			retroOnlineStyle.Render("[玩家] ") + retroTextStyle.Render("来一把经典场"),
		}
	}

	chatInputView := lobby.ChatInput().View()
	if !lobby.ChatInput().Focused() {
		chatInputView = retroDimStyle.Render("ENTER 输入聊天 / ↑↓ 选择菜单 ") + retroBrightStyle.Render("█")
	}

	chatHeader := retroGlowStyle.Render("WORLD CHAT")
	innerHeight := height - 2
	usedLines := 1 + len(chatLines) + 1 // header + chat + input
	emptyLines := max(innerHeight-usedLines, 0)
	contentLines := make([]string, 0, 1+len(chatLines)+emptyLines+1)
	contentLines = append(contentLines, chatHeader)
	contentLines = append(contentLines, chatLines...)
	for range emptyLines {
		contentLines = append(contentLines, "")
	}
	contentLines = append(contentLines, chatInputView)

	chatBoxWidth := 50
	chatBoxContent := lipgloss.JoinVertical(lipgloss.Left, contentLines...)
	return retroPanelStyle.Width(chatBoxWidth).Height(innerHeight).Render(chatBoxContent)
}

// LobbyView renders the lobby view.
func LobbyView(m model.Model) string {
	lobby := m.Lobby()
	var sb strings.Builder

	screenWidth := max(min(m.Width()-4, 108), 84)
	screenHeight := max(min(m.Height()-4, 32), 24)
	playerName := m.PlayerName()
	if playerName == "" {
		playerName = "GUEST"
	}

	title := retroGlowStyle.Render("╔════════════════════════════════╗\n║        欢乐斗地主 DDZ         ║\n╚════════════════════════════════╝")
	sb.WriteString(lipgloss.PlaceHorizontal(screenWidth, lipgloss.Center, title))
	sb.WriteString("\n")

	topBar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		retroTextStyle.Render("PLAYER: "),
		retroBrightStyle.Render(playerName),
		retroDimStyle.Render("  |  "),
		retroTextStyle.Render("ONLINE: "),
		retroOnlineStyle.Render(fmt.Sprintf("%03d", lobby.OnlineCount())),
		retroDimStyle.Render("  |  "),
		retroWarningStyle.Render("DOS/CRT MODE"),
	)
	sb.WriteString(lipgloss.PlaceHorizontal(screenWidth, lipgloss.Center, topBar))
	sb.WriteString("\n")

	if notification := m.GetCurrentNotification(); notification != nil {
		notificationStyle := getNotificationStyle(notification.Type)
		sb.WriteString(lipgloss.PlaceHorizontal(screenWidth, lipgloss.Center,
			notificationStyle.Render(notification.Message)))
		sb.WriteString("\n")
	}
	sb.WriteString(retroDimStyle.Render(strings.Repeat("░", screenWidth)))
	sb.WriteString("\n")

	menuItems := []string{
		"[01] 快速匹配",
		"[02] 创建房间",
		"[03] 加入房间",
		"[04] 人机练习",
		"[05] 排行榜",
		"[06] 我的战绩",
		"[07] 游戏规则",
	}

	lobbyModel := m.Lobby()
	menuLines := make([]string, 0, 2+len(menuItems))
	menuLines = append(menuLines, retroGlowStyle.Render("MAIN MENU"), "")
	for i, item := range menuItems {
		prefix := "  "
		style := retroTextStyle
		if i == lobbyModel.SelectedIndex() {
			prefix = "> "
			style = retroBrightStyle
		}
		menuLines = append(menuLines, style.Render(prefix+item))
	}
	menuLines = append(menuLines, "", retroDimStyle.Render("SCANLINE: ACTIVE"))

	menu := retroPanelStyle.Width(30).Render(lipgloss.JoinVertical(lipgloss.Left, menuLines...))
	menuHeight := lipgloss.Height(menu)

	// Chat box
	chatBox := renderChatBox(lobby, menuHeight)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, menu, "  ", chatBox)
	sb.WriteString(lipgloss.PlaceHorizontal(screenWidth, lipgloss.Center, mainContent))
	sb.WriteString("\n")

	// Only show blinking cursor on lobby input when chat is not focused
	var inputView string
	if lobby.ChatInput().Focused() {
		m.Input().Blur()
		inputView = retroWarningStyle.Render("ENTER 输入聊天 / ↑↓ 选择菜单")
	} else {
		m.Input().Focus()
		m.Input().Placeholder = "ENTER 确认 / ↑↓ 选择菜单 / 输入房间号"
		inputView = retroTextStyle.Render(m.Input().View()) + retroBrightStyle.Render(" █")
	}
	footer := lipgloss.NewStyle().
		Background(lipgloss.Color("#0a1e10")).
		Foreground(lipgloss.Color("#c7ff7a")).
		Width(screenWidth-4).
		Padding(0, 2).
		Render(inputView)
	sb.WriteString(lipgloss.PlaceHorizontal(screenWidth, lipgloss.Center, footer))

	sb.WriteString("\n")
	vignette := retroDimStyle.Render("╚" + strings.Repeat("═", screenWidth-2) + "╝")
	sb.WriteString(vignette)

	content := lipgloss.NewStyle().
		Width(screenWidth).
		Height(screenHeight).
		Background(retroBgColor).
		Foreground(lipgloss.Color("#6fa87d")).
		Padding(1, 2).
		Render(sb.String())
	return lipgloss.Place(m.Width(), m.Height(), lipgloss.Center, lipgloss.Center, content)
}

// RoomListView renders the room list view.
func RoomListView(m model.Model) string {
	lobby := m.Lobby()
	var sb strings.Builder

	title := common.TitleStyle("📋 可加入的房间")
	sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, title))
	sb.WriteString("\n\n")

	rooms := lobby.AvailableRooms()
	if len(rooms) == 0 {
		noRooms := "暂无可加入的房间\n\n按 ESC 返回大厅"
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, noRooms))
	} else {
		var roomList strings.Builder
		roomList.WriteString("房间列表:\n\n")

		for i, room := range rooms {
			prefix := "  "
			if i == lobby.SelectedRoomIdx() {
				prefix = "▶ "
			}
			fmt.Fprintf(&roomList, "%s房间 %s  (%d/3)\n", prefix, room.RoomCode, room.PlayerCount)
		}

		roomList.WriteString("\n↑↓ 选择  回车加入  ESC 返回")

		roomBox := common.BoxStyle.Render(roomList.String())
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, roomBox))
		sb.WriteString("\n\n")
	}

	inputView := lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, m.Input().View())
	sb.WriteString(inputView)

	return sb.String()
}

// LeaderboardView renders the leaderboard view.
func LeaderboardView(m model.Model) string {
	lobby := m.Lobby()
	var sb strings.Builder

	title := common.TitleStyle("🏆 排行榜")
	sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, title))
	sb.WriteString("\n\n")

	entries := lobby.Leaderboard()
	if len(entries) > 0 {
		leaderboard := renderLeaderboardTable(entries)
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, leaderboard))
	} else {
		noData := "正在加载排行榜..."
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, noData))
	}

	sb.WriteString("\n\n")
	hint := "按 ESC 返回大厅"
	sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, hint))

	return sb.String()
}

func renderLeaderboardTable(entries []protocol.LeaderboardEntry) string {
	var sb strings.Builder

	title := "🏆 排行榜 TOP 10"
	titleLine := lipgloss.PlaceHorizontal(50, lipgloss.Center, title)
	sb.WriteString(titleLine + "\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n")

	sb.WriteString("排名\t玩家\t\t积分\t胜场\t胜率\n")
	sb.WriteString(strings.Repeat("─", 50) + "\n")

	for _, e := range entries {
		rankStr := fmt.Sprintf("%2d.", e.Rank)
		fmt.Fprintf(&sb, "%s\t%s\t\t%d\t%d\t%.1f%%\n",
			rankStr, common.TruncateName(e.PlayerName, 10), e.Score, e.Wins, e.WinRate)
	}

	return common.BoxStyle.Render(sb.String())
}

// StatsView renders the stats view.
func StatsView(m model.Model) string {
	lobby := m.Lobby()
	var sb strings.Builder

	title := common.TitleStyle("📊 我的战绩")
	sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, title))
	sb.WriteString("\n\n")

	stats := lobby.MyStats()
	if stats != nil && stats.TotalGames > 0 {
		statsTable := renderStatsTable(stats)
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, statsTable))
	} else {
		noData := "暂无战绩数据"
		sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, noData))
	}

	sb.WriteString("\n\n")
	hint := "按 ESC 返回大厅"
	sb.WriteString(lipgloss.PlaceHorizontal(m.Width(), lipgloss.Center, hint))

	return sb.String()
}

func renderStatsTable(s *protocol.StatsResultPayload) string {
	var sb strings.Builder
	sb.WriteString("📊 我的战绩\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	rankStr := "未上榜"
	if s.Rank > 0 {
		rankStr = fmt.Sprintf("#%d", s.Rank)
	}
	fmt.Fprintf(&sb, "排名: %s  |  积分: %d\n", rankStr, s.Score)
	sb.WriteString(strings.Repeat("─", 40) + "\n")

	fmt.Fprintf(&sb, "总场次: %d  胜: %d  负: %d  胜率: %.1f%%\n",
		s.TotalGames, s.Wins, s.Losses, s.WinRate)

	landlordRate := 0.0
	if s.LandlordGames > 0 {
		landlordRate = float64(s.LandlordWins) / float64(s.LandlordGames) * 100
	}
	farmerRate := 0.0
	if s.FarmerGames > 0 {
		farmerRate = float64(s.FarmerWins) / float64(s.FarmerGames) * 100
	}

	fmt.Fprintf(&sb, "地主: %d胜/%d场 (%.1f%%)  |  农民: %d胜/%d场 (%.1f%%)\n",
		s.LandlordWins, s.LandlordGames, landlordRate,
		s.FarmerWins, s.FarmerGames, farmerRate)

	streakStr := ""
	if s.CurrentStreak > 0 {
		streakStr = fmt.Sprintf("🔥 %d 连胜!", s.CurrentStreak)
	} else if s.CurrentStreak < 0 {
		streakStr = fmt.Sprintf("💔 %d 连败", -s.CurrentStreak)
	}
	if s.MaxWinStreak > 0 {
		streakStr += fmt.Sprintf("  最高连胜: %d", s.MaxWinStreak)
	}
	if streakStr != "" {
		sb.WriteString(streakStr + "\n")
	}

	return common.BoxStyle.Render(sb.String())
}
