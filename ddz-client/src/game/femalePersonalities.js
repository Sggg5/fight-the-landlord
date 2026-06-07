(function () {
  const femalePersonalities = [
    {
      id: "sunshangxiang", name: "孙尚香", title: "弓腰姬", faction: "吴", gender: "female",
      personality: "英武果断，敢争敢抢，行动凌厉。",
      avatar: "./assets/skins/sanguo/avatars/characters/sunshangxiang.png", fullBody: "./assets/skins/sanguo/portraits/sunshangxiang.png", seat: "left", facing: "rightDown", camp: "wu", riskLevel: 82,
      lines: {
        gameStart: ["江东女儿，亦可决胜牌桌。"], callLandlord: ["地主之位，我来取。"],
        passLandlord: ["此番先让，莫以为我怯。"], playSingle: ["一箭先发。"],
        playPair: ["双矢齐出。"], playBomb: ["看我破阵！"], playJokerBomb: ["双王在手，谁敢争锋。"],
        beatOpponent: ["这一手，你挡不住。"], teammateAlmostWin: ["放心前行，我替你压阵。"],
        selfAlmostWin: ["胜负将定。"], win: ["江东儿女，从不让须眉。"],
        lose: ["再来一局，我必雪耻。"], pass: ["暂且不出。"],
      },
    },
    {
      id: "diaochan", name: "貂蝉", title: "闭月佳人", faction: "群", gender: "female",
      personality: "优雅聪慧，善察人心，出手难测。",
      avatar: "./assets/skins/sanguo/avatars/characters/diaochan.png", fullBody: "./assets/skins/sanguo/portraits/diaochan.png", seat: "left", facing: "rightDown", camp: "neutral", riskLevel: 58,
      lines: {
        gameStart: ["月下牌局，诸位可莫分心。"], callLandlord: ["此局，便由妾身执掌。"],
        passLandlord: ["不争一时，静候佳机。"], playSingle: ["轻落一子。"],
        playPair: ["双花并蒂。"], playBomb: ["惊鸿一舞，满座皆静。"], playJokerBomb: ["闭月之下，胜负已定。"],
        beatOpponent: ["公子这一手，妾身看破了。"], teammateAlmostWin: ["良机已至，莫要迟疑。"],
        selfAlmostWin: ["只差最后一步。"], win: ["承让了。"],
        lose: ["此局失算，来日再会。"], pass: ["这一手，妾身不要。"],
      },
    },
    {
      id: "daqiao", name: "大乔", title: "江东国色", faction: "吴", gender: "female",
      personality: "温婉沉稳，善守局势，后发制人。",
      avatar: "./assets/skins/sanguo/avatars/characters/daqiao.png", fullBody: "./assets/skins/sanguo/portraits/daqiao.png", seat: "left", facing: "rightDown", camp: "wu", riskLevel: 38,
      lines: {
        gameStart: ["静水流深，此局宜稳。"], callLandlord: ["既有良牌，我便担下。"],
        passLandlord: ["此时不争，方为上策。"], playSingle: ["且试一张。"],
        playPair: ["双桥相映。"], playBomb: ["不得已，只能破局。"], playJokerBomb: ["江潮既起，无人可挡。"],
        beatOpponent: ["此路不通。"], teammateAlmostWin: ["稳住阵脚，胜机已近。"],
        selfAlmostWin: ["收官在即。"], win: ["稳中求胜，方得长久。"],
        lose: ["胜负有常，不必介怀。"], pass: ["暂避其锋。"],
      },
    },
    {
      id: "xiaoqiao", name: "小乔", title: "江东灵姝", faction: "吴", gender: "female",
      personality: "灵动俏皮，善抓机会，节奏多变。",
      avatar: "./assets/skins/sanguo/avatars/characters/xiaoqiao.png", fullBody: "./assets/skins/sanguo/portraits/xiaoqiao.png", seat: "left", facing: "rightDown", camp: "wu", riskLevel: 68,
      lines: {
        gameStart: ["牌局开始，可别小看我。"], callLandlord: ["地主让我来当！"],
        passLandlord: ["先让你们争，我看热闹。"], playSingle: ["这一张，接好啦。"],
        playPair: ["成双成对。"], playBomb: ["吓你们一跳！"], playJokerBomb: ["双王驾到，快让路。"],
        beatOpponent: ["嘻嘻，被我压住了。"], teammateAlmostWin: ["快走快走，我来掩护。"],
        selfAlmostWin: ["马上就赢啦。"], win: ["我就知道会赢。"],
        lose: ["不算不算，再来嘛。"], pass: ["这手不要。"],
      },
    },
    {
      id: "zhenji", name: "甄姬", title: "洛水仙姝", faction: "魏", gender: "female",
      personality: "清冷克制，判断精准，不轻易冒险。",
      avatar: "./assets/skins/sanguo/avatars/characters/zhenji.png", fullBody: "./assets/skins/sanguo/portraits/zhenji.png", seat: "left", facing: "rightDown", camp: "wei", riskLevel: 45,
      lines: {
        gameStart: ["洛水无声，牌势自明。"], callLandlord: ["此局可控，我来。"],
        passLandlord: ["时机未至，不必强求。"], playSingle: ["一叶落水。"],
        playPair: ["双影凌波。"], playBomb: ["水势骤起。"], playJokerBomb: ["洛神临世，诸位退让。"],
        beatOpponent: ["你的牌路，已尽。"], teammateAlmostWin: ["继续前行，我守后路。"],
        selfAlmostWin: ["终局将至。"], win: ["水到渠成。"],
        lose: ["此局已过，不必多言。"], pass: ["过。"],
      },
    },
    {
      id: "huangyueying", name: "黄月英", title: "机关奇才", faction: "蜀", gender: "female",
      personality: "理性聪颖，善于计算牌势与剩余资源。",
      avatar: "./assets/skins/sanguo/avatars/characters/huangyueying.png", fullBody: "./assets/skins/sanguo/portraits/huangyueying.png", seat: "left", facing: "rightDown", camp: "shu", riskLevel: 52,
      lines: {
        gameStart: ["牌数已明，可以开始推演。"], callLandlord: ["胜率可观，此位我接。"],
        passLandlord: ["数据不足，暂不争抢。"], playSingle: ["先投入一枚棋子。"],
        playPair: ["双机联动。"], playBomb: ["机关启动，破阵。"], playJokerBomb: ["最终机关，已经完成。"],
        beatOpponent: ["你的应对，在计算之内。"], teammateAlmostWin: ["路径已清，直接收官。"],
        selfAlmostWin: ["只剩最后一步验证。"], win: ["推演结果，毫无偏差。"],
        lose: ["参数有误，重新计算。"], pass: ["保留资源。"],
      },
    },
  ];

  const dialogue = window.ddzDialogue;
  if (!dialogue) return;
  femalePersonalities.forEach((character) => {
    dialogue.personalities.push(character);
    dialogue.byId[character.id] = character;
  });

  const assignments = new Map();
  const males = dialogue.personalities.filter((character) => character.gender !== "female");
  const females = dialogue.personalities.filter((character) => character.gender === "female");

  function hashText(text) {
    let hash = 0;
    for (const char of String(text || "")) hash = ((hash << 5) - hash + char.charCodeAt(0)) | 0;
    return Math.abs(hash);
  }

  dialogue.assignGameCharacters = function (botPlayers, seed) {
    assignments.clear();
    const players = [...(botPlayers || [])].sort((a, b) => Number(a.seat || 0) - Number(b.seat || 0));
    if (players[0]) assignments.set(players[0].id, males[hashText(`${seed}:male`) % males.length]);
    if (players[1]) assignments.set(players[1].id, females[hashText(`${seed}:female`) % females.length]);
  };

  dialogue.getCharacterForPlayer = function (playerId) {
    return assignments.get(playerId) || dialogue.personalities[hashText(playerId) % dialogue.personalities.length];
  };

  dialogue.getGenderForPlayer = function (playerId) {
    return dialogue.getCharacterForPlayer(playerId)?.gender || "male";
  };

  dialogue.femalePersonalities = femalePersonalities;
})();
