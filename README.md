# Movie Discovery Web App

A comprehensive entertainment discovery platform built with Go, where users can search for movies and TV shows, view detailed information, manage personal watchlists, and discover trending content.

![Movie Discovery App](https://img.shields.io/badge/Go-1.24-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

## 🎬 Features

### Core Features
- **Real-time Search**: Search for movies and TV shows with debounced input and auto-suggestions
- **Detailed Information**: View comprehensive details including ratings from multiple sources (TMDB, IMDb, Rotten Tomatoes)
- **Personal Watchlist**: Add/remove titles, mark as watched, and rate content
- **Trending Content**: Discover popular movies and shows (daily/weekly trends)
- **Genre Filtering**: Browse content by genre with advanced filtering options
- **Responsive Design**: Optimized for both desktop and mobile devices

### Advanced Features
- **Recommendation Engine**: Personalized recommendations based on watchlist preferences
- **Multi-source Data**: Combines data from TMDB and OMDB APIs for comprehensive information
- **Caching System**: Intelligent caching for improved performance
- **Rate Limiting**: Graceful API rate limiting to prevent service disruption
- **Dark/Light Theme**: Toggle between themes with persistent preference storage
- **Export/Import**: Export watchlist data as JSON

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- TMDB API key ([Get one here](https://www.themoviedb.org/settings/api))
- OMDB API key ([Get one here](http://www.omdbapi.com/apikey.aspx))

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/movie-discovery-app.git
   cd movie-discovery-app
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and add your API keys:
   ```env
   TMDB_API_KEY=your_tmdb_api_key_here
   OMDB_API_KEY=your_omdb_api_key_here
   PORT=8080
   HOST=localhost
   ```

   **API Integration Status:**
   - ✅ **TMDB API**: Fully integrated for movie/TV search, trending content, and detailed information
   - ✅ **OMDB API**: Integrated for additional movie metadata (ratings, awards, plot details)
   - ✅ **Real Data**: Application now returns real movie data with actual posters and information

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

5. **Open your browser**
   Navigate to `http://localhost:8080`

## 📁 Project Structure

```
movie-discovery-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── configs/
│   └── config.go                # Configuration management
├── internal/
│   ├── api/
│   │   ├── handlers.go          # HTTP request handlers
│   │   ├── handlers_test.go     # Handler tests
│   │   └── router.go            # Route definitions
│   ├── models/
│   │   └── movie.go             # Data models
│   └── services/
│       ├── discovery.go         # Main discovery service
│       ├── discovery_test.go    # Discovery service tests
│       ├── tmdb.go              # TMDB API client
│       ├── omdb.go              # OMDB API client
│       ├── watchlist.go         # Watchlist management
│       ├── watchlist_test.go    # Watchlist tests
│       ├── recommendations.go   # Recommendation engine
│       └── genres.go            # Genre filtering
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── styles.css       # Application styles
│   │   ├── js/
│   │   │   ├── app.js           # Main application logic
│   │   │   ├── search.js        # Search functionality
│   │   │   ├── watchlist.js     # Watchlist management
│   │   │   ├── trending.js      # Trending content
│   │   │   └── details.js       # Movie details modal
│   │   └── images/              # Static images
│   └── templates/
│       └── index.html           # Main HTML template
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── README.md                    # This file
```

## 🔧 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Search
- `GET /search/movies?q={query}&page={page}` - Search movies
- `GET /search/tv?q={query}&page={page}` - Search TV shows

#### Content Details
- `GET /movies/{id}` - Get movie details
- `GET /trending/movies?time_window={day|week}&page={page}` - Get trending movies

#### Genres
- `GET /genres/movies` - Get movie genres
- `GET /genres/tv` - Get TV show genres
- `GET /discover/genre/{genreId}?page={page}&sort_by={sort}&min_rating={rating}` - Discover by genre

#### Watchlist
- `GET /watchlist` - Get user's watchlist
- `POST /watchlist` - Add item to watchlist
- `DELETE /watchlist/{type}/{id}` - Remove item from watchlist
- `PUT /watchlist/{type}/{id}/watched` - Mark item as watched
- `GET /watchlist/stats` - Get watchlist statistics

#### Recommendations
- `GET /recommendations?limit={limit}` - Get personalized recommendations

#### Health Check
- `GET /health` - Service health status

### Example API Usage

**Search for movies:**
```bash
curl "http://localhost:8080/api/v1/search/movies?q=inception&page=1"
```

**Add to watchlist:**
```bash
curl -X POST "http://localhost:8080/api/v1/watchlist" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "27205",
    "type": "movie",
    "title": "Inception",
    "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg"
  }'
```

## 🧪 Testing

### Run all tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test -cover ./...
```

### Run specific test package
```bash
go test ./internal/services
```

### Run benchmarks
```bash
go test -bench=. ./...
```

## 🎨 Frontend Features

### User Interface
- **Responsive Grid Layout**: Adaptive movie/TV show cards
- **Search with Suggestions**: Real-time search suggestions with debouncing
- **Modal Details View**: Comprehensive movie/show information
- **Theme Toggle**: Dark/light mode with system preference detection
- **Loading States**: Smooth loading indicators and error handling

### JavaScript Architecture
- **Modular Design**: Separate files for different functionalities
- **Event-Driven**: Efficient event handling and DOM manipulation
- **Local Storage**: Persistent theme and preference storage
- **API Integration**: Clean separation between frontend and backend

## 🔒 Security Features

- **API Key Management**: Secure environment variable handling
- **Rate Limiting**: Protection against API abuse
- **Input Validation**: Comprehensive input sanitization
- **CORS Support**: Configurable cross-origin resource sharing
- **Error Handling**: Graceful error responses without sensitive data exposure

## 🚀 Performance Optimizations

- **Caching Strategy**: Multi-level caching for API responses
- **Debounced Search**: Reduced API calls with intelligent debouncing
- **Lazy Loading**: Efficient image loading and pagination
- **Compression**: Optimized asset delivery
- **Connection Pooling**: Efficient HTTP client management

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `HOST` | Server host | `localhost` | No |
| `TMDB_API_KEY` | TMDB API key | - | Yes |
| `OMDB_API_KEY` | OMDB API key | - | Yes |
| `CACHE_DURATION_MINUTES` | Cache duration | `30` | No |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limit | `60` | No |

## 📝 Development

### Adding New Features

1. **Backend**: Add service logic in `internal/services/`
2. **API**: Add handlers in `internal/api/handlers.go`
3. **Routes**: Update `internal/api/router.go`
4. **Frontend**: Add JavaScript in `web/static/js/`
5. **Tests**: Add tests for new functionality

### Code Style
- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comprehensive comments for public functions
- Write tests for all new functionality

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [The Movie Database (TMDB)](https://www.themoviedb.org/) for movie/TV data
- [Open Movie Database (OMDB)](http://www.omdbapi.com/) for additional ratings
- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- [Font Awesome](https://fontawesome.com/) for icons


**Happy movie discovering! 🎬🍿**
