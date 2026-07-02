FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
  && addgroup -S app \
  && adduser -S app -G app

WORKDIR /app

COPY --from=builder /api /app/api
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY sql/migrations /app/sql/migrations

USER app

EXPOSE 8080

ENTRYPOINT ["/app/api"]
