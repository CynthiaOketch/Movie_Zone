const trendingSection = document.getElementById('trending-section');
const trendingResults = document.getElementById('trending-results');
const searchSection = document.getElementById('search-section');
const searchResults = document.getElementById('search-results');
const searchForm = document.getElementById('search-form');
const searchInput = document.getElementById('search-input');
const searchType = document.getElementById('search-type');
const watchlistSection = document.getElementById('watchlist-section');
const watchlistResults = document.getElementById('watchlist-results');
const navTrending = document.getElementById('nav-trending');
const navSearch = document.getElementById('nav-search');
const navWatchlist = document.getElementById('nav-watchlist');

function createCard(item, inWatchlist, isWatchlistPage) {
  const title = item.title || item.name || 'Untitled';
  const poster = item.poster_path
    ? `https://image.tmdb.org/t/p/w300${item.poster_path}`
    : (item.poster || 'https://via.placeholder.com/140x210?text=No+Image');
  const ratings = item.omdb_ratings && item.omdb_ratings.length
    ? item.omdb_ratings.map(r => `${r.Source}: ${r.Value}`).join(' | ')
    : 'No ratings';
  let btn = '';
  let watchedBadge = '';
  let cardClass = '';
  let detailsBtn = `<button class="details-btn" data-id="${item.id || item.tmdb_id}" data-type="${item.media_type || item.type}" aria-label="Show details for ${title}">Details</button>`;
  if (isWatchlistPage) {
    if (item.watched) {
      cardClass = 'watched-card';
      watchedBadge = '<span class="watched-badge" title="Watched" aria-label="Watched">✔</span>';
    }
    btn = `<button class="watchlist-btn remove" data-action="remove" data-id="${item.tmdb_id}" data-type="${item.type}" aria-label="Remove ${title} from watchlist">Remove</button>`;
    btn += `<button class="watched-btn${item.watched ? '' : ' unwatched'}" data-action="watched" data-id="${item.tmdb_id}" data-type="${item.type}" data-watched="${item.watched ? '1' : '0'}" aria-label="${item.watched ? 'Mark as unwatched' : 'Mark as watched'} for ${title}">${item.watched ? 'Watched' : 'Mark as Watched'}</button>`;
  } else {
    btn = inWatchlist
      ? `<button class="watchlist-btn remove" data-action="remove" data-id="${item.id}" data-type="${item.media_type}" aria-label="Remove ${title} from watchlist">Remove</button>`
      : `<button class="watchlist-btn" data-action="add" data-id="${item.id}" data-type="${item.media_type}" data-title="${title}" aria-label="Add ${title} to watchlist">Add to Watchlist</button>`;
  }
  return `
    <div class="card ${cardClass}">
      ${watchedBadge}
      <img src="${poster}" alt="${title}">
      <h3>${title}</h3>
      <div class="ratings">${ratings}</div>
      ${btn}
      ${detailsBtn}
    </div>
  `;
}

function renderResults(container, items, watchlist=[]) {
  try {
    if (!items.length) {
      container.innerHTML = '<div class="centered" style="color:#aaa;font-size:1.2rem;">🎬<br>No results found.</div>';
      return;
    }
    const watchlistIds = new Set(watchlist.map(i => i.tmdb_id));
    container.innerHTML = items.map(item => createCard(item, watchlistIds.has(String(item.id)), false)).join('');
    container.querySelectorAll('.card').forEach(card => card.classList.add('fade-in'));
    container.querySelectorAll('.watchlist-btn').forEach(btn => {
      btn.onclick = async e => {
        const id = btn.getAttribute('data-id');
        const type = btn.getAttribute('data-type');
        const title = btn.getAttribute('data-title');
        btn.disabled = true;
        if (btn.classList.contains('remove')) {
          btn.textContent = 'Removed!';
          await removeFromWatchlist({ id, media_type: type });
          setTimeout(() => {
            fetchWatchlist();
            fetchTrending();
          }, 1000);
        } else {
          btn.textContent = 'Added!';
          await addToWatchlist({ id, media_type: type, title });
          setTimeout(() => {
            fetchWatchlist();
            fetchTrending();
          }, 1000);
        }
      };
    });
    container.querySelectorAll('.details-btn').forEach(btn => {
      btn.onclick = async e => {
        const id = btn.getAttribute('data-id');
        const type = btn.getAttribute('data-type');
        showMovieDetails(id, type);
      };
    });
  } catch (e) {
    console.error('Error in renderResults:', e);
    container.innerHTML = '<div class="centered error-message">Error rendering results.</div>';
  }
}