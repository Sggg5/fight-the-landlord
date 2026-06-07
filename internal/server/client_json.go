package server

import (
	"encoding/json"

	"github.com/palemoky/fight-the-landlord/internal/protocol"
	"github.com/palemoky/fight-the-landlord/internal/protocol/codec"
	payloadconv "github.com/palemoky/fight-the-landlord/internal/protocol/convert/payload"
)

var jsonPayloadDecoders = map[protocol.MessageType]func(json.RawMessage) ([]byte, error){
	protocol.MsgPing:           decodeJSONToPingPayload,
	protocol.MsgReconnect:      decodeJSONToReconnectPayload,
	protocol.MsgJoinRoom:       decodeJSONToJoinRoomPayload,
	protocol.MsgBid:            decodeJSONToBidPayload,
	protocol.MsgPlayCards:      decodeJSONToPlayCardsPayload,
	protocol.MsgChat:           decodeJSONToChatPayload,
	protocol.MsgPurchaseItem:   decodeJSONToPurchaseItemPayload,
	protocol.MsgClaimDailyTask: decodeJSONToClaimDailyTaskPayload,
	protocol.MsgGetLeaderboard: decodeJSONToGetLeaderboardPayload,
}

func decodeJSONToPingPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.PingPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	result, err := payloadconv.EncodePayload(protocol.MsgPing, payload)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func decodeJSONToReconnectPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.ReconnectPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgReconnect, payload)
}

func decodeJSONToJoinRoomPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.JoinRoomPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgJoinRoom, payload)
}

func decodeJSONToBidPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.BidPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgBid, payload)
}

func decodeJSONToPlayCardsPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.PlayCardsPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgPlayCards, payload)
}

func decodeJSONToChatPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.ChatPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgChat, payload)
}

func decodeJSONToPurchaseItemPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.PurchaseItemPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgPurchaseItem, payload)
}

func decodeJSONToClaimDailyTaskPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.ClaimDailyTaskPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgClaimDailyTask, payload)
}

func decodeJSONToGetLeaderboardPayload(data json.RawMessage) ([]byte, error) {
	var payload protocol.GetLeaderboardPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payloadconv.EncodePayload(protocol.MsgGetLeaderboard, payload)
}

func decodeJSONMessage(data []byte) (*protocol.Message, error) {

	var raw struct {
		Type    protocol.MessageType "json:\"type\""
		Payload json.RawMessage      "json:\"payload\""
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	msg := codec.GetMessage()
	msg.Type = raw.Type
	if len(raw.Payload) == 0 {
		return msg, nil
	}
	if decoder, ok := jsonPayloadDecoders[raw.Type]; ok {
		pbBytes, err := decoder(raw.Payload)
		if err != nil {
			codec.PutMessage(msg)
			return nil, err
		}
		msg.Payload = pbBytes
	} else {
		msg.Payload = raw.Payload
	}
	return msg, nil
}

func msgToJSON(msg *protocol.Message) ([]byte, error) {

	type jsonMsg struct {
		Type    protocol.MessageType "json:\"type\""
		Payload interface{}          "json:\"payload,omitempty\""
	}
	result := jsonMsg{Type: msg.Type}
	if len(msg.Payload) == 0 {
		return json.Marshal(result)
	}
	payload, err := decodePayloadToStruct(msg.Type, msg.Payload)
	if err != nil {
		return nil, err
	}
	result.Payload = payload
	return json.Marshal(result)
}

func decodePayloadToStruct(msgType protocol.MessageType, data []byte) (interface{}, error) {

	type entry struct {
		msgType protocol.MessageType
		target  interface{}
	}
	knownTypes := []entry{
		{protocol.MsgConnected, &protocol.ConnectedPayload{}},
		{protocol.MsgReconnected, &protocol.ReconnectedPayload{}},
		{protocol.MsgPong, &protocol.PongPayload{}},
		{protocol.MsgOnlineCount, &protocol.OnlineCountPayload{}},
		{protocol.MsgRoomCreated, &protocol.RoomCreatedPayload{}},
		{protocol.MsgRoomJoined, &protocol.RoomJoinedPayload{}},
		{protocol.MsgPlayerJoined, &protocol.PlayerJoinedPayload{}},
		{protocol.MsgPlayerLeft, &protocol.PlayerLeftPayload{}},
		{protocol.MsgDealCards, &protocol.DealCardsPayload{}},
		{protocol.MsgBidTurn, &protocol.BidTurnPayload{}},
		{protocol.MsgBidResult, &protocol.BidResultPayload{}},
		{protocol.MsgPlayTurn, &protocol.PlayTurnPayload{}},
		{protocol.MsgError, &protocol.ErrorPayload{}},
		{protocol.MsgPlayerOffline, &protocol.PlayerOfflinePayload{}},
		{protocol.MsgPlayerOnline, &protocol.PlayerOnlinePayload{}},
		{protocol.MsgMaintenancePush, &protocol.MaintenancePayload{}},
		{protocol.MsgMaintenancePull, &protocol.MaintenancePayload{}},
		{protocol.MsgChat, &protocol.ChatPayload{}},
		{protocol.MsgStatsResult, &protocol.StatsResultPayload{}},
		{protocol.MsgLeaderboardResult, &protocol.LeaderboardResultPayload{}},
		{protocol.MsgRoomListResult, &protocol.RoomListResultPayload{}},
		{protocol.MsgSignInResult, &protocol.SignInResultPayload{}},
		{protocol.MsgAchievementsResult, &protocol.AchievementsResultPayload{}},
		{protocol.MsgShopListResult, &protocol.ShopListResultPayload{}},
		{protocol.MsgPurchaseItemResult, &protocol.PurchaseItemResultPayload{}},
		{protocol.MsgDailyTasksResult, &protocol.DailyTasksResultPayload{}},
		{protocol.MsgClaimDailyTaskResult, &protocol.ClaimDailyTaskResultPayload{}},
		{protocol.MsgGameStart, &protocol.GameStartPayload{}},
		{protocol.MsgLandlord, &protocol.LandlordPayload{}},
		{protocol.MsgCardPlayed, &protocol.CardPlayedPayload{}},
		{protocol.MsgPlayerPass, &protocol.PlayerPassPayload{}},
		{protocol.MsgGameOver, &protocol.GameOverPayload{}},
		{protocol.MsgHintResult, &protocol.HintResultPayload{}},
		{protocol.MsgPlayerReady, &protocol.PlayerReadyPayload{}},
		{protocol.MsgMatchFound, &protocol.RoomJoinedPayload{}},
	}
	for _, kt := range knownTypes {
		if kt.msgType == msgType {
			if err := payloadconv.DecodePayload(msgType, data, kt.target); err != nil {
				return nil, err
			}
			return kt.target, nil
		}
	}
	return nil, nil
}
