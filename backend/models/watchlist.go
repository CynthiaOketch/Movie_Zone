package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Watchlist model and DB logic will go here

type WatchlistItem struct {
	ID      int    `json:"id"`
	TMDBID  string `json:"tmdb_id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Watched bool   `json:"watched"`
}

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	createTable := `CREATE TABLE IF NOT EXISTS watchlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tmdb_id TEXT NOT NULL,
		type TEXT NOT NULL,
		title TEXT,
		watched INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createTable)
	return err
}

func AddToWatchlist(item WatchlistItem) error {
	_, err := db.Exec(`INSERT INTO watchlist (tmdb_id, type, title, watched) VALUES (?, ?, ?, ?)`,
		item.TMDBID, item.Type, item.Title, item.Watched)
	return err
}

func RemoveFromWatchlist(tmdbID, mediaType string) error {
	_, err := db.Exec(`DELETE FROM watchlist WHERE tmdb_id = ? AND type = ?`, tmdbID, mediaType)
	return err
}

func ListWatchlist() ([]WatchlistItem, error) {
	rows, err := db.Query(`SELECT id, tmdb_id, type, title, watched FROM watchlist`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WatchlistItem
	for rows.Next() {
		var item WatchlistItem
		var watchedInt int
		if err := rows.Scan(&item.ID, &item.TMDBID, &item.Type, &item.Title, &watchedInt); err != nil {
			return nil, err
		}
		item.Watched = watchedInt != 0
		items = append(items, item)
	}
	return items, nil
}

func SetWatched(tmdbID, mediaType string, watched bool) error {
	watchedInt := 0
	if watched {
		watchedInt = 1
	}
	_, err := db.Exec(`UPDATE watchlist SET watched = ? WHERE tmdb_id = ? AND type = ?`, watchedInt, tmdbID, mediaType)
	return err
} 