// Movie/TV show details functionality
MovieDiscoveryApp.prototype.showMovieDetails = async function(itemId, itemType = 'movie') {
    try {
        // Store current media info for advanced features
        this.setCurrentMedia(itemId, itemType);

        let data;
        if (itemType === 'movie') {
            data = await this.apiRequest(`/api/v1/movies/${itemId}`);
        } else if (itemType === 'tv') {
            // Use the proper TV show endpoint
            data = await this.apiRequest(`/api/v1/tv/${itemId}`);
        } else {
            throw new Error(`Unknown item type: ${itemType}`);
        }

        this.displayMovieDetails(data, itemType);
        document.getElementById('movie-modal').style.display = 'block';

        // Load advanced features silently in the background (no error popups)
        this.loadAdvancedFeaturesSilently(itemId, itemType);

    } catch (error) {
        console.error('Failed to load details:', error);
        this.showToast(`Failed to load ${itemType === 'tv' ? 'TV show' : 'movie'} details`, 'error');
    }
};

// Load advanced features without showing error popups
MovieDiscoveryApp.prototype.loadAdvancedFeaturesSilently = async function(itemId, itemType) {
    // Load trailers silently
    try {
        const trailers = await this.silentApiRequest(`/api/v1/${itemType}/${itemId}/trailers`);
        this.displayTrailers(trailers, itemId, itemType);
    } catch (error) {
        console.log('Trailers not available:', error.message);
        // Show "no trailers" message instead of error
        this.displayTrailers([], itemId, itemType);
    }

    // Load watch providers silently
    try {
        const providers = await this.silentApiRequest(`/api/v1/${itemType}/${itemId}/providers`);
        this.displayWatchProviders(providers, 'US');
    } catch (error) {
        console.log('Watch providers not available:', error.message);
        // Show "no providers" message instead of error
        this.displayWatchProviders(null, 'US');
    }
};

