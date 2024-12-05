package handler

import (
	"log"
	"net/http"
	"os"

	"github-met/domain"
	"github-met/infrastructure/utils"
	"github-met/types"
)


var token string
var githubGraphQLURL string

func init() {
	token = os.Getenv("GITHUB_TOKEN")
	githubGraphQLURL = os.Getenv("GITHUB_GRAPHQL_URL")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}
	if githubGraphQLURL == "" {
		githubGraphQLURL = "https://api.github.com/graphql"
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	userCreatedAt := domain.GetUserCreatedAt(&username)

	totalContributions, calculatedStreakData := domain.GetAllContributions(username, &userCreatedAt)

	svg := utils.StreakSVG(&types.RenderData{
		CalculatedStreakData: calculatedStreakData,
		StartedDate:          userCreatedAt,
		TotalContributions:   totalContributions,
	})

	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(svg))
}


