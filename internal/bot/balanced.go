package bot

import (
	"context"
	"math/rand/v2"

	"github.com/palemoky/fight-the-landlord/internal/game/card"
)

// BalancedEngine mixes stronger DouZero play with simple heuristic play.
type BalancedEngine struct {
	strong       DecisionEngine
	weak         DecisionEngine
	strongChance int
}

func NewBalancedEngine(strong DecisionEngine, strongChance int) *BalancedEngine {
	if strongChance < 0 {
		strongChance = 0
	}
	if strongChance > 100 {
		strongChance = 100
	}
	return &BalancedEngine{
		strong:       strong,
		weak:         NewHeuristicEngine(),
		strongChance: strongChance,
	}
}

func (e *BalancedEngine) DecideBid(ctx context.Context, botName string, hand []card.Card, prevBid *bool) bool {
	return e.weak.DecideBid(ctx, botName, hand, prevBid)
}

func (e *BalancedEngine) DecidePlay(ctx context.Context, botName string, gctx GameContext) []card.Card {
	if e.strong != nil && rand.IntN(100) < e.strongChance {
		return e.strong.DecidePlay(ctx, botName, gctx)
	}
	return e.weak.DecidePlay(ctx, botName, gctx)
}
