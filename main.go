package main

import (
	languages "github-met/api/languages"
	languagesJSON "github-met/api/languages/json"
	streaks "github-met/api/streaks"
	streaksJSON "github-met/api/streaks/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/streaks", streaks.Handler)
	mux.HandleFunc("/api/streaks/json", streaksJSON.Handler)
	mux.HandleFunc("/api/languages", languages.Handler)
	mux.HandleFunc("/api/languages/json", languagesJSON.Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
