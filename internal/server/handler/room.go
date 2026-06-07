package handler

import (
	"encoding/json"
	"errors"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	"github.com/palemoky/fight-the-landlord/internal/types"

	"github.com/palemoky/fight-the-landlord/internal/apperrors"
	"github.com/palemoky/fight-the-landlord/internal/game/room"
)

type practiceMatchPayload struct {
	Difficulty string `json:"difficulty"`
}

// handleCreateRoom 处理创建房间
func (h *Handler) handleCreateRoom(client types.ClientInterface) {
	// 维护模式检查
	if h.server.IsMaintenanceMode() {
		client.SendMessage(codec.NewErrorMessageWithText(
			protocol.ErrCodeServerMaintenance, "服务器维护中，暂停创建房间"))
		return
	}

	// 如果已在房间中，先离开
	if client.GetRoom() != "" {
		h.roomManager.LeaveRoom(client)
	}

	room, err := h.roomManager.CreateRoom(client)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, err.Error()))
		return
	}

	if room == nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "创建房间失败"))
		return
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgRoomCreated, protocol.RoomCreatedPayload{
		RoomCode: room.Code,
		Player:   room.GetPlayerInfo(client.GetID()),
	}))
}

// handleJoinRoom 处理加入房间
func (h *Handler) handleJoinRoom(client types.ClientInterface, msg *protocol.Message) {
	// 维护模式检查
	if h.server.IsMaintenanceMode() {
		client.SendMessage(codec.NewErrorMessageWithText(
			protocol.ErrCodeServerMaintenance, "服务器维护中，暂停加入房间"))
		return
	}

	payload, err := codec.ParsePayload[protocol.JoinRoomPayload](msg)
	if err != nil {
		client.SendMessage(codec.NewErrorMessage(protocol.ErrCodeInvalidMsg))
		return
	}

	// 如果已在房间中，先离开
	if client.GetRoom() != "" {
		h.roomManager.LeaveRoom(client)
	}

	room, err := h.roomManager.JoinRoom(client, payload.RoomCode)
	if err != nil {
		var gameErr *apperrors.GameError
		if errors.As(err, &gameErr) {
			client.SendMessage(codec.NewErrorMessage(gameErr.Code))
		} else {
			client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, err.Error()))
		}
		return
	}

	if room == nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, "加入房间失败"))
		return
	}

	client.SendMessage(codec.MustNewMessage(protocol.MsgRoomJoined, protocol.RoomJoinedPayload{
		RoomCode: room.Code,
		Player:   room.GetPlayerInfo(client.GetID()),
		Players:  room.GetAllPlayersInfo(),
	}))
}

// handleLeaveRoom 处理离开房间
func (h *Handler) handleLeaveRoom(client types.ClientInterface) {
	h.roomManager.LeaveRoom(client)
}

// handleQuickMatch 处理快速匹配
func (h *Handler) handleQuickMatch(client types.ClientInterface) {
	// 维护模式检查
	if h.server.IsMaintenanceMode() {
		client.SendMessage(codec.NewErrorMessageWithText(
			protocol.ErrCodeServerMaintenance, "服务器维护中，暂停快速匹配"))
		return
	}

	// 如果已在房间中，先离开
	if client.GetRoom() != "" {
		h.roomManager.LeaveRoom(client)
	}

	h.matcher.AddToQueue(client)
}

// handlePracticeMatch 处理人机练习
func (h *Handler) handlePracticeMatch(client types.ClientInterface, msg *protocol.Message) {
	if h.server.IsMaintenanceMode() {
		client.SendMessage(codec.NewErrorMessageWithText(
			protocol.ErrCodeServerMaintenance, "服务器维护中，暂停人机练习"))
		return
	}
	difficulty := parsePracticeDifficulty(msg)

	if client.GetRoom() != "" {
		r := h.roomManager.GetRoom(client.GetRoom())
		if r != nil && r.State == room.RoomStateEnded {
			h.matcher.PracticeMatchWithDifficulty(client, difficulty)
			return
		}
		if err := h.matcher.AddBotToRoomWithDifficulty(client, difficulty); err != nil {
			client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, err.Error()))
		}
		return
	}

	h.matcher.PracticeMatchWithDifficulty(client, difficulty)
}

func parsePracticeDifficulty(msg *protocol.Message) string {
	if msg == nil || len(msg.Payload) == 0 {
		return ""
	}
	var payload practiceMatchPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return ""
	}
	switch payload.Difficulty {
	case "easy", "normal", "hard":
		return payload.Difficulty
	default:
		return ""
	}
}

// handleReady 处理准备
func (h *Handler) handleReady(client types.ClientInterface, ready bool) {
	if gameSession := h.GetGameSession(client.GetRoom()); gameSession != nil && gameSession.SetTrusteeship(client.GetID(), ready) {
		return
	}
	err := h.roomManager.SetPlayerReady(client, ready)
	if err != nil {
		var gameErr *apperrors.GameError
		if errors.As(err, &gameErr) {
			client.SendMessage(codec.NewErrorMessage(gameErr.Code))
		} else {
			client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, err.Error()))
		}
	}
}

func (h *Handler) handleReplay(client types.ClientInterface) {
	err := h.roomManager.RestartRoom(client)
	if err != nil {
		client.SendMessage(codec.NewErrorMessageWithText(protocol.ErrCodeUnknown, err.Error()))
	}
}
