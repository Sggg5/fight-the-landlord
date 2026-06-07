package storage

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis key
	playerStatsKey    = "player:stats:"
	leaderboardKey    = "leaderboard:score"
	dailyLeaderboard  = "leaderboard:daily:"
	weeklyLeaderboard = "leaderboard:weekly:"
)

// PlayerStats ç©å®¶ç»è®¡æ°æ®
type PlayerStats struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`

	// æ»è®¡
	TotalGames int `json:"total_games"` // æ»åºæ¬?
	Wins       int `json:"wins"`        // èåº
	Losses     int `json:"losses"`      // è´¥åº

	// å°ä¸»/åæ°åå¼ç»è®¡
	LandlordGames int `json:"landlord_games"` // å°ä¸»åºæ¬¡
	LandlordWins  int `json:"landlord_wins"`  // å°ä¸»èåº
	FarmerGames   int `json:"farmer_games"`   // åæ°åºæ¬¡
	FarmerWins    int `json:"farmer_wins"`    // åæ°èåº

	// ç§¯å
	Score int `json:"score"` // å½åç§¯å

	// ç»æµ
	Coins            int `json:"coins"`             // å½åè±å­
	LastCoinChange   int `json:"last_coin_change"`  // ä¸ä¸å±è±å­åå
	BankruptcyGrants int `json:"bankruptcy_grants"` // ç ´äº§è¡¥å©æ¬¡æ°

	// è¿è/è¿è´¥
	CurrentStreak int `json:"current_streak"` // æ­£æ°ä¸ºè¿èï¼è´æ°ä¸ºè¿è´?
	MaxWinStreak  int `json:"max_win_streak"` // æå¤§è¿è?

	// æ¶é´
	LastPlayedAt int64 `json:"last_played_at"` // æåæ¸¸ææ¶é?
	CreatedAt    int64 `json:"created_at"`     // é¦æ¬¡æ¸¸ææ¶é´
	// Ç©µ½
	LastSignInDate string `json:"last_sign_in_date"`      // ×îºóÇ©µ½ÈÕÆÚ (2006-01-02)
	ConsecutiveSignIns int `json:"consecutive_sign_ins"`   // Á¬ÐøÇ©µ½ÌìÊý

	// ³É¾Í¼ÆÊý
	BombsPlayed int `json:"bombs_played"`           // ÀÛ¼ÆÕ¨µ¯Êý
	SpringWins int `json:"spring_wins"`            // ´ºÌì/·´´ºÌìÊ¤ÀûÊý
	AchievedAchievements []string `json:"achieved_achievements"`

	// 商埆
	Inventory          []string `json:"inventory"`

	// 每日任务
	DailyTaskDate     string   `json:"daily_task_date"`
	DailyTaskProgress map[string]int `json:"daily_task_progress"`
	ClaimedDailyTasks []string `json:"claimed_daily_tasks"` // ÒÑ»ñµÃ³É¾ÍIDÁÐ±í
}

// ç§¯åè§å
const (
	WinAsLandlord  = 30  // å°ä¸»è·è
	WinAsFarmer    = 15  // åæ°è·è
	LoseAsLandlord = -20 // å°ä¸»å¤±è´¥
	LoseAsFarmer   = -10 // åæ°å¤±è´¥

	// è¿èå æ
	StreakBonus3  = 5  // 3 è¿èå æ
	StreakBonus5  = 10 // 5 è¿èå æ
	StreakBonus10 = 20 // 10 è¿èå æ
)

// ç»æµè§å
const (
	InitialCoins        = 1000
	BaseStake           = 10
	BankruptcyThreshold = 100
	BankruptcySubsidy   = 1000
)

// LeaderboardEntry æè¡æ¦æ¡ç?
type LeaderboardEntry struct {
	Rank       int     `json:"rank"`
	PlayerID   string  `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Score      int     `json:"score"`
	Wins       int     `json:"wins"`
	WinRate    float64 `json:"win_rate"`
}

