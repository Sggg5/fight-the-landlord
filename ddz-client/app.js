const schema = `
syntax = "proto3";
package protocol;

enum MessageType {
  MSG_UNKNOWN = 0;
  MSG_RECONNECT = 1;
  MSG_PING = 2;
  MSG_CREATE_ROOM = 3;
  MSG_JOIN_ROOM = 4;
  MSG_LEAVE_ROOM = 5;
  MSG_QUICK_MATCH = 6;
  MSG_READY = 7;
  MSG_CANCEL_READY = 8;
  MSG_BID = 9;
  MSG_PLAY_CARDS = 10;
  MSG_PASS = 11;
  MSG_GET_STATS = 12;
  MSG_GET_LEADERBOARD = 13;
  MSG_GET_ROOM_LIST = 14;
  MSG_GET_ONLINE_COUNT = 15;
  MSG_CHAT = 16;
  MSG_GET_MAINTENANCE_STATUS = 17;
  MSG_CONNECTED = 100;
  MSG_RECONNECTED = 101;
  MSG_PONG = 102;
  MSG_PLAYER_OFFLINE = 103;
  MSG_PLAYER_ONLINE = 104;
  MSG_ONLINE_COUNT = 105;
  MSG_ROOM_CREATED = 106;
  MSG_ROOM_JOINED = 107;
  MSG_PLAYER_JOINED = 108;
  MSG_PLAYER_LEFT = 109;
  MSG_PLAYER_READY = 110;
  MSG_MATCH_FOUND = 111;
  MSG_GAME_START = 112;
  MSG_DEAL_CARDS = 113;
  MSG_BID_TURN = 114;
  MSG_BID_RESULT = 115;
  MSG_LANDLORD = 116;
  MSG_PLAY_TURN = 117;
  MSG_CARD_PLAYED = 118;
  MSG_PLAYER_PASS = 119;
  MSG_GAME_OVER = 120;
  MSG_ROUND_RESULT = 121;
  MSG_STATS_RESULT = 122;
  MSG_LEADERBOARD_RESULT = 123;
  MSG_ROOM_LIST_RESULT = 124;
  MSG_MAINTENANCE_STATUS = 125;
  MSG_MAINTENANCE = 126;
  MSG_ERROR = 200;
  MSG_PRACTICE_MATCH = 201;
}

message Message { MessageType type = 1; bytes payload = 2; }
message CardInfo { int64 suit = 1; int64 rank = 2; int64 color = 3; }
message PlayerInfo { string id = 1; string name = 2; int64 seat = 3; bool ready = 4; bool is_landlord = 5; int64 cards_count = 6; bool online = 7; }
message PlayerHand { string player_id = 1; string player_name = 2; repeated CardInfo cards = 3; }
message PlayerScore { string player_id = 1; string player_name = 2; bool is_landlord = 3; int64 score = 4; }
message GameStateDTO { string phase = 1; repeated PlayerInfo players = 2; repeated CardInfo hand = 3; repeated CardInfo bottom_cards = 4; string current_turn = 5; repeated CardInfo last_played = 6; string last_player_id = 7; bool must_play = 8; bool can_beat = 9; }
message LeaderboardEntry { int64 rank = 1; string player_id = 2; string player_name = 3; int64 score = 4; int64 wins = 5; double win_rate = 6; }
message RoomListItem { string room_code = 1; int64 player_count = 2; int64 max_players = 3; }
message ReconnectPayload { string token = 1; string player_id = 2; }
message PingPayload { int64 timestamp = 1; }
message JoinRoomPayload { string room_code = 1; }
message BidPayload { bool bid = 1; }
message PlayCardsPayload { repeated CardInfo cards = 1; }
message GetLeaderboardPayload { string type = 1; int64 offset = 2; int64 limit = 3; }
message ConnectedPayload { string player_id = 1; string player_name = 2; string reconnect_token = 3; }
message ReconnectedPayload { string player_id = 1; string player_name = 2; string room_code = 3; GameStateDTO game_state = 4; }
message PongPayload { int64 client_timestamp = 1; int64 server_timestamp = 2; }
message PlayerOfflinePayload { string player_id = 1; string player_name = 2; int64 timeout = 3; }
message PlayerOnlinePayload { string player_id = 1; string player_name = 2; }
message OnlineCountPayload { int64 count = 1; }
message MaintenanceStatusPayload { bool maintenance = 1; }
message MaintenancePayload { bool maintenance = 1; }
message ErrorPayload { int64 code = 1; string message = 2; }
message StatsResultPayload { string player_id = 1; string player_name = 2; int64 total_games = 3; int64 wins = 4; int64 losses = 5; double win_rate = 6; int64 landlord_games = 7; int64 landlord_wins = 8; int64 farmer_games = 9; int64 farmer_wins = 10; int64 score = 11; int64 rank = 12; int64 current_streak = 13; int64 max_win_streak = 14; }
message LeaderboardResultPayload { string type = 1; repeated LeaderboardEntry entries = 2; }
message RoomListResultPayload { repeated RoomListItem rooms = 1; }
message RoomCreatedPayload { string room_code = 1; PlayerInfo player = 2; }
message RoomJoinedPayload { string room_code = 1; PlayerInfo player = 2; repeated PlayerInfo players = 3; }
message PlayerJoinedPayload { PlayerInfo player = 1; }
message PlayerLeftPayload { string player_id = 1; string player_name = 2; }
message PlayerReadyPayload { string player_id = 1; bool ready = 2; }
message GameStartPayload { repeated PlayerInfo players = 1; }
message DealCardsPayload { repeated CardInfo cards = 1; repeated CardInfo bottom_cards = 2; }
message BidTurnPayload { string player_id = 1; int64 timeout = 2; bool is_grab = 3; int64 multiplier = 4; }
message BidResultPayload { string player_id = 1; string player_name = 2; bool bid = 3; bool is_grab = 4; int64 multiplier = 5; }
message LandlordPayload { string player_id = 1; string player_name = 2; repeated CardInfo bottom_cards = 3; int64 multiplier = 4; }
message PlayTurnPayload { string player_id = 1; int64 timeout = 2; bool must_play = 3; bool can_beat = 4; }
message CardPlayedPayload { string player_id = 1; string player_name = 2; repeated CardInfo cards = 3; int64 cards_left = 4; string hand_type = 5; }
message PlayerPassPayload { string player_id = 1; string player_name = 2; }
message GameOverPayload { string winner_id = 1; string winner_name = 2; bool is_landlord = 3; repeated PlayerHand player_hands = 4; int64 multiplier = 5; repeated PlayerScore scores = 6; }
`;

const root = protobuf.parse(schema).root;
const Message = root.lookupType("protocol.Message");

const typeNameToEnum = {
  reconnect: 1, ping: 2, create_room: 3, join_room: 4, leave_room: 5,
  quick_match: 6, ready: 7, cancel_ready: 8, replay: 202, bid: 9, play_cards: 10,
  pass: 11, get_stats: 12, get_leaderboard: 13, get_room_list: 14,
  get_online_count: 15, chat: 16, get_maintenance_status: 17, hint: 18,
  connected: 100, reconnected: 101, pong: 102, player_offline: 103,
  player_online: 104, online_count: 105, room_created: 106,
  room_joined: 107, player_joined: 108, player_left: 109,
  player_ready: 110, match_found: 111, game_start: 112, deal_cards: 113,
  bid_turn: 114, bid_result: 115, landlord: 116, play_turn: 117,
  card_played: 118, player_pass: 119, game_over: 120, round_result: 121,
  stats_result: 122, leaderboard_result: 123, room_list_result: 124,
  maintenance_status: 125, maintenance: 126, hint_result: 127, error: 200, practice_match: 201,
};
const enumToTypeName = Object.fromEntries(Object.entries(typeNameToEnum).map(([k, v]) => [v, k]));
const payloadType = {
  reconnect: "ReconnectPayload", ping: "PingPayload", join_room: "JoinRoomPayload",
  bid: "BidPayload", play_cards: "PlayCardsPayload", get_leaderboard: "GetLeaderboardPayload",
  connected: "ConnectedPayload", reconnected: "ReconnectedPayload", pong: "PongPayload",
  player_offline: "PlayerOfflinePayload", player_online: "PlayerOnlinePayload",
  online_count: "OnlineCountPayload", maintenance_status: "MaintenanceStatusPayload",
  maintenance: "MaintenancePayload", error: "ErrorPayload",
  leaderboard_result: "LeaderboardResultPayload", room_list_result: "RoomListResultPayload",
  room_created: "RoomCreatedPayload", room_joined: "RoomJoinedPayload",
  player_joined: "PlayerJoinedPayload", player_left: "PlayerLeftPayload",
  player_ready: "PlayerReadyPayload", game_start: "GameStartPayload",
  deal_cards: "DealCardsPayload", bid_turn: "BidTurnPayload", bid_result: "BidResultPayload",
  landlord: "LandlordPayload", play_turn: "PlayTurnPayload", card_played: "CardPlayedPayload",
  player_pass: "PlayerPassPayload", game_over: "GameOverPayload",
};

