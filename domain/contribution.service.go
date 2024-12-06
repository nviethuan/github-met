package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github-met/infrastructure/utils"
	"github-met/types"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var token string
var githubGraphQLURL string

func init() {
	token = os.Getenv("GITHUB_TOKEN")
	githubGraphQLURL = os.Getenv("GITHUB_GRAPHQL_URL")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set")
		os.Exit(1)
	}
	if githubGraphQLURL == "" {
		githubGraphQLURL = "https://api.github.com/graphql"
	}
}

func GetUserCreatedAt(username *string) time.Time {
	query := `
	query {
	  user(login: "` + *username + `") {
	    createdAt
	  }
	}`

	graphqlQuery := types.GraphQLQuery{
		Query: query,
	}

	payload, err := json.Marshal(graphqlQuery)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", githubGraphQLURL, bytes.NewBuffer(payload))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var data types.User
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		os.Exit(1)
	}

	createdAt, err := time.Parse(time.RFC3339, data.Data.User.CreatedAt)
	if err != nil {
		fmt.Println("Error parsing created at:", err)
		os.Exit(1)
	}
	return createdAt
}

// func GetContributionsStreak(username string, startDate *time.Time) types.StreakData {
// 	query := `
// 	query {
// 	  user(login: "` + username + `") {
// 	    contributionsCollection {
// 	      contributionCalendar {
// 	        weeks {
// 	          contributionDays {
// 	            date
// 	            contributionCount
// 	          }
// 	        }
// 	      }
// 	    }
// 	  }
// 	}`

// 	graphqlQuery := types.GraphQLQuery{
// 		Query: query,
// 	}

// 	payload, err := json.Marshal(graphqlQuery)
// 	if err != nil {
// 		fmt.Println("Error marshaling query:", err)
// 		os.Exit(1)
// 	}

// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", githubGraphQLURL, bytes.NewBuffer(payload))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		os.Exit(1)
// 	}

// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		os.Exit(1)
// 	}
// 	defer resp.Body.Close()

// 	// Read the response
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		os.Exit(1)
// 	}

// 	// Unmarshal the response into ContributionData
// 	var data types.ContributionData
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling response:", err)
// 		os.Exit(1)
// 	}
// 	streak, startedDateStreak, endDateStreak := utils.CalculateStreak(data.Data.User.ContributionsCollection.ContributionCalendar.Weeks)
// 	startedDate, err := time.Parse(time.RFC3339, data.Data.User.CreatedAt)
// 	if err != nil {
// 		fmt.Println("Error parsing started date:", err)
// 		os.Exit(1)
// 	}

// 	return types.StreakData{
// 		Streak:          streak,
// 		StreakStartDate: startedDateStreak,
// 		StreakEndDate:   endDateStreak,
// 		StartedDate:     startedDate,
// 	}
// }

func GetContributionsForYear(username string, start *time.Time, end *time.Time) types.ContributionData {
	contributionsCollectionParams := ""

	if start != nil && end != nil {
		contributionsCollectionParams = `(from: "` + start.Format(time.DateTime) + `", to: "` + end.Format(time.DateTime) + `")`
	} else if start != nil {
		contributionsCollectionParams = `(from: "` + start.Format(time.DateTime) + `")`
	}

	query := `
	query {
		user(login: "` + username + `") {
			contributionsCollection` + contributionsCollectionParams + ` {
				contributionCalendar {
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
					totalContributions
				}
			}
		}
	}`

	fmt.Println("query:", query)

	graphqlQuery := types.GraphQLQuery{
		Query: query,
	}

	payload, err := json.Marshal(graphqlQuery)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", githubGraphQLURL, bytes.NewBuffer(payload))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var data types.ContributionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		os.Exit(1)
	}

	prettyJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error generating pretty JSON:", err)
	} else {
		fmt.Printf("data: %v\n", string(prettyJSON))
	}

	return data
}

func GetAllContributions(username string, start *time.Time) (int, types.CalculatedStreakData) {
	yearRanges := utils.RangeOfYears(start)

	contributionsDataChan := make(chan types.ContributionData, len(yearRanges))

	for _, yearRange := range yearRanges {
		go func(start, end time.Time) {
			contributionsDataChan <- GetContributionsForYear(username, &start, &end)
		}(yearRange[0], yearRange[1])
	}

	totalContributions := 0
	weeks := []types.ContributionWeek{}
	for i := 0; i < len(yearRanges); i++ {
		contributionData := <-contributionsDataChan

		totalContributions += contributionData.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions
		weeks = append(weeks, contributionData.Data.User.ContributionsCollection.ContributionCalendar.Weeks...)
	}

	calculatedStreakData := utils.CalculateStreak(weeks)

	return totalContributions, calculatedStreakData
}
