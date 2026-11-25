FROM golang:1.27.4-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/kunime-core ./cmd/api

FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/kunime-core /app/kunime-core

ENV APP_PORT=8080

EXPOSE 8080

CMD ["/app/kunime-core"]