const state = {
  ws: null,
  connected: false,
  playerId: "",
  playerName: "",
  reconnectToken: "",
  roomCode: "",
  players: [],
  hand: [],
  selected: new Set(),
  bottomCards: [],
  lastPlayed: [],
  lastPlayerId: "",
  lastPlayerName: "",
  currentTurn: "",
  turnMustPlay: false,
  turnCanBeat: true,
  bidTurn: "",
  phase: "lobby",
  multiplier: 1,
  yuanzi: 0,
  botYuanzi: {},
  trusteeship: false,
  onlineCount: null,
  latency: null,
};

const uiState = {
  currentView: "lobby",
};
const dialogueTimers = new Map();
let practiceReplayTimer = null;
const DEFAULT_CENTER_POETRY = "\u6620\u9636\u78a7\u8349\u81ea\u6625\u8272\uff0c\u9694\u53f6\u9ec4\u9e42\u7a7a\u597d\u97f3\u3002";

const poetryFallbackLines = [
  "醉卧沙场君莫笑，古来征战几人回。",
  "黄沙百战穿金甲，不破楼兰终不还。",
  "会当凌绝顶，一览众山小。",
  "长风破浪会有时，直挂云帆济沧海。",
  "大鹏一日同风起，扶摇直上九万里。",
  "男儿何不带吴钩，收取关山五十州。",
];

const poetryExpandedLines = [
  "滚滚长江东逝水，浪花淘尽英雄。",
  "青山依旧在，几度夕阳红。",
  "大江东去，浪淘尽，千古风流人物。",
  "江山如画，一时多少豪杰。",
  "羽扇纶巾，谈笑间，樯橹灰飞烟灭。",
  "故垒西边，人道是，三国周郎赤壁。",
  "出师未捷身先死，长使英雄泪满襟。",
  "三顾频烦天下计，两朝开济老臣心。",
  "功盖三分国，名成八阵图。",
  "江流石不转，遗恨失吞吴。",
  "丞相祠堂何处寻，锦官城外柏森森。",
  "映阶碧草自春色，隔叶黄鹂空好音。",
  "白帝高为三峡镇，瞿塘险过百牢关。",
  "星垂平野阔，月涌大江流。",
  "无边落木萧萧下，不尽长江滚滚来。",
  "会当凌绝顶，一览众山小。",
  "长风破浪会有时，直挂云帆济沧海。",
  "大鹏一日同风起，扶摇直上九万里。",
  "黄沙百战穿金甲，不破楼兰终不还。",
  "醉卧沙场君莫笑，古来征战几人回。",
  "男儿何不带吴钩，收取关山五十州。",
  "但使龙城飞将在，不教胡马度阴山。",
  "秦时明月汉时关，万里长征人未还。",
  "愿将腰下剑，直为斩楼兰。",
  "报君黄金台上意，提携玉龙为君死。",
  "黑云压城城欲摧，甲光向日金鳞开。",
  "风萧萧兮易水寒，壮士一去兮不复还。",
  "老骥伏枥，志在千里。",
  "烈士暮年，壮心不已。",
  "对酒当歌，人生几何。",
  "山不厌高，海不厌深。",
  "周公吐哺，天下归心。",
  "秋风萧瑟，洪波涌起。",
  "日月之行，若出其中。",
  "星汉灿烂，若出其里。",
  "捐躯赴国难，视死忽如归。",
  "丈夫志四海，万里犹比邻。",
  "刑天舞干戚，猛志固常在。",
  "壮志饥餐胡虏肉，笑谈渴饮匈奴血。",
  "莫等闲，白了少年头，空悲切。",
  "人生自古谁无死，留取丹心照汗青。",
  "千磨万击还坚劲，任尔东西南北风。",
  "我自横刀向天笑，去留肝胆两昆仑。",
  "十步杀一人，千里不留行。",
  "事了拂衣去，深藏身与名。",
  "天生我材必有用，千金散尽还复来。",
  "仰天大笑出门去，我辈岂是蓬蒿人。",
  "俱怀逸兴壮思飞，欲上青天揽明月。",
  "人生得意须尽欢，莫使金樽空对月。",
  "海内存知己，天涯若比邻。",
  "潮平两岸阔，风正一帆悬。",
  "沉舟侧畔千帆过，病树前头万木春。",
  "不畏浮云遮望眼，自缘身在最高层。",
  "山重水复疑无路，柳暗花明又一村。",
  "纸上得来终觉浅，绝知此事要躬行。",
  "路漫漫其修远兮，吾将上下而求索。",
  "亦余心之所善兮，虽九死其犹未悔。",
  "宝剑锋从磨砺出，梅花香自苦寒来。",
  "千淘万漉虽辛苦，吹尽狂沙始到金。",
  "莫愁前路无知己，天下谁人不识君。",
];
let poetryQueue = [];
let lastPoetryLine = "";

let poetryLoading = false;

const $ = (id) => document.getElementById(id);
const developerMode = new URLSearchParams(window.location.search).get("dev") === "1";
document.body.dataset.developerMode = String(developerMode);
let reconnectTimer = null;
let activeServerUrl = "";

function selectedPracticeDifficulty() {
  const botControl = $("botDifficultyControl");
  const value = botControl && !botControl.hidden ? $("botDifficulty")?.value : $("practiceDifficulty")?.value;
  return ["easy", "normal", "hard"].includes(value) ? value : "normal";
}

function practiceMatchPayload() {
  return { difficulty: selectedPracticeDifficulty() };
}

function lookupPayload(type) {
  const name = payloadType[type];
  return name ? root.lookupType(`protocol.${name}`) : null;
}

function normalize(obj) {
  return JSON.parse(JSON.stringify(obj ?? {}));
}

function encodePayload(type, payload) {
  if (!payload) return new Uint8Array();
  const Proto = lookupPayload(type);
  if (Proto) return Proto.encode(Proto.create(payload)).finish();
  return new TextEncoder().encode(JSON.stringify(payload));
}

function decodePayload(type, bytes) {
  if (!bytes || bytes.length === 0) return {};
  const Proto = lookupPayload(type);
  if (Proto) return normalize(Proto.toObject(Proto.decode(bytes), { longs: Number, defaults: true }));
  try {
    return JSON.parse(new TextDecoder().decode(bytes));
  } catch {
    return {};
  }
}

function send(type, payload) {
  if (!state.ws || state.ws.readyState !== WebSocket.OPEN) {
    addLog("系统", "还没有连接服务器");
    return;
  }
  const packet = Message.encode(Message.create({
    type: typeNameToEnum[type],
    payload: encodePayload(type, payload),
  })).finish();
  state.ws.send(packet);
}

function connect(url) {
  if (!url) return;
  activeServerUrl = url;
  window.clearTimeout(reconnectTimer);
  if ((state.ws?.readyState === WebSocket.OPEN || state.ws?.readyState === WebSocket.CONNECTING) && state.ws.url === url) return;
  if (state.ws) state.ws.close();
  state.ws = new WebSocket(url);
  state.ws.binaryType = "arraybuffer";
  setStatus("连接中");

  state.ws.onopen = () => {
    state.connected = true;
    setStatus("已连接");
    addLog("系统", developerMode ? `已连接 ${url}` : "已连接服务器");
    send("get_online_count");
    send("get_maintenance_status");
  };

  state.ws.onclose = () => {
    state.connected = false;
    setStatus("已断开");
    addLog("系统", "连接已断开");
    render();
    if (!developerMode && activeServerUrl) {
      setStatus("重新连接中");
      reconnectTimer = window.setTimeout(() => connect(activeServerUrl), 3000);
    }
  };

  state.ws.onerror = () => addLog("系统", "连接发生错误");
  state.ws.onmessage = (event) => {
    const msg = Message.decode(new Uint8Array(event.data));
    const type = enumToTypeName[msg.type] || "unknown";
    handleMessage(type, decodePayload(type, msg.payload));
  };
}

