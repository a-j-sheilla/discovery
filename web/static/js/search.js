// Search functionality
MovieDiscoveryApp.prototype.searchMovies = async function(query, page = 1) {
    try {
        const data = await this.apiRequest(`/api/v1/search/movies?q=${encodeURIComponent(query)}&page=${page}`);
        this.displaySearchResults(data, 'movies');
        this.currentPage = page;
    } catch (error) {
        console.error('Movie search failed:', error);
        // Show demo data when API fails
        this.showDemoSearchResults(query, 'movies');
    }
};

MovieDiscoveryApp.prototype.searchTVShows = async function(query, page = 1) {
    try {
        const data = await this.apiRequest(`/api/v1/search/tv?q=${encodeURIComponent(query)}&page=${page}`);
        this.displaySearchResults(data, 'tv');
        this.currentPage = page;
    } catch (error) {
        console.error('TV show search failed:', error);
    }
};

MovieDiscoveryApp.prototype.displaySearchResults = function(data, type) {
    const resultsContainer = document.getElementById('search-results');
    const resultsGrid = document.getElementById('results-grid');
    const resultsTitle = document.getElementById('results-title');
    const resultsCount = document.getElementById('results-count');

    // Update title and count
    const typeLabel = type === 'movies' ? 'Movies' : 'TV Shows';
    resultsTitle.textContent = `${typeLabel} - "${this.searchQuery}"`;
    resultsCount.textContent = `${data.total_results} results found`;

    // Clear previous results
    resultsGrid.innerHTML = '';

    if (data.results && data.results.length > 0) {
        // Create movie cards
        data.results.forEach(item => {
            const card = this.createMovieCard(item, type);
            resultsGrid.appendChild(card);
        });

        // Create pagination
        this.createPagination(data.page, data.total_pages, (page) => {
            if (type === 'movies') {
                this.searchMovies(this.searchQuery, page);
            } else {
                this.searchTVShows(this.searchQuery, page);
            }
        });

        this.showSearchResults();
    } else {
        // No results found or API error
        const isAPIError = data.error || (data.results && data.results.length === 0 && data.total_results === 0);
        const errorMessage = isAPIError ?
            'Search temporarily unavailable. Please check API configuration.' :
            'No results found. Try searching with different keywords.';

        resultsGrid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px;">
                <i class="fas fa-${isAPIError ? 'exclamation-triangle' : 'search'}" style="font-size: 3rem; color: var(--text-secondary); margin-bottom: 20px;"></i>
                <h3 style="color: var(--text-secondary); margin-bottom: 10px;">${isAPIError ? 'Service Unavailable' : 'No Results Found'}</h3>
                <p style="color: var(--text-secondary);">${errorMessage}</p>
                ${isAPIError ? '<p style="color: var(--text-secondary); font-size: 0.9rem; margin-top: 10px;">Note: This demo requires valid TMDB API keys</p>' : ''}
            </div>
        `;
        document.getElementById('pagination').innerHTML = '';
        this.showSearchResults();
    }
};

MovieDiscoveryApp.prototype.getSearchSuggestions = async function(query) {
    if (query.length < 2) return [];
    
    try {
        // Get suggestions from both movies and TV shows
        const [movieData, tvData] = await Promise.all([
            this.apiRequest(`/api/v1/search/movies?q=${encodeURIComponent(query)}&page=1`),
            this.apiRequest(`/api/v1/search/tv?q=${encodeURIComponent(query)}&page=1`)
        ]);

        const suggestions = [];
        
        // Add top 3 movies
        if (movieData.results) {
            movieData.results.slice(0, 3).forEach(movie => {
                suggestions.push({
                    title: movie.title,
                    year: movie.release_date ? movie.release_date.substring(0, 4) : '',
                    type: 'movie',
                    id: movie.id
                });
            });
        }

        // Add top 3 TV shows
        if (tvData.results) {
            tvData.results.slice(0, 3).forEach(show => {
                suggestions.push({
                    title: show.name,
                    year: show.first_air_date ? show.first_air_date.substring(0, 4) : '',
                    type: 'tv',
                    id: show.id
                });
            });
        }

        return suggestions;
    } catch (error) {
        console.error('Failed to get suggestions:', error);
        return [];
    }
};

MovieDiscoveryApp.prototype.setupSearchSuggestions = function() {
    const searchInput = document.getElementById('search-input');
    const searchContainer = document.querySelector('.search-container');
    
    // Create suggestions dropdown
    const suggestionsDropdown = document.createElement('div');
    suggestionsDropdown.className = 'search-suggestions';
    suggestionsDropdown.style.cssText = `
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        background-color: var(--surface-color);
        border: 1px solid var(--border-color);
        border-top: none;
        border-radius: 0 0 var(--border-radius) var(--border-radius);
        box-shadow: var(--shadow);
        max-height: 300px;
        overflow-y: auto;
        z-index: 100;
        display: none;
    `;
    
    searchContainer.style.position = 'relative';
    searchContainer.appendChild(suggestionsDropdown);

    let suggestionTimer;

    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        
        clearTimeout(suggestionTimer);
        
        if (query.length >= 2) {
            suggestionTimer = setTimeout(async () => {
                const suggestions = await this.getSearchSuggestions(query);
                this.displaySearchSuggestions(suggestions, suggestionsDropdown);
            }, 300);
        } else {
            suggestionsDropdown.style.display = 'none';
        }
    });

    // Hide suggestions when clicking outside
    document.addEventListener('click', (e) => {
        if (!searchContainer.contains(e.target)) {
            suggestionsDropdown.style.display = 'none';
        }
    });
};

MovieDiscoveryApp.prototype.displaySearchSuggestions = function(suggestions, dropdown) {
    dropdown.innerHTML = '';
    
    if (suggestions.length === 0) {
        dropdown.style.display = 'none';
        return;
    }

    suggestions.forEach(suggestion => {
        const item = document.createElement('div');
        item.className = 'suggestion-item';
        item.style.cssText = `
            padding: 12px 20px;
            cursor: pointer;
            border-bottom: 1px solid var(--border-color);
            transition: var(--transition);
            display: flex;
            justify-content: space-between;
            align-items: center;
        `;
        
        item.innerHTML = `
            <div>
                <div style="font-weight: 500; color: var(--text-color);">${suggestion.title}</div>
                <div style="font-size: 0.8rem; color: var(--text-secondary);">
                    ${suggestion.year} • ${suggestion.type === 'movie' ? 'Movie' : 'TV Show'}
                </div>
            </div>
            <i class="fas fa-arrow-right" style="color: var(--text-secondary);"></i>
        `;

        item.addEventListener('mouseenter', () => {
            item.style.backgroundColor = 'var(--background-color)';
        });

        item.addEventListener('mouseleave', () => {
            item.style.backgroundColor = 'transparent';
        });

        item.addEventListener('click', () => {
            document.getElementById('search-input').value = suggestion.title;
            this.searchQuery = suggestion.title;
            this.currentSearchType = suggestion.type === 'movie' ? 'movies' : 'tv';
            
            // Update search type button
            document.querySelectorAll('.search-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            document.querySelector(`[data-type="${this.currentSearchType}"]`).classList.add('active');
            
            this.performSearch();
            dropdown.style.display = 'none';
        });

        dropdown.appendChild(item);
    });

    dropdown.style.display = 'block';
};

// Advanced search filters
MovieDiscoveryApp.prototype.setupAdvancedSearch = function() {
    // This could be expanded to include genre filters, year ranges, rating filters, etc.
    // For now, we'll keep it simple with the basic search functionality
};

// Demo data for when API is not available
MovieDiscoveryApp.prototype.showDemoSearchResults = function(query, type) {
    const demoMovies = [
        {
            id: 550,
            title: "Fight Club",
            overview: "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
            release_date: "1999-10-15",
            poster_path: null,
            vote_average: 8.4,
            popularity: 61.416
        },
        {
            id: 13,
            title: "Forrest Gump",
            overview: "A man with a low IQ has accomplished great things in his life and been present during significant historic events.",
            release_date: "1994-07-06",
            poster_path: null,
            vote_average: 8.5,
            popularity: 75.123
        },
        {
            id: 27205,
            title: "Inception",
            overview: "Cobb, a skilled thief who commits corporate espionage by infiltrating the subconscious of his targets.",
            release_date: "2010-07-16",
            poster_path: null,
            vote_average: 8.4,
            popularity: 151.489
        }
    ];

    const demoData = {
        page: 1,
        results: demoMovies.filter(movie =>
            movie.title.toLowerCase().includes(query.toLowerCase())
        ),
        total_pages: 1,
        total_results: demoMovies.length
    };

    // If no matches, show all demo movies
    if (demoData.results.length === 0) {
        demoData.results = demoMovies;
    }

    this.displaySearchResults(demoData, type);

    // Show demo notice
    const resultsContainer = document.getElementById('search-results');
    const demoNotice = document.createElement('div');
    demoNotice.style.cssText = `
        background-color: #fff3cd;
        border: 1px solid #ffeaa7;
        color: #856404;
        padding: 10px;
        border-radius: 5px;
        margin-bottom: 20px;
        text-align: center;
    `;
    demoNotice.innerHTML = `
        <i class="fas fa-info-circle"></i>
        <strong>Demo Mode:</strong> Showing sample data. Add valid API keys to see real movie data.
    `;

    const resultsGrid = document.getElementById('results-grid');
    resultsContainer.insertBefore(demoNotice, resultsGrid);
};

// Initialize search suggestions when the app loads
document.addEventListener('DOMContentLoaded', () => {
    // Wait for app to be initialized
    setTimeout(() => {
        if (window.app) {
            app.setupSearchSuggestions();
        }
    }, 100);
});
