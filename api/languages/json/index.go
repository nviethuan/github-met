package handler

import (
	"encoding/json"
	"github-met/domain"
	"github-met/types"
	"net/http"
	"strings"
)

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

	username := strings.TrimSpace(r.URL.Query().Get("username"))
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	languages, err := domain.GetUserLanguages(username)
	if err != nil {
		http.Error(w, "Failed to fetch language statistics", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	_ = json.NewEncoder(w).Encode(struct {
		Username  string             `json:"username"`
		Languages []languageResponse `json:"languages"`
	}{
		Username:  username,
		Languages: formatLanguages(languages),
	})
}

type languageResponse struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Bytes int64  `json:"bytes"`
}

func formatLanguages(languages []types.LanguageData) []languageResponse {
	result := make([]languageResponse, len(languages))
	for i, language := range languages {
		result[i] = languageResponse{
			Name:  language.Name,
			Color: language.Color,
			Bytes: language.Bytes,
		}
	}
	return result
}
