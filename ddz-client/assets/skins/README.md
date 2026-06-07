皮肤资源目录说明

新增皮肤时复制一个现有目录，并在 `src/skin/skinConfig.js` 增加配置。

推荐结构：

- `backgrounds/battlefield.jpg`：牌桌背景或夜景背景
- `backgrounds/table.png`：桌面/牌桌面板
- `ui/panel.png`：通用面板
- `ui/avatar-frame.png`：玩家头像框
- `ui/chat-panel.png`：聊天抽屉面板
- `ui/btn-primary.png`：主按钮
- `ui/btn-secondary.png`：次按钮
- `ui/choose-glow.png`：选中牌高亮特效
- `cards/back.png`：牌背
- `cards/{suit}-{rank}.png`：牌面

没有图片时，当前 CSS 变量会自动使用兜底样式。

启用图片资源时，在 `src/skin/skinConfig.js` 里把对应 `assets` 开关改成 `true`，并填写路径：

- `avatarFrame`
- `chatPanel`
- `chooseGlow`
