package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)


// SearchResult represents a unified movie/TV search result
type SearchResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Name        string `json:"name"`
	MediaType   string `json:"media_type"`
	Overview    string `json:"overview"`
	PosterPath  string `json:"poster_path"`
	ReleaseDate string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
	OMDBRatings []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"omdb_ratings,omitempty"`
	OMDBPlot string `json:"omdb_plot,omitempty"`
}

// SearchTMDB searches TMDB for movies or TV shows
func SearchTMDB(query, mediaType, apiKey string, page int) ([]SearchResult, int, error) {
	baseURL := "https://api.themoviedb.org/3/search/"
	if mediaType != "movie" && mediaType != "tv" {
		mediaType = "movie"
	}
	endpoint := fmt.Sprintf("%s%s", baseURL, mediaType)

	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("query", query)
	params.Set("language", "en-US")
	params.Set("page", fmt.Sprintf("%d", page))

	fullURL := endpoint + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, 1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 1, fmt.Errorf("TMDB API error: %s", resp.Status)
	}

	var tmdbResp struct {
		Results    []SearchResult `json:"results"`
		TotalPages int            `json:"total_pages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, 1, err
	}

	// Set MediaType for each result
	for i := range tmdbResp.Results {
		tmdbResp.Results[i].MediaType = mediaType
	}

	return tmdbResp.Results, tmdbResp.TotalPages, nil
}

func FetchTMDBData() error {
	
	return nil
}

// FetchTMDBDetails fetches details for a movie or TV show by TMDB ID and type
func FetchTMDBDetails(id, mediaType, apiKey string) (map[string]interface{}, error) {
	if mediaType != "movie" && mediaType != "tv" {
		mediaType = "movie"
	}
	endpoint := "https://api.themoviedb.org/3/" + mediaType + "/" + id
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("language", "en-US")
	fullURL := endpoint + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %s", resp.Status)
	}
	var details map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}
	return details, nil
}

// FetchTMDBTrending fetches trending movies or TV shows from TMDB
func FetchTMDBTrending(mediaType, apiKey string, page int) ([]SearchResult, int, error) {
	if mediaType != "movie" && mediaType != "tv" {
		mediaType = "movie"
	}
	endpoint := "https://api.themoviedb.org/3/trending/" + mediaType + "/day"
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("language", "en-US")
	params.Set("page", fmt.Sprintf("%d", page))
	fullURL := endpoint + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, 1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, 1, fmt.Errorf("TMDB API error: %s", resp.Status)
	}
	var tmdbResp struct {
		Results    []SearchResult `json:"results"`
		TotalPages int            `json:"total_pages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tmdbResp); err != nil {
		return nil, 1, err
	}
	for i := range tmdbResp.Results {
		tmdbResp.Results[i].MediaType = mediaType
	}
	return tmdbResp.Results, tmdbResp.TotalPages, nil
}

// FetchTMDBTrailer fetches the YouTube trailer key for a movie or TV show by TMDB ID and type
func FetchTMDBTrailer(id, mediaType, apiKey string) (string, error) {
	if mediaType != "movie" && mediaType != "tv" {
		mediaType = "movie"
	}
	endpoint := "https://api.themoviedb.org/3/" + mediaType + "/" + id + "/videos"
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("language", "en-US")
	fullURL := endpoint + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("TMDB API error: %s", resp.Status)
	}
	var videoResp struct {
		Results []struct {
			Key  string `json:"key"`
			Site string `json:"site"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&videoResp); err != nil {
		return "", err
	}
	for _, v := range videoResp.Results {
		if v.Site == "YouTube" && v.Type == "Trailer" {
			return v.Key, nil
		}
	}
	return "", fmt.Errorf("No YouTube trailer found")
} 