// Main application JavaScript
class MovieDiscoveryApp {
    constructor() {
        this.currentSection = 'home';
        this.currentSearchType = 'movies';
        this.currentPage = 1;
        this.searchQuery = '';
        this.debounceTimer = null;
        
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadTheme();
        this.showHome();
    }

    setupEventListeners() {
        // Search input with debouncing
        const searchInput = document.getElementById('search-input');
        searchInput.addEventListener('input', (e) => {
            clearTimeout(this.debounceTimer);
            this.debounceTimer = setTimeout(() => {
                this.handleSearch(e.target.value);
            }, 500);
        });

        // Enter key for search
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                clearTimeout(this.debounceTimer);
                this.handleSearch(e.target.value);
            }
        });

        // Navigation links
        document.querySelectorAll('.nav-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                this.updateActiveNavLink(link);
            });
        });

        // Search type buttons
        document.querySelectorAll('.search-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                this.updateActiveSearchBtn(btn);
                if (this.searchQuery) {
                    this.performSearch();
                }
            });
        });
    }

    handleSearch(query) {
        this.searchQuery = query.trim();
        if (this.searchQuery.length >= 2) {
            this.currentPage = 1;
            this.performSearch();
        } else if (this.searchQuery.length === 0) {
            this.hideSearchResults();
        }
    }

    performSearch() {
        if (this.currentSearchType === 'movies') {
            this.searchMovies(this.searchQuery, this.currentPage);
        } else {
            this.searchTVShows(this.searchQuery, this.currentPage);
        }
    }

    updateActiveNavLink(activeLink) {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        activeLink.classList.add('active');
    }

    updateActiveSearchBtn(activeBtn) {
        document.querySelectorAll('.search-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        activeBtn.classList.add('active');
        this.currentSearchType = activeBtn.dataset.type;
    }

    showSection(sectionId) {
        document.querySelectorAll('.section').forEach(section => {
            section.classList.remove('active');
        });
        document.getElementById(sectionId).classList.add('active');
        this.currentSection = sectionId.replace('-section', '');
    }

    showHome() {
        this.showSection('home-section');
        this.hideSearchResults();
    }

    showTrending() {
        this.showSection('trending-section');
        this.loadTrending('week');
    }

    showWatchlist() {
        this.showSection('watchlist-section');
        this.loadWatchlist();
    }

    hideSearchResults() {
        document.getElementById('search-results').style.display = 'none';
    }

    showSearchResults() {
        document.getElementById('search-results').style.display = 'block';
    }

    showLoading() {
        document.getElementById('loading').style.display = 'flex';
    }

    hideLoading() {
        document.getElementById('loading').style.display = 'none';
    }

    showToast(message, type = 'info') {
        console.log('showToast called:', { message, type });

        const toastContainer = document.getElementById('toast-container');
        if (!toastContainer) {
            console.error('Toast container not found!');
            return;
        }

        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.textContent = message;

        console.log('Toast element created:', toast);
        toastContainer.appendChild(toast);
        console.log('Toast added to container');

        setTimeout(() => {
            toast.remove();
            console.log('Toast removed');
        }, 5000);
    }

    async apiRequest(url, options = {}) {
        try {
            this.showLoading();
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            this.showToast('Request failed. Please try again.', 'error');
            throw error;
        } finally {
            this.hideLoading();
        }
    }

    createMovieCard(item, type = 'movie') {
        const card = document.createElement('div');
        card.className = 'movie-card';
        card.style.position = 'relative';

        const posterUrl = item.poster_path
            ? `https://image.tmdb.org/t/p/w300${item.poster_path}`
            : '/static/images/placeholder.svg';

        const title = item.title || item.name || 'Unknown Title';
        const year = item.release_date || item.first_air_date || '';
        const rating = item.vote_average || 0;

        card.innerHTML = `
            <img src="${posterUrl}" alt="${title}" class="movie-poster"
                 onerror="this.src='/static/images/placeholder.svg'">
            <button class="watchlist-btn" onclick="app.toggleWatchlist('${item.id}', '${type}', '${title}', '${posterUrl}')">
                <i class="fas fa-bookmark"></i>
            </button>
            <div class="movie-info">
                <div class="movie-title">${title}</div>
                <div class="movie-year">${year ? year.substring(0, 4) : ''}</div>
                ${rating > 0 ? `
                    <div class="movie-rating">
                        <i class="fas fa-star rating-star"></i>
                        <span class="rating-value">${rating.toFixed(1)}</span>
                    </div>
                ` : ''}
            </div>
        `;

        card.addEventListener('click', (e) => {
            if (!e.target.closest('.watchlist-btn')) {
                this.showMovieDetails(item.id, type);
            }
        });

        return card;
    }

    createPagination(currentPage, totalPages, onPageChange) {
        const pagination = document.getElementById('pagination');
        pagination.innerHTML = '';

        if (totalPages <= 1) return;

        // Previous button
        const prevBtn = document.createElement('button');
        prevBtn.textContent = 'Previous';
        prevBtn.disabled = currentPage === 1;
        prevBtn.addEventListener('click', () => onPageChange(currentPage - 1));
        pagination.appendChild(prevBtn);

        // Page numbers
        const startPage = Math.max(1, currentPage - 2);
        const endPage = Math.min(totalPages, currentPage + 2);

        if (startPage > 1) {
            const firstBtn = document.createElement('button');
            firstBtn.textContent = '1';
            firstBtn.addEventListener('click', () => onPageChange(1));
            pagination.appendChild(firstBtn);

            if (startPage > 2) {
                const ellipsis = document.createElement('span');
                ellipsis.textContent = '...';
                ellipsis.style.padding = '10px';
                pagination.appendChild(ellipsis);
            }
        }

        for (let i = startPage; i <= endPage; i++) {
            const pageBtn = document.createElement('button');
            pageBtn.textContent = i;
            pageBtn.className = i === currentPage ? 'active' : '';
            pageBtn.addEventListener('click', () => onPageChange(i));
            pagination.appendChild(pageBtn);
        }

        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                const ellipsis = document.createElement('span');
                ellipsis.textContent = '...';
                ellipsis.style.padding = '10px';
                pagination.appendChild(ellipsis);
            }

            const lastBtn = document.createElement('button');
            lastBtn.textContent = totalPages;
            lastBtn.addEventListener('click', () => onPageChange(totalPages));
            pagination.appendChild(lastBtn);
        }

        // Next button
        const nextBtn = document.createElement('button');
        nextBtn.textContent = 'Next';
        nextBtn.disabled = currentPage === totalPages;
        nextBtn.addEventListener('click', () => onPageChange(currentPage + 1));
        pagination.appendChild(nextBtn);
    }

    loadTheme() {
        const savedTheme = localStorage.getItem('theme') || 'light';
        document.documentElement.setAttribute('data-theme', savedTheme);
        this.updateThemeIcon(savedTheme);
    }

    toggleTheme() {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        
        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        this.updateThemeIcon(newTheme);
    }

    updateThemeIcon(theme) {
        const icon = document.querySelector('#theme-toggle-btn i');
        icon.className = theme === 'dark' ? 'fas fa-sun' : 'fas fa-moon';
    }
}

// Global functions for HTML onclick handlers
function showHome() {
    app.showHome();
}

function showTrending() {
    app.showTrending();
}

function showWatchlist() {
    app.showWatchlist();
}

function toggleTheme() {
    app.toggleTheme();
}

function searchMovies() {
    document.querySelector('[data-type="movies"]').click();
}

function searchTVShows() {
    document.querySelector('[data-type="tv"]').click();
}

function closeModal() {
    document.getElementById('movie-modal').style.display = 'none';
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.app = new MovieDiscoveryApp();
});
