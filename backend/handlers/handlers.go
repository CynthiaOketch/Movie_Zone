package handlers

import(
	"encoding/json"
	"net/http"
	"os"
	"moviezone/api"
	"moviezone/models"
	"io/ioutil"
	"strconv"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query().Get("q")
	mediaType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if q == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing search query parameter 'q'"}`))
		return
	}
	tmdbKey := os.Getenv("TMDB_API_KEY")
	omdbKey := os.Getenv("OMDB_API_KEY")
	results, totalPages, err := api.SearchTMDB(q, mediaType, tmdbKey, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// OMDB enrichment
	for i := range results {
		title := results[i].Title
		if title == "" {
			title = results[i].Name
		}
		year := ""
		if results[i].ReleaseDate != "" {
			year = results[i].ReleaseDate[:4]
		} else if results[i].FirstAirDate != "" {
			year = results[i].FirstAirDate[:4]
		}
		if title != "" && omdbKey != "" {
			omdbInfo, err := api.FetchOMDBInfo(title, year, omdbKey)
			if err == nil && omdbInfo != nil {
				results[i].OMDBRatings = omdbInfo.Ratings
				results[i].OMDBPlot = omdbInfo.Plot
			}
			// If OMDB fails, just skip enrichment for this result
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"page": page,
		"total_pages": totalPages,
	})
}

func HandleDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	mediaType := r.URL.Query().Get("type")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing 'id' query parameter"}`))
		return
	}
	tmdbKey := os.Getenv("TMDB_API_KEY")
	omdbKey := os.Getenv("OMDB_API_KEY")
	details, err := api.FetchTMDBDetails(id, mediaType, tmdbKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// Try to enrich with OMDB
	title := ""
	year := ""
	if t, ok := details["title"].(string); ok && t != "" {
		title = t
	} else if n, ok := details["name"].(string); ok && n != "" {
		title = n
	}
	if rd, ok := details["release_date"].(string); ok && len(rd) >= 4 {
		year = rd[:4]
	} else if fad, ok := details["first_air_date"].(string); ok && len(fad) >= 4 {
		year = fad[:4]
	}
	if title != "" && omdbKey != "" {
		omdbInfo, err := api.FetchOMDBInfo(title, year, omdbKey)
		if err == nil && omdbInfo != nil {
			details["omdb_ratings"] = omdbInfo.Ratings
			details["omdb_plot"] = omdbInfo.Plot
		}
	}
	json.NewEncoder(w).Encode(details)
}

func HandleTrending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mediaType := r.URL.Query().Get("type")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	tmdbKey := os.Getenv("TMDB_API_KEY")
	omdbKey := os.Getenv("OMDB_API_KEY")
	results, totalPages, err := api.FetchTMDBTrending(mediaType, tmdbKey, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	// OMDB enrichment (optional, can be slow)
	for i := range results {
		title := results[i].Title
		if title == "" {
			title = results[i].Name
		}
		year := ""
		if results[i].ReleaseDate != "" {
			year = results[i].ReleaseDate[:4]
		} else if results[i].FirstAirDate != "" {
			year = results[i].FirstAirDate[:4]
		}
		if title != "" && omdbKey != "" {
			omdbInfo, err := api.FetchOMDBInfo(title, year, omdbKey)
			if err == nil && omdbInfo != nil {
				results[i].OMDBRatings = omdbInfo.Ratings
				results[i].OMDBPlot = omdbInfo.Plot
			}
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
		"page": page,
		"total_pages": totalPages,
	})
}

func HandleGenres(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "genres endpoint placeholder"}`))
}

func HandleWatchlist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := models.ListWatchlist()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if items == nil {
			items = []models.WatchlistItem{}
		}
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}
		var item models.WatchlistItem
		if err := json.Unmarshal(body, &item); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}
		if item.TMDBID == "" || item.Type == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing tmdb_id or type"})
			return
		}
		if err := models.AddToWatchlist(item); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Added to watchlist"})
	case http.MethodDelete:
		tmdbID := r.URL.Query().Get("tmdb_id")
		mediaType := r.URL.Query().Get("type")
		if tmdbID == "" || mediaType == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing tmdb_id or type"})
			return
		}
		if err := models.RemoveFromWatchlist(tmdbID, mediaType); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "Removed from watchlist"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
	}
}

func HandleWatchlistWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	var req struct {
		TMDBID   string `json:"tmdb_id"`
		Type     string `json:"type"`
		Watched  bool   `json:"watched"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}
	if req.TMDBID == "" || req.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing tmdb_id or type"})
		return
	}
	if err := models.SetWatched(req.TMDBID, req.Type, req.Watched); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Watch status updated"})
} 