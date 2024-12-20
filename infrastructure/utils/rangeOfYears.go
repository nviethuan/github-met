package utils

import "time"

func RangeOfYears(start *time.Time, timezone *time.Location) [][]time.Time {
	if start == nil {
		now := time.Now()
		start = &now
	}
	currentYear := start.Year()
	currentTime := time.Now()
	endYear := currentTime.Year()

	var yearRanges [][]time.Time

	for year := currentYear; year <= endYear; year++ {
		startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, timezone)
		endOfYear := time.Date(year, time.December, 31, 23, 59, 59, 0, timezone)
		if year == currentYear {
			startOfYear = *start
		}
		if year == endYear {
			endOfYear = currentTime
		}
		yearRanges = append(yearRanges, []time.Time{startOfYear, endOfYear})
	}

	return yearRanges
}
