# Kunime Core

Kunime Core serves the data that the [Kunime app](https://github.com/kudanilll/kunime) doesn't get from the catalog: banners, genre artwork, recommendations, popularity, and app version settings. It records watch events in Firestore and uses the completed-anime feed from [Kunime API](https://github.com/kudanilll/kunime-api) to pick up to five top-rated titles for recommendations.

## Run Locally

Install Go 1.27.4 or newer (see [`go.mod`](go.mod)). Set up a Firestore project for development, start Kunime API, and copy [`.env.example`](.env.example) to `.env`. Use a development database because Core attempts to rebuild recommendations at startup and schedules hourly jobs to refresh recommendations and remove expired watch events.

| Variable | What to set |
| --- | --- |
| `APP_PORT` | HTTP port. Defaults to `8080`. Set it to `8081` if Kunime API is already on `8080`. |
| `APP_ENV` | Defaults to `development` for console logs. Use `production` for JSON logs. |
| `CLIENT_API_KEY` | Key for the client routes. |
| `ADMIN_API_KEY` | Separate key for administrative routes. Don't ship it with the app. |
| `KUNIME_API_BASE_URL` | Kunime API URL ending in `/api/v1`, without a trailing slash (for example, `http://localhost:8080/api/v1` when both services run on the host). |
| `KUNIME_API_KEY` | Must match `API_KEY` in Kunime API. |
| `FIREBASE_CREDENTIALS_BASE64` | Base64-encoded service account JSON for the Firestore project. Required at startup. |

```sh
go run ./cmd/api
```

In a second terminal, check the HTTP server. Use `8081` instead if you changed `APP_PORT`.

```sh
curl --fail http://localhost:8080/health
```

Expect `{"status":"ok"}`. To check the Firestore read path, set `CLIENT_API_KEY` in your shell to the value used by Core and request the recommendations:

```sh
curl --fail -H "X-API-Key: $CLIENT_API_KEY" \
  http://localhost:8080/api/v1/recommendations
```

An empty `data` array is valid before there are any recommendations. If you expected titles, check the Core logs for errors from the initial rebuild.

## Tests

```sh
go test -mod=readonly ./...
```

The watch-event route test uses an in-memory repository; it can run without Firebase credentials or Kunime API.

## Deployment to Dokploy

The [GitHub Actions workflow](.github/workflows/docker-build-push.yml) runs Go checks and a Docker build on pull requests to `master`. A push to `master` runs the Go checks, pushes a Linux/amd64 image to GitHub Container Registry (GHCR) with `latest` and `sha-<full-commit>` tags, then calls Dokploy. [Dependabot](.github/dependabot.yml) checks Go modules, Actions, and base images weekly.

1. In Dokploy, select the **Docker** source. Set the image to `ghcr.io/kudanilll/kunime-core:latest` and point the domain to container port `8080`. If the package is private, configure GHCR credentials with `read:packages` and access to the package.
2. Set the runtime variables from the [table above](#run-locally), using `APP_ENV=production` and `APP_PORT=8080`. `KUNIME_API_BASE_URL` must be reachable from the Core container and end in `/api/v1`.
3. Enable **Auto Deploy**. Copy the HTTPS webhook URL from the application's **Deployments** tab into the GitHub Actions secret `DOKPLOY_WEBHOOK_URL`, at repository level or in the `production` environment. Use one deployment trigger per image to avoid duplicate runs.
4. After a push to `master`, confirm the rollout in Dokploy. Check `GET /health` on the deployed domain, then request `GET /api/v1/recommendations` with the client key. The webhook response only confirms Dokploy accepted the request.

To roll back, point Dokploy to a known-good `sha-<full-commit>` tag and deploy it manually. Switch back to `latest` when you want automatic deployment to resume.

## Folder Structure

- [`cmd/api`](cmd/api) starts the HTTP server and scheduled jobs.
- [`internal/bootstrap`](internal/bootstrap) wires up Fiber, Firestore, and the jobs.
- [`internal/transport/http`](internal/transport/http) contains the routes, handlers, and API key middleware.
- [`internal/usecase`](internal/usecase) contains application logic. [`internal/infra`](internal/infra) contains the Firestore repositories and the Kunime API client.

## API Endpoints

`GET /health` is public. Send `X-API-Key` on every `/api/v1` request: use `CLIENT_API_KEY` on client routes and `ADMIN_API_KEY` under `/api/v1/admin`. The watch-event route also requires `X-Device-Id`.

| Method | Path | Access | Purpose |
| --- | --- | --- | --- |
| GET | `/health` | Public | Check that the HTTP server responds. |
| POST | `/api/v1/events/watch` | Client | Record an episode watch event. |
| GET | `/api/v1/anime/popular` | Client | List anime by watch count. |
| GET | `/api/v1/recommendations` | Client | Read the current recommendation list. |
| GET | `/api/v1/banners` | Client | Read banners active now. |
| GET | `/api/v1/genres` | Client | Read genre slugs and image URLs. |
| GET | `/api/v1/app/version` | Client | Read mobile app version settings. |
| POST | `/api/v1/admin/app/version` | Admin | Update mobile app version settings. |
| POST | `/api/v1/admin/recommendations/refresh` | Admin | Rebuild recommendations on demand. |

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
      "AnimeID": "one-piece",
      "WatchCount": 1234,
      "LastUpdated": "0001-01-01T00:00:00Z"
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
      "source": "completed_top_rated",
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

## Example: Record a Watch Event

Use this against your development instance. The request writes a watch event and increments the popularity counter.

```sh
curl --fail -X POST http://localhost:8080/api/v1/events/watch \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $CLIENT_API_KEY" \
  -H "X-Device-Id: device-123" \
  -d '{"anime_id":"one-piece","episode":1097}'
```

## License

[GNU Affero General Public License v3.0](LICENSE).
