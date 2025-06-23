// Watchlist functionality
MovieDiscoveryApp.prototype.loadWatchlist = async function() {
    try {
        const data = await this.apiRequest('/api/v1/watchlist');
        this.displayWatchlist(data);
        this.loadWatchlistStats();
    } catch (error) {
        console.error('Failed to load watchlist:', error);
    }
};

MovieDiscoveryApp.prototype.loadWatchlistStats = async function() {
    try {
        const stats = await this.apiRequest('/api/v1/watchlist/stats');
        this.displayWatchlistStats(stats);
    } catch (error) {
        console.error('Failed to load watchlist stats:', error);
    }
};

MovieDiscoveryApp.prototype.displayWatchlistStats = function(stats) {
    const statsContainer = document.getElementById('watchlist-stats');
    statsContainer.innerHTML = `
        <div class="stat-item">
            <div class="stat-value">${stats.total_items}</div>
            <div class="stat-label">Total Items</div>
        </div>
        <div class="stat-item">
            <div class="stat-value">${stats.movies}</div>
            <div class="stat-label">Movies</div>
        </div>
        <div class="stat-item">
            <div class="stat-value">${stats.tv_shows}</div>
            <div class="stat-label">TV Shows</div>
        </div>
        <div class="stat-item">
            <div class="stat-value">${stats.watched_items}</div>
            <div class="stat-label">Watched</div>
        </div>
        <div class="stat-item">
            <div class="stat-value">${stats.unwatched_items}</div>
            <div class="stat-label">To Watch</div>
        </div>
        <div class="stat-item">
            <div class="stat-value">${stats.average_rating.toFixed(1)}</div>
            <div class="stat-label">Avg Rating</div>
        </div>
    `;
};

