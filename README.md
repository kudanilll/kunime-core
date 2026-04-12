# Kunime Core

Backend API for Kunime mobile app data (watch events, popular anime, recommendations, banners, genres, and app versioning).

## Run Locally

1. Copy `.env.example` to `.env`
2. Fill in required environment variables
3. Run:

```bash
go run cmd/api/main.go
```

Server default: `http://localhost:8080`

---

## Folder Structure

```text
kunime-core/
├── cmd/
│   └── api/
│       └── main.go                     # entry point (HTTP + cron)
├── internal/
│   ├── bootstrap/
│   │   └── app.go                      # init firestore, fcm, cron, router
│   ├── config/
│   │   └── config.go                   # env/config loader
│   ├── domain/
│   │   ├── anime_popular.go
│   │   ├── app.go
│   │   ├── banner.go
│   │   ├── genre.go
│   │   ├── recommendation.go
│   │   └── watch_event.go
│   ├── infra/
│   │   ├── firestore/
│   │   │   ├── app_repo.go
│   │   │   ├── banner_repo.go
│   │   │   ├── client.go
│   │   │   ├── genre_repo.go
│   │   │   ├── popular_repo.go
│   │   │   ├── recommendation_repo.go
│   │   │   └── watch_repo.go
│   │   └── http/
│   │       └── kunime_api_client.go
│   ├── job/
│   │   └── cleanup_watch_event.go
│   ├── transport/
│   │   └── http/
│   │       ├── router.go
│   │       ├── handler/
│   │       │   ├── app_handler.go
│   │       │   ├── banner_handler.go
│   │       │   ├── genre_handler.go
│   │       │   ├── popular_handler.go
│   │       │   ├── recommendation_handler.go
│   │       │   └── watch_handler.go
│   │       └── middleware/
│   │           └── apikey.go
│   └── usecase/
│       ├── app/
│       │   ├── get_version.go
│       │   └── update_version.go
│       ├── banner/
│       │   └── get_active.go
│       ├── genre/
│       │   └── get_all.go
│       ├── popular/
│       │   └── get_popular.go
│       ├── recommendation/
│       │   ├── build_recommendation.go
│       │   └── get_recommendation.go
│       └── watch/
│           └── record_watch.go
├── .env
├── .env.example
├── .gitignore
├── firebase-service-account.json
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

---

## API Endpoints

### Base URL

- Local: `http://localhost:8080`
- API prefix: `/api/v1`

### Authentication

All `/api/v1/**` routes require header:

- `X-API-Key: <key>`

Route groups:

- Client routes use `CLIENT_API_KEY`
- Admin routes use `ADMIN_API_KEY`

Additional header for watch event:

- `X-Device-Id: <device-id>` (required for `POST /api/v1/events/watch`)

---

### 1) Health Check

**GET** `/health`

**Response 200**

```json
{
  "status": "ok"
}
```

---

### 2) Record Watch Event (Client)

**POST** `/api/v1/events/watch`

Headers:

- `X-API-Key` (client key)
- `X-Device-Id` (required)

Body:

```json
{
  "anime_id": "one-piece",
  "episode": 1097
}
```

Success response `200`:

```json
{
  "status": "recorded"
}
```

Possible error responses:

- `400` invalid body
- `400` `anime_id is required`
- `400` `episode must be greater than 0`
- `400` `device_id is required`
- `401` missing/invalid API key

---

### 3) Get Popular Anime (Client)

**GET** `/api/v1/anime/popular`

Query params:

- `limit` (optional, integer, default `10`)

Success response `200`:

```json
{
  "data": [
    {
      "anime_id": "one-piece",
      "watch_count": 1234
    }
  ]
}
```

Possible error responses:

- `500` `failed to fetch popular anime`
- `401` missing/invalid API key

---

### 4) Get Recommendations (Client)

**GET** `/api/v1/recommendations`

Success response `200`:

```json
{
  "data": [
    {
      "anime_id": "frieren",
      "title": "Sousou no Frieren",
      "rating": 9.1,
      "image": "https://...",
      "source": "jikan",
      "updated_at": "2026-02-20T15:00:00Z"
    }
  ]
}
```

Notes:

- Returns empty array if no data (`"data": []`).

Possible error responses:

- `500` `failed to fetch recommendations`
- `401` missing/invalid API key

---

### 4b) Refresh Recommendations (Admin)

**POST** `/api/v1/admin/recommendations/refresh`

Headers:

- `X-API-Key` (admin key)

Success response `200`:

```json
{
  "message": "recommendations refreshed"
}
```

Possible error responses:

- `500` `failed to refresh recommendations`
- `401` missing/invalid API key

---

### 5) Get Banners (Client)

**GET** `/api/v1/banners`

Success response `200`:

```json
{
  "data": [
    {
      "id": "banner_1",
      "title": "Promo",
      "image_url": "https://...",
      "redirect_url": "https://..."
    }
  ]
}
```

Notes:

- Returns only active banners in valid time window.
- Returns empty array if no data.

Possible error responses:

- `500` `failed to fetch banners`
- `401` missing/invalid API key

---

### 6) Get Genres (Client)

**GET** `/api/v1/genres`

Success response `200`:

```json
{
  "data": [
    {
      "slug": "action",
      "image_url": "https://..."
    }
  ]
}
```

Notes:

- `slug` is derived from Firestore document ID.
- Returns empty array if no data.

Possible error responses:

- `500` `failed to fetch genres`
- `401` missing/invalid API key

---

### 7) Get App Version (Client)

**GET** `/api/v1/app/version`

Success response `200`:

```json
{
  "data": {
    "latest_version": "1.2.3",
    "min_supported_version": "1.0.0",
    "force_update": false,
    "update_url": "https://play.google.com/..."
  }
}
```

Possible error responses:

- `500` `failed to get app version`
- `401` missing/invalid API key

---

### 8) Update App Version (Admin)

**POST** `/api/v1/admin/app/version`

Headers:

- `X-API-Key` (admin key)

Body:

```json
{
  "latest_version": "1.2.4",
  "min_supported_version": "1.1.0",
  "force_update": true,
  "update_url": "https://play.google.com/..."
}
```

Success response `200`:

```json
{
  "message": "app version updated"
}
```

Possible error responses:

- `400` invalid request body
- `400` `latest_version is required`
- `500` `failed to update version`
- `401` missing/invalid API key

---

## Quick cURL Examples

```bash
# Health
curl http://localhost:8080/health

# Popular (client)
curl -H "X-API-Key: $CLIENT_API_KEY" \
  "http://localhost:8080/api/v1/anime/popular?limit=10"

# Record watch (client)
curl -X POST "http://localhost:8080/api/v1/events/watch" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $CLIENT_API_KEY" \
  -H "X-Device-Id: device-123" \
  -d '{"anime_id":"one-piece","episode":1097}'

# Update app version (admin)
curl -X POST "http://localhost:8080/api/v1/admin/app/version" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $ADMIN_API_KEY" \
  -d '{"latest_version":"1.2.4","min_supported_version":"1.1.0","force_update":true,"update_url":"https://play.google.com/..."}'
```
