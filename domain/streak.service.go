package domain

import (
	"fmt"
	"github-met/types"
	"os"
	"time"
)

func CalculateStreak(days []types.ContributionDay) types.CalculatedStreakData {
	fmt.Println("days length", len(days))
	fmt.Println("days", days)
	var currentStreakData types.StreakData
	var longestStreakData types.StreakData

	currentStreakLength := 0
	longestStreakLength := 0

	for _, day := range days {
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			os.Exit(1)
		}

		if day.ContributionCount > 0 {
			if currentStreakLength == 0 {
				currentStreakData.StreakStartDate = date
			}
			currentStreakLength++
			currentStreakData.StreakEndDate = date
		} else {
			if currentStreakLength > longestStreakLength {
				longestStreakLength = currentStreakLength
				longestStreakData = currentStreakData
			}
			currentStreakLength = 0
			currentStreakData = types.StreakData{} // Reset current streak data
		}
	}

	if currentStreakLength > longestStreakLength {
		longestStreakData = currentStreakData
	}

	currentStreakData.Streak = currentStreakLength
	longestStreakData.Streak = longestStreakLength

	return types.CalculatedStreakData{
		CurrentStreak: currentStreakData,
		LongestStreak: longestStreakData,
	}
}

func FlattenContributionDays(weeks []types.ContributionWeek) []types.ContributionDay {
	var contributionDays []types.ContributionDay

	for _, week := range weeks {
		contributionDays = append(contributionDays, week.ContributionDays...)
	}

	return contributionDays
}

func SortContributionDays(contributionDays []types.ContributionDay) []types.ContributionDay {
	quickSort(contributionDays, 0, len(contributionDays)-1)
	return contributionDays
}

func quickSort(contributionDays []types.ContributionDay, low, high int) {
	if low < high {
		pi := partition(contributionDays, low, high)
		quickSort(contributionDays, low, pi-1)
		quickSort(contributionDays, pi+1, high)
	}
}

func partition(contributionDays []types.ContributionDay, low, high int) int {
	pivot := contributionDays[high]
	i := low - 1

	for j := low; j < high; j++ {
		if contributionDays[j].Date < pivot.Date {
			i++
			contributionDays[i], contributionDays[j] = contributionDays[j], contributionDays[i]
		}
	}
	contributionDays[i+1], contributionDays[high] = contributionDays[high], contributionDays[i+1]
	return i + 1
}
