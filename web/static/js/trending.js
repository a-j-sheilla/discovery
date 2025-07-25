// Trending content functionality
MovieDiscoveryApp.prototype.loadTrending = async function(timeWindow = 'week') {
    try {
        const data = await this.apiRequest(`/api/v1/trending/movies?time_window=${timeWindow}&page=1`);
        this.displayTrending(data);
        this.updateTrendingFilter(timeWindow);
    } catch (error) {
        console.error('Failed to load trending content:', error);
        this.showDemoTrending(timeWindow);
    }
};

MovieDiscoveryApp.prototype.displayTrending = function(data) {
    const trendingGrid = document.getElementById('trending-grid');
    trendingGrid.innerHTML = '';

    if (data.results && data.results.length > 0) {
        data.results.forEach(item => {
            const card = this.createTrendingCard(item);
            trendingGrid.appendChild(card);
        });
    } else {
        trendingGrid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px;">
                <i class="fas fa-exclamation-triangle" style="font-size: 3rem; color: var(--text-secondary); margin-bottom: 20px;"></i>
                <h3 style="color: var(--text-secondary); margin-bottom: 10px;">Trending Content Unavailable</h3>
                <p style="color: var(--text-secondary);">This feature requires valid TMDB API keys</p>
                <p style="color: var(--text-secondary); font-size: 0.9rem; margin-top: 10px;">
                    To enable trending content, add your TMDB API key to the .env file
                </p>
            </div>
        `;
    }
};

MovieDiscoveryApp.prototype.createTrendingCard = function(item) {
    const card = document.createElement('div');
    card.className = 'movie-card trending-card';
    card.style.position = 'relative';

    const posterUrl = item.poster_path
        ? `https://image.tmdb.org/t/p/w300${item.poster_path}`
        : '/static/images/placeholder.svg';

    const title = item.title || item.name || 'Unknown Title';
    const year = item.release_date || item.first_air_date || '';
    const rating = item.vote_average || 0;
    const popularity = item.popularity || 0;

    // Determine if it's a movie or TV show
    const type = item.title ? 'movie' : 'tv';

    card.innerHTML = `
        <div class="trending-badge">
            <i class="fas fa-fire"></i>
            <span>${Math.round(popularity)}</span>
        </div>
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
            <div class="trending-info">
                <span class="trending-type">${type === 'movie' ? 'Movie' : 'TV Show'}</span>
                <span class="popularity-score">Popularity: ${Math.round(popularity)}</span>
            </div>
        </div>
    `;

    // Add click handler for details
    card.addEventListener('click', (e) => {
        if (!e.target.closest('.watchlist-btn')) {
            this.showMovieDetails(item.id, type);
        }
    });

    return card;
};

MovieDiscoveryApp.prototype.updateTrendingFilter = function(activeFilter) {
    document.querySelectorAll('.trending-filters .filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Find and activate the correct button
    const buttons = document.querySelectorAll('.trending-filters .filter-btn');
    buttons.forEach(btn => {
        const onclick = btn.getAttribute('onclick');
        if (onclick && onclick.includes(`'${activeFilter}'`)) {
            btn.classList.add('active');
        }
    });
};

// Enhanced trending functionality with categories
MovieDiscoveryApp.prototype.loadTrendingByCategory = async function(category, timeWindow = 'week') {
    try {
        let endpoint;
        switch (category) {
            case 'movies':
                endpoint = `/api/v1/trending/movies?time_window=${timeWindow}&page=1`;
                break;
            case 'tv':
                endpoint = `/api/v1/trending/tv?time_window=${timeWindow}&page=1`;
                break;
            default:
                endpoint = `/api/v1/trending/movies?time_window=${timeWindow}&page=1`;
        }

        const data = await this.apiRequest(endpoint);
        this.displayTrending(data);
        this.updateTrendingFilter(timeWindow);
    } catch (error) {
        console.error('Failed to load trending content:', error);
        this.showToast('Failed to load trending content', 'error');
    }
};

// Add trending-specific CSS styles
MovieDiscoveryApp.prototype.addTrendingStyles = function() {
    const style = document.createElement('style');
    style.textContent = `
        .trending-card {
            position: relative;
        }
        
        .trending-badge {
            position: absolute;
            top: 10px;
            left: 10px;
            background: linear-gradient(45deg, #ff6b6b, #ffa500);
            color: white;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 0.8rem;
            font-weight: bold;
            display: flex;
            align-items: center;
            gap: 5px;
            z-index: 5;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
        }
        
        .trending-badge i {
            animation: flicker 2s infinite alternate;
        }
        
        @keyframes flicker {
            0% { opacity: 1; }
            100% { opacity: 0.7; }
        }
        
        .trending-info {
            margin-top: 8px;
            font-size: 0.75rem;
            color: var(--text-secondary);
        }
        
        .trending-type {
            background-color: var(--secondary-color);
            color: white;
            padding: 2px 6px;
            border-radius: 10px;
            font-size: 0.7rem;
            margin-right: 5px;
        }
        
        .popularity-score {
            font-weight: 500;
        }
        
        .trending-filters .filter-btn {
            position: relative;
            overflow: hidden;
        }
        
        .trending-filters .filter-btn.active::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: linear-gradient(45deg, var(--secondary-color), var(--accent-color));
            z-index: -1;
        }
        
        .trending-filters .filter-btn.active {
            color: white;
            border-color: transparent;
        }
    `;
    
    if (!document.querySelector('#trending-styles')) {
        style.id = 'trending-styles';
        document.head.appendChild(style);
    }
};

// Demo trending data
MovieDiscoveryApp.prototype.showDemoTrending = function(timeWindow) {
    const demoTrendingMovies = [
        {
            id: 550,
            title: "Fight Club",
            overview: "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression.",
            release_date: "1999-10-15",
            poster_path: null,
            vote_average: 8.4,
            popularity: 95.5
        },
        {
            id: 13,
            title: "Forrest Gump",
            overview: "A man with a low IQ has accomplished great things in his life.",
            release_date: "1994-07-06",
            poster_path: null,
            vote_average: 8.5,
            popularity: 88.2
        },
        {
            id: 27205,
            title: "Inception",
            overview: "Cobb, a skilled thief who commits corporate espionage by infiltrating the subconscious.",
            release_date: "2010-07-16",
            poster_path: null,
            vote_average: 8.4,
            popularity: 92.1
        },
        {
            id: 278,
            title: "The Shawshank Redemption",
            overview: "Two imprisoned men bond over a number of years, finding solace and eventual redemption.",
            release_date: "1994-09-23",
            poster_path: null,
            vote_average: 9.3,
            popularity: 87.9
        }
    ];

    const demoData = {
        page: 1,
        results: demoTrendingMovies,
        total_pages: 1,
        total_results: demoTrendingMovies.length
    };

    this.displayTrending(demoData);
    this.updateTrendingFilter(timeWindow);

    // Add demo notice
    const trendingGrid = document.getElementById('trending-grid');
    const demoNotice = document.createElement('div');
    demoNotice.style.cssText = `
        grid-column: 1 / -1;
        background-color: #fff3cd;
        border: 1px solid #ffeaa7;
        color: #856404;
        padding: 15px;
        border-radius: 8px;
        margin-bottom: 20px;
        text-align: center;
    `;
    demoNotice.innerHTML = `
        <i class="fas fa-info-circle"></i>
        <strong>Demo Mode:</strong> Showing sample trending data. Add valid TMDB API keys to see real trending movies.
    `;

    trendingGrid.insertBefore(demoNotice, trendingGrid.firstChild);
};

// Global functions for trending filters
function loadTrending(timeWindow) {
    app.loadTrending(timeWindow);
}

// Initialize trending styles when the app loads
document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        if (window.app) {
            app.addTrendingStyles();
        }
    }, 100);
});
