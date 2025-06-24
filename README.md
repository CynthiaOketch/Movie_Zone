# Movie_Zone
A comprehensive entertainment discovery platform where users can search for movies and TV shows, view detailed information, manage personal watchlists, and discover trending content.

## Features
- Search for movies and TV shows
- View trending content
- View detailed information (including OMDB ratings and plot)
- Add/remove items from your personal watchlist
- Mark items as watched/unwatched
- Pagination for trending and search results
- View trailers for movies and TV shows directly in the app

## Setup Instructions

### Backend (Go)
1. Make sure you have Go installed (version 1.18+ recommended).
2. Clone the repository and navigate to the `backend` directory:
   ```bash
   git clone https://github.com/CynthiaOketch/Movie_Zone.git
   cd Movie_Zone/backend
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Set up your environment variables (create a `.env` file or export them):
   - `TMDB_API_KEY` (required)
   - `OMDB_API_KEY` (required)
5. Run the backend server:
   ```bash
   go run .
   ```

### Frontend
- The frontend is served automatically from the `backend/static` directory. Open [http://localhost:8080](http://localhost:8080) in your browser after starting the backend.

## Troubleshooting

### Go Module Cache Permission Error
If you see an error like:
```
go: could not create module cache: mkdir /home/bocal/go: permission denied
```
This means your user does not have permission to write to the Go module cache directory. To fix this, set your `GOPATH` to a directory you own:
```bash
export GOPATH=$HOME/go
```
You can add this line to your `~/.bashrc` or `~/.profile` to make it permanent.

## Notes
- Trailers are fetched from YouTube via TMDB. If no trailer is available, an error message will be shown.
- Pagination controls appear below trending and search results.

## Contributors
- Cynthia Oketch - ([CynthiaOketch](https://github.com/CynthiaOketch))


## License
This project is licensed under the [MIT License](LICENSE).

---
Feel free to open issues or contribute!