function handleMessage(type, payload) {
  switch (type) {
    case "connected":
      state.playerId = payload.playerId;
      state.playerName = payload.playerName;
      state.reconnectToken = payload.reconnectToken;
      localStorage.setItem("ddzReconnect", JSON.stringify({
        playerId: state.playerId,
        token: state.reconnectToken,
      }));
      send("get_stats");
      addLog("系统", `欢迎 ${state.playerName}`);
      break;
    case "reconnected":
      state.playerId = payload.playerId;
      state.playerName = payload.playerName;
      state.roomCode = payload.roomCode || "";
      if (payload.gameState) applyGameState(payload.gameState);
      if (state.roomCode) showGame();
      addLog("系统", "重连成功");
      break;
    case "pong":
      state.latency = Math.max(0, Date.now() - Number(payload.clientTimestamp || Date.now()));
      break;
    case "online_count":
      state.onlineCount = payload.count;
      break;
    case "maintenance_status":
    case "maintenance":
      if (payload.maintenance) addLog("系统", "服务器正在维护");
      break;
    case "room_created":
      state.roomCode = payload.roomCode;
      state.players = [payload.player].filter(Boolean);
      state.phase = "waiting";
      showGame();
      addLog("房间", `已创建房间 ${state.roomCode}`);
      break;
    case "room_joined":
      state.roomCode = payload.roomCode;
      state.players = payload.players || [];
      assignBotCharacters();
      state.phase = "waiting";
      showGame();
      addLog("房间", `已加入房间 ${state.roomCode}`);
      break;
    case "player_joined":
      upsertPlayer(payload.player);
      addLog("房间", `${payload.player?.name || "玩家"} 加入房间`);
      break;
    case "player_left":
      const leavingName = playerName(payload.playerId) || payload.playerName || "玩家";
      state.players = state.players.filter((p) => p.id !== payload.playerId);
      addLog("房间", `${leavingName} 离开房间`);
      break;
    case "player_ready":
      state.players = state.players.map((p) => p.id === payload.playerId ? { ...p, ready: payload.ready } : p);
      break;
    case "game_start":
      state.players = (payload.players || state.players).map((player) => ({ ...player, isLandlord: false }));
      assignBotCharacters();
      state.phase = "bidding";
      state.hand = [];
      state.selected.clear();
      state.bottomCards = [];
      state.lastPlayed = [];
      state.lastPlayerId = "";
      state.lastPlayerName = "";
      state.currentTurn = "";
      state.bidTurn = "";
      state.turnMustPlay = false;
      state.turnCanBeat = true;
      state.multiplier = 1;
      addLog("牌局", "游戏开始");
      window.setTimeout(() => {
        state.players.filter(isBotPlayer).forEach((player, index) => {
          window.setTimeout(() => triggerDialogue(player.id, "gameStart"), index * 450);
        });
      }, 300);
      break;
    case "deal_cards":
      state.hand = sortCards(payload.cards || []);
      state.selected.clear();
      state.bottomCards = payload.bottomCards || [];
      refreshPoetryLine();
      break;
    case "bid_turn":
      state.phase = "bidding";
      state.bidTurn = payload.playerId;
      state.multiplier = payload.multiplier || state.multiplier;
      addLog("叫地主", `${playerName(payload.playerId)} ${payload.isGrab ? "抢地主" : "叫地主"} 回合`);
      break;
    case "bid_result":
      state.multiplier = payload.multiplier || state.multiplier;
      addLog("叫地主", `${playerName(payload.playerId)} ${payload.bid ? (payload.isGrab ? "抢" : "叫") : "不叫"}`);
      triggerDialogue(payload.playerId, payload.bid ? "callLandlord" : "passLandlord");
      break;
    case "landlord":
      state.phase = "playing";
      state.multiplier = payload.multiplier || state.multiplier;
      state.bottomCards = payload.bottomCards || [];
      state.players = state.players.map((p) => ({ ...p, isLandlord: p.id === payload.playerId }));
      if (payload.playerId === state.playerId) state.hand = sortCards([...state.hand, ...state.bottomCards]);
      addLog("地主", `${playerName(payload.playerId)} 成为地主`);
      break;
    case "play_turn":
      state.phase = "playing";
      state.currentTurn = payload.playerId;
      state.turnMustPlay = Boolean(payload.mustPlay);
      state.turnCanBeat = payload.canBeat !== false;
      addLog("出牌", `轮到 ${playerName(payload.playerId)}${payload.mustPlay ? "，必须出牌" : ""}`);
      break;
    case "card_played":
      state.lastPlayed = payload.cards || [];
      state.lastPlayerId = payload.playerId;
      state.lastPlayerName = playerName(payload.playerId);
      updateCardsLeft(payload.playerId, payload.cardsLeft);
      if (payload.playerId === state.playerId) removePlayedCards(payload.cards || []);
      addLog("出牌", `${playerName(payload.playerId)} 出了 ${payload.handType || "牌"}`);
      triggerDialogue(payload.playerId, dialogueEventForCards(payload.cards || [], payload.handType));
      if (payload.cardsLeft > 0 && payload.cardsLeft <= 2) triggerDialogue(payload.playerId, "selfAlmostWin");
      state.players
        .filter((player) => player.id !== payload.playerId && isBotPlayer(player) && player.cardsCount > 0 && player.cardsCount <= 2)
        .forEach((player) => triggerDialogue(player.id, "teammateAlmostWin"));
      break;
    case "player_pass":
      addLog("出牌", `${playerName(payload.playerId)} 不出`);
      triggerDialogue(payload.playerId, "pass");
      break;
    case "game_over":
      const shouldAutoReplay = isPracticeRoom();
      state.phase = "game_over";
      state.trusteeship = false;
      state.multiplier = payload.multiplier || state.multiplier;
      state.currentTurn = "";
      state.turnMustPlay = false;
      state.turnCanBeat = true;
      addLog("结算", `${playerName(payload.winnerId) || payload.winnerName} 获胜，倍数 ${payload.multiplier || 1}x`);
      state.players.forEach((player) => {
        triggerDialogue(player.id, player.id === payload.winnerId ? "win" : "lose");
      });
      updateBotYuanzi(payload.scores || []);
      renderResult(payload);
      window.setTimeout(() => send("get_stats"), 500);
      if (shouldAutoReplay) schedulePracticeReplay();
      break;
    case "room_list_result":
      renderRoomList(payload.rooms || []);
      break;
    case "leaderboard_result":
      renderLeaderboard(payload.entries || []);
      break;
    case "stats_result":
      state.yuanzi = Number(payload.coins || 0);
      renderStats(payload);
      break;
    case "hint_result":
      applyHintCards(payload.cards || []);
      addLog("提示", payload.message || "已给出推荐");
      break;
    case "chat":
      addLog(payload.senderName || "消息", payload.content || "");
      break;
    case "error":
      addLog("错误", payload.message || `错误 ${payload.code}`);
      break;
    default:
      addLog("收到", type);
  }
  render();
}

function applyGameState(gs) {
  state.phase = gs.phase || state.phase;
  state.players = gs.players || [];
  state.hand = sortCards(gs.hand || []);
  state.bottomCards = gs.bottomCards || [];
  state.currentTurn = gs.currentTurn || "";
  state.turnMustPlay = Boolean(gs.mustPlay);
  state.turnCanBeat = gs.canBeat !== false;
  state.lastPlayed = gs.lastPlayed || [];
  state.lastPlayerId = gs.lastPlayerId || "";
}

function upsertPlayer(player) {
  if (!player) return;
  const idx = state.players.findIndex((p) => p.id === player.id);
  if (idx >= 0) state.players[idx] = player;
  else state.players.push(player);
}

function updateCardsLeft(playerId, count) {
  state.players = state.players.map((p) => p.id === playerId ? { ...p, cardsCount: count } : p);
}

function removePlayedCards(cards) {
  const keys = cards.map(cardKey);
  state.hand = state.hand.filter((card) => {
    const idx = keys.indexOf(cardKey(card));
    if (idx >= 0) {
      keys.splice(idx, 1);
      return false;
    }
    return true;
  });
  state.selected.clear();
}

function sortCards(cards) {
  const suitOrder = [4, 1, 3, 0, 2];
  return [...cards].sort((a, b) => (b.rank - a.rank) || (suitOrder.indexOf(a.suit) - suitOrder.indexOf(b.suit)));
}

function cardKey(card) {
  return `${card.suit}:${card.rank}:${card.color}`;
}

