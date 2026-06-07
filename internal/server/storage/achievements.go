package storage

// AchievementDef defines a single achievement.
type AchievementDef struct {
	ID          string
	Name        string
	Description string
	// CheckProgress returns (achieved, progressPercent). Called after each game.
	CheckProgress func(stats *PlayerStats) (achieved bool, progress int)
}

// AllAchievements returns the full achievement list.
func AllAchievements() []AchievementDef {
	return []AchievementDef{
		{
			ID: "first_win", Name: "First Win", Description: "Win your first game",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.Wins >= 1, min(s.Wins*100, 100)
			},
		},
		{
			ID: "win_10", Name: "Win 10", Description: "Win 10 games",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.Wins >= 10, min(s.Wins*10, 100)
			},
		},
		{
			ID: "win_50", Name: "Win 50", Description: "Win 50 games",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.Wins >= 50, min(s.Wins*2, 100)
			},
		},
		{
			ID: "landlord_10", Name: "Landlord x10", Description: "Play 10 games as landlord",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.LandlordGames >= 10, min(s.LandlordGames*10, 100)
			},
		},
		{
			ID: "farmer_10", Name: "Farmer x10", Description: "Play 10 games as farmer",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.FarmerGames >= 10, min(s.FarmerGames*10, 100)
			},
		},
		{
			ID: "bomb_10", Name: "Bomb Master", Description: "Play 10 bombs in total",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.BombsPlayed >= 10, min(s.BombsPlayed*10, 100)
			},
		},
		{
			ID: "spring_1", Name: "Spring", Description: "Win a spring/anti-spring game",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.SpringWins >= 1, min(s.SpringWins*100, 100)
			},
		},
		{
			ID: "streak_5", Name: "Streak 5", Description: "Achieve a 5-win streak",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.MaxWinStreak >= 5, min(s.MaxWinStreak*20, 100)
			},
		},
		{
			ID: "streak_10", Name: "Streak 10", Description: "Achieve a 10-win streak",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.MaxWinStreak >= 10, min(s.MaxWinStreak*10, 100)
			},
		},
		{
			ID: "rich_5000", Name: "Rich", Description: "Accumulate 5000 coins",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.Coins >= 5000, min(s.Coins/50, 100)
			},
		},
		{
			ID: "signin_7", Name: "Loyal", Description: "Sign in for 7 consecutive days",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.ConsecutiveSignIns >= 7, min(s.ConsecutiveSignIns*15, 100)
			},
		},
		{
			ID: "bankruptcy", Name: "Bankrupt", Description: "Receive bankruptcy subsidy",
			CheckProgress: func(s *PlayerStats) (bool, int) {
				return s.BankruptcyGrants >= 1, min(s.BankruptcyGrants*100, 100)
			},
		},
	}
}

// GetAchievementByID finds an achievement by ID.
func GetAchievementByID(id string) *AchievementDef {
	for _, a := range AllAchievements() {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

// CheckNewAchievements returns newly unlocked achievement IDs.
func CheckNewAchievements(stats *PlayerStats) []string {
	achieved := make(map[string]bool)
	for _, id := range stats.AchievedAchievements {
		achieved[id] = true
	}
	var newOnes []string
	for _, a := range AllAchievements() {
		if achieved[a.ID] {
			continue
		}
		ok, _ := a.CheckProgress(stats)
		if ok {
			newOnes = append(newOnes, a.ID)
		}
	}
	return newOnes
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
