package domain

import (
	"fmt"
	"os"
	"time"
	types "github-met/types"
)

func CalculateStreak(weeks []types.ContributionWeek) (int, time.Time, time.Time) {
	var streak int
	var previousDate, startDate, endDate time.Time

	// Iterate backward through weeks and days to calculate streak
	for i := len(weeks) - 1; i >= 0; i-- {
		week := weeks[i]
		for j := len(week.ContributionDays) - 1; j >= 0; j-- {
			day := week.ContributionDays[j]
			if day.ContributionCount == 0 {
				return streak, startDate, endDate // Streak ends when a day with no contributions is encountered
			}

			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}

			// Check for consecutive days
			if !previousDate.IsZero() && !date.AddDate(0, 0, 1).Equal(previousDate) {
				return streak, startDate, endDate
			}

			if streak == 0 {
				endDate = date // Set end date of streak
			}
			streak++
			startDate = date // Update start date of streak
			previousDate = date
		}
	}

	return streak, startDate, endDate
}

func GetLastStreak(weeks []types.ContributionWeek) (int, time.Time, time.Time) {
	var streak int
	var startDate, endDate time.Time
	var currentStreak int
	var currentStartDate, currentEndDate time.Time

	for i := len(weeks) - 1; i >= 0; i-- {
		week := weeks[i]
		for j := len(week.ContributionDays) - 1; j >= 0; j-- {
			day := week.ContributionDays[j]
			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}

			if day.ContributionCount > 0 {
				if currentStreak == 0 {
					currentEndDate = date
				}
				currentStreak++
				currentStartDate = date
			} else {
				if currentStreak > streak {
					streak = currentStreak
					startDate = currentStartDate
					endDate = currentEndDate
				}
				currentStreak = 0
			}
		}
	}

	if currentStreak > streak {
		streak = currentStreak
		startDate = currentStartDate
		endDate = currentEndDate
	}

	return streak, startDate, endDate
}


