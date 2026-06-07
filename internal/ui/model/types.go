// Package model defines the core types and interfaces for the UI.
package model

import (
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"

	gameClient "github.com/palemoky/fight-the-landlord/internal/client"
	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/transport"
)

// GamePhase represents the current game phase.
type GamePhase int

const (
	PhaseConnecting GamePhase = iota
	PhaseReconnecting
	PhaseLobby
	PhaseRoomList
	PhaseMatching
	PhaseWaiting
	PhaseBidding
	PhasePlaying
	PhaseGameOver
	PhaseLeaderboard
	PhaseStats
	PhaseAchievements
	PhaseRules
)

// NotificationType represents types of system notifications.
type NotificationType int

const (
	NotifyError            NotificationType = iota // éè¯¯ä¿¡æ¯ï¼ä¸´æ¶ï¼
	NotifyRateLimit                                // éé¢æç¤ºï¼ä¸´æ¶ï¼
	NotifyReconnecting                             // éè¿ä¸­ï¼æä¹ï¼?
	NotifyReconnectSuccess                         // éè¿æåï¼ä¸´æ¶ï¼
	NotifyMaintenance                              // ç»´æ¤éç¥ï¼æä¹ï¼
	NotifyOnlineCount                              // å¨çº¿äººæ°ï¼æä¹ï¼
	NotifyInfo                                     // ä¸è¬ä¿¡æ¯æç¤ºï¼ä¸´æ¶ï¼?
)

// SystemNotification represents a system notification.
type SystemNotification struct {
	Message   string
	Type      NotificationType
	Temporary bool // æ¯å¦ä¸ºä¸´æ¶éç¥ï¼?ç§åèªå¨æ¶å¤±ï¼?
}

// --- Tea Messages ---

// ServerMessage wraps a protocol message for tea.Msg.
type ServerMessage struct {
	Msg *protocol.Message
}

// ConnectedMsg indicates successful connection.
type ConnectedMsg struct{}

// ConnectionErrorMsg indicates a connection error.
type ConnectionErrorMsg struct {
	Err error
}

// ReconnectingMsg indicates reconnection in progress.
type ReconnectingMsg struct {
	Attempt  int
	MaxTries int
}

// ReconnectSuccessMsg indicates successful reconnection.
type ReconnectSuccessMsg struct{}

// ClearReconnectMsg clears reconnection message.
type ClearReconnectMsg struct{}

// ClearErrorMsg clears error message.
type ClearErrorMsg struct{}

// ClearInputErrorMsg clears input error message.
type ClearInputErrorMsg struct{}

// ClearSystemNotificationMsg clears system notification.
type ClearSystemNotificationMsg struct{}

// GameOverDelayMsg triggers the transition to PhaseGameOver after a delay.
type GameOverDelayMsg struct{}

// --- Model Interface ---

// Model is the main interface for OnlineModel, used by handler/view/input packages.
type Model interface {
	// Phase management
	Phase() GamePhase
	SetPhase(GamePhase)

	// Player info
	PlayerID() string
	PlayerName() string
	SetPlayerInfo(id, name string)

	// Client access
	Client() *transport.Client

	// UI components
	Input() *textinput.Model
	Timer() *timer.Model
	SetTimer(timer.Model)

	// Sub-models
	Lobby() LobbyAccessor
	Game() GameAccessor

	// Notification management
	SetNotification(notifyType NotificationType, message string, temporary bool)
	ClearNotification(notifyType NotificationType)
	GetCurrentNotification() *SystemNotification

	// State management
	EnterLobby()
	IsMaintenanceMode() bool
	SetMaintenanceMode(bool)

	// Matching
	MatchingStartTime() time.Time
	SetMatchingStartTime(time.Time)

	// Sound
	PlaySound(name string)
	PlaySequence(names ...string)
	PlayBGM(name string)
	PlayBGMAnyOf(names ...string)
	StopBGM()
	ToggleMute() bool
	Muted() bool

	// Dimensions
	Width() int
	Height() int
}

// LobbyAccessor provides access to lobby data.
type LobbyAccessor interface {
	// Data
	OnlineCount() int
	SetOnlineCount(int)
	AvailableRooms() []protocol.RoomListItem
	SetAvailableRooms([]protocol.RoomListItem)
	SelectedRoomIdx() int
	SetSelectedRoomIdx(int)
	SelectedIndex() int
	SetSelectedIndex(int)
	Leaderboard() []protocol.LeaderboardEntry
	SetLeaderboard([]protocol.LeaderboardEntry)
	MyStats() *protocol.StatsResultPayload
	SetMyStats(*protocol.StatsResultPayload)

	// Sign-in
	SignInInfo() (consecutive int, canSignIn bool)
	SetSignInInfo(consecutive int, canSignIn bool, coins int)

	// Achievements
	Achievements() []protocol.AchievementInfo
	SetAchievements([]protocol.AchievementInfo)

	// Shop
	ShopItems() []protocol.ShopItem
	SetShopItems([]protocol.ShopItem)

	// Daily Tasks
	DailyTasks() ([]protocol.DailyTaskInfo, int64)
	SetDailyTasks([]protocol.DailyTaskInfo, int64)

	// Chat
	ChatHistory() []string
	AddChatMessage(string)
	ChatInput() *textinput.Model

	// Navigation
	HandleUpKey(phase GamePhase)
	HandleDownKey(phase GamePhase)

	// Dimensions
	Width() int
	Height() int
}

// GameAccessor provides access to game data.
type GameAccessor interface {
	// State - uses the existing client.GameState
	State() *gameClient.GameState

	// Bidding
	BidTurn() string
	SetBidTurn(string)

	// Turn indicators (setters only - getters unused)
	SetMustPlay(bool)
	SetCanBeat(bool)

	// Timer
	TimerDuration() time.Duration
	SetTimerDuration(time.Duration)
	TimerStartTime() time.Time
	SetTimerStartTime(time.Time)

	// Bell (setter only - getter unused)
	SetBellPlayed(bool)

	// Features
	CardCounterEnabled() bool
	SetCardCounterEnabled(bool)
	ShowingHelp() bool
	SetShowingHelp(bool)

	// Chat
	ChatHistory() []string
	AddChatMessage(string)
	ShowQuickMsgMenu() bool
	SetShowQuickMsgMenu(bool)
	QuickMsgInput() string
	SetQuickMsgInput(string)
	AppendQuickMsgInput(rune)
	ClearQuickMsgInput()
	QuickMsgScroll() int
	SetQuickMsgScroll(int)
}

// --- Handler Interface ---

// Handler processes server messages.
type Handler interface {
	HandleServerMessage(m Model, msg *protocol.Message) tea.Cmd
}

// InputHandler processes keyboard input.
type InputHandler interface {
	HandleKeyPress(m Model, msg tea.KeyMsg) (handled bool, cmd tea.Cmd)
}
