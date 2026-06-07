package storage

import (
	"context"

	"time"
)

// DailyTaskDef defines a daily task template.
type DailyTaskDef struct {
	ID          string
	Name        string
	Description string
	Requirement int
	RewardCoins int
	RewardScore int
	// CheckProgress returns the current progress given the player stats.
	CheckProgress func(stats *PlayerStats) int
}

// AllDailyTasks returns the full daily task list.
func AllDailyTasks() []DailyTaskDef {
	return []DailyTaskDef{
		{
			ID: "daily_win_3", Name: "小试牛刀", Description: "今日赢 3 局",
			Requirement: 3, RewardCoins: 100, RewardScore: 5,
			CheckProgress: func(s *PlayerStats) int { return s.Wins }, // Wins is total, but we track daily via progress map
		},
		{
			ID: "daily_landlord", Name: "称霸一方", Description: "今日做 2 次地主",
			Requirement: 2, RewardCoins: 80, RewardScore: 3,
			CheckProgress: func(s *PlayerStats) int { return s.LandlordGames },
		},
		{
			ID: "daily_bomb", Name: "炮火连天", Description: "今日打出 3 次炸弹",
			Requirement: 3, RewardCoins: 150, RewardScore: 8,
			CheckProgress: func(s *PlayerStats) int { return s.BombsPlayed },
		},
		{
			ID: "daily_play_5", Name: "乐此不疲", Description: "今日玩 5 局",
			Requirement: 5, RewardCoins: 60, RewardScore: 3,
			CheckProgress: func(s *PlayerStats) int { return s.TotalGames },
		},
		{
			ID: "daily_farmer", Name: "农民起义", Description: "今日做 3 次农民赢 2 局",
			Requirement: 2, RewardCoins: 90, RewardScore: 5,
			CheckProgress: func(s *PlayerStats) int { return s.FarmerWins },
		},
	}
}

// GetDailyTaskByID finds a daily task by ID.
func GetDailyTaskByID(id string) *DailyTaskDef {
	for _, t := range AllDailyTasks() {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

// DailyTaskStatus provides the full status of a daily task for a player.
type DailyTaskStatus struct {
	Def         DailyTaskDef
	Progress    int
	Completed   bool
	Claimed     bool
}

// GetDailyTasksStatus returns the status of all daily tasks for a player.
// It resets progress if the date has changed.
func (lm *LeaderboardManager) GetDailyTasksStatus(ctx context.Context, playerID, playerName string) ([]DailyTaskStatus, int64, error) {
	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return nil, 0, err
	}

	today := time.Now().Format("2006-01-02")
	now := time.Now()
	// Next reset is tomorrow at midnight (server local time)
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	nextReset := tomorrow.Unix()

	// Reset progress if a new day
	if stats.DailyTaskDate != today {
		stats.DailyTaskDate = today
		stats.DailyTaskProgress = make(map[string]int)
		stats.ClaimedDailyTasks = nil
		// Save the reset
		_ = lm.SavePlayerStats(ctx, stats)
	}

	claimed := make(map[string]bool)
	for _, id := range stats.ClaimedDailyTasks {
		claimed[id] = true
	}

	var tasks []DailyTaskStatus
	for _, def := range AllDailyTasks() {
		progress := stats.DailyTaskProgress[def.ID]
		if progress == 0 {
			// Initialize from stats
			progress = def.CheckProgress(stats)
		}
		completed := progress >= def.Requirement
		tasks = append(tasks, DailyTaskStatus{
			Def:       def,
			Progress:  progress,
			Completed: completed,
			Claimed:   claimed[def.ID],
		})
	}

	return tasks, nextReset, nil
}

// UpdateDailyTaskProgress updates the progress for a specific daily task after a game.
func (lm *LeaderboardManager) UpdateDailyTaskProgress(ctx context.Context, playerID string, taskID string, progress int) error {
	stats, err := lm.GetPlayerStats(ctx, playerID)
	if err != nil || stats == nil {
		return err
	}

	if stats.DailyTaskProgress == nil {
		stats.DailyTaskProgress = make(map[string]int)
	}

	today := time.Now().Format("2006-01-02")
	if stats.DailyTaskDate != today {
		stats.DailyTaskDate = today
		stats.DailyTaskProgress = make(map[string]int)
		stats.ClaimedDailyTasks = nil
	}

	if progress > stats.DailyTaskProgress[taskID] {
		stats.DailyTaskProgress[taskID] = progress
	}

	return lm.SavePlayerStats(ctx, stats)
}

// ClaimDailyTask claims the reward for a completed daily task.
func (lm *LeaderboardManager) ClaimDailyTask(ctx context.Context, playerID, playerName, taskID string) (*ClaimResult, error) {
	def := GetDailyTaskByID(taskID)
	if def == nil {
		return &ClaimResult{Success: false, Error: "任务不存在"}, nil
	}

	stats, err := lm.getOrCreateStats(ctx, playerID, playerName)
	if err != nil {
		return nil, err
	}

	// Check if already claimed
	for _, id := range stats.ClaimedDailyTasks {
		if id == taskID {
			return &ClaimResult{Success: false, Error: "已领取"}, nil
		}
	}

	// Check if completed
	tasks, _, _ := lm.GetDailyTasksStatus(ctx, playerID, playerName)
	var taskStatus *DailyTaskStatus
	for i := range tasks {
		if tasks[i].Def.ID == taskID {
			taskStatus = &tasks[i]
			break
		}
	}
	if taskStatus == nil || !taskStatus.Completed {
		return &ClaimResult{Success: false, Error: "任务未完成"}, nil
	}

	// Award rewards
	stats.Coins += def.RewardCoins
	stats.Score += def.RewardScore
	stats.LastCoinChange = def.RewardCoins
	stats.ClaimedDailyTasks = append(stats.ClaimedDailyTasks, taskID)

	if err := lm.SavePlayerStats(ctx, stats); err != nil {
		return nil, err
	}

	return &ClaimResult{
		Success:     true,
		TaskID:      taskID,
		RewardCoins: def.RewardCoins,
		RewardScore: def.RewardScore,
		Coins:       stats.Coins,
	}, nil
}

// ClaimResult holds the result of claiming a daily task reward.
type ClaimResult struct {
	Success     bool
	TaskID      string
	RewardCoins int
	RewardScore int
	Coins       int
	Error       string
}
