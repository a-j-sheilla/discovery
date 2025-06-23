# Movie Discovery App API Documentation

## Overview

The Movie Discovery App provides a RESTful API for searching movies and TV shows, managing watchlists, and getting personalized recommendations. The API combines data from TMDB (The Movie Database) and OMDB (Open Movie Database) to provide comprehensive entertainment information.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. In a production environment, you would implement proper authentication and authorization.

## Rate Limiting

The API implements rate limiting to prevent abuse:
- Default: 60 requests per minute per IP
- Rate limit headers are included in responses
- Exceeding the limit returns HTTP 429 (Too Many Requests)

## Response Format

All API responses are in JSON format. Successful responses include the requested data, while error responses follow this structure:

```json
{
  "error": "Error message",
  "status_code": 400
}
```

## Endpoints

### Health Check

#### GET /health

Check the health status of the API service.

**Response:**
```json
{
  "status": "healthy",
  "service": "movie-discovery-app"
}
```

### Search

#### GET /search/movies

Search for movies by title.

**Parameters:**
- `q` (required): Search query string
- `page` (optional): Page number (default: 1, max: 1000)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/search/movies?q=inception&page=1"
```

**Response:**
```json
{
  "page": 1,
  "results": [
    {
      "id": 27205,
      "title": "Inception",
      "overview": "Cobb, a skilled thief who commits corporate espionage...",
      "release_date": "2010-07-15",
      "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
      "backdrop_path": "/aej3LRUga5rhgkmRP6XMFw3ejbl.jpg",
      "vote_average": 8.4,
      "vote_count": 31546,
      "popularity": 151.489,
      "genre_ids": [28, 878, 53],
      "imdb_rating": "8.8",
      "rotten_tomatoes": "87%",
      "director": "Christopher Nolan",
      "actors": "Leonardo DiCaprio, Marion Cotillard, Tom Hardy"
    }
  ],
  "total_pages": 1,
  "total_results": 1
}
```

#### GET /search/tv

Search for TV shows by title.

**Parameters:**
- `q` (required): Search query string
- `page` (optional): Page number (default: 1, max: 1000)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/search/tv?q=breaking%20bad&page=1"
```

**Response:** Similar to movies search but with TV show specific fields.

### Movie Details

#### GET /movies/{id}

Get detailed information about a specific movie.

**Parameters:**
- `id` (required): Movie ID from TMDB

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/movies/27205"
```

**Response:**
```json
{
  "id": 27205,
  "title": "Inception",
  "overview": "Cobb, a skilled thief who commits corporate espionage...",
  "release_date": "2010-07-15",
  "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
  "backdrop_path": "/aej3LRUga5rhgkmRP6XMFw3ejbl.jpg",
  "vote_average": 8.4,
  "vote_count": 31546,
  "popularity": 151.489,
  "runtime": 148,
  "genres": [
    {"id": 28, "name": "Action"},
    {"id": 878, "name": "Science Fiction"},
    {"id": 53, "name": "Thriller"}
  ],
  "imdb_rating": "8.8",
  "rotten_tomatoes": "87%",
  "plot": "Dom Cobb is a skilled thief, the absolute best in the dangerous art of extraction...",
  "director": "Christopher Nolan",
  "writer": "Christopher Nolan",
  "actors": "Leonardo DiCaprio, Marion Cotillard, Tom Hardy, Elliot Page",
  "language": "English, Japanese, French",
  "country": "United States, United Kingdom",
  "awards": "Won 4 Oscars. Another 143 wins & 198 nominations.",
  "imdb_id": "tt1375666"
}
```

### Trending Content

#### GET /trending/movies

Get trending movies.

**Parameters:**
- `time_window` (optional): `day` or `week` (default: `week`)
- `page` (optional): Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/trending/movies?time_window=week&page=1"
```

**Response:** Similar to search results but with trending movies.

### Genres

#### GET /genres/movies

