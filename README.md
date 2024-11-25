# Music Library API

A RESTful API service for managing a music library, built with Go.

## Features

- CRUD operations for songs
- Filtering and pagination for song listing
- Integration with external music info API
- Automatic database migrations
- Swagger documentation
- Structured logging

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Docker (optional)

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=music_library
DB_SSL_MODE=disable

MUSIC_API_URL=http://localhost:8081
SERVER_PORT=8080
```

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Start PostgreSQL database
4. Run the application:
   ```bash
   go run cmd/main.go
   ```

## API Documentation

Once the server is running, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## API Endpoints

- `POST /api/v1/songs` - Create a new song
- `GET /api/v1/songs` - List songs (with filtering and pagination)
- `GET /api/v1/songs/:id` - Get a specific song
- `PUT /api/v1/songs/:id` - Update a song
- `DELETE /api/v1/songs/:id` - Delete a song

## Example Usage

Create a new song:
```bash
curl -X POST http://localhost:8080/api/v1/songs \
  -H "Content-Type: application/json" \
  -d '{
    "group": "Muse",
    "song": "Supermassive Black Hole"
  }'
```

List songs with filtering:
```bash
curl "http://localhost:8080/api/v1/songs?group=Muse&page=1&page_size=10"
```