function playerName(id) {
  const player = state.players.find((item) => item.id === id);
  if (player && isBotPlayer(player)) return characterForPlayer(id)?.name || player.name || "人机";
  return player?.name || (id === state.playerId ? state.playerName : "玩家");
}

function isBotPlayer(player) {
  return Boolean(player?.isBot || String(player?.name || "").startsWith("🤖"));
}

function botYuanzi(player) {
  if (!player?.id || !isBotPlayer(player)) return 0;
  if (state.botYuanzi[player.id] == null) {
    let hash = 0;
    for (const char of player.id) hash = ((hash * 31) + char.charCodeAt(0)) >>> 0;
    state.botYuanzi[player.id] = 1000 + (hash % 4001);
  }
  return state.botYuanzi[player.id];
}

function updateBotYuanzi(scores) {
  scores.forEach((score) => {
    const player = state.players.find((item) => item.id === score.playerId);
    if (!isBotPlayer(player)) return;
    state.botYuanzi[player.id] = Math.max(0, botYuanzi(player) + (Number(score.score || 0) * 10));
  });
}

function characterForPlayer(playerId) {
  return window.ddzDialogue?.getCharacterForPlayer(playerId) || null;
}

function visualCharacterForPlayer(player, self) {
  if (self || !player?.id) return null;
  return characterForPlayer(player.id);
}

function portraitPathForPlayer(player, self, character) {
  if (character?.fullBody) return character.fullBody;
  if (character?.portrait) return character.portrait;
  if (character?.id) return `./assets/skins/sanguo/portraits/${character.id}.png`;
  return "";
}

function campForCharacter(character) {
  const value = String(character?.camp || character?.faction || "").toLowerCase();
  if (value === "wei" || value.includes("魏") || value.includes("榄")) return "wei";
  if (value === "shu" || value.includes("蜀") || value.includes("铚")) return "shu";
  if (value === "wu" || value.includes("吴") || value.includes("鍚")) return "wu";
  return "neutral";
}

function visualSeatForPlayer(self, el) {
  if (self) return "self";
  if (el?.id === "seat2") return "right";
  return "left";
}

function facingForSeat(seatName, character) {
  if (seatName === "right") return character?.facing === "front" ? "front" : "leftDown";
  if (seatName === "left") return character?.facing === "front" ? "front" : "rightDown";
  if (character?.facing) return character.facing;
  return "front";
}

function campLabel(camp) {
  return ({ wei: "\u9b4f", shu: "\u8700", wu: "\u5434", neutral: "\u7fa4" })[camp] || "\u7fa4";
}

function verticalNameHtml(name) {
  return Array.from(String(name || "\u73a9\u5bb6")).map((char) => `<span>${escapeHtml(char)}</span>`).join("");
}

function assignBotCharacters() {
  const bots = state.players.filter(isBotPlayer);
  if (bots.length >= 2) window.ddzDialogue?.assignGameCharacters?.(bots, state.roomCode || bots.map((bot) => bot.id).join(":"));
}

isBotPlayer = function(player) {
  const name = String(player?.name || "");
  return Boolean(player?.isBot || name.startsWith("🤖") || name.startsWith("馃"));
};

function isPracticeRoom() {
  return state.phase === "practice" || state.players.some(isBotPlayer);
}

function resetRoundState() {
  window.clearTimeout(practiceReplayTimer);
  state.roomCode = "";
  state.players = [];
  state.hand = [];
  state.selected.clear();
  state.bottomCards = [];
  state.lastPlayed = [];
  state.lastPlayerId = "";
  state.lastPlayerName = "";
  state.currentTurn = "";
  state.bidTurn = "";
  state.turnMustPlay = false;
  state.turnCanBeat = true;
  state.multiplier = 1;
  if (typeof ddzCardCounter !== "undefined") ddzCardCounter.reset();
}

function startPracticeMatch() {
  const replayingInSameRoom = state.phase === "game_over" && isPracticeRoom() && Boolean(state.roomCode);
  if (replayingInSameRoom) {
    window.clearTimeout(practiceReplayTimer);
    if (state.connected) send("practice_match", practiceMatchPayload());
    return;
  }
  resetRoundState();
  state.phase = "practice";
  state.playerName = state.playerName || "练习玩家";
  if (state.connected) send("practice_match", practiceMatchPayload());
  showGame();
  render();
}

function schedulePracticeReplay() {
  if (!isPracticeRoom()) return;
  window.clearTimeout(practiceReplayTimer);
  let seconds = 5;
  const tick = () => {
    if (state.phase !== "game_over" || !isPracticeRoom()) return;
    $("listPanel").innerHTML = `<div class="list-item">人机房 ${seconds} 秒后自动再来一局</div>` + $("listPanel").innerHTML;
    if (seconds <= 0) {
      startPracticeMatch();
      return;
    }
    seconds -= 1;
    practiceReplayTimer = window.setTimeout(tick, 1000);
  };
  practiceReplayTimer = window.setTimeout(tick, 500);
}

function dialogueEventForCards(cards, handType) {
  const lowerType = String(handType || "").toLowerCase();
  if (lowerType.includes("rocket") || lowerType.includes("joker")) return "playJokerBomb";
  if (lowerType.includes("bomb") || (cards?.length === 4 && new Set(cards.map((card) => card.rank)).size === 1)) return "playBomb";
  if (lowerType.includes("pair") || (cards?.length === 2 && cards[0]?.rank === cards[1]?.rank)) return "playPair";
  if (cards?.length === 1) return "playSingle";
  return "beatOpponent";
}

function showDialogue(playerId, eventName) {
  const player = state.players.find((item) => item.id === playerId);
  if (!player || !isBotPlayer(player) || !window.ddzDialogue) return "";
  const character = characterForPlayer(playerId);
  if (!character) return "";
  const line = window.ddzDialogue.getLine(character.id, eventName);
  if (!line) return "";

  const seat = [...document.querySelectorAll(".seat")].find((item) => item.dataset.playerId === playerId);
  if (!seat) return line;
  let bubble = seat.querySelector(".dialogue-bubble");
  if (!bubble) {
    bubble = document.createElement("div");
    bubble.className = "dialogue-bubble";
    seat.appendChild(bubble);
  }
  bubble.innerHTML = `<strong>${escapeHtml(character.name)}</strong><span>${escapeHtml(line)}</span>`;
  bubble.classList.add("show");

  window.clearTimeout(dialogueTimers.get(playerId));
  dialogueTimers.set(playerId, window.setTimeout(() => {
    bubble.classList.remove("show");
  }, 3000));
  return line;
}

window.showDialogue = showDialogue;

function triggerDialogue(playerId, eventName) {
  window.setTimeout(() => showDialogue(playerId, eventName), 0);
}

function setStatus(text) {
  $("statusText").textContent = text;
}

function showLobby() {
  uiState.currentView = "lobby";
  $("lobbyView").hidden = false;
  $("gameView").hidden = true;
  document.body.dataset.screen = "lobby";
}

function showGame() {
  uiState.currentView = "game";
  $("lobbyView").hidden = true;
  $("gameView").hidden = false;
  document.body.dataset.screen = "game";
  refreshPoetryLine();
}

window.showLobby = showLobby;
window.showGame = showGame;

function localPoetryLine() {
  if (poetryQueue.length === 0) {
    poetryQueue = [...poetryExpandedLines];
    for (let index = poetryQueue.length - 1; index > 0; index--) {
      const randomIndex = Math.floor(Math.random() * (index + 1));
      [poetryQueue[index], poetryQueue[randomIndex]] = [poetryQueue[randomIndex], poetryQueue[index]];
    }
    if (poetryQueue[poetryQueue.length - 1] === lastPoetryLine && poetryQueue.length > 1) {
      [poetryQueue[0], poetryQueue[poetryQueue.length - 1]] = [poetryQueue[poetryQueue.length - 1], poetryQueue[0]];
    }
  }
  lastPoetryLine = poetryQueue.pop() || poetryFallbackLines[0];
  return lastPoetryLine;
}

function normalizePoetryLine(payload) {
  const poem = payload?.data;
  const lines = Array.isArray(poem?.content) ? poem.content : [];
  const cleanLines = lines
    .map((line) => String(line || "").replace(/[（）()「」]/g, "").trim())
    .filter((line) => line && line.length <= 28 && !line.includes("注"));
  const line = cleanLines[Math.floor(Math.random() * cleanLines.length)] || localPoetryLine();
  const author = poem?.author?.name ? ` · ${poem.author.name}` : "";
  return `${line}${author}`;
}

