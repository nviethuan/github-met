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
	"sort"
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

func GetUserLanguages(username string) ([]types.LanguageData, error) {
	const query = `
	query($login: String!, $after: String) {
	  user(login: $login) {
	    repositories(ownerAffiliations: OWNER, first: 10, after: $after, isFork: false) {
	      nodes {
	        languages(first: 40, orderBy: {field: SIZE, direction: DESC}) {
	          edges {
	            size
	            node {
	              name
	              color
	            }
	          }
	        }
	      }
	      pageInfo {
	        endCursor
	        hasNextPage
	      }
	    }
	  }
	}`

	var cursor *string
	languagesByName := make(map[string]types.LanguageData)
	client := &http.Client{Timeout: 15 * time.Second}

	for {
		graphqlQuery := struct {
			Query     string `json:"query"`
			Variables struct {
				Login string  `json:"login"`
				After *string `json:"after"`
			} `json:"variables"`
		}{
			Query: query,
			Variables: struct {
				Login string  `json:"login"`
				After *string `json:"after"`
			}{
				Login: username,
				After: cursor,
			},
		}

		payload, err := json.Marshal(graphqlQuery)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", githubGraphQLURL, bytes.NewBuffer(payload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GitHub API returned status %s", resp.Status)
		}

		var result struct {
			Data struct {
				User *struct {
					Repositories struct {
						Nodes []struct {
							Languages struct {
								Edges []struct {
									Size int64 `json:"size"`
									Node struct {
										Name  string `json:"name"`
										Color string `json:"color"`
									} `json:"node"`
								} `json:"edges"`
							} `json:"languages"`
						} `json:"nodes"`
						PageInfo struct {
							EndCursor   *string `json:"endCursor"`
							HasNextPage bool    `json:"hasNextPage"`
						} `json:"pageInfo"`
					} `json:"repositories"`
				} `json:"user"`
			} `json:"data"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("GitHub GraphQL error: %s", result.Errors[0].Message)
		}
		if result.Data.User == nil {
			return nil, fmt.Errorf("GitHub user %q was not found", username)
		}

		for _, repository := range result.Data.User.Repositories.Nodes {
			for _, edge := range repository.Languages.Edges {
				language := languagesByName[edge.Node.Name]
				language.Name = edge.Node.Name
				if language.Color == "" {
					language.Color = edge.Node.Color
				}
				language.Bytes += edge.Size
				languagesByName[edge.Node.Name] = language
			}
		}

		pageInfo := result.Data.User.Repositories.PageInfo
		if !pageInfo.HasNextPage {
			break
		}
		if pageInfo.EndCursor == nil || (cursor != nil && *pageInfo.EndCursor == *cursor) {
			return nil, fmt.Errorf("GitHub returned an invalid repository pagination cursor")
		}
		cursor = pageInfo.EndCursor
	}

	languages := make([]types.LanguageData, 0, len(languagesByName))
	for _, language := range languagesByName {
		languages = append(languages, language)
	}
	sort.Slice(languages, func(i, j int) bool {
		if languages[i].Bytes == languages[j].Bytes {
			return languages[i].Name < languages[j].Name
		}
		return languages[i].Bytes > languages[j].Bytes
	})

	return languages, nil
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

func GetContributionsForYear(username string, start *time.Time, end *time.Time) (types.ContributionData, error) {
	contributionsCollectionParams := ""

	if start != nil && end != nil {
		contributionsCollectionParams = `(from: "` + start.Format(time.RFC3339) + `", to: "` + end.Format(time.RFC3339) + `")`
	} else if start != nil {
		contributionsCollectionParams = `(from: "` + start.Format(time.RFC3339) + `")`
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

	graphqlQuery := types.GraphQLQuery{
		Query: query,
	}

	payload, err := json.Marshal(graphqlQuery)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		return types.ContributionData{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", githubGraphQLURL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return types.ContributionData{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return types.ContributionData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return types.ContributionData{}, err
	}

	var data types.ContributionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return types.ContributionData{}, err
	}

	return data, nil
}

func GetAllContributions(username string, start *time.Time, location *time.Location) (int, types.CalculatedStreakData) {
	yearRanges := utils.RangeOfYears(start, location)

	contributionsDataChan := make(chan types.ContributionData, len(yearRanges))
	errorsChan := make(chan error, len(yearRanges))

	for _, yearRange := range yearRanges {
		go func(start, end time.Time) {
			contributionData, err := GetContributionsForYear(username, &start, &end)
			if err != nil {
				errorsChan <- err
				return
			}
			contributionsDataChan <- contributionData
		}(yearRange[0], yearRange[1])
	}

	totalContributions := 0
	weeks := []types.ContributionWeek{}
	for i := 0; i < len(yearRanges); i++ {
		contributionData := <-contributionsDataChan

		totalContributions += contributionData.Data.User.ContributionsCollection.ContributionCalendar.TotalContributions
		weeks = append(weeks, contributionData.Data.User.ContributionsCollection.ContributionCalendar.Weeks...)
	}

	contributionDays := SortContributionDays(FlattenContributionDays(weeks))
	calculatedStreakData := CalculateStreak(contributionDays, location)

	return totalContributions, calculatedStreakData
}