// LeaderboardManager æè¡æ¦ç®¡çå¨
type LeaderboardManager struct {
	redis *redis.Client
}

// NewLeaderboardManager åå»ºæè¡æ¦ç®¡çå¨
func NewLeaderboardManager(client *redis.Client) *LeaderboardManager {
	return &LeaderboardManager{redis: client}
}

// IsReady æ£æ?Redis å®¢æ·ç«¯æ¯å¦å¯ç?
func (lm *LeaderboardManager) IsReady() bool {
	return lm != nil && lm.redis != nil
}

// GetPlayerStats è·åç©å®¶ç»è®¡
func (lm *LeaderboardManager) GetPlayerStats(ctx context.Context, playerID string) (*PlayerStats, error) {
	key := playerStatsKey + playerID
	data, err := lm.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var stats PlayerStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// SavePlayerStats ä¿å­ç©å®¶ç»è®¡
func (lm *LeaderboardManager) SavePlayerStats(ctx context.Context, stats *PlayerStats) error {
	key := playerStatsKey + stats.PlayerID
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}
	return lm.redis.Set(ctx, key, data, 0).Err()
}

// getOrCreateStats è·åæåå»ºç©å®¶ç»è®?
func (lm *LeaderboardManager) getOrCreateStats(ctx context.Context, playerID, playerName string) (*PlayerStats, error) {
	stats, err := lm.GetPlayerStats(ctx, playerID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return &PlayerStats{
			PlayerID:   playerID,
			PlayerName: playerName,
			Coins:      InitialCoins,
			CreatedAt:  time.Now().Unix(),
		}, nil
	}

	NormalizeEconomy(stats)

	return stats, nil
}

// NormalizeEconomy initializes economy fields for old stats records.
func NormalizeEconomy(stats *PlayerStats) {
	if stats == nil {
		return
	}
	if stats.Coins == 0 && stats.LastCoinChange == 0 && stats.BankruptcyGrants == 0 {
		stats.Coins = InitialCoins
	}
}

// updateRoleStats æ´æ°è§è²ç¸å³ç»è®¡å¹¶è¿ååºç¡ç§¯ååå
func updateRoleStats(stats *PlayerStats, isLandlord, isWinner bool) int {
	switch {
	case isLandlord && isWinner:
		stats.LandlordGames++
		stats.LandlordWins++
		return WinAsLandlord
	case isLandlord && !isWinner:
		stats.LandlordGames++
		return LoseAsLandlord
	case !isLandlord && isWinner:
		stats.FarmerGames++
		stats.FarmerWins++
		return WinAsFarmer
	default: // !isLandlord && !isWinner
		stats.FarmerGames++
		return LoseAsFarmer
	}
}

// updateWinLossStats æ´æ°èè´ç»è®¡åè¿è?è¿è´¥
func updateWinLossStats(stats *PlayerStats, isWinner bool) {
	if isWinner {
		stats.Wins++
		stats.CurrentStreak = max(1, stats.CurrentStreak+1)
	} else {
		stats.Losses++
		stats.CurrentStreak = min(-1, stats.CurrentStreak-1)
	}

	if stats.CurrentStreak > stats.MaxWinStreak {
		stats.MaxWinStreak = stats.CurrentStreak
	}
}

// calculateStreakBonus è®¡ç®è¿èå æ
func calculateStreakBonus(streak int) int {
	switch {
	case streak >= 10:
		return StreakBonus10
	case streak >= 5:
		return StreakBonus5
	case streak >= 3:
		return StreakBonus3
	default:
		return 0
	}
}

func updateEconomyStats(stats *PlayerStats, roundScore int) {
	coinChange := roundScore * BaseStake
	stats.LastCoinChange = coinChange
	stats.Coins = max(0, stats.Coins+coinChange)

	if stats.Coins < BankruptcyThreshold {
		stats.Coins = BankruptcySubsidy
		stats.BankruptcyGrants++
	}
}