async function refreshPoetryLine() {
  const el = $("poetryLine");
  if (!el || poetryLoading) return;
  poetryLoading = true;
  el.textContent = DEFAULT_CENTER_POETRY;

  const controller = new AbortController();
  const timeout = window.setTimeout(() => controller.abort(), 4500);
  try {
    const response = await fetch("https://poetry.palemoky.com/api/poems/random?lang=zh-Hans", {
      signal: controller.signal,
      cache: "no-store",
    });
    if (!response.ok) throw new Error(`poetry ${response.status}`);
    el.textContent = normalizePoetryLine(await response.json());
  } catch {
    el.textContent = localPoetryLine() || DEFAULT_CENTER_POETRY;
  } finally {
    window.clearTimeout(timeout);
    poetryLoading = false;
  }
}

function addLog(who, text) {
  const line = document.createElement("div");
  line.className = "log-line";
  line.textContent = `[${new Date().toLocaleTimeString()}] ${who}: ${text}`;
  $("log").appendChild(line);
  $("log").scrollTop = $("log").scrollHeight;
}

function render() {
  $("playerLabel").textContent = state.playerName ? `${state.playerName} (${state.playerId.slice(0, 6)})` : "未连接";
  $("roomText").textContent = state.roomCode || "-";
  $("onlineText").textContent = state.onlineCount ?? "-";
  const yuanziText = Number(state.yuanzi || 0).toLocaleString("zh-CN");
  $("yuanziText").textContent = yuanziText;
  $("gameYuanziText").textContent = yuanziText;
  $("latencyText").textContent = state.latency == null ? "-- ms" : `${state.latency} ms`;
  $("multiplierText").textContent = `${state.multiplier || 1}x`;
  $("roundText").textContent = state.currentTurn || state.bidTurn ? playerName(state.currentTurn || state.bidTurn) : phaseLabel(state.phase);
  const landlord = state.players.find((player) => player.isLandlord);
  $("landlordStatus").textContent = landlord ? `地主 ${landlord.name || "玩家"}` : "地主未确定";
  renderSeats();
  renderCards($("handCards"), state.hand, true);
  renderCards($("bottomCards"), state.bottomCards, false);
  renderTurnPrompt();
  renderPlayedCardsBySeat();
  renderActionButtons();
  if (typeof ddzCardCounter !== "undefined") ddzCardCounter.render();
}

function renderPlayedCardsBySeat() {
  const center = $("lastPlayedCards");
  const left = $("seat1Played");
  const right = $("seat2Played");
  renderCards(center, [], false);
  renderCards(left, [], false);
  renderCards(right, [], false);
  if (!state.lastPlayed.length || !state.lastPlayerId) return;
  if (state.lastPlayerId === state.playerId) {
    renderCards(center, state.lastPlayed, false);
    return;
  }
  const others = [...state.players]
    .filter((player) => player.id !== state.playerId)
    .sort((a, b) => Number(a.seat || 0) - Number(b.seat || 0));
  renderCards(others[0]?.id === state.lastPlayerId ? left : right, state.lastPlayed, false);
}

function renderActionButtons() {
  const isMyBidTurn = state.phase === "bidding" && state.bidTurn === state.playerId;
  const isMyPlayTurn = state.phase === "playing" && state.currentTurn === state.playerId;
  const hasBot = state.players.some(isBotPlayer);
  const setVisible = (action, visible) => {
    document.querySelectorAll(`[data-action='${action}']`).forEach((button) => {
      button.hidden = !visible;
      button.disabled = !visible;
    });
  };

  setVisible("bid_yes", isMyBidTurn);
  setVisible("bid_no", isMyBidTurn);
  setVisible("play_cards", isMyPlayTurn);
  setVisible("pass", isMyPlayTurn && !state.turnMustPlay);
  setVisible("trusteeship", state.phase === "bidding" || state.phase === "playing");
  setVisible("ready", state.phase === "waiting");
  setVisible("add_bot", state.phase === "waiting" && Boolean(state.roomCode) && state.players.length < 3);
  const botDifficultyControl = $("botDifficultyControl");
  if (botDifficultyControl) {
    botDifficultyControl.hidden = !(state.phase === "waiting" && Boolean(state.roomCode) && state.players.length < 3);
  }
  setVisible("replay", state.phase === "game_over" && Boolean(state.roomCode));
  document.querySelectorAll("[data-action='trusteeship']").forEach((button) => {
    button.textContent = state.trusteeship ? "取消托管" : "托管";
    button.classList.toggle("active", state.trusteeship);
  });

  document.querySelectorAll("[data-action='leave_room']").forEach((button) => {
    button.hidden = !state.roomCode && state.phase !== "practice" && state.phase !== "game_over";
  });
}

function phaseLabel(phase) {
  return ({
    lobby: "等待开局",
    practice: "人机练习",
    waiting: "等待准备",
    bidding: "叫地主",
    playing: "出牌阶段",
    game_over: "本局结束",
  })[phase] || "等待开局";
}

function renderTurnPrompt() {
  const banner = $("turnBanner");
  const hint = $("turnHint");
  if (!banner) return;
  if (hint) hint.hidden = true;

  const activePlayerId = state.currentTurn || state.bidTurn;
  const isMe = activePlayerId && activePlayerId === state.playerId;
  const name = activePlayerId ? playerName(activePlayerId) : "";
  banner.classList.toggle("my-turn", Boolean(isMe));
  banner.classList.toggle("waiting-turn", !activePlayerId);

  if (!activePlayerId) {
    banner.textContent = "等待玩家操作";
    hint.textContent = state.lastPlayed.length ? "当前出牌显示在中央" : "等待出牌";
    return;
  }

  if (state.phase === "bidding") {
    banner.textContent = isMe ? "轮到你叫地主" : `轮到 ${name} 叫地主`;
    hint.textContent = isMe ? "请选择：叫 / 抢地主 或 不叫" : "等待对手叫地主";
    return;
  }

  banner.textContent = isMe ? "轮到你出牌" : `轮到 ${name} 出牌`;
  hint.textContent = isMe
    ? (state.turnMustPlay ? "必须出牌，不能不出" : "可以出牌或选择不出")
    : "等待对手操作";
}

function renderSeats() {
  const sorted = [...state.players].sort((a, b) => a.seat - b.seat);
  const me = sorted.find((p) => p.id === state.playerId);
  const others = sorted.filter((p) => p.id !== state.playerId);
  renderSeat($("seat0"), me || { name: state.playerName || "我", id: state.playerId, cardsCount: state.hand.length }, true);
  renderSeat($("seat1"), others[0], false);
  renderSeat($("seat2"), others[1], false);
}

function renderSeat(el, player, self) {
  if (!player) {
    el.innerHTML = `<div><div class="player-name">空座</div><div class="player-meta">Lv. -- · 等待玩家</div></div>`;
    el.classList.remove("turn");
    el.removeAttribute("data-player-id");
    el.removeAttribute("data-visual-seat");
    el.removeAttribute("data-facing");
    el.classList.remove("camp-wei", "camp-shu", "camp-wu", "camp-neutral");
    el.style.removeProperty("--avatar-image");
    el.style.removeProperty("--portrait-image");
    el.classList.remove("has-portrait");
    return;
  }
  el.dataset.playerId = player.id || "";
  el.classList.toggle("turn", player.id && (player.id === state.currentTurn || player.id === state.bidTurn));
  const character = isBotPlayer(player) ? characterForPlayer(player.id) : null;
  const visualCharacter = visualCharacterForPlayer(player, self);
  const avatarKind = player.isLandlord ? "landlord" : (Number(player.seat || 0) % 2 === 0 ? "farmerA" : "farmerB");
  const avatarPath = character?.avatar || window.ddzSkinManager?.getAvatarPath(avatarKind);
  const portraitPath = portraitPathForPlayer(player, self, visualCharacter);
  const visualSeat = visualSeatForPlayer(self, el);
  const camp = campForCharacter(visualCharacter);
  const facing = facingForSeat(visualSeat, visualCharacter);
  el.dataset.visualSeat = visualSeat;
  el.dataset.facing = facing;
  el.classList.remove("camp-wei", "camp-shu", "camp-wu", "camp-neutral");
  el.classList.add(`camp-${camp}`);
  if (avatarPath) el.style.setProperty("--avatar-image", `url("${avatarPath}")`);
  else el.style.removeProperty("--avatar-image");
  if (portraitPath) el.style.setProperty("--portrait-image", `url("${portraitPath}")`);
  else el.style.removeProperty("--portrait-image");
  el.classList.toggle("has-portrait", Boolean(portraitPath));
  const displayName = visualCharacter?.name || character?.name || player.name || (self ? "\u6211" : "\u73a9\u5bb6");
  const cardsCount = Number(player.cardsCount || (self ? state.hand.length : 0));
  if (self) {
    el.innerHTML = `
      <div class="self-player-chip" aria-label="${escapeHtml(displayName)}"></div>
      <div class="dialogue-bubble" aria-live="polite"></div>
    `;
    return;
  }
  el.innerHTML = `
    <div class="general-plaque" aria-label="${escapeHtml(displayName)}">
      <div class="plaque-crown">${player.isLandlord ? "\ud83d\udc51" : ""}</div>
      <div class="plaque-camp">${campLabel(camp)}</div>
      <div class="plaque-name">${verticalNameHtml(displayName)}</div>
      <div class="plaque-count">${cardsCount > 0 ? `${cardsCount}\u5f20` : ""}</div>
      <div class="plaque-role">${player.isLandlord ? "\u5730\u4e3b" : ""}</div>
    </div>
    <div class="dialogue-bubble" aria-live="polite"></div>
  `;
}

