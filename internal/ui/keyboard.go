// Package input handles keyboard input processing.
package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/palemoky/fight-the-landlord/internal/game/card"
	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	"github.com/palemoky/fight-the-landlord/internal/protocol/convert"
	"github.com/palemoky/fight-the-landlord/internal/ui/model"
	"github.com/palemoky/fight-the-landlord/internal/ui/view"
)

func clearSystemNotification() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return model.ClearSystemNotificationMsg{}
	})
}

// sendChatMessage sends a chat message and returns error command if failed
func sendChatMessage(m model.Model, content, scope string) tea.Cmd {
	chatMsg := codec.MustNewMessage(protocol.MsgChat, protocol.ChatPayload{
		Content: content,
		Scope:   scope,
	})
	if err := m.Client().SendMessage(chatMsg); err != nil {
		m.SetNotification(model.NotifyError, fmt.Sprintf("â ï¸ åéæ¶æ¯å¤±è´? %v", err), true)
		return clearSystemNotification()
	}
	return nil
}

// handleLobbyChatInput handles chat input in the lobby
// Returns (handled, cmd) where handled indicates if the key was processed
func handleLobbyChatInput(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.Phase() != model.PhaseLobby {
		return false, nil
	}

	chatInput := m.Lobby().ChatInput()

	// "/" key focuses chat input
	if !chatInput.Focused() {
		if msg.String() == "/" {
			chatInput.Focus()
			return true, nil
		}
		return false, nil
	}

	// Chat input is focused - handle input
	switch msg.Key().Code {
	case tea.KeyEnter:
		if content := chatInput.Value(); content != "" {
			if cmd := sendChatMessage(m, content, "lobby"); cmd != nil {
				return true, cmd
			}
			chatInput.SetValue("")
		}
		return true, nil
	case tea.KeyEsc:
		chatInput.Blur()
		return true, nil
	default:
		var cmd tea.Cmd
		*chatInput, cmd = chatInput.Update(msg)
		return true, cmd
	}
}

// handleQuickMessageMenu handles the quick message menu in-game
func handleQuickMessageMenu(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if m.Phase() != model.PhaseBidding && m.Phase() != model.PhasePlaying {
		return false, nil
	}

	// Toggle quick message menu with 'T' key
	if tryToggleQuickMenu(m, msg) {
		return true, nil
	}

	// Handle menu interactions
	if !m.Game().ShowQuickMsgMenu() {
		return false, nil
	}

	return processQuickMenuKey(m, msg)
}

func tryToggleQuickMenu(m model.Model, msg tea.KeyMsg) bool {
	if msg.String() == "t" || msg.String() == "T" {
		if !m.Game().ShowQuickMsgMenu() {
			m.Game().SetShowQuickMsgMenu(true)
			return true
		}
	}
	return false
}

func processQuickMenuKey(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	switch msg.Key().Code {
	case tea.KeyEsc:
		m.Game().SetShowQuickMsgMenu(false)
		return true, nil
	case tea.KeyUp, tea.KeyDown:
		return handleQuickMsgScroll(m, msg)
	case tea.KeyEnter:
		return handleQuickMsgEnter(m)
	case tea.KeyBackspace:
		return handleQuickMsgInput(m, msg)
	default:
		if msg.Key().Text != "" {
			return handleQuickMsgInput(m, msg)
		}
	}
	return true, nil
}

func handleQuickMsgScroll(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	scroll := m.Game().QuickMsgScroll()
	if msg.Key().Code == tea.KeyUp {
		if scroll > 0 {
			m.Game().SetQuickMsgScroll(scroll - 1)
		}
	} else {
		maxScroll := max(len(view.QuickMessages)-10, 0)
		if scroll < maxScroll {
			m.Game().SetQuickMsgScroll(scroll + 1)
		}
	}
	return true, nil
}

