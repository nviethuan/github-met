package handler

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

type GraphQLQuery struct {
	Query string `json:"query"`
}

type ContributionData struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
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
	        totalContributions
	      }
	    }
	  }
	}`

	graphqlQuery := GraphQLQuery{
		Query: query,
	}

	payload, err := json.Marshal(graphqlQuery)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		os.Exit(1)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	var data ContributionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