function renderCards(container, cards, selectable) {
  container.innerHTML = "";
  const mid = (cards.length - 1) / 2;
  cards.forEach((card, index) => {
    const el = document.createElement("button");
    el.type = "button";
    el.className = `card ${card.color === 1 ? "red" : ""} ${card.suit === 4 ? "joker" : ""}`;
    if (selectable) {
      const fan = Math.max(-22, Math.min(22, (index - mid) * 2.2));
      el.style.setProperty("--fan", `${fan}deg`);
      el.style.zIndex = String(index + 1);
    }
    el.innerHTML = `<span class="card-rank">${rankLabel(card.rank)}</span><span class="card-suit">${suitLabel(card)}</span><span></span>`;
    const cardFacePath = window.ddzSkinManager?.getCardFacePath(card);
    if (cardFacePath) {
      el.dataset.hasFace = "true";
      el.style.setProperty("--card-face-url", `url("${cardFacePath}")`);
    }
    if (selectable) {
      el.classList.toggle("selected", state.selected.has(index));
      el.addEventListener("click", () => {
        if (state.selected.has(index)) state.selected.delete(index);
        else state.selected.add(index);
        renderCards(container, cards, selectable);
      });
    }
    container.appendChild(el);
  });
}

function rankLabel(rank) {
  return ({ 11: "J", 12: "Q", 13: "K", 14: "A", 15: "2", 16: "小王", 17: "大王" })[rank] || String(rank);
}

function suitLabel(card) {
  if (card.rank === 16) return "JOKER";
  if (card.rank === 17) return "JOKER";
  return ["♠", "♥", "♣", "♦", "★"][card.suit] || "";
}

function selectedCards() {
  return [...state.selected].sort((a, b) => a - b).map((i) => state.hand[i]).filter(Boolean);
}

function applyHintCards(cards) {
  state.selected.clear();
  const used = new Set();
  for (const hintCard of cards) {
    const index = state.hand.findIndex((card, i) =>
      !used.has(i) && card.suit === hintCard.suit && card.rank === hintCard.rank && card.color === hintCard.color
    );
    if (index >= 0) {
      used.add(index);
      state.selected.add(index);
    }
  }
  renderCards($("handCards"), state.hand, true);
}

function renderRoomList(rooms) {
  $("listPanel").innerHTML = rooms.length
    ? rooms.map((room) => `<button class="list-item" data-room="${escapeHtml(room.roomCode)}">房间 ${escapeHtml(room.roomCode)} · ${room.playerCount}/${room.maxPlayers}</button>`).join("")
    : `<div class="list-item">暂无房间</div>`;
  $("listPanel").querySelectorAll("[data-room]").forEach((btn) => {
    btn.addEventListener("click", () => send("join_room", { roomCode: btn.dataset.room }));
  });
  scheduleListPanelHide();
}

function renderLeaderboard(entries) {
  $("listPanel").innerHTML = entries.length
    ? entries.map((e) => `<div class="list-item">#${e.rank} ${escapeHtml(e.playerName)} · ${e.score} 分 · ${Number(e.winRate).toFixed(1)}%</div>`).join("")
    : `<div class="list-item">暂无排行</div>`;
  scheduleListPanelHide();
}

function renderStats(s) {
  $("listPanel").innerHTML = `<div class="list-item">${escapeHtml(s.playerName)} · ${s.score} 分 · 排名 ${s.rank || "-"}</div>
    <div class="list-item">${s.totalGames} 局，${s.wins} 胜 ${s.losses} 负，胜率 ${Number(s.winRate || 0).toFixed(1)}%</div>`;
  scheduleListPanelHide();
}

function renderResult(payload) {
  $("listPanel").innerHTML = (payload.scores || []).map((s) =>
    `<div class="list-item">${escapeHtml(playerName(s.playerId) || s.playerName)} ${s.isLandlord ? "地主" : "农民"} · ${s.score > 0 ? "+" : ""}${s.score}分 · 元子 ${s.score > 0 ? "+" : ""}${Number(s.score || 0) * 10}</div>`
  ).join("");
  scheduleListPanelHide(12000);
}

function escapeHtml(value) {
  return String(value ?? "").replace(/[&<>"']/g, (ch) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#039;" }[ch]));
}

let listPanelTimer = null;
function scheduleListPanelHide(delay = 8000) {
  window.clearTimeout(listPanelTimer);
  listPanelTimer = window.setTimeout(() => {
    $("listPanel").innerHTML = "";
  }, delay);
}

$("connectForm").addEventListener("submit", (event) => {
  event.preventDefault();
  connect($("serverInput").value.trim());
});

document.querySelectorAll("[data-action]").forEach((button) => {
  button.addEventListener("click", () => {
    const action = button.dataset.action;
    if (action === "leave_room") {
      window.clearTimeout(practiceReplayTimer);
      send("leave_room");
      state.roomCode = "";
      state.phase = "lobby";
      state.players = [];
      state.hand = [];
      state.bottomCards = [];
  state.lastPlayed = [];
  state.lastPlayerId = "";
      showLobby();
      return render();
    }
    if (action === "practice_match") {
      startPracticeMatch();
      return;
    }
    if (action === "replay") {
      if (isPracticeRoom()) startPracticeMatch();
      else send("replay");
      return;
    }
    if (action === "add_bot") {
      send("practice_match", practiceMatchPayload());
      return;
    }
    if (action === "trusteeship") {
      state.trusteeship = !state.trusteeship;
      send(state.trusteeship ? "ready" : "cancel_ready");
      addLog("托管", state.trusteeship ? "已开启托管" : "已取消托管");
      return render();
    }
    if (action === "practice_match" || action === "replay") {
      state.phase = "practice";
      state.playerName = state.playerName || "练习玩家";
      state.roomCode = "";
      state.players = [];
      state.hand = [];
      state.selected.clear();
      state.bottomCards = [];
      state.lastPlayed = [];
      state.lastPlayerId = "";
      state.lastPlayerName = "";
      state.currentTurn = "";
      state.bidTurn = "";
      if (typeof ddzCardCounter !== "undefined") ddzCardCounter.reset();
      if (state.connected) send("practice_match", practiceMatchPayload());
      showGame();
      return render();
    }
    if (action === "quick_match" || action === "create_room") {
      send(action);
      showGame();
      return render();
    }
    if (action === "join_room") return send("join_room", { roomCode: $("roomInput").value.trim() });
    if (action === "bid_yes") return send("bid", { bid: true });
    if (action === "bid_no") return send("bid", { bid: false });
    if (action === "play_cards") return send("play_cards", { cards: selectedCards() });
    if (action === "get_leaderboard") return send("get_leaderboard", { type: "total", offset: 0, limit: 20 });
    send(action);
  });
});

document.querySelector(".hint-button")?.addEventListener("click", () => {
  if (state.phase !== "playing" || state.currentTurn !== state.playerId) {
    addLog("提示", "轮到你出牌时才能提示");
    return;
  }
  send("hint");
});

$("chatForm").addEventListener("submit", (event) => {
  event.preventDefault();
  const content = $("chatInput").value.trim();
  if (!content) return;
  send("chat", { content, scope: state.roomCode ? "room" : "lobby" });
  $("chatInput").value = "";
});

$("chatToggle")?.addEventListener("click", () => {
  $("chatPanel")?.classList.toggle("collapsed");
});

