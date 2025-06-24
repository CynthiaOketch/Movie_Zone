package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"moviezone/api"
	"moviezone/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Load environment variables (API keys)
	tmdbKey := os.Getenv("TMDB_API_KEY")
	omdbKey := os.Getenv("OMDB_API_KEY")

	if tmdbKey == "" || omdbKey == "" {
		log.Println("Warning: TMDB_API_KEY or OMDB_API_KEY not set in environment.")
	}

	// Initialize SQLite3 DB
	if err := models.InitDB("moviezone.db"); err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	// Serve static files (frontend)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/details", handleDetails)
	http.HandleFunc("/api/trending", handleTrending)
	http.HandleFunc("/api/genres", handleGenres)
	http.HandleFunc("/api/watchlist", handleWatchlist)
	http.HandleFunc("/api/watchlist/watched", handleWatchlistWatched)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}