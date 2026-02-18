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

```
kunime-core/
├── cmd/
│   └── api/
│       └── main.go          # entry point (HTTP + cron)
├── internal/
│   ├── bootstrap/
│   │   └── app.go           # init firestore, fcm, cron, router
│   ├── config/
│   │   └── config.go        # env, secrets, config loader
│   ├── domain/
│   │   ├── watch_event.go
│   │   ├── anime_popular.go
│   │   ├── recommendation.go
│   │   ├── banner.go
│   │   ├── genre.go
│   │   └── app.go
│   ├── infra/
│   │   ├── firestore/
│   │   │   ├── client.go
│   │   │   ├── watch_repo.go
│   │   │   ├── popular_repo.go
│   │   │   ├── banner_repo.go
│   │   │   ├── genre_repo.go
│   │   │   ├── recommendation_repo.go
│   │   │   └── app_repo.go
│   │   ├── fcm/
│   │   │   └── client.go
│   │   └── http/
│   │       └── kunime_api_client.go
│   ├── job/
│   │   └── cleanup_watch_event.go
│   ├── transport/
│   │   └── http/
│   │       ├── router.go
│   │       ├── middleware/
│   │       │   ├── apikey.go
│   │       │   └── device.go
│   │       └── handler/
│   │           ├── watch_handler.go
│   │           ├── popular_handler.go
│   │           ├── recommendation_handler.go
│   │           ├── banner_handler.go
│   │           ├── genre_handler.go
│   │           └── app_handler.go
│   ├── usecase/
│   │   ├── watch/
│   │   │   └── record_watch.go
│   │   ├── popular/
│   │   │   └── get_popular.go
│   │   ├── recommendation/
│   │   │   ├── build_recommendation.go
│   │   │   └── get_recommendation.go
│   │   ├── banner/
│   │   │   └── get_active.go
│   │   ├── genre/
│   │   │   └── get_genres.go
│   │   └── app/
│   │       ├── get_version.go
│   │       └── update_version.go
│   └── util/
│       ├── time.go
│       └── response.go
├── go.mod
├── go.sum
└── README.md
```
