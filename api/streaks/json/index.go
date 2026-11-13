package handler

import (
	"encoding/json"
	"github-met/domain"
	"net/http"
	"time"
)

type streakPeriod struct {
	Streak    int       `json:"streak"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type streakResponse struct {
	Username           string       `json:"username"`
	CurrentStreak      streakPeriod `json:"currentStreak"`
	LongestStreak      streakPeriod `json:"longestStreak"`
	TotalContributions int          `json:"totalContributions"`
	StartedDate        time.Time    `json:"startedDate"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, OPTIONS")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	timezone := r.URL.Query().Get("timezone")
	if timezone == "" {
		timezone = "Asia/Ho_Chi_Minh"
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		http.Error(w, "Failed to load location", http.StatusBadRequest)
		return
	}

	startedDate := domain.GetUserCreatedAt(&username)
	totalContributions, streaks := domain.GetAllContributions(username, &startedDate, location)

	response := streakResponse{
		Username: username,
		CurrentStreak: streakPeriod{
			Streak:    streaks.CurrentStreak.Streak,
			StartDate: streaks.CurrentStreak.StreakStartDate,
			EndDate:   streaks.CurrentStreak.StreakEndDate,
		},
		LongestStreak: streakPeriod{
			Streak:    streaks.LongestStreak.Streak,
			StartDate: streaks.LongestStreak.StreakStartDate,
			EndDate:   streaks.LongestStreak.StreakEndDate,
		},
		TotalContributions: totalContributions,
		StartedDate:        startedDate,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	_ = json.NewEncoder(w).Encode(response)
}
