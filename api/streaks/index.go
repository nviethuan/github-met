package handler

import (
	"net/http"

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

	theme := r.URL.Query().Get("theme")
	if theme == "" || (theme != "dark" && theme != "light") {
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