MovieDiscoveryApp.prototype.displayMovieDetails = function(movie, itemType) {
    const detailsContainer = document.getElementById('movie-details');
    
    const backdropUrl = movie.backdrop_path
        ? `https://image.tmdb.org/t/p/w1280${movie.backdrop_path}`
        : '/static/images/no-backdrop.jpg';
    
    const posterUrl = movie.poster_path
        ? `https://image.tmdb.org/t/p/w500${movie.poster_path}`
        : '/static/images/placeholder.svg';

    const title = movie.title || movie.name || 'Unknown Title';
    const releaseDate = movie.release_date || movie.first_air_date || '';
    const year = releaseDate ? releaseDate.substring(0, 4) : '';
    const runtime = movie.runtime ? `${movie.runtime} min` : '';
    const rating = movie.vote_average || 0;
    const voteCount = movie.vote_count || 0;

    // TV show specific fields
    const seasons = movie.number_of_seasons ? `${movie.number_of_seasons} season${movie.number_of_seasons > 1 ? 's' : ''}` : '';
    const episodes = movie.number_of_episodes ? `${movie.number_of_episodes} episode${movie.number_of_episodes > 1 ? 's' : ''}` : '';
    const lastAirDate = movie.last_air_date || '';
    const isOngoing = itemType === 'tv' && lastAirDate && new Date(lastAirDate) > new Date();
    
    // Format genres
    const genres = movie.genres ? movie.genres.map(g => g.name).join(', ') : '';
    
    // OMDB data
    const imdbRating = movie.imdb_rating || '';
    const rottenTomatoes = movie.rotten_tomatoes || '';
    const plot = movie.plot || movie.overview || '';
    const director = movie.director || '';
    const writer = movie.writer || '';
    const actors = movie.actors || '';
    const language = movie.language || '';
    const country = movie.country || '';
    const awards = movie.awards || '';

    detailsContainer.innerHTML = `
        <img src="${backdropUrl}" alt="${title}" class="movie-backdrop" 
             onerror="this.style.display='none'">
        
        <div class="movie-details-content">
            <div class="movie-details-header">
                <img src="${posterUrl}" alt="${title}" class="movie-details-poster"
                     onerror="this.src='/static/images/placeholder.svg'">
                
                <div class="movie-details-info">
                    <h1 class="movie-details-title">${title}</h1>
                    
                    <div class="movie-details-meta">
                        ${year ? `<div class="meta-item"><i class="fas fa-calendar"></i> ${year}${isOngoing ? ' - Ongoing' : ''}</div>` : ''}
                        ${itemType === 'movie' && runtime ? `<div class="meta-item"><i class="fas fa-clock"></i> ${runtime}</div>` : ''}
                        ${itemType === 'tv' && seasons ? `<div class="meta-item"><i class="fas fa-tv"></i> ${seasons}</div>` : ''}
                        ${itemType === 'tv' && episodes ? `<div class="meta-item"><i class="fas fa-list"></i> ${episodes}</div>` : ''}
                        ${genres ? `<div class="meta-item"><i class="fas fa-tags"></i> ${genres}</div>` : ''}
                        ${language ? `<div class="meta-item"><i class="fas fa-language"></i> ${language}</div>` : ''}
                    </div>
                    
                    <div class="ratings-section">
                        ${rating > 0 ? `
                            <div class="rating-item">
                                <span class="rating-label">TMDB</span>
                                <div class="rating-value">
                                    <i class="fas fa-star"></i>
                                    ${rating.toFixed(1)}/10
                                    <span class="vote-count">(${voteCount.toLocaleString()} votes)</span>
                                </div>
                            </div>
                        ` : ''}
                        
                        ${imdbRating && imdbRating !== 'N/A' ? `
                            <div class="rating-item">
                                <span class="rating-label">IMDb</span>
                                <div class="rating-value">
                                    <i class="fas fa-star"></i>
                                    ${imdbRating}/10
                                </div>
                            </div>
                        ` : ''}
                        
                        ${rottenTomatoes && rottenTomatoes !== 'N/A' ? `
                            <div class="rating-item">
                                <span class="rating-label">Rotten Tomatoes</span>
                                <div class="rating-value">
                                    <i class="fas fa-tomato"></i>
                                    ${rottenTomatoes}
                                </div>
                            </div>
                        ` : ''}
                    </div>
                    
                    <div class="action-buttons">
                        <button class="action-btn primary" onclick="app.toggleWatchlist('${movie.id}', '${itemType}', '${title}', '${posterUrl}')">
                            <i class="fas fa-bookmark"></i>
                            Add to Watchlist
                        </button>
                        <button class="action-btn secondary" onclick="app.shareMovie('${title}', '${movie.id}')">
                            <i class="fas fa-share"></i>
                            Share
                        </button>
                    </div>
                </div>
            </div>
            
            <!-- Detailed information moved to Overview tab -->
        </div>
    `;
    
    // Update watchlist button state
    this.updateDetailsWatchlistButton(movie.id, itemType);

    // Populate the overview tab with detailed information
    this.populateOverviewTab(movie, itemType);
};

MovieDiscoveryApp.prototype.updateDetailsWatchlistButton = async function(itemId, itemType) {
    const isInWatchlist = await this.isInWatchlist(itemId, itemType);
    const button = document.querySelector('.movie-details .action-btn.primary');

    if (button) {
        if (isInWatchlist) {
            button.innerHTML = '<i class="fas fa-check"></i> In Watchlist';
            button.classList.add('in-watchlist');
        } else {
            button.innerHTML = '<i class="fas fa-bookmark"></i> Add to Watchlist';
            button.classList.remove('in-watchlist');
        }
    }
};

