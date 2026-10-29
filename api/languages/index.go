package handler

import (
	"github-met/domain"
	"github-met/infrastructure/utils"
	"net/http"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.URL.Query().Get("username"))
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

	currentHour := time.Now().In(location).Hour()
	theme := "dark"
	if currentHour >= 6 && currentHour < 18 {
		theme = "light"
	}

	languages, err := domain.GetUserLanguages(username)
	if err != nil {
		http.Error(w, "Failed to fetch language statistics", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(utils.LanguagesSVG(username, languages, theme)))
}
