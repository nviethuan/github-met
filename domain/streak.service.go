package domain

import (
	"fmt"
	"github-met/types"
	"os"
	"sync"
	"time"
)

func CalculateStreak(days []types.ContributionDay, location *time.Location) types.CalculatedStreakData {
	currentStreakData := make(chan types.StreakData)
	longestStreakData := make(chan types.StreakData)

	go CurrentStreak(&days, currentStreakData, location)
	go LongestStreak(&days, longestStreakData)

	return types.CalculatedStreakData{
		CurrentStreak: <-currentStreakData,
		LongestStreak: <-longestStreakData,
	}
}

func CurrentStreak(days *[]types.ContributionDay, currentStreakDataChan chan types.StreakData, location *time.Location) {
	lastDay := (*days)[len(*days)-1]
	date, err := time.Parse("2006-01-02", lastDay.Date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		os.Exit(1)
	}

	today := time.Now().In(location).Format("2006-01-02")
	isToday := lastDay.Date == today

	// Nếu hôm nay chưa có contribution, chúng ta cần kiểm tra streak từ hôm qua
	if lastDay.ContributionCount == 0 && isToday {
		// Tìm streak từ ngày hôm qua trở về trước
		currentStreakData := types.StreakData{
			Streak:          0,
			StreakStartDate: date,
			StreakEndDate:   date,
		}

		// Bắt đầu từ ngày hôm qua (index len-2)
		for i := len(*days) - 2; i >= 0; i-- {
			currentDay := (*days)[i]
			currentDate, err := time.Parse("2006-01-02", currentDay.Date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}

			if currentDay.ContributionCount == 0 {
				// Nếu gặp ngày không có contribution, dừng lại và trả về streak hiện tại
				currentStreakDataChan <- currentStreakData
				return
			}

			// Nếu có contribution, tăng streak
			if currentStreakData.Streak == 0 {
				currentStreakData.StreakEndDate = currentDate
			}
			currentStreakData.StreakStartDate = currentDate
			currentStreakData.Streak++
		}

		currentStreakDataChan <- currentStreakData
		return
	}

	// Nếu hôm nay có contribution hoặc không phải hôm nay, tính toán bình thường
	if lastDay.ContributionCount == 0 {
		currentStreakData := types.StreakData{
			Streak:          0,
			StreakStartDate: date,
			StreakEndDate:   date,
		}

		currentStreakDataChan <- currentStreakData
		return
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
			return
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
	contributionDays = blockSort(contributionDays, 10)
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

func blockSort(contributionDays []types.ContributionDay, blockSize int) []types.ContributionDay {
	n := len(contributionDays)
	if n <= blockSize {
		quickSort(contributionDays, 0, n-1)
		return contributionDays
	}

	numBlocks := (n + blockSize - 1) / blockSize
	blocks := make([][]types.ContributionDay, numBlocks)
	var wg sync.WaitGroup

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := start + blockSize
		if end > n {
			end = n
		}
		blocks[i] = contributionDays[start:end]
		wg.Add(1)
		go func(block []types.ContributionDay) {
			defer wg.Done()
			quickSort(block, 0, len(block)-1)
		}(blocks[i])
	}

	wg.Wait()
	sortedDays := mergeBlocks(blocks)
	return sortedDays
}

func mergeBlocks(blocks [][]types.ContributionDay) []types.ContributionDay {
	var sortedDays []types.ContributionDay
	indices := make([]int, len(blocks))

	for {
		minIndex := -1
		for i := 0; i < len(blocks); i++ {
			if indices[i] < len(blocks[i]) {
				if minIndex == -1 || blocks[i][indices[i]].Date < blocks[minIndex][indices[minIndex]].Date {
					minIndex = i
				}
			}
		}

		if minIndex == -1 {
			break
		}

		sortedDays = append(sortedDays, blocks[minIndex][indices[minIndex]])
		indices[minIndex]++
	}

	return sortedDays
}
