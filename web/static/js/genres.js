// Genre filtering functionality
class GenreManager {
    constructor() {
        this.currentGenreType = 'movies';
        this.currentGenreId = null;
        this.currentPage = 1;
        this.currentSort = 'popularity.desc';
        this.movieGenres = [];
        this.tvGenres = [];
        this.isLoading = false;
    }

    // Initialize genre functionality
    async init() {
        await this.loadGenres();
        this.setupEventListeners();
    }

    // Load genres from API
    async loadGenres() {
        try {
            // Load movie genres
            const movieResponse = await fetch('/api/v1/genres/movies');
            if (movieResponse.ok) {
                this.movieGenres = await movieResponse.json();
            }

            // Load TV genres
            const tvResponse = await fetch('/api/v1/genres/tv');
            if (tvResponse.ok) {
                this.tvGenres = await tvResponse.json();
            }

            this.renderGenreButtons();
        } catch (error) {
            console.error('Error loading genres:', error);
            app.showToast('Failed to load genres', 'error');
        }
    }

    // Setup event listeners
    setupEventListeners() {
        // Genre type buttons
        document.querySelectorAll('[data-genre-type]').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const type = e.target.dataset.genreType;
                this.setGenreType(type);
            });
        });

        // Sort dropdown
        const sortSelect = document.getElementById('genre-sort');
        if (sortSelect) {
            sortSelect.addEventListener('change', () => {
                this.applyGenreSort();
            });
        }
    }

    // Set genre type (movies or tv)
    setGenreType(type) {
        this.currentGenreType = type;
        this.currentGenreId = null;
        this.currentPage = 1;

        // Update active button
        document.querySelectorAll('[data-genre-type]').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.genreType === type);
        });

        this.renderGenreButtons();
        this.clearResults();
    }

    // Render genre buttons
    renderGenreButtons() {
        const container = document.getElementById('genre-buttons');
        if (!container) return;

        const genres = this.currentGenreType === 'movies' ? this.movieGenres : this.tvGenres;
        
        if (genres.length === 0) {
            container.innerHTML = '<div class="genre-loading">Loading genres...</div>';
            return;
        }

        container.innerHTML = genres.map(genre => `
            <button class="genre-btn" data-genre-id="${genre.id}" onclick="genreManager.selectGenre(${genre.id}, '${genre.name}')">
                ${genre.name}
            </button>
        `).join('');
    }

    // Select a genre
    async selectGenre(genreId, genreName) {
        if (this.isLoading) return;

        this.currentGenreId = genreId;
        this.currentPage = 1;

        // Update active genre button
        document.querySelectorAll('.genre-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.genreId == genreId);
        });

        // Update results header
        const titleElement = document.getElementById('genre-results-title');
        const countElement = document.getElementById('genre-results-count');
        
        if (titleElement) titleElement.textContent = `${genreName} ${this.currentGenreType === 'movies' ? 'Movies' : 'TV Shows'}`;
        if (countElement) countElement.textContent = 'Loading...';

        await this.loadGenreResults();
    }

    // Load results for selected genre
    async loadGenreResults() {
        if (!this.currentGenreId || this.isLoading) return;

        this.isLoading = true;
        this.showLoading();

        try {
            const params = new URLSearchParams({
                page: this.currentPage,
                sort_by: this.currentSort
            });

            const response = await fetch(`/api/v1/discover/genre/${this.currentGenreId}?${params}`);
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            this.renderResults(data);
            this.renderPagination(data);

        } catch (error) {
            console.error('Error loading genre results:', error);
            app.showToast('Failed to load results', 'error');
            this.clearResults();
        } finally {
            this.isLoading = false;
            this.hideLoading();
        }
    }

    // Render results
    renderResults(data) {
        const grid = document.getElementById('genre-results-grid');
        const countElement = document.getElementById('genre-results-count');
        
        if (!grid) return;

        if (!data.results || data.results.length === 0) {
            grid.innerHTML = '<div class="no-results">No results found for this genre.</div>';
            if (countElement) countElement.textContent = 'No results found';
            return;
        }

        // Update count
        if (countElement) {
            countElement.textContent = `${data.total_results.toLocaleString()} results found`;
        }

        // Render movie/TV cards
        grid.innerHTML = data.results.map(item => {
            const isMovie = this.currentGenreType === 'movies';
            const title = isMovie ? item.title : item.name;
            const date = isMovie ? item.release_date : item.first_air_date;
            const year = date ? new Date(date).getFullYear() : 'N/A';
            const posterUrl = item.poster_path
                ? `https://image.tmdb.org/t/p/w500${item.poster_path}`
                : '/static/images/placeholder.svg';

            return `
                <div class="movie-card" onclick="showMovieDetails(${item.id}, '${this.currentGenreType === 'movies' ? 'movie' : 'tv'}')">
                    <img src="${posterUrl}" alt="${title}" class="movie-poster" loading="lazy">
                    <div class="movie-info">
                        <div class="movie-title">${title}</div>
                        <div class="movie-year">${year}</div>
                        ${item.vote_average ? `
                            <div class="movie-rating">
                                <i class="fas fa-star rating-star"></i>
                                <span class="rating-value">${item.vote_average.toFixed(1)}</span>
                            </div>
                        ` : ''}
                    </div>
                    <button class="watchlist-btn" onclick="event.stopPropagation(); toggleWatchlist(${item.id}, '${this.currentGenreType === 'movies' ? 'movie' : 'tv'}', '${title}', '${posterUrl}')">
                        <i class="fas fa-plus"></i>
                    </button>
                </div>
            `;
        }).join('');
    }

    // Render pagination
    renderPagination(data) {
        const container = document.getElementById('genre-pagination');
        if (!container) return;

        const totalPages = Math.min(data.total_pages, 500); // TMDB limit
        const currentPage = data.page;

        if (totalPages <= 1) {
            container.innerHTML = '';
            return;
        }

        let paginationHTML = '';

        // Previous button
        if (currentPage > 1) {
            paginationHTML += `<button onclick="genreManager.goToPage(${currentPage - 1})" class="pagination-btn">Previous</button>`;
        }

        // Page numbers
        const startPage = Math.max(1, currentPage - 2);
        const endPage = Math.min(totalPages, currentPage + 2);

        if (startPage > 1) {
            paginationHTML += `<button onclick="genreManager.goToPage(1)" class="pagination-btn">1</button>`;
            if (startPage > 2) {
                paginationHTML += `<span class="pagination-ellipsis">...</span>`;
            }
        }

        for (let i = startPage; i <= endPage; i++) {
            paginationHTML += `<button onclick="genreManager.goToPage(${i})" class="pagination-btn ${i === currentPage ? 'active' : ''}">${i}</button>`;
        }

        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                paginationHTML += `<span class="pagination-ellipsis">...</span>`;
            }
            paginationHTML += `<button onclick="genreManager.goToPage(${totalPages})" class="pagination-btn">${totalPages}</button>`;
        }

        // Next button
        if (currentPage < totalPages) {
            paginationHTML += `<button onclick="genreManager.goToPage(${currentPage + 1})" class="pagination-btn">Next</button>`;
        }

        container.innerHTML = paginationHTML;
    }

    // Go to specific page
    async goToPage(page) {
        if (this.isLoading || page === this.currentPage) return;
        
        this.currentPage = page;
        await this.loadGenreResults();
        
        // Scroll to top of results
        document.getElementById('genre-results').scrollIntoView({ behavior: 'smooth' });
    }

    // Apply sort filter
    async applyGenreSort() {
        const sortSelect = document.getElementById('genre-sort');
        if (!sortSelect) return;

        this.currentSort = sortSelect.value;
        this.currentPage = 1;

        if (this.currentGenreId) {
            await this.loadGenreResults();
        }
    }

    // Clear results
    clearResults() {
        const grid = document.getElementById('genre-results-grid');
        const titleElement = document.getElementById('genre-results-title');
        const countElement = document.getElementById('genre-results-count');
        const pagination = document.getElementById('genre-pagination');

        if (grid) grid.innerHTML = '';
        if (titleElement) titleElement.textContent = 'Select a Genre';
        if (countElement) countElement.textContent = 'Choose a genre to explore';
        if (pagination) pagination.innerHTML = '';
    }

    // Show loading state
    showLoading() {
        const grid = document.getElementById('genre-results-grid');
        if (grid) {
            grid.innerHTML = '<div class="loading-placeholder">Loading...</div>';
        }
    }

    // Hide loading state
    hideLoading() {
        // Loading state is cleared when results are rendered
    }
}

// Global functions for navigation
function setGenreType(type) {
    if (window.genreManager) {
        window.genreManager.setGenreType(type);
    }
}

function applyGenreSort() {
    if (window.genreManager) {
        window.genreManager.applyGenreSort();
    }
}

function showGenres() {
    // Hide all sections
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // Show genre section
    const genreSection = document.getElementById('genre-section');
    if (genreSection) {
        genreSection.classList.add('active');
    }
    
    // Update navigation
    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
    });
    
    // Find and activate the genres nav link
    const genresLink = document.querySelector('.nav-link[onclick="showGenres()"]');
    if (genresLink) {
        genresLink.classList.add('active');
    }
    
    // Initialize genre manager if not already done
    if (!window.genreManager) {
        window.genreManager = new GenreManager();
        window.genreManager.init();
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.genreManager = new GenreManager();
    window.genreManager.init();
});
