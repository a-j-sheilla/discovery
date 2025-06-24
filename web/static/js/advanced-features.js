// Advanced Features JavaScript - Trailers and Watch Providers

// Trailer functionality
MovieDiscoveryApp.prototype.loadTrailers = async function(mediaId, mediaType) {
    try {
        const trailers = await this.apiRequest(`/api/v1/${mediaType}/${mediaId}/trailers`);
        this.displayTrailers(trailers, mediaId, mediaType);
    } catch (error) {
        console.log('Failed to load trailers:', error);
        // Display no trailers message instead of error toast
        this.displayTrailers([], mediaId, mediaType);
    }
};

MovieDiscoveryApp.prototype.loadOfficialTrailer = async function(mediaId, mediaType) {
    try {
        const trailer = await this.apiRequest(`/api/v1/${mediaType}/${mediaId}/trailer`);
        return trailer;
    } catch (error) {
        console.error('Failed to load official trailer:', error);
        return null;
    }
};

MovieDiscoveryApp.prototype.displayTrailers = function(trailers, mediaId, mediaType) {
    const trailersContainer = document.getElementById('trailers-container');
    if (!trailersContainer) return;

    trailersContainer.innerHTML = '';

    if (!trailers || trailers.length === 0) {
        trailersContainer.innerHTML = `
            <div class="no-trailers">
                <i class="fas fa-video-slash"></i>
                <p>No trailers available</p>
            </div>
        `;
        return;
    }

    const trailersHTML = trailers.map(trailer => `
        <div class="trailer-card" onclick="playTrailer('${trailer.video_id}', '${trailer.title}')">
            <div class="trailer-thumbnail">
                <img src="${trailer.thumbnail}" alt="${trailer.title}" loading="lazy">
                <div class="play-overlay">
                    <i class="fas fa-play"></i>
                </div>
            </div>
            <div class="trailer-info">
                <h4 class="trailer-title">${trailer.title}</h4>
                <p class="trailer-channel">${trailer.channel_title}</p>
                <p class="trailer-date">${this.formatDate(trailer.published_at)}</p>
            </div>
        </div>
    `).join('');

    trailersContainer.innerHTML = `
        <div class="trailers-grid">
            ${trailersHTML}
        </div>
    `;
};

MovieDiscoveryApp.prototype.playTrailer = function(videoId, title) {
    const modal = document.getElementById('trailer-modal');
    const iframe = document.getElementById('trailer-iframe');
    const modalTitle = document.getElementById('trailer-modal-title');

    if (modal && iframe && modalTitle) {
        modalTitle.textContent = title;
        iframe.src = `https://www.youtube.com/embed/${videoId}?autoplay=1&rel=0`;
        modal.style.display = 'block';
        document.body.style.overflow = 'hidden';
    }
};

MovieDiscoveryApp.prototype.closeTrailerModal = function() {
    const modal = document.getElementById('trailer-modal');
    const iframe = document.getElementById('trailer-iframe');

    if (modal && iframe) {
        modal.style.display = 'none';
        iframe.src = '';
        document.body.style.overflow = 'auto';
    }
};

// Watch Providers functionality
MovieDiscoveryApp.prototype.loadWatchProviders = async function(mediaId, mediaType, region = 'US') {
    try {
        const providers = await this.apiRequest(`/api/v1/${mediaType}/${mediaId}/providers`);
        this.displayWatchProviders(providers, region);
    } catch (error) {
        console.log('Failed to load watch providers:', error);
        // Display no providers message instead of error toast
        this.displayWatchProviders(null, region);
    }
};

MovieDiscoveryApp.prototype.loadStreamingServices = async function(mediaId, mediaType, region = 'US') {
    try {
        const services = await this.apiRequest(`/api/v1/${mediaType}/${mediaId}/streaming?region=${region}`);
        return services;
    } catch (error) {
        console.error('Failed to load streaming services:', error);
        return [];
    }
};

