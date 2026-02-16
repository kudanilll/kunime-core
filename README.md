# Kunime Core

## Run Locally

1. Copy `.env.example` to `.env`
2. Fill in the environment variables
3. Run the application

```bash
go run cmd/api/main.go
```

The application will be available at `http://localhost:8080`

## Folder Structure

```kunime-core/
├── cmd/
│   └── api/
│       └── main.go          # entry point (HTTP + cron)
│
├── internal/
│   ├── config/
│   │   └── config.go        # env, secrets, config loader
│   │
│   ├── bootstrap/
│   │   └── app.go           # init firestore, fcm, cron, router
│   │
│   ├── transport/
│   │   └── http/
│   │       ├── router.go
│   │       ├── middleware/
│   │       │   └── device.go
│   │       └── handler/
│   │           ├── watch_handler.go
│   │           ├── popular_handler.go
│   │           ├── recommendation_handler.go
│   │           ├── banner_handler.go
│   │           ├── genre_handler.go
│   │           └── app_handler.go
│   │
│   ├── usecase/
│   │   ├── watch/
│   │   │   └── record_watch.go
│   │   ├── popular/
│   │   │   └── aggregate_popular.go
│   │   ├── recommendation/
│   │   │   └── global_top_rated.go
│   │   ├── banner/
│   │   │   └── get_banners.go
│   │   ├── genre/
│   │   │   └── get_genres.go
│   │   └── app/
│   │       └── get_version.go
│   │
│   ├── job/
│   │   ├── notify_release.go
│   │   └── recommendation_cron.go
│   │
│   ├── domain/
│   │   ├── watch_event.go
│   │   ├── anime_popular.go
│   │   ├── recommendation.go
│   │   ├── banner.go
│   │   ├── genre.go
│   │   └── release_state.go
│   │
│   ├── infra/
│   │   ├── firestore/
│   │   │   ├── client.go
│   │   │   ├── watch_repo.go
│   │   │   ├── popular_repo.go
│   │   │   ├── banner_repo.go
│   │   │   ├── genre_repo.go
│   │   │   ├── recommendation_repo.go
│   │   │   └── release_repo.go
│   │   │
│   │   ├── fcm/
│   │   │   └── client.go
│   │   │
│   │   └── animeapi/
│   │       └── client.go    # client ke kunime-api
│   │
│   └── util/
│       ├── time.go
│       └── response.go
│
├── go.mod
├── go.sum
└── README.md
```
