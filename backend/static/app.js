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
  let detailsBtn = `<button class="details-btn" data-id="${item.id || item.tmdb_id}" data-type="${item.media_type || item.type}" aria-label="Show details for ${title}">+</button>`;
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

function renderWatchlist(container, items) {
  if (!items.length) {
    container.innerHTML = '<div class="centered" style="color:#aaa;font-size:1.2rem;">🍿<br>Your watchlist is empty.</div>';
    return;
  }
  container.innerHTML = items.map(item => createCard(item, true, true)).join('');
  container.querySelectorAll('.watchlist-btn.remove').forEach(btn => {
    btn.onclick = async e => {
      const id = btn.getAttribute('data-id');
      const type = btn.getAttribute('data-type');
      btn.disabled = true;
      btn.textContent = 'Removed!';
      await removeFromWatchlist({ tmdb_id: id, type });
      setTimeout(() => {
        fetchWatchlist();
        fetchTrending();
      }, 1000);
    };
  });
  container.querySelectorAll('.watched-btn').forEach(btn => {
    btn.onclick = async e => {
      const id = btn.getAttribute('data-id');
      const type = btn.getAttribute('data-type');
      const watched = btn.getAttribute('data-watched') === '1';
      btn.disabled = true;
      btn.textContent = watched ? 'Mark as Watched' : 'Watched!';
      await setWatched({ tmdb_id: id, type }, !watched);
      setTimeout(() => {
        fetchWatchlist();
      }, 1000);
    };
  });
  container.querySelectorAll('.card').forEach(card => card.classList.add('fade-in'));
}

// --- Navigation logic ---
function showSection(section) {
  trendingSection.style.display = 'none';
  searchSection.style.display = 'none';
  watchlistSection.style.display = 'none';
  searchForm.style.display = 'none';
  navTrending.classList.remove('active');
  navSearch.classList.remove('active');
  navWatchlist.classList.remove('active');
  trendingSection.classList.remove('fade-in');
  searchSection.classList.remove('fade-in');
  watchlistSection.classList.remove('fade-in');
  if (section === 'trending') {
    trendingSection.style.display = '';
    trendingSection.classList.add('fade-in');
    navTrending.classList.add('active');
  } else if (section === 'search') {
    searchSection.style.display = '';
    searchSection.classList.add('fade-in');
    searchForm.style.display = '';
    navSearch.classList.add('active');
  } else if (section === 'watchlist') {
    watchlistSection.style.display = '';
    watchlistSection.classList.add('fade-in');
    navWatchlist.classList.add('active');
    fetchWatchlist();
  }
}
navTrending.onclick = () => showSection('trending');
navSearch.onclick = () => showSection('search');
navWatchlist.onclick = () => showSection('watchlist');

// --- Watchlist API ---
function showSpinner(container) {
  container.innerHTML = '<div class="centered"><div class="spinner"></div></div>';
}
function showError(container, message, retryFn) {
  container.innerHTML = `<div class="centered"><div class="error-message">${message}</div><button class="retry-btn">Retry</button></div>`;
  container.querySelector('.retry-btn').onclick = retryFn;
}

async function fetchWatchlist() {
  showSpinner(watchlistResults);
  try {
    const res = await fetch('/api/watchlist');
    const data = await res.json();
    renderWatchlist(watchlistResults, data);
  } catch (err) {
    showError(watchlistResults, 'Error loading watchlist.', fetchWatchlist);
  }
}
async function addToWatchlist(item) {
  await fetch('/api/watchlist', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tmdb_id: item.id, type: item.media_type, title: item.title || item.name })
  });
}
async function removeFromWatchlist(item) {
  await fetch(`/api/watchlist?tmdb_id=${encodeURIComponent(item.tmdb_id || item.id)}&type=${item.type || item.media_type}`, {
    method: 'DELETE'
  });
}
async function setWatched(item, watched) {
  await fetch('/api/watchlist/watched', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tmdb_id: item.tmdb_id, type: item.type, watched })
  });
}

// --- Trending and Search logic ---

let currentWatchlist = [];
let trendingPage = 1;
let trendingTotalPages = 1;
let searchPage = 1;
let searchTotalPages = 1;

function renderPagination(container, page, totalPages, onPageChange) {
  const pagination = document.createElement('div');
  pagination.className = 'pagination-controls';
  pagination.style.display = 'flex';
  pagination.style.justifyContent = 'center';
  pagination.style.alignItems = 'center';
  pagination.style.gap = '1rem';
  pagination.style.margin = '1.2rem 0 0.5rem 0';

  const prevBtn = document.createElement('button');
  prevBtn.textContent = 'Prev';
  prevBtn.disabled = page <= 1;
  prevBtn.onclick = () => onPageChange(page - 1);

  const nextBtn = document.createElement('button');
  nextBtn.textContent = 'Next';
  nextBtn.disabled = page >= totalPages;
  nextBtn.onclick = () => onPageChange(page + 1);

  const pageInfo = document.createElement('span');
  pageInfo.textContent = `Page ${page} of ${totalPages}`;

  pagination.appendChild(prevBtn);
  pagination.appendChild(pageInfo);
  pagination.appendChild(nextBtn);

  container.parentNode.insertBefore(pagination, container.nextSibling);
}

