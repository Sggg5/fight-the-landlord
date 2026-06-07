export type DialogueEvent =
  | "gameStart"
  | "callLandlord"
  | "passLandlord"
  | "playSingle"
  | "playPair"
  | "playBomb"
  | "playJokerBomb"
  | "beatOpponent"
  | "teammateAlmostWin"
  | "selfAlmostWin"
  | "win"
  | "lose"
  | "pass";

export interface PersonalityLines {
  gameStart: string[];
  callLandlord: string[];
  passLandlord: string[];
  playSingle: string[];
  playPair: string[];
  playBomb: string[];
  playJokerBomb: string[];
  beatOpponent: string[];
  teammateAlmostWin: string[];
  selfAlmostWin: string[];
  win: string[];
  lose: string[];
  pass: string[];
}

export interface CharacterPersonality {
  id: string;
  name: string;
  title: string;
  faction: "蜀" | "魏" | "晋";
  personality: string;
  avatar: string;
  riskLevel: number;
  lines: PersonalityLines;
}

export const personalities: CharacterPersonality[] = [
  {
    id: "liubei",
    name: "刘备",
    title: "昭烈仁君",
    faction: "蜀",
    personality: "仁义稳重，重视同伴，不轻易冒进。",
    avatar: "./assets/skins/sanguo/avatars/characters/liubei.png",
    riskLevel: 35,
    lines: {
      gameStart: ["诸君同心，此局当以稳为先。", "牌局虽小，亦见仁义。"],
      callLandlord: ["民心所向，备愿担此重任。", "既有良机，备便一试。"],
      passLandlord: ["此时不争，静观其变。", "让贤于诸位，备且守势。"],
      playSingle: ["先出一子，以探虚实。", "小牌铺路，莫急。"],
      playPair: ["双骑并进。", "两翼齐出，稳住阵脚。"],
      playBomb: ["不得已，只能雷霆一击。", "此炸非为逞强，为护全局。"],
      playJokerBomb: ["天命在此，诸君退让。", "双王齐至，局势已明。"],
      beatOpponent: ["承让，此手我接下了。", "以和为贵，但牌不可让。"],
      teammateAlmostWin: ["同伴将成，备当相助。", "再稳一步，胜机已近。"],
      selfAlmostWin: ["胜负将定，仍不可骄。", "只差一步，诸君共勉。"],
      win: ["此胜非备一人之功。", "仁者不争，争则必成。"],
      lose: ["胜败常事，备再修德行。", "此局有失，来日再战。"],
      pass: ["暂且不出。", "让一手，未必是退。"],
    },
  },
  {
    id: "guanyu",
    name: "关羽",
    title: "武圣云长",
    faction: "蜀",
    personality: "骄傲强势，出手果断，压制感强。",
    avatar: "./assets/skins/sanguo/avatars/characters/guanyu.png",
    riskLevel: 72,
    lines: {
      gameStart: ["关某在此，尔等且看。", "青龙偃月，今日试牌。"],
      callLandlord: ["地主之位，关某取了。", "此局当由关某统兵。"],
      passLandlord: ["区区地主，暂且让你。", "关某不屑此时相争。"],
      playSingle: ["一骑足矣。", "单刀赴会。"],
      playPair: ["双刀并斩。", "两路齐下。"],
      playBomb: ["看我破阵。", "此炸一出，谁敢争锋。"],
      playJokerBomb: ["双王压阵，万军辟易。", "尔等可服？"],
      beatOpponent: ["此牌，也配挡我？", "关某接了。"],
      teammateAlmostWin: ["速战速决，莫负良机。", "兄弟将胜，关某护阵。"],
      selfAlmostWin: ["胜势已成。", "再一手，定乾坤。"],
      win: ["关某胜之，理所当然。", "云长在此，何惧群雄。"],
      lose: ["胜败暂论，关某不服。", "来日再战，必雪此局。"],
      pass: ["不出。", "此手暂让。"],
    },
  },
  {
    id: "zhangfei",
    name: "张飞",
    title: "燕人翼德",
    faction: "蜀",
    personality: "暴躁豪爽，喜欢强攻，气势压人。",
    avatar: "./assets/skins/sanguo/avatars/characters/zhangfei.png",
    riskLevel: 88,
    lines: {
      gameStart: ["俺张飞来也，谁敢一战。", "痛快些，别磨磨蹭蹭。"],
      callLandlord: ["地主俺来当。", "都闪开，俺要抢。"],
      passLandlord: ["哼，这回先饶你们。", "不抢不抢，俺等大牌。"],
      playSingle: ["先来一张。", "小的也能吓你一跳。"],
      playPair: ["一对，接好了。", "两张一起上。"],
      playBomb: ["炸你个痛快。", "轰开阵脚。"],
      playJokerBomb: ["王炸来了，怕不怕。", "哈哈，这才叫痛快。"],
      beatOpponent: ["压住你了。", "就这也敢出？"],
      teammateAlmostWin: ["兄弟快走，俺挡着。", "快赢了，别怂。"],
      selfAlmostWin: ["俺马上就赢。", "再来一手，收工。"],
      win: ["哈哈，痛快。", "俺老张赢了。"],
      lose: ["不算不算，再来。", "俺刚才手滑了。"],
      pass: ["不要。", "这手俺先忍了。"],
    },
  },
  {
    id: "zhugeliang",
    name: "诸葛亮",
    title: "卧龙军师",
    faction: "蜀",
    personality: "冷静谋略，善于观察节奏和牌势。",
    avatar: "./assets/skins/sanguo/avatars/characters/zhugeliang.png",
    riskLevel: 48,
    lines: {
      gameStart: ["观此牌势，已有三分端倪。", "且看风向，再定奇谋。"],
      callLandlord: ["亮愿借东风一用。", "此局可谋，亮便接下。"],
      passLandlord: ["时机未至，不必强争。", "此位让出，另有后手。"],
      playSingle: ["投石问路。", "一子落，局势明。"],
      playPair: ["双子成势。", "此对正合阵法。"],
      playBomb: ["火计已成。", "借一声惊雷。"],
      playJokerBomb: ["天时地利，皆在掌中。", "双王落定，胜负可知。"],
      beatOpponent: ["此计已被亮识破。", "先生此手，亮可解。"],
      teammateAlmostWin: ["护住此势，胜机将至。", "同伴将成，不可轻失。"],
      selfAlmostWin: ["只差最后一计。", "收官之势已定。"],
      win: ["不过略施小计。", "此局，尽在推演之中。"],
      lose: ["谋事在人，成事在天。", "此局失算，亮记下了。"],
      pass: ["不急。", "暂避其锋。"],
    },
  },
  {
    id: "caocao",
    name: "曹操",
    title: "魏武孟德",
    faction: "魏",
    personality: "多疑霸气，敢压敢控，喜欢掌握主动。",
    avatar: "./assets/skins/sanguo/avatars/characters/caocao.png",
    riskLevel: 78,
    lines: {
      gameStart: ["宁教我控牌，不教牌控我。", "群雄在座，且看孟德手段。"],
      callLandlord: ["地主之权，孤取了。", "此局归孤调度。"],
      passLandlord: ["孤暂不争，且看尔等。", "疑兵之计，不必明说。"],
      playSingle: ["一牌足以试胆。", "孤先落子。"],
      playPair: ["成双入阵。", "此对正可压势。"],
      playBomb: ["乱世当用重典。", "孤一炸，群雄皆静。"],
      playJokerBomb: ["王权在手，谁敢不从。", "双王既出，天下归心。"],
      beatOpponent: ["想压孤？还早。", "此路已被孤断。"],
      teammateAlmostWin: ["快收，孤替你镇场。", "胜机在前，不可迟疑。"],
      selfAlmostWin: ["大势归孤。", "最后一手，天下可定。"],
      win: ["天下英雄，唯孤与诸君耳。", "此胜，孤笑纳了。"],
      lose: ["此败有诈。", "孤输一局，不输天下。"],
      pass: ["不出。", "孤且观望。"],
    },
  },
  {
    id: "simayi",
    name: "司马懿",
    title: "冢虎仲达",
    faction: "晋",
    personality: "隐忍阴狠，少言冷静，等待致命机会。",
    avatar: "./assets/skins/sanguo/avatars/characters/simayi.png",
    riskLevel: 62,
    lines: {
      gameStart: ["静待其变。", "牌局如朝局，忍者胜。"],
      callLandlord: ["既如此，便由我执局。", "此位，我收下。"],
      passLandlord: ["不急，先让他们争。", "锋芒太露，非智者所为。"],
      playSingle: ["一枚暗子。", "落子无声。"],
      playPair: ["双线并行。", "此对，恰到好处。"],
      playBomb: ["忍到此刻，正为这一击。", "惊雷之后，万籁俱寂。"],
      playJokerBomb: ["双王既现，局已入瓮。", "此刻，才是终局。"],
      beatOpponent: ["你已入局。", "此手，我等很久了。"],
      teammateAlmostWin: ["莫急，我会替你挡住。", "胜机已近，静收即可。"],
      selfAlmostWin: ["终于等到了。", "收网。"],
      win: ["能忍者，终得天下。", "此局，不过迟早。"],
      lose: ["输一局，无妨。", "我还活着，便未败。"],
      pass: ["过。", "让他先走。"],
    },
  },
];

export const personalityById = Object.fromEntries(personalities.map((item) => [item.id, item]));
