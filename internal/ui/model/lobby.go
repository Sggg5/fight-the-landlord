// Package model contains the UI model implementations.
package model

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/transport"
)

const (
	// chatBoxWidth is the width of the chat box.
	chatBoxWidth = 50
	// chatInputWidth is the width of the chat input.
	chatInputWidth = chatBoxWidth - 5
)

// LobbyModel handles the lobby interface.
type LobbyModel struct {
	client *transport.Client
	width  int
	height int

	// Navigation
	selectedIndex int

	// Data
	onlineCount     int
	availableRooms  []protocol.RoomListItem
	selectedRoomIdx int
	leaderboard     []protocol.LeaderboardEntry
	myStats         *protocol.StatsResultPayload

	// Chat
	chatHistory []string
	chatInput   textinput.Model

	// Input reference
	input *textinput.Model
}

// NewLobbyModel creates a new LobbyModel.
func NewLobbyModel(c *transport.Client, input *textinput.Model) *LobbyModel {
	chatInput := textinput.New()
	chatInput.Placeholder = "ENTER è¾å¥èå¤© / ââ éæ©èå"
	chatInput.CharLimit = 50
	chatInput.SetWidth(chatInputWidth)

	return &LobbyModel{
		client:    c,
		input:     input,
		chatInput: chatInput,
	}
}

func (m *LobbyModel) Init() tea.Cmd {
	return nil
}

func (m *LobbyModel) View() tea.View {
	return tea.NewView("") // Not used directly, managed by OnlineModel
}

func (m *LobbyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// --- LobbyAccessor implementation ---

func (m *LobbyModel) OnlineCount() int                        { return m.onlineCount }
func (m *LobbyModel) SetOnlineCount(count int)                { m.onlineCount = count }
func (m *LobbyModel) AvailableRooms() []protocol.RoomListItem { return m.availableRooms }
func (m *LobbyModel) SetAvailableRooms(rooms []protocol.RoomListItem) {
	m.availableRooms = rooms
	m.selectedRoomIdx = 0
}
func (m *LobbyModel) SelectedRoomIdx() int                               { return m.selectedRoomIdx }
func (m *LobbyModel) SetSelectedRoomIdx(idx int)                         { m.selectedRoomIdx = idx }
func (m *LobbyModel) Leaderboard() []protocol.LeaderboardEntry           { return m.leaderboard }
func (m *LobbyModel) SetLeaderboard(entries []protocol.LeaderboardEntry) { m.leaderboard = entries }
func (m *LobbyModel) MyStats() *protocol.StatsResultPayload              { return m.myStats }
func (m *LobbyModel) SetMyStats(stats *protocol.StatsResultPayload)      { m.myStats = stats }

func (m *LobbyModel) ChatHistory() []string { return m.chatHistory }
func (m *LobbyModel) AddChatMessage(msg string) {
	m.chatHistory = append(m.chatHistory, msg)
	if len(m.chatHistory) > 50 {
		m.chatHistory = m.chatHistory[len(m.chatHistory)-50:]
	}
}
func (m *LobbyModel) ChatInput() *textinput.Model { return &m.chatInput }

// HandleNavigationKey å¤çä¸ä¸é®å¯¼è?
// direction: -1 è¡¨ç¤ºåä¸ï¼? è¡¨ç¤ºåä¸
func (m *LobbyModel) HandleNavigationKey(phase GamePhase, direction int) {
	switch phase {
	case PhaseRoomList:
		if len(m.availableRooms) > 0 {
			m.selectedRoomIdx += direction
			if m.selectedRoomIdx < 0 {
				m.selectedRoomIdx = len(m.availableRooms) - 1
			} else if m.selectedRoomIdx >= len(m.availableRooms) {
				m.selectedRoomIdx = 0
			}
		}
	case PhaseLobby:
		m.selectedIndex += direction
		if m.selectedIndex < 0 {
			m.selectedIndex = 6
		} else if m.selectedIndex > 7 {
			m.selectedIndex = 0
		}
	}
}

// HandleUpKey å¤çåä¸é?
func (m *LobbyModel) HandleUpKey(phase GamePhase) {
	m.HandleNavigationKey(phase, -1)
}

// HandleDownKey å¤çåä¸é?
func (m *LobbyModel) HandleDownKey(phase GamePhase) {
	m.HandleNavigationKey(phase, 1)
}

func (m *LobbyModel) Width() int  { return m.width }
func (m *LobbyModel) Height() int { return m.height }
func (m *LobbyModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
func (m *LobbyModel) Input() *textinput.Model   { return m.input }
func (m *LobbyModel) SelectedIndex() int        { return m.selectedIndex }
func (m *LobbyModel) SetSelectedIndex(idx int)  { m.selectedIndex = idx }
func (m *LobbyModel) Client() *transport.Client { return m.client }


// --- Sign-in ---
func (m *LobbyModel) SignInInfo() (int, bool) { return m.signInConsecutive, m.signInCanSignIn }
func (m *LobbyModel) SetSignInInfo(consecutive int, canSignIn bool, coins int) {
	m.signInConsecutive = consecutive
	m.signInCanSignIn = canSignIn
}

// --- Achievements ---
func (m *LobbyModel) Achievements() []protocol.AchievementInfo { return m.achievements }
func (m *LobbyModel) SetAchievements(a []protocol.AchievementInfo) { m.achievements = a }

// --- Shop ---
func (m *LobbyModel) ShopItems() []protocol.ShopItem { return m.shopItems }
func (m *LobbyModel) SetShopItems(items []protocol.ShopItem) { m.shopItems = items }

// --- Daily Tasks ---
func (m *LobbyModel) DailyTasks() ([]protocol.DailyTaskInfo, int64) { return m.dailyTasks, m.dailyTasksReset }
func (m *LobbyModel) SetDailyTasks(tasks []protocol.DailyTaskInfo, reset int64) {
	m.dailyTasks = tasks
	m.dailyTasksReset = reset
}