function clearPagination(container) {
  let next = container.nextSibling;
  while (next && next.classList && next.classList.contains('pagination-controls')) {
    const toRemove = next;
    next = next.nextSibling;
    toRemove.remove();
  }
}

async function fetchTrending(page = 1) {
  trendingPage = page;
  showSpinner(trendingResults);
  clearPagination(trendingResults);
  try {
    const [res, wlRes] = await Promise.all([
      fetch(`/api/trending?type=movie&page=${page}`),
      fetch('/api/watchlist')
    ]);
    if (!res.ok) throw new Error('Trending fetch failed: ' + res.status + ' ' + res.statusText);
    if (!wlRes.ok) throw new Error('Watchlist fetch failed: ' + wlRes.status + ' ' + wlRes.statusText);
    const data = await res.json();
    trendingTotalPages = data.total_pages || 1;
    currentWatchlist = await wlRes.json();
    renderResults(trendingResults, data.results, currentWatchlist);
    renderPagination(trendingResults, data.page, trendingTotalPages, fetchTrending);
  } catch (err) {
    showError(trendingResults, 'Error loading trending movies.', () => fetchTrending(page));
  }
}

async function fetchSearch(query, type, page = 1) {
  searchPage = page;
  showSpinner(searchResults);
  clearPagination(searchResults);
  try {
    const [res, wlRes] = await Promise.all([
      fetch(`/api/search?q=${encodeURIComponent(query)}&type=${type}&page=${page}`),
      fetch('/api/watchlist')
    ]);
    const data = await res.json();
    searchTotalPages = data.total_pages || 1;
    const watchlist = await wlRes.json();
    renderResults(searchResults, data.results, watchlist);
    renderPagination(searchResults, data.page, searchTotalPages, p => fetchSearch(query, type, p));
  } catch (err) {
    showError(searchResults, 'Error loading search results.', () => fetchSearch(query, type, page));
  }
}

searchForm.addEventListener('submit', e => {
  e.preventDefault();
  const query = searchInput.value.trim();
  const type = searchType.value;
  if (!query) return;
  showSection('search');
  searchPage = 1;
  fetchSearch(query, type, 1);
});

navTrending.onclick = () => {
  showSection('trending');
  trendingPage = 1;
  fetchTrending(1);
};

// On load
showSection('trending');
fetchTrending(); 
fetchTrending(); 

async function showMovieDetails(id, type) {
  const modal = document.getElementById('details-modal');
  const content = document.getElementById('details-content');
  content.innerHTML = '<div class="centered"><div class="spinner"></div></div>';
  modal.style.display = 'flex';
  try {
    const res = await fetch(`/api/details?id=${id}&type=${type}`);
    const data = await res.json();
    content.innerHTML = `
      <h2>${data.title || data.name}</h2>
      <img src="https://image.tmdb.org/t/p/w300${data.poster_path}" alt="${data.title || data.name}">
      <div class="ratings">${(data.omdb_ratings || []).map(r => `${r.Source}: ${r.Value}`).join(' | ')}</div>
      <div class="plot">${data.omdb_plot || data.overview || ''}</div>
      <button id="trailer-btn" class="trailer-btn">Watch Trailer</button>
      <div id="trailer-container"></div>
    `;
    document.getElementById('trailer-btn').onclick = async () => {
      const trailerContainer = document.getElementById('trailer-container');
      trailerContainer.innerHTML = '<div class="centered"><div class="spinner"></div></div>';
      try {
        const trailerRes = await fetch(`/api/trailer?id=${id}&type=${type}`);
        const trailerData = await trailerRes.json();
        if (trailerData.key) {
          trailerContainer.innerHTML = `<iframe width="100%" height="315" src="https://www.youtube.com/embed/${trailerData.key}" frameborder="0" allowfullscreen></iframe>`;
        } else {
          trailerContainer.innerHTML = '<div class="error-message">Trailer not found.</div>';
        }
      } catch (err) {
        trailerContainer.innerHTML = '<div class="error-message">Error loading trailer.</div>';
      }
    };
  } catch (err) {
    content.innerHTML = '<div class="error-message">Error loading details.</div>';
  }
}
document.addEventListener('DOMContentLoaded', function() {
  document.getElementById('close-modal').onclick = () => {
    document.getElementById('details-modal').style.display = 'none';
    document.getElementById('details-content').innerHTML = '';
  };
}); 