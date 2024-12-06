package utils

import (
	"fmt"
	"time"

	types "github-met/types"
)

func CalculateStreak(weeks []types.ContributionWeek) types.CalculatedStreakData {
	fmt.Println("length of weeks in calculate streak", len(weeks))

	var currentStreak types.StreakData
	var longestStreak types.StreakData

	currentStreakLength := 0
	longestStreakLength := 0

	for _, week := range weeks {
		for _, day := range week.ContributionDays {
			if day.ContributionCount > 0 {
				if currentStreakLength == 0 {
					currentStreak.StreakStartDate, _ = time.Parse("2006-01-02", day.Date)
				}
				currentStreakLength++
				currentStreak.StreakEndDate, _ = time.Parse("2006-01-02", day.Date)
			} else {
				if currentStreakLength > longestStreakLength {
					longestStreakLength = currentStreakLength
					longestStreak = currentStreak
				}
				currentStreakLength = 0
			}
		}
	}

	if currentStreakLength > longestStreakLength {
		longestStreak = currentStreak
	}

	currentStreak.Streak = currentStreakLength
	longestStreak.Streak = longestStreakLength

	return types.CalculatedStreakData{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}
}
