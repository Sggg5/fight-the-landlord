import { DialogueEvent, personalityById } from "./personalities";

const TALK_COOLDOWN_MS = 8000;
const forceEvents = new Set<DialogueEvent>(["playBomb", "playJokerBomb", "win"]);
const lastTalkAt = new Map<string, number>();

export function canSpeak(characterId: string, eventName: DialogueEvent, now = Date.now()): boolean {
  if (forceEvents.has(eventName)) return true;
  const last = lastTalkAt.get(characterId) || 0;
  return now - last >= TALK_COOLDOWN_MS;
}

export function getLine(characterId: string, eventName: DialogueEvent): string {
  const character = personalityById[characterId];
  const lines = character?.lines?.[eventName];
  if (!lines || lines.length === 0) return "";
  if (!canSpeak(characterId, eventName)) return "";
  lastTalkAt.set(characterId, Date.now());
  return lines[Math.floor(Math.random() * lines.length)];
}

export function resetDialogueCooldown(characterId?: string): void {
  if (characterId) {
    lastTalkAt.delete(characterId);
    return;
  }
  lastTalkAt.clear();
}
