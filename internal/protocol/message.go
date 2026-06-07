package protocol

import "encoding/json"

// Message 氓聼潞莽隆聙忙露聢忙聛炉莽禄聯忙聻聞
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// MessageType 忙露聢忙聛炉莽卤禄氓聻聥
type MessageType string

// 氓庐垄忙聢路莽芦?芒聠?忙聹聧氓聤隆莽芦?忙露聢忙聛炉莽卤禄氓聻聥
const (
	// 猫驴聻忙聨楼忙聯聧盲陆聹
	MsgReconnect MessageType = "reconnect" // 忙聳颅莽潞驴茅聡聧猫驴聻
	MsgPing      MessageType = "ping"      // 氓驴聝猫路鲁 ping

	// 忙聢驴茅聴麓忙聯聧盲陆聹
	MsgCreateRoom    MessageType = "create_room"    // 氓聢聸氓禄潞忙聢驴茅聴麓
	MsgJoinRoom      MessageType = "join_room"      // 氓聤聽氓聟楼忙聢驴茅聴麓
	MsgLeaveRoom     MessageType = "leave_room"     // 莽娄禄氓录聙忙聢驴茅聴麓
	MsgQuickMatch    MessageType = "quick_match"    // 氓驴芦茅聙聼氓聦鹿茅聟?
	MsgPracticeMatch MessageType = "practice_match" // 盲潞潞忙聹潞莽禄聝盲鹿聽
	MsgReady         MessageType = "ready"          // 氓聡聠氓陇聡氓掳卤莽禄陋
	MsgCancelReady   MessageType = "cancel_ready"   // 氓聫聳忙露聢氓聡聠氓陇聡
	MsgReplay        MessageType = "replay"

	// 忙赂赂忙聢聫忙聯聧盲陆聹
	MsgBid       MessageType = "bid"        // 氓聫芦氓聹掳盲赂?
	MsgPlayCards MessageType = "play_cards" // 氓聡潞莽聣聦
	MsgPass      MessageType = "pass"       // 盲赂聧氓聡潞
	MsgHint      MessageType = "hint"

	// 忙聨聮猫隆聦忙娄?
	MsgGetStats             MessageType = "get_stats"              // 猫聨路氓聫聳盲赂陋盲潞潞莽禄聼猫庐隆
	MsgGetLeaderboard       MessageType = "get_leaderboard"        // 猫聨路氓聫聳忙聨聮猫隆聦忙娄?
	MsgGetRoomList          MessageType = "get_room_list"          // 猫聨路氓聫聳忙聢驴茅聴麓氓聢聴猫隆篓
	MsgGetOnlineCount       MessageType = "get_online_count"       // 猫聨路氓聫聳氓聹篓莽潞驴盲潞潞忙聲掳
	MsgGetMaintenanceStatus MessageType = "get_maintenance_status" // 猫聨路氓聫聳莽禄麓忙聤陇莽聤露忙聙?
	MsgChat                 MessageType = "chat"                   // 猫聛聤氓陇漏忙露聢忙聛炉

	// ?????
	MsgSignIn          MessageType = "sign_in" // 脟漏碌陆
	MsgGetAchievements MessageType = "get_achievements"

	// 鍟嗗煄
	MsgShopList     MessageType = "shop_list"
	MsgPurchaseItem MessageType = "purchase_item"

	// 姣忔棩浠诲姟
	MsgGetDailyTasks  MessageType = "get_daily_tasks"
	MsgClaimDailyTask MessageType = "claim_daily_task" // 禄帽脠隆鲁脡戮脥脕脨卤铆

	// 脡脤鲁脟

	// 脙驴脠脮脠脦脦帽
)

