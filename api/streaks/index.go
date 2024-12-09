package handler

import (
	"net/http"
	"time"
	"github-met/domain"
	"github-met/infrastructure/utils"
	"github-met/types"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	startDay := 6
	endDay := 18

	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		http.Error(w, "Failed to load location", http.StatusInternalServerError)
		return
	}

	currentDay := time.Now().In(location)

	var theme string
	if currentDay.Hour() >= startDay && currentDay.Hour() < endDay {
		theme = "light"
	} else {
		theme = "dark"
	}

	userCreatedAt := domain.GetUserCreatedAt(&username)

	totalContributions, calculatedStreakData := domain.GetAllContributions(username, &userCreatedAt)

	svg := utils.StreakSVG(&types.RenderData{
		CalculatedStreakData: calculatedStreakData,
		StartedDate:          userCreatedAt,
		TotalContributions:   totalContributions,
		Background:           theme,
	})

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(svg))
}


