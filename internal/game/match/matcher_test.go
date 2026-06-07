package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palemoky/fight-the-landlord/internal/bot"
	"github.com/palemoky/fight-the-landlord/internal/config"
	"github.com/palemoky/fight-the-landlord/internal/game/room"
	"github.com/palemoky/fight-the-landlord/internal/server/storage"
	"github.com/palemoky/fight-the-landlord/internal/testutil"
)

func TestMatcher_QueueOps(t *testing.T) {
	// As long as we keep queue size < 3, it won't call CreateRoom.
	matcher := NewMatcher(MatcherDeps{}) // nil dependencies for testing

	c1 := &testutil.SimpleClient{ID: "p1", Name: "Player1"}
	c2 := &testutil.SimpleClient{ID: "p2", Name: "Player2"}

	// Add c1
	matcher.AddToQueue(c1)
	assert.Equal(t, 1, matcher.GetQueueLength())

	// Add c1 again (should be ignored)
	matcher.AddToQueue(c1)
	assert.Equal(t, 1, matcher.GetQueueLength())

	// Add c2
	matcher.AddToQueue(c2)
	assert.Equal(t, 2, matcher.GetQueueLength())

	// Remove c1
	matcher.RemoveFromQueue(c1)
	assert.Equal(t, 1, matcher.GetQueueLength())

	// Remove c1 again (should be no-op)
	matcher.RemoveFromQueue(c1)
	assert.Equal(t, 1, matcher.GetQueueLength())

	// Remove c2
	matcher.RemoveFromQueue(c2)
	assert.Equal(t, 0, matcher.GetQueueLength())
}

func TestMatcher_AddBotToWaitingRoom(t *testing.T) {
	t.Parallel()

	rm := room.NewRoomManager(storage.NewRedisStore(nil), config.GameConfig{RoomTimeout: 10})
	matcher := NewMatcher(MatcherDeps{
		RoomManager: rm,
		BotEngine:   bot.NewHeuristicEngine(),
	})
	client := testutil.NewSimpleClient("p1", "Player1")
	r, err := rm.CreateRoom(client)
	require.NoError(t, err)

	require.NoError(t, matcher.AddBotToRoom(client))

	assert.Len(t, r.Players, 2)
	assert.True(t, r.HasBot())
	for _, player := range r.Players {
		if player.Client.IsBot() {
			assert.True(t, player.Ready)
		}
	}
}