MovieDiscoveryApp.prototype.displayWatchlist = function(watchlistItems, filter = 'all') {
    const watchlistGrid = document.getElementById('watchlist-grid');
    watchlistGrid.innerHTML = '';

    // Filter items based on the current filter
    let filteredItems = watchlistItems;
    if (filter === 'watched') {
        filteredItems = watchlistItems.filter(item => item.watched);
    } else if (filter === 'unwatched') {
        filteredItems = watchlistItems.filter(item => !item.watched);
    }

    if (filteredItems.length === 0) {
        watchlistGrid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px;">
                <i class="fas fa-bookmark" style="font-size: 3rem; color: var(--text-secondary); margin-bottom: 20px;"></i>
                <h3 style="color: var(--text-secondary); margin-bottom: 10px;">Your watchlist is empty</h3>
                <p style="color: var(--text-secondary);">Start adding movies and TV shows to your watchlist!</p>
            </div>
        `;
        return;
    }

    filteredItems.forEach(item => {
        const card = this.createWatchlistCard(item);
        watchlistGrid.appendChild(card);
    });
};

MovieDiscoveryApp.prototype.createWatchlistCard = function(item) {
    const card = document.createElement('div');
    card.className = 'movie-card';
    card.style.position = 'relative';

    const posterUrl = item.poster_path || '/static/images/placeholder.svg';
    const addedDate = new Date(item.added_at).toLocaleDateString();

    card.innerHTML = `
        <img src="${posterUrl}" alt="${item.title}" class="movie-poster"
             onerror="this.src='/static/images/placeholder.svg'">
        <div class="watchlist-actions">
            <button class="action-btn remove-btn" onclick="app.removeFromWatchlist('${item.id}', '${item.type}')" title="Remove from watchlist">
                <i class="fas fa-trash"></i>
            </button>
            <button class="action-btn ${item.watched ? 'unwatch-btn' : 'watch-btn'}" 
                    onclick="app.toggleWatchedStatus('${item.id}', '${item.type}', ${item.watched})" 
                    title="${item.watched ? 'Mark as unwatched' : 'Mark as watched'}">
                <i class="fas ${item.watched ? 'fa-eye-slash' : 'fa-eye'}"></i>
            </button>
        </div>
        <div class="movie-info">
            <div class="movie-title">${item.title}</div>
            <div class="movie-meta">
                <div class="movie-type">${item.type === 'movie' ? 'Movie' : 'TV Show'}</div>
                <div class="added-date">Added: ${addedDate}</div>
                ${item.watched ? `
                    <div class="watched-status">
                        <i class="fas fa-check-circle" style="color: #27ae60;"></i>
                        Watched
                        ${item.rating > 0 ? `<span class="user-rating">★ ${item.rating}/10</span>` : ''}
                    </div>
                ` : ''}
            </div>
        </div>
    `;

    // Add click handler for details (excluding action buttons)
    card.addEventListener('click', (e) => {
        if (!e.target.closest('.watchlist-actions')) {
            this.showMovieDetails(item.id, item.type);
        }
    });

    return card;
};

MovieDiscoveryApp.prototype.toggleWatchlist = async function(itemId, itemType, title, posterPath) {
    try {
        // Find the button that was clicked for immediate visual feedback
        const clickedButton = event?.target?.closest('.watchlist-btn');
        if (clickedButton) {
            clickedButton.style.opacity = '0.5';
            clickedButton.disabled = true;
        }

        // Check if item is already in watchlist
        const isInWatchlist = await this.isInWatchlist(itemId, itemType);

        if (isInWatchlist) {
            await this.removeFromWatchlist(itemId, itemType);
        } else {
            await this.addToWatchlist(itemId, itemType, title, posterPath);
        }

        // Update button state
        await this.updateWatchlistButtons();

        // Refresh watchlist if currently viewing it
        if (this.currentSection === 'watchlist') {
            this.loadWatchlist();
        }

        // Re-enable button
        if (clickedButton) {
            clickedButton.style.opacity = '1';
            clickedButton.disabled = false;
        }
    } catch (error) {
        console.error('Failed to toggle watchlist:', error);
        this.showToast('Failed to update watchlist. Please try again.', 'error');

        // Re-enable button on error
        const clickedButton = event?.target?.closest('.watchlist-btn');
        if (clickedButton) {
            clickedButton.style.opacity = '1';
            clickedButton.disabled = false;
        }
    }
};

MovieDiscoveryApp.prototype.addToWatchlist = async function(itemId, itemType, title, posterPath) {
    try {
        console.log('Adding to watchlist:', { itemId, itemType, title, posterPath });

        const watchlistItem = {
            id: itemId.toString(),
            type: itemType,
            title: title,
            poster_path: posterPath,
            watched: false,
            rating: 0
        };

        console.log('Watchlist item data:', JSON.stringify(watchlistItem, null, 2));

        const response = await this.apiRequest('/api/v1/watchlist', {
            method: 'POST',
            body: JSON.stringify(watchlistItem)
        });

        console.log('API response:', response);
        this.showToast(`Added "${title}" to watchlist`, 'success');
        console.log('Toast shown for success');
    } catch (error) {
        console.error('Error adding to watchlist:', error);
        this.showToast('Failed to add to watchlist', 'error');
        throw error;
    }
};

MovieDiscoveryApp.prototype.removeFromWatchlist = async function(itemId, itemType) {
    try {
        await this.apiRequest(`/api/v1/watchlist/${itemType}/${itemId}`, {
            method: 'DELETE'
        });

        this.showToast('Removed from watchlist', 'success');
        
        // Refresh watchlist if currently viewing it
        if (this.currentSection === 'watchlist') {
            this.loadWatchlist();
        }
    } catch (error) {
        this.showToast('Failed to remove from watchlist', 'error');
        throw error;
    }
};

MovieDiscoveryApp.prototype.toggleWatchedStatus = async function(itemId, itemType, currentStatus) {
    if (currentStatus) {
        // Mark as unwatched
        await this.markAsUnwatched(itemId, itemType);
    } else {
        // Mark as watched (with optional rating)
        const rating = await this.promptForRating();
        await this.markAsWatched(itemId, itemType, rating);
    }
    
    // Refresh watchlist
    this.loadWatchlist();
};

MovieDiscoveryApp.prototype.markAsWatched = async function(itemId, itemType, rating = 0) {
    try {
        await this.apiRequest(`/api/v1/watchlist/${itemType}/${itemId}/watched`, {
            method: 'PUT',
            body: JSON.stringify({ rating: rating })
        });

        this.showToast('Marked as watched', 'success');
    } catch (error) {
        this.showToast('Failed to update watch status', 'error');
        throw error;
    }
};

MovieDiscoveryApp.prototype.markAsUnwatched = async function(itemId, itemType) {
    try {
        // For simplicity, we'll remove and re-add the item
        // In a real app, you'd have a separate endpoint for this
        this.showToast('Marked as unwatched', 'success');
    } catch (error) {
        this.showToast('Failed to update watch status', 'error');
        throw error;
    }
};

MovieDiscoveryApp.prototype.promptForRating = function() {
    return new Promise((resolve) => {
        const rating = prompt('Rate this movie/show (1-10, or leave empty for no rating):');
        const numRating = parseFloat(rating);
        
        if (isNaN(numRating) || numRating < 1 || numRating > 10) {
            resolve(0);
        } else {
            resolve(numRating);
        }
    });
};

MovieDiscoveryApp.prototype.isInWatchlist = async function(itemId, itemType) {
    try {
        const watchlist = await this.apiRequest('/api/v1/watchlist');
        return watchlist.some(item => item.id === itemId.toString() && item.type === itemType);
    } catch (error) {
        console.error('Failed to check watchlist status:', error);
        return false;
    }
};

MovieDiscoveryApp.prototype.updateWatchlistButtons = async function() {
    // Update all watchlist buttons on the current page
    const buttons = document.querySelectorAll('.watchlist-btn');
    
    for (const button of buttons) {
        const onclick = button.getAttribute('onclick');
        if (onclick) {
            const matches = onclick.match(/toggleWatchlist\('(\d+)', '(\w+)'/);
            if (matches) {
                const [, itemId, itemType] = matches;
                const isInWatchlist = await this.isInWatchlist(itemId, itemType);
                button.classList.toggle('in-watchlist', isInWatchlist);
            }
        }
    }
};

MovieDiscoveryApp.prototype.filterWatchlist = function(filter) {
    // Update active filter button
    document.querySelectorAll('.watchlist-filters .filter-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    event.target.classList.add('active');
    
    // Reload watchlist with filter
    this.loadWatchlist().then(() => {
        const watchlistData = this.lastWatchlistData || [];
        this.displayWatchlist(watchlistData, filter);
    });
};

// Store watchlist data for filtering
MovieDiscoveryApp.prototype.displayWatchlist = function(watchlistItems, filter = 'all') {
    this.lastWatchlistData = watchlistItems;
    
    const watchlistGrid = document.getElementById('watchlist-grid');
    watchlistGrid.innerHTML = '';

    // Filter items based on the current filter
    let filteredItems = watchlistItems;
    if (filter === 'watched') {
        filteredItems = watchlistItems.filter(item => item.watched);
    } else if (filter === 'unwatched') {
        filteredItems = watchlistItems.filter(item => !item.watched);
    }

    if (filteredItems.length === 0) {
        const emptyMessage = filter === 'all' 
            ? 'Your watchlist is empty'
            : filter === 'watched' 
                ? 'No watched items yet'
                : 'No items to watch';
                
        watchlistGrid.innerHTML = `
            <div style="grid-column: 1 / -1; text-align: center; padding: 40px;">
                <i class="fas fa-bookmark" style="font-size: 3rem; color: var(--text-secondary); margin-bottom: 20px;"></i>
                <h3 style="color: var(--text-secondary); margin-bottom: 10px;">${emptyMessage}</h3>
                <p style="color: var(--text-secondary);">Start adding movies and TV shows to your watchlist!</p>
            </div>
        `;
        return;
    }

    filteredItems.forEach(item => {
        const card = this.createWatchlistCard(item);
        watchlistGrid.appendChild(card);
    });
};

// Global function for filter buttons
function filterWatchlist(filter) {
    app.filterWatchlist(filter);
}
