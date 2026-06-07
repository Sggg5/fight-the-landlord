package bot

import (
	"context"
	"log"
	"slices"

	"github.com/palemoky/fight-the-landlord/internal/game/card"
	"github.com/palemoky/fight-the-landlord/internal/game/rule"
)

// StrongHeuristicEngine is a stronger local fallback for hard bots.
// It stays deterministic and legal, but leads with combinations instead of
// always leaking the smallest single card like the beginner engine.
type StrongHeuristicEngine struct{}

func NewStrongHeuristicEngine() *StrongHeuristicEngine {
	return &StrongHeuristicEngine{}
}

func (e *StrongHeuristicEngine) DecideBid(_ context.Context, _ string, hand []card.Card, _ *bool) bool {
	return scoredBid(hand)
}

func (e *StrongHeuristicEngine) DecidePlay(_ context.Context, botName string, gctx GameContext) []card.Card {
	if !gctx.MustPlay && !gctx.CanBeat {
		return nil
	}

	var cards []card.Card
	if gctx.MustPlay || gctx.RecentPlays[0].Played.IsEmpty() {
		cards = bestLead(gctx.Hand)
	} else {
		cards = rule.FindSmallestBeatingCards(gctx.Hand, gctx.RecentPlays[0].Played)
		if shouldSaveBomb(cards, gctx) {
			cards = nil
		}
	}

	if cards == nil {
		log.Printf("🤖 %s strong pass", botName)
	} else {
		log.Printf("🤖 %s strong play: %s", botName, cardsToStr(cards))
	}
	return cards
}

func bestLead(hand []card.Card) []card.Card {
	if len(hand) == 0 {
		return nil
	}
	if parsed, err := rule.ParseHand(hand); err == nil && parsed.Type != rule.Invalid {
		return hand
	}
	if straight := longestStraight(hand); len(straight) >= 5 {
		return straight
	}
	if trio := lowestTrioWithKicker(hand); trio != nil {
		return trio
	}
	if pair := lowestNOfRank(hand, 2); pair != nil {
		return pair
	}
	return []card.Card{lowestCard(hand)}
}

func shouldSaveBomb(cards []card.Card, gctx GameContext) bool {
	if len(cards) == 0 || gctx.MustPlay {
		return false
	}
	parsed, err := rule.ParseHand(cards)
	if err != nil {
		return false
	}
	if parsed.Type != rule.Bomb && parsed.Type != rule.Rocket {
		return false
	}
	lastType := gctx.RecentPlays[0].Played.Type
	if lastType == rule.Bomb || lastType == rule.Rocket {
		return false
	}
	return gctx.PlayerCounts[0] > 2
}

func longestStraight(hand []card.Card) []card.Card {
	byRank := cardsByRank(hand)
	ranks := make([]card.Rank, 0, len(byRank))
	for rank := range byRank {
		if rank < card.Rank2 {
			ranks = append(ranks, rank)
		}
	}
	slices.Sort(ranks)

	var best []card.Rank
	for i := 0; i < len(ranks); {
		j := i + 1
		for j < len(ranks) && ranks[j] == ranks[j-1]+1 {
			j++
		}
		if j-i >= 5 && j-i > len(best) {
			best = ranks[i:j]
		}
		i = j
	}
	if len(best) < 5 {
		return nil
	}

	result := make([]card.Card, 0, len(best))
	for _, rank := range best {
		result = append(result, byRank[rank][0])
	}
	return result
}

func lowestTrioWithKicker(hand []card.Card) []card.Card {
	byRank := cardsByRank(hand)
	ranks := sortedRanks(byRank)
	for _, rank := range ranks {
		if len(byRank[rank]) < 3 {
			continue
		}
		result := append([]card.Card{}, byRank[rank][:3]...)
		if pair := lowestKicker(byRank, rank, 2); pair != nil {
			return append(result, pair...)
		}
		if single := lowestKicker(byRank, rank, 1); single != nil {
			return append(result, single...)
		}
		return result
	}
	return nil
}

func lowestNOfRank(hand []card.Card, n int) []card.Card {
	byRank := cardsByRank(hand)
	for _, rank := range sortedRanks(byRank) {
		if len(byRank[rank]) >= n {
			return append([]card.Card{}, byRank[rank][:n]...)
		}
	}
	return nil
}

func lowestKicker(byRank map[card.Rank][]card.Card, exclude card.Rank, n int) []card.Card {
	for _, rank := range sortedRanks(byRank) {
		if rank == exclude || len(byRank[rank]) < n {
			continue
		}
		return append([]card.Card{}, byRank[rank][:n]...)
	}
	return nil
}

func lowestCard(hand []card.Card) card.Card {
	lowest := hand[0]
	for _, c := range hand[1:] {
		if c.Rank < lowest.Rank {
			lowest = c
		}
	}
	return lowest
}

func cardsByRank(hand []card.Card) map[card.Rank][]card.Card {
	byRank := make(map[card.Rank][]card.Card)
	for _, c := range hand {
		byRank[c.Rank] = append(byRank[c.Rank], c)
	}
	return byRank
}

func sortedRanks(byRank map[card.Rank][]card.Card) []card.Rank {
	ranks := make([]card.Rank, 0, len(byRank))
	for rank := range byRank {
		ranks = append(ranks, rank)
	}
	slices.Sort(ranks)
	return ranks
}