func handleQuickMsgEnter(m model.Model) (bool, tea.Cmd) {
	input := m.Game().QuickMsgInput()
	if input != "" {
		idx := 0
		for _, c := range input {
			idx = idx*10 + int(c-'0')
		}
		idx-- // Convert to 0-indexed
		if idx >= 0 && idx < len(view.QuickMessages) {
			if cmd := sendChatMessage(m, view.QuickMessages[idx], "room"); cmd != nil {
				return true, cmd
			}
			m.Game().SetShowQuickMsgMenu(false)
			return true, nil
		}
	}
	m.Game().ClearQuickMsgInput()
	return true, nil
}

func handleQuickMsgInput(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	if msg.Key().Code == tea.KeyBackspace {
		input := m.Game().QuickMsgInput()
		if input != "" {
			m.Game().SetQuickMsgInput(input[:len(input)-1])
		}
		return true, nil
	}

	if msg.String() == "t" || msg.String() == "T" {
		m.Game().SetShowQuickMsgMenu(false)
		return true, nil
	}
	// Accumulate digits for message selection
	runes := []rune(msg.Key().Text)
	if len(runes) == 1 && runes[0] >= '0' && runes[0] <= '9' {
		input := m.Game().QuickMsgInput()
		if len(input) < 2 {
			m.Game().AppendQuickMsgInput(runes[0])
		}
	}
	return true, nil
}

// HandleKeyPress handles keyboard input and returns whether it was handled.
func HandleKeyPress(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	// Try lobby chat handling first
	if handled, cmd := handleLobbyChatInput(m, msg); handled {
		return true, cmd
	}

	// Try quick message menu handling
	if handled, cmd := handleQuickMessageMenu(m, msg); handled {
		return true, cmd
	}

	// General key handling
	switch msg.Key().Code {
	case tea.KeyEsc:
		return handleEscKey(m)
	case tea.KeyUp:
		playMenuFeedback(m)
		m.Lobby().HandleUpKey(m.Phase())
		return false, nil
	case tea.KeyDown:
		playMenuFeedback(m)
		m.Lobby().HandleDownKey(m.Phase())
		return false, nil
	case tea.KeyEnter:
		playMenuFeedback(m)
		cmd := handleEnter(m)
		return false, cmd
	default:
		if msg.String() == "ctrl+c" {
			return handleEscKey(m)
		}
		if msg.Key().Text != "" {
			return handleRuneKey(m, msg)
		}
	}
	return false, nil
}

// playMenuFeedback å¨å¤§å?/ æ¿é´åè¡¨ç¨ä¸ä¸é®å¯¼èªæåè½¦éæ©æ¶ç»åºæé®é³åé¦
func playMenuFeedback(m model.Model) {
	switch m.Phase() {
	case model.PhaseLobby, model.PhaseRoomList:
		m.PlaySound("menu")
	}
}

func handleEscKey(m model.Model) (bool, tea.Cmd) {
	if m.Game().ShowingHelp() {
		m.Game().SetShowingHelp(false)
		return true, nil
	}

	switch m.Phase() {
	case model.PhaseRoomList, model.PhaseMatching, model.PhaseLeaderboard, model.PhaseStats, model.PhaseAchievements, model.PhaseRules, model.PhaseGameOver:
		m.EnterLobby()
		return true, nil
	case model.PhaseWaiting:
		_ = m.Client().LeaveRoom()
		m.EnterLobby()
		return true, nil
	case model.PhaseBidding, model.PhasePlaying:
		m.SetNotification(model.NotifyError, "â ï¸ æ¸¸æè¿è¡ä¸­ï¼æ æ³éåºï¼", true)
		return true, clearSystemNotification()
	}

	m.Client().Close()
	return true, tea.Quit
}

