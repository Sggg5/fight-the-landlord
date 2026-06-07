(function () {
  const storageKey = "ddz.currentSkin";
  const presetKey = "ddz.skinPreset.v3";
  const configs = Array.isArray(window.DDZ_SKIN_CONFIGS) ? window.DDZ_SKIN_CONFIGS : [];
  const fallbackSkinId = configs[0]?.id || "classic";

  function toCssVarName(key) {
    return `--${key.replace(/[A-Z]/g, (match) => `-${match.toLowerCase()}`)}`;
  }

  function cssUrl(path) {
    return path ? `url("${path}")` : "none";
  }

  function findSkin(id) {
    return configs.find((skin) => skin.id === id) || configs.find((skin) => skin.id === fallbackSkinId) || configs[0];
  }

  function applyCssVariables(skin) {
    if (!skin) return;

    const root = document.documentElement;
    const colors = skin.themeColors || {};
    Object.entries(colors).forEach(([key, value]) => {
      root.style.setProperty(toCssVarName(key), value);
    });

    root.style.setProperty("--bg-image", cssUrl(skin.background));
    root.style.setProperty("--card-back", cssUrl(skin.assets?.cardBack ? skin.cardBack : ""));
    root.style.setProperty("--avatar-frame", cssUrl(skin.assets?.avatarFrame ? skin.avatarFrame : ""));
    root.style.setProperty("--chat-panel-image", cssUrl(skin.assets?.chatPanel ? skin.chatPanel : ""));
    root.style.setProperty("--choose-glow-image", cssUrl(skin.assets?.chooseGlow ? skin.chooseGlow : ""));
    root.style.setProperty("--button-radius", skin.buttonStyle?.radius || "0");
    root.style.setProperty("--button-border-width", skin.buttonStyle?.borderWidth || "1px");
    root.style.setProperty("--button-shadow", skin.buttonStyle?.shadow || "none");
    root.style.setProperty("--button-sprite", cssUrl(skin.buttonStyle?.sprite));
    root.style.setProperty("--table-radius", skin.tableStyle?.radius || "0");
    root.style.setProperty("--table-border-width", skin.tableStyle?.borderWidth || "1px");
    root.style.setProperty("--table-texture", cssUrl(skin.tableStyle?.texture));
    const generated = skin.generatedAssets || {};
    [
      "lobbyBg", "modeTaoyuan", "modeGuandu", "modeChibi", "modeQunxiong",
      "aiPracticeGeneral", "friendBattleGeneral", "rank", "record", "shop",
      "task", "gift", "settings",
    ].forEach((key) => root.style.setProperty(toCssVarName(`generated${key[0].toUpperCase()}${key.slice(1)}`), cssUrl(generated[key])));
    renderGeneratedImages(skin);
    root.dataset.skin = skin.id;
  }

  function GeneratedImage({ src, className = "", alt = "" }) {
    const image = document.createElement("img");
    image.src = src;
    image.className = className;
    image.alt = alt;
    image.loading = "lazy";
    image.draggable = false;
    return image;
  }

  function renderGeneratedImages(skin) {
    const generated = skin.generatedAssets || {};
    document.querySelectorAll("[data-generated-asset]").forEach((host) => {
      const src = generated[host.dataset.generatedAsset] || "";
      host.replaceChildren();
      host.classList.toggle("has-generated-image", Boolean(src));
      if (src) host.appendChild(GeneratedImage({ src, className: "generated-image", alt: "" }));
    });
  }

  function getCurrentSkinId() {
    return localStorage.getItem(storageKey) || fallbackSkinId;
  }

  function getCurrentSkin() {
    return findSkin(getCurrentSkinId());
  }

  function setSkin(id) {
    const skin = findSkin(id);
    if (!skin) return null;
    localStorage.setItem(storageKey, skin.id);
    applyCssVariables(skin);
    window.dispatchEvent(new CustomEvent("ddz:skin-change", { detail: { skin } }));
    return skin;
  }

  function initSkinSelect(selectId = "skinSelect") {
    const select = document.getElementById(selectId);
    if (!select) return;

    select.innerHTML = configs
      .map((skin) => `<option value="${skin.id}">${skin.name}</option>`)
      .join("");
    select.value = getCurrentSkin()?.id || fallbackSkinId;
    select.addEventListener("change", () => setSkin(select.value));
  }

  function getCardFacePath(card) {
    const skin = getCurrentSkin();
    if (!skin?.assets?.cardFaces || !skin.cardFacePath || !card) return "";
    const suit = Number(card.suit);
    const rank = Number(card.rank);
    return `${skin.cardFacePath}/${suit}-${rank}.png`;
  }

  function getAvatarFramePath() {
    return getCurrentSkin()?.avatarFrame || "";
  }

  function getAvatarPath(kind) {
    const skin = getCurrentSkin();
    if (!skin?.assets?.avatars) return "";
    return skin.avatars?.[kind] || "";
  }

  window.ddzSkinManager = {
    get skins() {
      return [...configs];
    },
    getCurrentSkin,
    getCurrentSkinId,
    setSkin,
    applyCssVariables,
    getCardFacePath,
    getAvatarFramePath,
    getAvatarPath,
    GeneratedImage,
    renderGeneratedImages,
    initSkinSelect,
  };

  document.addEventListener("DOMContentLoaded", () => {
    if (!localStorage.getItem(presetKey)) {
      localStorage.setItem(storageKey, "sanguo");
      localStorage.setItem(presetKey, "true");
    }
    setSkin(getCurrentSkinId());
    initSkinSelect();
  });
})();
