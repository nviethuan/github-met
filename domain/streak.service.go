package domain

import (
	"fmt"
	"github-met/types"
	"os"
	"time"
)

func CalculateStreak(days []types.ContributionDay) types.CalculatedStreakData {
	currentStreakData := make(chan types.StreakData)
	longestStreakData := make(chan types.StreakData)

	go CurrentStreak(&days, currentStreakData)
	go LongestStreak(&days, longestStreakData)

	return types.CalculatedStreakData{
		CurrentStreak: <-currentStreakData,
		LongestStreak: <-longestStreakData,
	}
}

func CurrentStreak(days *[]types.ContributionDay, currentStreakDataChan chan types.StreakData) {
	lastDay := (*days)[len(*days)-1]

	if lastDay.ContributionCount == 0 {
		date, err := time.Parse("2006-01-02", lastDay.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			os.Exit(1)
		}

		currentStreakDataChan <- types.StreakData{
			Streak:          0,
			StreakStartDate: date,
			StreakEndDate:   date,
		}
	}

	currentStreakData := types.StreakData{}

	for i := len(*days) - 1; i >= 0; i-- {
		if (*days)[i].ContributionCount > 0 {
			date, err := time.Parse("2006-01-02", (*days)[i].Date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}
			currentStreakData.StreakStartDate = date
			currentStreakData.StreakEndDate = date
			currentStreakData.Streak++
		} else {
			break
		}
	}

	currentStreakDataChan <- currentStreakData
}

func LongestStreak(days *[]types.ContributionDay, longestStreakDataChan chan types.StreakData) {
	var longestStreakData types.StreakData
	var currentStreakData types.StreakData

	for i, day := range *days {
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			os.Exit(1)
		}

		if day.ContributionCount > 0 {
			if currentStreakData.Streak == 0 {
				currentStreakData.StreakStartDate = date
			}
			currentStreakData.StreakEndDate = date
			currentStreakData.Streak++
		} else {
			if currentStreakData.Streak > longestStreakData.Streak {
				longestStreakData = currentStreakData
			}
			currentStreakData = types.StreakData{}
		}

		if i == len(*days)-1 && currentStreakData.Streak > longestStreakData.Streak {
			longestStreakData = currentStreakData
		}
	}

	longestStreakDataChan <- longestStreakData
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
