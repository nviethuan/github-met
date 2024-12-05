package utils

import (
	"fmt"
	"os"
	"time"

	types "github-met/types"
)

func CalculateStreak(weeks []types.ContributionWeek) types.CalculatedStreakData {
	var currentStreak, longestStreak int
	var previousDate time.Time
	var startDate, endDate, longestStartDate, longestEndDate time.Time

	// Iterate backward through weeks and days to calculate streak
	for i := len(weeks) - 1; i >= 0; i-- {
		week := weeks[i]
		for j := len(week.ContributionDays) - 1; j >= 0; j-- {
			day := week.ContributionDays[j]
			if day.ContributionCount == 0 {
				if currentStreak > longestStreak {
					longestStreak = currentStreak
					longestStartDate = startDate
					longestEndDate = endDate
				}
				return types.CalculatedStreakData{
					CurrentStreak: types.StreakData{
						Streak:          currentStreak,
						StreakStartDate: startDate,
						StreakEndDate:   endDate,
					},
					LongestStreak: types.StreakData{
						Streak:          longestStreak,
						StreakStartDate: longestStartDate,
						StreakEndDate:   longestEndDate,
					},
				}
			}

			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}

			// Check for consecutive days
			if !previousDate.IsZero() && !date.AddDate(0, 0, 1).Equal(previousDate) {
				if currentStreak > longestStreak {
					longestStreak = currentStreak
					longestStartDate = startDate
					longestEndDate = endDate
				}
				currentStreak = 0
			}

			if currentStreak == 0 {
				endDate = date
			}
			currentStreak++
			startDate = date
			previousDate = date
		}
	}

	if currentStreak > longestStreak {
		longestStreak = currentStreak
		longestStartDate = startDate
		longestEndDate = endDate
	}

	return types.CalculatedStreakData{
		CurrentStreak: types.StreakData{
			Streak:          currentStreak,
			StreakStartDate: startDate,
			StreakEndDate:   endDate,
		},
		LongestStreak: types.StreakData{
			Streak:          longestStreak,
			StreakStartDate: longestStartDate,
			StreakEndDate:   longestEndDate,
		},
	}
}