// RecordGameResult è®°å½æ¸¸æç»æ
func (lm *LeaderboardManager) RecordGameResult(ctx context.Context, playerID, playerName string, isLandlord, isWinner bool) error {
	roundScore := 1
	if isLandlord {
		roundScore = 2
	}
	if !isWinner {
		roundScore = -roundScore
	}
	return lm.RecordGameResultWithScore(ctx, playerID, playerName, isLandlord, isWinner, roundScore)
}

// RecordGameResultWithScore records a game result and applies economy changes.
func (lm *LeaderboardManager) RecordGameResultWithScore(ctx context.Context, playerID, playerName string, isLandlord, isWinner bool, roundScore int) error {
	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return err
	}

	// æ´æ°åºæ¬ä¿¡æ¯
	stats.PlayerName = playerName
	stats.TotalGames++
	stats.LastPlayedAt = time.Now().Unix()

	// æ´æ°è§è²åèè´ç»è®?
	scoreChange := updateRoleStats(stats, isLandlord, isWinner)
	updateWinLossStats(stats, isWinner)

	// è®¡ç®è¿èå æå¹¶æ´æ°ç§¯å?
	scoreChange += calculateStreakBonus(stats.CurrentStreak)
	stats.Score = max(0, stats.Score+scoreChange)
	updateEconomyStats(stats, roundScore)

	// ä¿å­å¹¶æ´æ°æè¡æ¦
	if err := lm.SavePlayerStats(ctx, stats); err != nil {
		return err
	}
	return lm.UpdateLeaderboard(ctx, stats)
}

// UpdateLeaderboard æ´æ°æè¡æ¦?


// SignIn performs daily sign-in and returns the reward.
func (lm *LeaderboardManager) SignIn(ctx context.Context, playerID, playerName string) (reward int, consecutive int, err error) {
	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return 0, 0, err
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if stats.LastSignInDate == today {
		return 0, stats.ConsecutiveSignIns, nil // already signed in today
	}

	if stats.LastSignInDate == yesterday {
		stats.ConsecutiveSignIns++
	} else {
		stats.ConsecutiveSignIns = 1
	}
	stats.LastSignInDate = today

	// Reward: day 1=50, 2=80, 3=100, 4=150, 5+=200
	switch {
	case stats.ConsecutiveSignIns >= 5:
		reward = 200
	case stats.ConsecutiveSignIns == 4:
		reward = 150
	case stats.ConsecutiveSignIns == 3:
		reward = 100
	case stats.ConsecutiveSignIns == 2:
		reward = 80
	default:
		reward = 50
	}

	stats.Coins += reward
	consecutive = stats.ConsecutiveSignIns

	if err := lm.SavePlayerStats(ctx, stats); err != nil {
		return 0, 0, err
	}
	return reward, consecutive, nil
}

// CanSignIn checks if the player can sign in today.
func (lm *LeaderboardManager) CanSignIn(ctx context.Context, playerID string) (bool, int) {
	stats, err := lm.GetPlayerStats(ctx, playerID)
	if err != nil || stats == nil {
		return true, 0
	}
	today := time.Now().Format("2006-01-02")
	return stats.LastSignInDate != today, stats.ConsecutiveSignIns
}

// GetAchievementsStatus returns the current status of all achievements for a player.
func (lm *LeaderboardManager) GetAchievementsStatus(ctx context.Context, playerID string) ([]AchievementStatus, error) {
	stats, err := lm.GetPlayerStats(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		return nil, nil
	}

	achieved := make(map[string]bool)
	for _, id := range stats.AchievedAchievements {
		achieved[id] = true
	}

	var result []AchievementStatus
	for _, a := range AllAchievements() {
		ok, prog := a.CheckProgress(stats)
		result = append(result, AchievementStatus{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			Achieved:    achieved[a.ID] || ok,
			Progress:    prog,
		})
	}
	return result, nil
}

