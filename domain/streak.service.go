package domain

import (
	"fmt"
	"github-met/types"
	"os"
	"time"
	"sync"
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
	date, err := time.Parse("2006-01-02", lastDay.Date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		os.Exit(1)
	}
	
	if lastDay.ContributionCount == 0 {
		currentStreakData := types.StreakData{
			Streak:          0,
			StreakStartDate: date,
			StreakEndDate:   date,
		}
		
		currentStreakDataChan <- currentStreakData
	}

	currentStreakData := types.StreakData{
		Streak:          1,
		StreakStartDate: date,
		StreakEndDate:   date,
	}

	for i := len(*days) - 2; i >= 0; i-- {
		currentDay := (*days)[i]
		date, err := time.Parse("2006-01-02", currentDay.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			os.Exit(1)
		}
		
		if currentDay.ContributionCount == 0 {
			currentStreakDataChan <- currentStreakData
		}
				
		if currentDay.ContributionCount > 0 {
		
			currentStreakData.StreakStartDate = date
			currentStreakData.Streak++
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
		
		var wg sync.WaitGroup
		wg.Add(2)
		
		go func() {
			defer wg.Done()
			quickSort(contributionDays, low, pi-1)
		}()
		
		go func() {
			defer wg.Done() 
			quickSort(contributionDays, pi+1, high)
		}()
		
		wg.Wait()
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