func handleRuneKey(m model.Model, msg tea.KeyMsg) (bool, tea.Cmd) {
	runes := []rune(msg.Key().Text)
	if len(runes) == 0 {
		return false, nil
	}

	if m.Phase() == model.PhaseGameOver && (runes[0] == 'r' || runes[0] == 'R') {
		return true, handleGameOverPracticeRestart(m)
	}

	// Mute toggle works in any phase
	if m.Phase() == model.PhaseLobby && (runes[0] == 's' || runes[0] == 'S') {
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgSignIn, nil))
		m.SetNotification(model.NotifyInfo, "???...", true)
		return true, clearSystemNotification()
	}

	if runes[0] == 'm' || runes[0] == 'M' {
		if muted := m.ToggleMute(); muted {
			m.SetNotification(model.NotifyInfo, "ð å·²éé?, true)
		} else {
			m.SetNotification(model.NotifyInfo, "ð å£°é³å·²å¼å?, true)
		}
		return true, clearSystemNotification()
	}

	// Handle game toggles (only during bidding/playing)
	if m.Phase() == model.PhaseBidding || m.Phase() == model.PhasePlaying {
		switch runes[0] {
		case 'c', 'C':
			m.Game().SetCardCounterEnabled(!m.Game().CardCounterEnabled())
			return true, nil
		case 'h', 'H':
			m.Game().SetShowingHelp(!m.Game().ShowingHelp())
			return true, nil
		}
	}

	return false, nil
}

func handleEnter(m model.Model) tea.Cmd {
	input := strings.TrimSpace(m.Input().Value())
	m.Input().Reset()

	switch m.Phase() {
	case model.PhaseLobby:
		return handleLobbyEnter(m, input)
	case model.PhaseRoomList:
		return handleRoomListEnter(m, input)
	case model.PhaseWaiting:
		return handleWaitingEnter(m, input)
	case model.PhaseBidding:
		return handleBiddingEnter(m, input)
	case model.PhasePlaying:
		return handlePlayingEnter(m, input)
	case model.PhaseGameOver:
		return handleGameOverEnter(m)
	}

	return nil
}

// checkServerAvailability æ£æ¥æå¡å¨æ¯å¦å¯ç¨äºæ¸¸ææä½?
// è¿å true åéè¯¯å½ä»¤å¦ææå¡å¨ä¸å¯ç¨ï¼è¿å false å?nil å¦æå¯ç¨
func checkServerAvailability(m model.Model) (blocked bool, cmd tea.Cmd) {
	if blocked, cmd := checkMaintenanceMode(m); blocked {
		return blocked, cmd
	}
	if m.Client().IsReconnecting() {
		m.SetNotification(model.NotifyError, "â ï¸ æ­£å¨éè¿ä¸­ï¼è¯·ç¨ååè¯?, true)
		return true, clearSystemNotification()
	}
	if !m.Client().IsConnected() {
		m.SetNotification(model.NotifyError, "â ï¸ æªè¿æ¥å°æå¡å?, true)
		return true, clearSystemNotification()
	}
	return false, nil
}

// checkMaintenanceMode æ£æ¥æå¡å¨æ¯å¦å¤äºç»´æ¤æ¨¡å¼
// è¿å true åéè¯¯å½ä»¤å¦æå¨ç»´æ¤æ¨¡å¼ï¼è¿å?false å?nil å¦ææ­£å¸¸
func checkMaintenanceMode(m model.Model) (blocked bool, cmd tea.Cmd) {
	if m.IsMaintenanceMode() {
		m.SetNotification(model.NotifyError, "â ï¸ æå¡å¨ç»´æ¤ä¸­ï¼æåæ¥åæ°è¿æ¥", true)
		return true, clearSystemNotification()
	}
	return false, nil
}

func handleLobbyEnter(m model.Model, input string) tea.Cmd {
	if input == "" {
		input = fmt.Sprintf("%d", m.Lobby().SelectedIndex()+1)
	}

	switch input {
	case "1": // å¿«éå¹é?
		if blocked, cmd := checkServerAvailability(m); blocked {
			return cmd
		}
		m.SetPhase(model.PhaseMatching)
		m.SetMatchingStartTime(time.Now())
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgQuickMatch, nil))

	case "2": // åå»ºæ¿é´
		if blocked, cmd := checkMaintenanceMode(m); blocked {
			return cmd
		}
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgCreateRoom, nil))

	case "3": // æ¿é´åè¡¨
		if blocked, cmd := checkMaintenanceMode(m); blocked {
			return cmd
		}
		m.SetPhase(model.PhaseRoomList)
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgGetRoomList, nil))
		m.Input().Placeholder = "è¾å¥æ¿é´å·ææ?ESC è¿å"

	case "4": // äººæºç»ä¹ 
		if blocked, cmd := checkServerAvailability(m); blocked {
			return cmd
		}
		m.SetPhase(model.PhaseMatching)
		m.SetMatchingStartTime(time.Now())
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgPracticeMatch, nil))

	case "5": // æè¡æ¦?
		m.SetPhase(model.PhaseLeaderboard)
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgGetLeaderboard, nil))

	case "6": // ç»è®¡ä¿¡æ¯
		m.SetPhase(model.PhaseStats)
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgGetStats, nil))

	case "7": // æ¸¸æè§å
	case "7": // ÓÎÏ·¹æÔò
		m.SetPhase(model.PhaseRules)

	case "8": // ³É¾Í
		m.SetPhase(model.PhaseAchievements)
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgGetAchievements, nil))

	default: // å å¥æ¿é´
		if blocked, cmd := checkMaintenanceMode(m); blocked {
			return cmd
		}
		_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgJoinRoom, protocol.JoinRoomPayload{
			RoomCode: input,
		}))
	}

	return nil
}

