package room

import (
	"errors"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	"github.com/palemoky/fight-the-landlord/internal/types"
)

// --- Room 方法 ---

// Broadcast 广播消息给房间内所有玩家
func (r *Room) Broadcast(msg *protocol.Message) {
	if r == nil {
		return
	}
	for _, player := range r.Players {
		if player == nil || player.Client == nil {
			continue
		}
		player.Client.SendMessage(msg)
	}
}

// broadcastExcept 广播消息给除指定玩家外的所有玩家
func (r *Room) BroadcastExcept(excludeID string, msg *protocol.Message) {
	if r == nil {
		return
	}
	for id, player := range r.Players {
		if id != excludeID {
			if player == nil || player.Client == nil {
				continue
			}
			player.Client.SendMessage(msg)
		}
	}
}

// checkAllReady 检查是否所有玩家都准备好
func (r *Room) checkAllReady() bool {
	if len(r.Players) < 3 {
		return false
	}
	for _, player := range r.Players {
		if !player.Ready {
			return false
		}
	}
	return true
}

// GetPlayerInfo 获取玩家信息
func (r *Room) GetPlayerInfo(playerID string) protocol.PlayerInfo {
	player := r.Players[playerID]
	cardsCount := 0
	// 游戏会话由外部调用方管理，此处暂不传入
	return protocol.PlayerInfo{
		ID:         player.Client.GetID(),
		Name:       player.Client.GetName(),
		Seat:       player.Seat,
		Ready:      player.Ready,
		IsLandlord: player.IsLandlord,
		CardsCount: cardsCount,
		IsBot:      player.Client.IsBot(),
	}
}

// GetAllPlayersInfo 获取所有玩家信息
func (r *Room) GetAllPlayersInfo() []protocol.PlayerInfo {
	infos := make([]protocol.PlayerInfo, 0, len(r.Players))
	for _, id := range r.PlayerOrder {
		infos = append(infos, r.GetPlayerInfo(id))
	}
	return infos
}

// StartGame 准备开始游戏（不创建GameSession，由外部管理）
// 注意：调用者负责保存到 Redis
func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.startGameLocked()
}

// RestartGame starts another round with the players already in an ended room.
func (r *Room) RestartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State != RoomStateEnded || len(r.Players) < 3 {
		return errors.New("cannot restart game: room has not ended or not enough players")
	}
	for _, player := range r.Players {
		if player == nil || player.Client == nil {
			return errors.New("cannot restart game: player is offline")
		}
		player.Ready = true
		player.IsLandlord = false
	}
	r.State = RoomStateReady
	r.Broadcast(codec.MustNewMessage(protocol.MsgGameStart, protocol.GameStartPayload{
		Players: r.GetAllPlayersInfo(),
	}))
	return nil
}

func (r *Room) HasBot() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, player := range r.Players {
		if player != nil && player.Client != nil && player.Client.IsBot() {
			return true
		}
	}
	return false
}

func (r *Room) GetClients() []types.ClientInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()
	clients := make([]types.ClientInterface, 0, len(r.PlayerOrder))
	for _, id := range r.PlayerOrder {
		if player := r.Players[id]; player != nil && player.Client != nil {
			clients = append(clients, player.Client)
		}
	}
	return clients
}

// startGameLocked 开始游戏（调用者已持有锁时使用）
func (r *Room) startGameLocked() error {
	if r.State != RoomStateWaiting || len(r.Players) < 3 {
		return errors.New("cannot start game: room not ready or not enough players")
	}

	r.State = RoomStateReady

	// 广播游戏开始
	r.Broadcast(codec.MustNewMessage(protocol.MsgGameStart, protocol.GameStartPayload{
		Players: r.GetAllPlayersInfo(),
	}))

	return nil
}