MovieDiscoveryApp.prototype.populateOverviewTab = function(movie, itemType) {
    const overviewTab = document.getElementById('overview-tab');
    if (!overviewTab) return;

    const plot = movie.plot || movie.overview || '';
    const director = movie.director || '';
    const writer = movie.writer || '';
    const actors = movie.actors || '';
    const country = movie.country || '';
    const awards = movie.awards || '';

    overviewTab.innerHTML = `
        ${plot ? `
            <div class="movie-overview">
                <h3>Overview</h3>
                <p>${plot}</p>
            </div>
        ` : ''}

        <div class="movie-details-sections">
            ${director ? `
                <div class="details-section">
                    <h3>Director</h3>
                    <p>${director}</p>
                </div>
            ` : ''}

            ${writer ? `
                <div class="details-section">
                    <h3>Writer</h3>
                    <p>${writer}</p>
                </div>
            ` : ''}

            ${actors ? `
                <div class="details-section">
                    <h3>Cast</h3>
                    <p>${actors}</p>
                </div>
            ` : ''}

            ${country ? `
                <div class="details-section">
                    <h3>Country</h3>
                    <p>${country}</p>
                </div>
            ` : ''}

            ${awards && awards !== 'N/A' ? `
                <div class="details-section">
                    <h3>Awards</h3>
                    <p>${awards}</p>
                </div>
            ` : ''}
        </div>
    `;
};

MovieDiscoveryApp.prototype.shareMovie = function(title, movieId) {
    if (navigator.share) {
        navigator.share({
            title: title,
            text: `Check out this movie: ${title}`,
            url: `${window.location.origin}/movie/${movieId}`
        }).catch(err => console.log('Error sharing:', err));
    } else {
        // Fallback: copy to clipboard
        const url = `${window.location.origin}/movie/${movieId}`;
        navigator.clipboard.writeText(url).then(() => {
            this.showToast('Link copied to clipboard!', 'success');
        }).catch(() => {
            this.showToast('Unable to copy link', 'error');
        });
    }
};

// Add details-specific CSS styles
MovieDiscoveryApp.prototype.addDetailsStyles = function() {
    const style = document.createElement('style');
    style.textContent = `
        .ratings-section {
            margin: 20px 0;
            display: flex;
            flex-wrap: wrap;
            gap: 20px;
        }
        
        .rating-item {
            display: flex;
            flex-direction: column;
            gap: 5px;
        }
        
        .rating-label {
            font-size: 0.8rem;
            color: var(--text-secondary);
            font-weight: 500;
            text-transform: uppercase;
        }
        
        .rating-value {
            display: flex;
            align-items: center;
            gap: 5px;
            font-weight: bold;
            color: var(--text-color);
        }
        
        .rating-value i {
            color: #f39c12;
        }
        
        .vote-count {
            font-size: 0.8rem;
            color: var(--text-secondary);
            font-weight: normal;
        }
        
        .action-buttons {
            margin-top: 20px;
            display: flex;
            gap: 15px;
            flex-wrap: wrap;
        }
        
        .action-btn {
            padding: 12px 24px;
            border: none;
            border-radius: var(--border-radius);
            font-weight: 500;
            cursor: pointer;
            transition: var(--transition);
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .action-btn.primary {
            background-color: var(--secondary-color);
            color: white;
        }
        
        .action-btn.primary:hover {
            background-color: #2980b9;
        }
        
        .action-btn.primary.in-watchlist {
            background-color: var(--accent-color);
        }
        
        .action-btn.secondary {
            background-color: var(--surface-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
        }
        
        .action-btn.secondary:hover {
            background-color: var(--background-color);
        }
        
        .movie-overview {
            margin: 30px 0;
        }
        
        .movie-overview h3 {
            margin-bottom: 15px;
            color: var(--primary-color);
        }
        
        .movie-overview p {
            line-height: 1.8;
            font-size: 1.1rem;
        }
        
        @media (max-width: 768px) {
            .ratings-section {
                flex-direction: column;
                gap: 15px;
            }
            
            .action-buttons {
                flex-direction: column;
            }
            
            .action-btn {
                justify-content: center;
            }
        }
    `;
    
    if (!document.querySelector('#details-styles')) {
        style.id = 'details-styles';
        document.head.appendChild(style);
    }
};

// Close modal when clicking outside
document.addEventListener('click', (e) => {
    const modal = document.getElementById('movie-modal');
    if (e.target === modal) {
        modal.style.display = 'none';
    }
});

// Close modal with Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        const modal = document.getElementById('movie-modal');
        if (modal.style.display === 'block') {
            modal.style.display = 'none';
        }
    }
});

// Initialize details styles when the app loads
document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        if (window.app) {
            app.addDetailsStyles();
        }
    }, 100);
});
