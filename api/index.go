package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
	"time"
)

type GraphQLQuery struct {
	Query string `json:"query"`
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
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		http.Error(w, "GITHUB_TOKEN is not set", http.StatusInternalServerError)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	query := `
	query {
	  user(login: "` + username + `") {
	    contributionsCollection {
	      contributionCalendar {
	        weeks {
	          contributionDays {
	            date
	            contributionCount
	          }
	        }
	      }
	    }
	  }
	}`

	graphqlQuery := GraphQLQuery{
		Query: query,
	}

	// Marshal the query into JSON
	payload, err := json.Marshal(graphqlQuery)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		os.Exit(1)
	}

	// Make the HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// Unmarshal the response into ContributionData
	var data ContributionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		os.Exit(1)
	}
	
	streak := calculateStreak(data.Data.User.ContributionsCollection.ContributionCalendar.Weeks)

	// Create a response map to include streak information
	response := map[string]interface{}{
		"contributionData": data,
		"streak":           streak,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func calculateStreak(weeks []ContributionWeek) int {
	var streak int
	var previousDate time.Time

	// Flatten the weeks into a single slice of days
	var contributionDays []ContributionDay
	for _, week := range weeks {
		contributionDays = append(contributionDays, week.ContributionDays...)
	}

	// Iterate backward to calculate streak
	for i := len(contributionDays) - 1; i >= 0; i-- {
		day := contributionDays[i]
		if day.ContributionCount == 0 {
			break // Streak ends when a day with no contributions is encountered
		}

		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			os.Exit(1)
		}

		// Check for consecutive days
		if !previousDate.IsZero() && !date.AddDate(0, 0, 1).Equal(previousDate) {
			break
		}

		streak++
		previousDate = date
	}

	return streak
}