const (
	// 忙聹聧氓聤隆莽芦?芒聠?氓庐垄忙聢路莽芦?忙露聢忙聛炉莽卤禄氓聻聥
	// 猫驴聻忙聨楼莽聸赂氓聟鲁
	MsgConnected     MessageType = "connected"      // 猫驴聻忙聨楼忙聢聬氓聤聼
	MsgReconnected   MessageType = "reconnected"    // 茅聡聧猫驴聻忙聢聬氓聤聼
	MsgPong          MessageType = "pong"           // 氓驴聝猫路鲁 pong
	MsgPlayerOffline MessageType = "player_offline" // 莽聨漏氓庐露忙聨聣莽潞驴茅聙職莽聼楼
	MsgPlayerOnline  MessageType = "player_online"  // 莽聨漏氓庐露盲赂聤莽潞驴茅聙職莽聼楼
	MsgOnlineCount   MessageType = "online_count"   // 氓聹篓莽潞驴盲潞潞忙聲掳忙聸麓忙聳掳

	// 忙聢驴茅聴麓莽聸赂氓聟鲁
	MsgRoomCreated  MessageType = "room_created"  // 忙聢驴茅聴麓氓聢聸氓禄潞忙聢聬氓聤聼
	MsgRoomJoined   MessageType = "room_joined"   // 氓聤聽氓聟楼忙聢驴茅聴麓忙聢聬氓聤聼
	MsgPlayerJoined MessageType = "player_joined" // 氓聟露盲禄聳莽聨漏氓庐露氓聤聽氓聟楼
	MsgPlayerLeft   MessageType = "player_left"   // 莽聨漏氓庐露莽娄禄氓录聙
	MsgPlayerReady  MessageType = "player_ready"  // 莽聨漏氓庐露氓聡聠氓陇聡
	MsgMatchFound   MessageType = "match_found"   // 氓聦鹿茅聟聧忙聢聬氓聤聼

	// 忙赂赂忙聢聫忙碌聛莽篓聥
	MsgGameStart   MessageType = "game_start"   // 忙赂赂忙聢聫氓录聙氓搂?
	MsgDealCards   MessageType = "deal_cards"   // 氓聫聭莽聣聦
	MsgBidTurn     MessageType = "bid_turn"     // 猫陆庐氓聢掳氓聫芦氓聹掳盲赂?
	MsgBidResult   MessageType = "bid_result"   // 氓聫芦氓聹掳盲赂禄莽禄聯忙聻?
	MsgLandlord    MessageType = "landlord"     // 氓聹掳盲赂禄莽隆庐氓庐職
	MsgPlayTurn    MessageType = "play_turn"    // 猫陆庐氓聢掳氓聡潞莽聣聦
	MsgCardPlayed  MessageType = "card_played"  // 忙聹聣盲潞潞氓聡潞莽聣聦
	MsgPlayerPass  MessageType = "player_pass"  // 忙聹聣盲潞潞盲赂聧氓聡潞
	MsgGameOver    MessageType = "game_over"    // 忙赂赂忙聢聫莽禄聯忙聺聼
	MsgRoundResult MessageType = "round_result" // 忙聹卢猫陆庐莽禄聯忙聻聹
	MsgHintResult  MessageType = "hint_result"

	// 忙聨聮猫隆聦忙娄?
	MsgStatsResult       MessageType = "stats_result"       // 盲赂陋盲潞潞莽禄聼猫庐隆莽禄聯忙聻聹
	MsgLeaderboardResult MessageType = "leaderboard_result" // 忙聨聮猫隆聦忙娄聹莽禄聯忙聻?
	MsgRoomListResult    MessageType = "room_list_result"   // 忙聢驴茅聴麓氓聢聴猫隆篓莽禄聯忙聻聹

	// 莽鲁禄莽禄聼茅聙職莽聼楼
	MsgMaintenancePush MessageType = "maintenance_push" // 盲赂禄氓聤篓忙聨篓茅聙?
	MsgMaintenancePull MessageType = "maintenance_pull" // 猫垄芦氓聤篓忙聥聣氓聫聳

	// 茅聰聶猫炉炉

	// ?????
	MsgSignInResult       MessageType = "sign_in_result" // 脟漏碌陆陆谩鹿没
	MsgAchievementsResult MessageType = "achievements_result"

	// 鍟嗗煄
	MsgShopListResult     MessageType = "shop_list_result"
	MsgPurchaseItemResult MessageType = "purchase_item_result"

	// 姣忔棩浠诲姟
	MsgDailyTasksResult     MessageType = "daily_tasks_result"
	MsgClaimDailyTaskResult MessageType = "claim_daily_task_result" // 鲁脡戮脥脕脨卤铆陆谩鹿没

	// 脡脤鲁脟

	// 脙驴脠脮脠脦脦帽

	MsgError MessageType = "error" // 茅聰聶猫炉炉忙露聢忙聛炉
)