Get all available movie genres.

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/genres/movies"
```

**Response:**
```json
[
  {"id": 28, "name": "Action"},
  {"id": 12, "name": "Adventure"},
  {"id": 16, "name": "Animation"},
  {"id": 35, "name": "Comedy"},
  {"id": 80, "name": "Crime"}
]
```

#### GET /genres/tv

Get all available TV show genres.

#### GET /discover/genre/{genreId}

Discover movies by genre with optional filters.

**Parameters:**
- `genreId` (required): Genre ID
- `page` (optional): Page number (default: 1)
- `sort_by` (optional): Sort order (default: `popularity.desc`)
  - Options: `popularity.desc`, `popularity.asc`, `vote_average.desc`, `vote_average.asc`, `release_date.desc`, `release_date.asc`
- `min_rating` (optional): Minimum vote average (0-10)
- `max_rating` (optional): Maximum vote average (0-10)
- `min_year` (optional): Minimum release year
- `max_year` (optional): Maximum release year

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/discover/genre/28?page=1&sort_by=vote_average.desc&min_rating=7.0"
```

### Watchlist Management

#### GET /watchlist

Get the user's watchlist.

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/watchlist"
```

**Response:**
```json
[
  {
    "id": "27205",
    "type": "movie",
    "title": "Inception",
    "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
    "added_at": "2023-12-01T10:30:00Z",
    "watched": true,
    "rating": 9.0
  }
]
```

#### POST /watchlist

Add an item to the watchlist.

**Request Body:**
```json
{
  "id": "27205",
  "type": "movie",
  "title": "Inception",
  "poster_path": "/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
  "watched": false,
  "rating": 0
}
```

**Response:**
```json
{
  "status": "success"
}
```

#### DELETE /watchlist/{type}/{id}

Remove an item from the watchlist.

**Parameters:**
- `type`: `movie` or `tv`
- `id`: Item ID

**Example Request:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/watchlist/movie/27205"
```

#### PUT /watchlist/{type}/{id}/watched

Mark an item as watched with optional rating.

**Parameters:**
- `type`: `movie` or `tv`
- `id`: Item ID

**Request Body:**
```json
{
  "rating": 8.5
}
```

**Example Request:**
```bash
curl -X PUT "http://localhost:8080/api/v1/watchlist/movie/27205/watched" \
  -H "Content-Type: application/json" \
  -d '{"rating": 8.5}'
```

#### GET /watchlist/stats

Get watchlist statistics.

**Response:**
```json
{
  "total_items": 25,
  "watched_items": 15,
  "unwatched_items": 10,
  "movies": 18,
  "tv_shows": 7,
  "average_rating": 7.8
}
```

### Recommendations

#### GET /recommendations

Get personalized recommendations based on the user's watchlist.

**Parameters:**
- `limit` (optional): Number of recommendations (default: 20, max: 50)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/recommendations?limit=10"
```

**Response:**
```json
[
  {
    "item": {
      "id": 550,
      "title": "Fight Club",
      "overview": "A ticking-time-bomb insomniac and a slippery soap salesman...",
      "vote_average": 8.4,
      "popularity": 61.416
    },
    "score": 8.7,
    "type": "movie"
  }
]
```

## Error Handling

The API uses standard HTTP status codes:

- `200 OK`: Successful request
- `400 Bad Request`: Invalid request parameters
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

Error responses include a descriptive message:

```json
{
  "error": "Query parameter 'q' is required",
  "status_code": 400
}
```

## Data Sources

The API combines data from multiple sources:

1. **TMDB (The Movie Database)**: Primary source for movie/TV data, images, and basic information
2. **OMDB (Open Movie Database)**: Additional ratings (IMDb, Rotten Tomatoes) and detailed plot information

## Caching

The API implements intelligent caching to improve performance:

- Search results: 30 minutes
- Movie details: 30 minutes
- Genre lists: 24 hours
- Trending content: 30 minutes

Cache headers are included in responses to indicate cache status.

## Best Practices

1. **Pagination**: Always use pagination for large result sets
2. **Rate Limiting**: Respect rate limits and implement exponential backoff
3. **Error Handling**: Always check response status codes and handle errors gracefully
4. **Caching**: Implement client-side caching to reduce API calls
5. **Search Debouncing**: Implement debouncing for search inputs to avoid excessive API calls

## SDK and Libraries

While no official SDK is provided, the API is designed to be easily consumed by any HTTP client. Popular libraries for different languages:

- **JavaScript**: `fetch()`, `axios`
- **Python**: `requests`, `httpx`
- **Go**: `net/http`
- **Java**: `OkHttp`, `HttpClient`
- **PHP**: `Guzzle`, `cURL`

## Support

For API support and questions:
1. Check this documentation
2. Review the [GitHub Issues](https://github.com/yourusername/movie-discovery-app/issues)
3. Create a new issue with detailed information about your problem
