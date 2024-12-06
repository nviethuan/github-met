package types

import "time"

type GraphQLQuery struct {
	Query string `json:"query"`
}

type User struct {
	Data struct {
		User struct {
			CreatedAt string `json:"createdAt"`
		} `json:"user"`
	} `json:"data"`
}

// ContributionDay represents a day's contributions
type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

// ContributionWeek represents a week's contributions
type ContributionWeek struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

// ContributionData is the structure for the response
type ContributionData struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []ContributionWeek `json:"weeks"`
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
			CreatedAt string `json:"createdAt"`
		} `json:"user"`
	} `json:"data"`
}

type StreakData struct {
	Streak          int
	StreakStartDate time.Time
	StreakEndDate   time.Time

	StartedDate time.Time
}

type CalculatedStreakData struct {
	CurrentStreak      StreakData
	LongestStreak      StreakData
}

type RenderData struct {
	CalculatedStreakData
	TotalContributions int
	StartedDate        time.Time
	Background         string
}