$("moreMenuToggle")?.addEventListener("click", () => {
  const menu = $("moreMenu");
  const expanded = !menu?.classList.contains("expanded");
  menu?.classList.toggle("expanded", expanded);
  $("moreMenuToggle")?.setAttribute("aria-expanded", String(expanded));
});

$("lobbyToggle")?.addEventListener("click", () => {
  showLobby();
});

$("joinModeFocus")?.addEventListener("click", () => $("roomInput")?.focus());

const bgmAudio = $("bgmAudio");
const bgmToggle = $("bgmToggle");
const bgmStorageKey = "ddz.bgmEnabled";
let bgmEnabled = localStorage.getItem(bgmStorageKey) === "true";

function renderBgmState() {
  if (!bgmToggle || !bgmAudio) return;
  bgmToggle.textContent = bgmEnabled ? "BGM 开" : "BGM";
  bgmToggle.setAttribute("aria-pressed", String(bgmEnabled));
  bgmToggle.classList.toggle("active", bgmEnabled);
}

bgmToggle?.addEventListener("click", async () => {
  if (typeof ddzAudio !== "undefined") {
    ddzAudio.toggleBGM();
    bgmEnabled = ddzAudio._bgmEnabled;
    return;
  }
  if (!bgmAudio) return;
  bgmEnabled = !bgmEnabled;
  localStorage.setItem(bgmStorageKey, String(bgmEnabled));
  bgmAudio.volume = 0.42;
  if (bgmEnabled) {
    try {
      await bgmAudio.play();
    } catch {
      bgmEnabled = false;
      localStorage.setItem(bgmStorageKey, "false");
      addLog("系统", "浏览器阻止了背景音乐，请再次点击 BGM");
    }
  } else {
    bgmAudio.pause();
  }
  renderBgmState();
});

renderBgmState();

window.addEventListener("ddz:skin-change", () => render());

setInterval(() => {
  if (state.ws?.readyState === WebSocket.OPEN) send("ping", { timestamp: Date.now() });
}, 20000);

// ===== Audio Manager =====
const ddzAudio = {
  _ctx: null, _buffers: {}, _bgmAudio: null,
  _bgmEnabled: false, _sfxEnabled: true, _initialized: false,
  async init() {
    if (this._initialized) return;
    try {
      this._ctx = new (window.AudioContext || window.webkitAudioContext)();
      this._bgmAudio = document.getElementById("bgmAudio");
      this._bgmEnabled = localStorage.getItem("ddz.bgmEnabled") === "true";
      this._sfxEnabled = localStorage.getItem("ddz.sfxEnabled") !== "false";
      this._initialized = true;
    } catch(e) {}
  },
  async preload(names) {
    for (const name of names) {
      if (this._buffers[name]) continue;
      try {
        const audio = new Audio("./assets/audio" + name);
        audio.preload = "auto";
        audio.volume = 0.65;
        this._buffers[name] = audio;
      } catch(e) {}
    }
  },
  play(name) {
    if (!this._sfxEnabled) return;
    const audio = this._buffers[name] || new Audio("./assets/audio" + name);
    try {
      audio.currentTime = 0;
      audio.volume = 0.65;
      audio.play().catch(function(){});
    } catch(e) {}
  },
  setBGM(track) {
    const bgm = this._bgmAudio;
    if (!bgm || !this._bgmEnabled) return;
    const m = {lobby:"./assets/audio/bgm_welcome.mp3",playing:"./assets/audio/bgm/bgm_normal.mp3"};
    const src = m[track] || m.lobby;
    if (bgm.src && bgm.src.indexOf(src) >= 0) return;
    bgm.src = src;
    bgm.volume = 0.42;
    bgm.loop = true;
    if (this._bgmEnabled) bgm.play().catch(function(){});
  },
  toggleBGM() {
    this._bgmEnabled = !this._bgmEnabled;
    bgmEnabled = this._bgmEnabled;
    localStorage.setItem("ddz.bgmEnabled", String(this._bgmEnabled));
    const bgm = this._bgmAudio;
    if (!bgm) return;
    if (this._bgmEnabled) { bgm.volume = 0.42; bgm.loop = true; bgm.play().catch(function(){}); }
    else { bgm.pause(); }
    renderBgmState();
  },
};

function voiceGender(playerId) {
  const characterGender = window.ddzDialogue?.getGenderForPlayer?.(playerId);
  if (characterGender === "male" || characterGender === "female") return characterGender;
  let hash = 0;
  for (const char of String(playerId || state.playerId)) hash = ((hash << 5) - hash + char.charCodeAt(0)) | 0;
  return Math.abs(hash) % 2 === 0 ? "male" : "female";
}

function voiceRank(rank) {
  return ({11:"J",12:"Q",13:"K",14:"A",15:"2",16:"joker_small",17:"joker_big"})[rank] || String(rank);
}

function randomVoice(folder, names) {
  const available = names.filter(Boolean);
  if (!available.length) return "";
  return `${folder}/${available[Math.floor(Math.random() * available.length)]}`;
}

function cardVoicePath(playerId, cards, handType) {
  const folder = `/voices/${voiceGender(playerId)}`;
  const type = String(handType || "");
  const rank = voiceRank(cards?.[0]?.rank);
  if (type === "单张" && rank) return `${folder}/single/single_${rank}.mp3`;
  if (type === "对子" && rank) return `${folder}/pair/pair_${rank}.mp3`;
  if (type === "三张" && rank) return `${folder}/trio/trio_${rank}.mp3`;
  const typeVoices = {
    "三带一": "type_trio_single.mp3",
    "三带二": "type_trio_pair.mp3",
    "顺子": "type_straight.mp3",
    "连对": "type_pairstraight.mp3",
    "飞机": "type_plane.mp3",
    "飞机带单": "type_plane.mp3",
    "飞机带对": "type_plane.mp3",
    "炸弹": "type_bomb.mp3",
    "王炸": "type_rocket.mp3",
    "四带二": "type_four_two.mp3",
    "四带两对": "type_four_twopair.mp3",
  };
  return typeVoices[type] ? `${folder}/type/${typeVoices[type]}` : randomVoice(folder, ["beat/beat.mp3", "beat/beat_bigger.mp3", "beat/beat_cover.mp3"]);
}

function bidVoicePath(payload) {
  const folder = `/voices/${voiceGender(payload.playerId)}/bid`;
  if (!payload.bid) return `${folder}/${payload.isGrab ? "bid_nograb.mp3" : "bid_nocall.mp3"}`;
  if (!payload.isGrab) return `${folder}/bid_call.mp3`;
  return randomVoice(folder, ["bid_grab.mp3", "bid_grab2.mp3", "bid_grab3.mp3"]);
}

function passVoicePath(playerId) {
  const gender = voiceGender(playerId);
  const folder = `/voices/${gender}/pass`;
  return randomVoice(folder, gender === "male"
    ? ["pass.mp3", "pass_buyao.mp3", "pass_guo.mp3", "pass_peng.mp3"]
    : ["pass_buyao.mp3", "pass_guo.mp3", "pass_kan.mp3", "pass_yaobuqi.mp3"]);
}