MovieDiscoveryApp.prototype.displayWatchProviders = function(providers, region = 'US') {
    const providersContainer = document.getElementById('providers-container');
    if (!providersContainer) return;

    providersContainer.innerHTML = '';

    if (!providers || !providers.results || !providers.results[region]) {
        providersContainer.innerHTML = `
            <div class="no-providers">
                <i class="fas fa-tv"></i>
                <p>No streaming information available for ${region}</p>
            </div>
        `;
        return;
    }

    const regionData = providers.results[region];
    let providersHTML = '';

    // Streaming services (subscription)
    if (regionData.flatrate && regionData.flatrate.length > 0) {
        providersHTML += `
            <div class="provider-section">
                <h4><i class="fas fa-play-circle"></i> Stream</h4>
                <div class="providers-list">
                    ${regionData.flatrate.map(provider => `
                        <div class="provider-item" title="${provider.provider_name}">
                            <img src="https://image.tmdb.org/t/p/original${provider.logo_path}" 
                                 alt="${provider.provider_name}" loading="lazy">
                            <span>${provider.provider_name}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    // Purchase options
    if (regionData.buy && regionData.buy.length > 0) {
        providersHTML += `
            <div class="provider-section">
                <h4><i class="fas fa-shopping-cart"></i> Buy</h4>
                <div class="providers-list">
                    ${regionData.buy.map(provider => `
                        <div class="provider-item" title="${provider.provider_name}">
                            <img src="https://image.tmdb.org/t/p/original${provider.logo_path}" 
                                 alt="${provider.provider_name}" loading="lazy">
                            <span>${provider.provider_name}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    // Rental options
    if (regionData.rent && regionData.rent.length > 0) {
        providersHTML += `
            <div class="provider-section">
                <h4><i class="fas fa-clock"></i> Rent</h4>
                <div class="providers-list">
                    ${regionData.rent.map(provider => `
                        <div class="provider-item" title="${provider.provider_name}">
                            <img src="https://image.tmdb.org/t/p/original${provider.logo_path}" 
                                 alt="${provider.provider_name}" loading="lazy">
                            <span>${provider.provider_name}</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    if (providersHTML) {
        providersContainer.innerHTML = providersHTML;
        
        // Add link to JustWatch if available
        if (regionData.link) {
            providersContainer.innerHTML += `
                <div class="provider-link">
                    <a href="${regionData.link}" target="_blank" rel="noopener noreferrer">
                        <i class="fas fa-external-link-alt"></i> View on JustWatch
                    </a>
                </div>
            `;
        }
    } else {
        providersContainer.innerHTML = `
            <div class="no-providers">
                <i class="fas fa-tv"></i>
                <p>No streaming information available</p>
            </div>
        `;
    }
};

// Enhanced movie details with advanced features
MovieDiscoveryApp.prototype.showMovieDetailsWithAdvanced = async function(movieId, mediaType = 'movie') {
    // Show existing movie details first
    await this.showMovieDetails(movieId, mediaType);
    
    // Load advanced features
    await Promise.all([
        this.loadTrailers(movieId, mediaType),
        this.loadWatchProviders(movieId, mediaType)
    ]);
};

// Utility functions
MovieDiscoveryApp.prototype.formatDate = function(dateString) {
    if (!dateString) return '';
    
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    } catch (error) {
        return dateString;
    }
};

// Global functions for HTML onclick handlers
function playTrailer(videoId, title) {
    app.playTrailer(videoId, title);
}

function closeTrailerModal() {
    app.closeTrailerModal();
}

function loadWatchProviders(mediaId, mediaType, region) {
    app.loadWatchProviders(mediaId, mediaType, region);
}

// Tab switching functionality
function showAdvancedTab(tabName) {
    // Hide all tab contents
    const tabContents = document.querySelectorAll('.tab-content');
    tabContents.forEach(content => content.classList.remove('active'));

    // Remove active class from all tab buttons
    const tabButtons = document.querySelectorAll('.tab-btn');
    tabButtons.forEach(btn => btn.classList.remove('active'));

    // Show selected tab content
    const selectedTab = document.getElementById(`${tabName}-tab`);
    if (selectedTab) {
        selectedTab.classList.add('active');
    }

    // Add active class to selected tab button
    const selectedButton = event.target;
    if (selectedButton) {
        selectedButton.classList.add('active');
    }

    // Load content for the selected tab only when user clicks on it
    const currentMediaId = app.currentMediaId;
    const currentMediaType = app.currentMediaType;

    if (currentMediaId && currentMediaType) {
        if (tabName === 'trailers') {
            // Only load if not already loaded
            const trailersContainer = document.getElementById('trailers-container');
            if (trailersContainer && trailersContainer.innerHTML.includes('loading-trailers')) {
                app.loadTrailers(currentMediaId, currentMediaType);
            }
        } else if (tabName === 'providers') {
            // Only load if not already loaded
            const providersContainer = document.getElementById('providers-container');
            if (providersContainer && providersContainer.innerHTML.includes('loading-providers')) {
                const region = document.getElementById('region-select')?.value || 'US';
                app.loadWatchProviders(currentMediaId, currentMediaType, region);
            }
        }
    }
}

// Region change functionality
function changeRegion() {
    const regionSelect = document.getElementById('region-select');
    const region = regionSelect?.value || 'US';

    const currentMediaId = app.currentMediaId;
    const currentMediaType = app.currentMediaType;

    if (currentMediaId && currentMediaType) {
        app.loadWatchProviders(currentMediaId, currentMediaType, region);
    }
}

// Store current media info for advanced features
MovieDiscoveryApp.prototype.setCurrentMedia = function(mediaId, mediaType) {
    this.currentMediaId = mediaId;
    this.currentMediaType = mediaType;
};

// Initialize advanced features when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Close trailer modal when clicking outside
    window.addEventListener('click', function(event) {
        const modal = document.getElementById('trailer-modal');
        if (event.target === modal) {
            app.closeTrailerModal();
        }
    });

    // Close trailer modal with Escape key
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape') {
            app.closeTrailerModal();
        }
    });
});