func handleRoomListEnter(m model.Model, input string) tea.Cmd {
	if m.IsMaintenanceMode() {
		m.SetNotification(model.NotifyError, "â ï¸ æå¡å¨ç»´æ¤ä¸­ï¼æåå å¥æ¿é?, true)
		return clearSystemNotification()
	}
	rooms := m.Lobby().AvailableRooms()
	if input == "" {
		if len(rooms) > 0 && m.Lobby().SelectedRoomIdx() < len(rooms) {
			roomCode := rooms[m.Lobby().SelectedRoomIdx()].RoomCode
			_ = m.Client().JoinRoom(roomCode)
		}
	} else {
		_ = m.Client().JoinRoom(input)
	}
	return nil
}

func handleWaitingEnter(m model.Model, input string) tea.Cmd {
	if strings.EqualFold(input, "r") || strings.EqualFold(input, "ready") {
		_ = m.Client().Ready()
	}
	return nil
}

func handleBiddingEnter(m model.Model, input string) tea.Cmd {
	if m.Game().BidTurn() == m.PlayerID() {
		switch strings.ToLower(input) {
		case "y", "yes", "1":
			_ = m.Client().Bid(true)
		case "n", "no", "0":
			_ = m.Client().Bid(false)
		}
	}
	return nil
}

func handlePlayingEnter(m model.Model, input string) tea.Cmd {
	if m.Game().State().CurrentTurn == m.PlayerID() {
		upperInput := strings.ToUpper(input)
		if upperInput == "PASS" || upperInput == "P" {
			_ = m.Client().Pass()
		} else if input != "" {
			cards, err := card.FindCardsInHand(m.Game().State().Hand, strings.ToUpper(input))
			if err != nil {
				m.Input().Placeholder = err.Error()
				return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
					return model.ClearInputErrorMsg{}
				})
			}
			_ = m.Client().PlayCards(convert.CardsToInfos(cards))
		}
	}
	return nil
}

func handleGameOverEnter(m model.Model) tea.Cmd {
	m.EnterLobby()
	m.Game().State().Reset()

	_ = m.Client().SendMessage(codec.MustNewMessage(protocol.MsgGetMaintenanceStatus, nil))

	return nil
}

func handleGameOverPracticeRestart(m model.Model) tea.Cmd {
	if !gameHadBot(m) {
		return nil
	}
	if blocked, cmd := checkServerAvailability(m); blocked {
		return cmd
	}

	m.Game().State().Reset()
	m.SetPhase(model.PhaseMatching)
	m.SetMatchingStartTime(time.Now())
	_ = m.Client().PracticeMatch()

	return nil
}

func gameHadBot(m model.Model) bool {
	for _, player := range m.Game().State().Players {
		if player.IsBot {
			return true
		}
	}
	return false
}