// ===== Card Counter =====
const ddzCardCounter = {
  _counts: {}, _visible: false, _dragReady: false,
  reset() {
    for (let r = 3; r <= 15; r++) this._counts[r] = 4;
    this._counts[16] = 1; this._counts[17] = 1;
  },
  syncFromHand(cards) {
    this.reset();
    this.deduct(cards);
  },
  deduct(cards) {
    if (!cards) return;
    for (const c of cards) {
      if (Number.isInteger(c?.rank) && this._counts[c.rank] > 0) this._counts[c.rank]--;
    }
    this.saveCounts();
  },
  saveCounts() {
    if (!state.roomCode) return;
    try {
      localStorage.setItem(`ddz.cardCounterCounts.${state.roomCode}`, JSON.stringify(this._counts));
    } catch(e) {}
  },
  restoreCounts(roomCode) {
    if (!roomCode) return false;
    try {
      const saved = JSON.parse(localStorage.getItem(`ddz.cardCounterCounts.${roomCode}`) || "null");
      if (!saved) return false;
      for (let rank = 3; rank <= 17; rank++) {
        const maxCount = rank <= 15 ? 4 : 1;
        if (!Number.isInteger(saved[rank]) || saved[rank] < 0 || saved[rank] > maxCount) return false;
      }
      this._counts = saved;
      return true;
    } catch(e) {
      return false;
    }
  },
  render() {
    const el = document.getElementById("cardCounter");
    if (!el) return;
    this.initDrag(el);
    if (!this._visible || state.phase === "lobby") { el.style.display = "none"; return; }
    el.style.display = "";
    const labels = {3:"3",4:"4",5:"5",6:"6",7:"7",8:"8",9:"9",10:"10",11:"J",12:"Q",13:"K",14:"A",15:"2",16:"\u5c0f\u738b",17:"\u5927\u738b"};
    let html = "<div class=\"cc-head\"><strong>记牌器</strong><span>余牌 / 手牌</span></div><div class=\"cc-grid\">";
    for (let r = 3; r <= 17; r++) {
      const left = this._counts[r] || 0;
      const cls = left >= 3 ? "cc-ok" : (left >= 1 ? "cc-warn" : "cc-gone");
      const hold = state.hand.filter(function(c) { return c.rank === r; }).length;
      html += "<div class=\"cc-cell " + cls + (hold > 0 ? " cc-hold" : "") + "\">"
        + "<span class=\"cc-rank\">" + (labels[r]||r) + "</span>"
        + "<span class=\"cc-left\">" + left + "</span>"
        + "<span class=\"cc-held\">" + hold + "</span>"
        + "</div>";
    }
    html += "</div>";
    el.innerHTML = html;
    this.applySavedPosition(el);
  },
  toggle() {
    this._visible = !this._visible;
    localStorage.setItem("ddz.cardCounter", String(this._visible));
    this.render();
  },
  applySavedPosition(el) {
    if (!el || el.dataset.positionApplied === "1") return;
    el.dataset.positionApplied = "1";
    try {
      const pos = JSON.parse(localStorage.getItem("ddz.cardCounterPos") || "null");
      if (!pos || !Number.isFinite(pos.x) || !Number.isFinite(pos.y)) return;
      const maxX = Math.max(8, window.innerWidth - el.offsetWidth - 8);
      const maxY = Math.max(8, window.innerHeight - el.offsetHeight - 8);
      el.style.left = Math.min(Math.max(8, pos.x), maxX) + "px";
      el.style.top = Math.min(Math.max(8, pos.y), maxY) + "px";
      el.style.right = "auto";
    } catch(e) {}
  },
  savePosition(el) {
    if (!el) return;
    try {
      localStorage.setItem("ddz.cardCounterPos", JSON.stringify({
        x: Math.round(el.offsetLeft),
        y: Math.round(el.offsetTop),
      }));
    } catch(e) {}
  },
  initDrag(el) {
    if (!el || this._dragReady) return;
    this._dragReady = true;
    let dragging = false;
    let startX = 0;
    let startY = 0;
    let baseX = 0;
    let baseY = 0;
    const move = (event) => {
      if (!dragging) return;
      const point = event.touches?.[0] || event;
      const width = el.offsetWidth;
      const height = el.offsetHeight;
      const nextX = Math.min(Math.max(8, baseX + point.clientX - startX), Math.max(8, window.innerWidth - width - 8));
      const nextY = Math.min(Math.max(8, baseY + point.clientY - startY), Math.max(8, window.innerHeight - height - 8));
      el.style.left = nextX + "px";
      el.style.top = nextY + "px";
      el.style.right = "auto";
      event.preventDefault?.();
    };
    const end = () => {
      if (!dragging) return;
      dragging = false;
      el.classList.remove("dragging");
      this.savePosition(el);
      window.removeEventListener("pointermove", move);
      window.removeEventListener("pointerup", end);
      window.removeEventListener("touchmove", move);
      window.removeEventListener("touchend", end);
    };
    el.addEventListener("pointerdown", (event) => {
      if (!event.target.closest(".cc-head")) return;
      dragging = true;
      startX = event.clientX;
      startY = event.clientY;
      baseX = el.offsetLeft;
      baseY = el.offsetTop;
      el.classList.add("dragging");
      window.addEventListener("pointermove", move, { passive: false });
      window.addEventListener("pointerup", end);
      event.preventDefault();
    });
    el.addEventListener("touchstart", (event) => {
      if (!event.target.closest(".cc-head")) return;
      const point = event.touches[0];
      dragging = true;
      startX = point.clientX;
      startY = point.clientY;
      baseX = el.offsetLeft;
      baseY = el.offsetTop;
      el.classList.add("dragging");
      window.addEventListener("touchmove", move, { passive: false });
      window.addEventListener("touchend", end);
      event.preventDefault();
    }, { passive: false });
  },
};

ddzCardCounter.reset();
try { ddzCardCounter._visible = localStorage.getItem("ddz.cardCounter") === "true"; } catch(e) {}

ddzAudio.init().then(function() {
  ddzAudio.preload(["/effects/deal.mp3","/effects/play.mp3","/effects/win.mp3","/effects/lose.mp3","/effects/reveal.mp3","/effects/bomb.mp3","/voices/male/bid/bid_call.mp3","/voices/male/bid/bid_grab.mp3","/voices/male/bid/bid_nocall.mp3","/voices/male/pass/pass_buyao.mp3"]);
}).catch(function(){});

// Update BGM toggle to use ddzAudio
bgmEnabled = ddzAudio._bgmEnabled;
renderBgmState();

// Wrap handleMessage to add sounds and card counter
{
  const orig = handleMessage;
  handleMessage = function(type, payload) {
    if (type === "deal_cards") {
      ddzAudio.play("/effects/deal.mp3");
    }
    if (type === "card_played") {
      ddzAudio.play("/effects/play.mp3");
      ddzAudio.play(cardVoicePath(payload.playerId, payload.cards, payload.handType));
      if (payload.handType === "炸弹" || payload.handType === "王炸") ddzAudio.play("/effects/bomb.mp3");
      if (payload.cardsLeft === 1 || payload.cardsLeft === 2) {
        ddzAudio.play(`/voices/${voiceGender(payload.playerId)}/alert/last${payload.cardsLeft}.mp3`);
      }
      if (payload.playerId !== state.playerId) { ddzCardCounter.deduct(payload.cards); ddzCardCounter.render(); }
    }
    if (type === "bid_result") { ddzAudio.play(bidVoicePath(payload)); }
    if (type === "landlord") { ddzAudio.play("/effects/reveal.mp3"); ddzAudio.setBGM("playing"); ddzCardCounter.render(); }
    if (type === "player_pass") { ddzAudio.play(passVoicePath(payload.playerId)); }
    if (type === "game_over") { ddzAudio.play(payload.winnerId===state.playerId?"/effects/win.mp3":"/effects/lose.mp3"); ddzAudio.setBGM("lobby"); }
    if (type === "game_start") { ddzAudio.setBGM("playing"); }
    
    orig(type, payload);

    if (type === "deal_cards" || (type === "landlord" && payload.playerId === state.playerId)) {
      ddzCardCounter.syncFromHand(state.hand);
    }
    if (type === "reconnected" && !ddzCardCounter.restoreCounts(state.roomCode)) {
      ddzCardCounter.syncFromHand(state.hand);
    }
    
    if (["deal_cards","card_played","player_pass","game_over","game_start","landlord","play_turn","reconnected"].indexOf(type) >= 0) {
      ddzCardCounter.render();
    }
  };
}

// Counter & SFX toggle buttons handler
document.addEventListener("DOMContentLoaded", function() {
  setTimeout(function() {
    var ct = document.getElementById("counterToggle");
    if (ct) {
      ct.addEventListener("click", function() { ddzCardCounter.toggle(); this.classList.toggle("active", ddzCardCounter._visible); });
      ct.classList.toggle("active", ddzCardCounter._visible);
    }
    var st = document.getElementById("sfxToggle");
    if (st) {
      st.addEventListener("click", function() { ddzAudio._sfxEnabled = !ddzAudio._sfxEnabled; localStorage.setItem("ddz.sfxEnabled", String(ddzAudio._sfxEnabled)); this.textContent = ddzAudio._sfxEnabled ? "\ud83d\udd0a" : "\ud83d\udd07"; this.classList.toggle("active", ddzAudio._sfxEnabled); });
      st.textContent = ddzAudio._sfxEnabled ? "\ud83d\udd0a" : "\ud83d\udd07";
      st.classList.toggle("active", ddzAudio._sfxEnabled);
    }
  }, 500);
});


showLobby();
render();
refreshPoetryLine();

if (!developerMode) {
  const defaultServerUrl = $("serverInput")?.value.trim();
  if (defaultServerUrl) connect(defaultServerUrl);
}
