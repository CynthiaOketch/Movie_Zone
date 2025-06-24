package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// OMDBInfo holds ratings and plot from OMDB
type OMDBInfo struct {
	Title   string `json:"Title"`
	Year    string `json:"Year"`
	Plot    string `json:"Plot"`
	Ratings []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
}

// FetchOMDBInfo fetches ratings and plot from OMDB by title and year
func FetchOMDBInfo(title, year, apiKey string) (*OMDBInfo, error) {
	baseURL := "https://www.omdbapi.com/"
	params := url.Values{}
	params.Set("apikey", apiKey)
	params.Set("t", title)
	if year != "" {
		params.Set("y", year)
	}
	params.Set("plot", "short")

	fullURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDB API error: %s", resp.Status)
	}

	var info OMDBInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if info.Title == "" {
		return nil, fmt.Errorf("OMDB: No data found for title '%s'", title)
	}
	return &info, nil
}

func FetchOMDBData() error {
	
	return nil
} 