// AchievementStatus represents a player''s progress on one achievement.
type AchievementStatus struct {
	ID          string
	Name        string
	Description string
	Achieved    bool
	Progress    int
}

func (lm *LeaderboardManager) UpdateLeaderboard(ctx context.Context, stats *PlayerStats) error {
	// æ´æ°æ»æè¡æ¦
	if err := lm.redis.ZAdd(ctx, leaderboardKey, redis.Z{
		Score:  float64(stats.Score),
		Member: stats.PlayerID,
	}).Err(); err != nil {
		return err
	}

	// æ´æ°æ¯æ¥æè¡æ¦?
	today := time.Now().Format("2006-01-02")
	dailyKey := dailyLeaderboard + today
	if err := lm.redis.ZAdd(ctx, dailyKey, redis.Z{
		Score:  float64(stats.Score),
		Member: stats.PlayerID,
	}).Err(); err != nil {
		return err
	}
	// è®¾ç½®è¿ææ¶é´ï¼?å¤©ï¼
	lm.redis.Expire(ctx, dailyKey, 48*time.Hour)

	// æ´æ°æ¯å¨æè¡æ¦?
	year, week := time.Now().ISOWeek()
	weeklyKey := fmt.Sprintf("%s%d-W%02d", weeklyLeaderboard, year, week)
	if err := lm.redis.ZAdd(ctx, weeklyKey, redis.Z{
		Score:  float64(stats.Score),
		Member: stats.PlayerID,
	}).Err(); err != nil {
		return err
	}
	// è®¾ç½®è¿ææ¶é´ï¼?å¤©ï¼
	lm.redis.Expire(ctx, weeklyKey, 8*24*time.Hour)

	return nil
}

// GetLeaderboard è·åæè¡æ¦?
func (lm *LeaderboardManager) GetLeaderboard(ctx context.Context, limit int) ([]*LeaderboardEntry, error) {
	leaderboardType := "total"
	offset := 0
	// ç¡®å®ä½¿ç¨åªä¸ªæè¡æ¦?
	key := leaderboardKey
	switch leaderboardType {
	case "daily":
		today := time.Now().Format("2006-01-02")
		key = dailyLeaderboard + today
	case "weekly":
		year, week := time.Now().ISOWeek()
		key = fmt.Sprintf("%s%d-W%02d", weeklyLeaderboard, year, week)
	}

	// è·åæè¡æ¦ï¼ä»é«å°ä½ï¼?
	results, err := lm.redis.ZRevRangeWithScores(ctx, key, int64(offset), int64(offset+limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]*LeaderboardEntry, 0, len(results))
	for i, result := range results {
		playerID := result.Member.(string)

		// è·åç©å®¶è¯¦ç»ç»è®¡
		stats, err := lm.GetPlayerStats(ctx, playerID)
		if err != nil || stats == nil {
			continue
		}

		winRate := 0.0
		if stats.TotalGames > 0 {
			winRate = float64(stats.Wins) / float64(stats.TotalGames) * 100
		}

		entries = append(entries, &LeaderboardEntry{
			Rank:       offset + i + 1,
			PlayerID:   playerID,
			PlayerName: stats.PlayerName,
			Score:      int(result.Score),
			Wins:       stats.Wins,
			WinRate:    winRate,
		})
	}

	return entries, nil
}

// GetPlayerRank è·åç©å®¶æå
func (lm *LeaderboardManager) GetPlayerRank(ctx context.Context, playerID string) (int64, error) {
	rank, err := lm.redis.ZRevRank(ctx, leaderboardKey, playerID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil // æªä¸æ¦?
		}
		return -1, err
	}
	return rank + 1, nil // Redis æåä»?0 å¼å§?
}

// SortByScore æç§¯åæåº?
func SortByScore(entries []LeaderboardEntry) {
	slices.SortFunc(entries, func(a, b LeaderboardEntry) int {
		return cmp.Compare(b.Score, a.Score)
	})
}
