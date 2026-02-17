# ---------- Stage 1: build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

# ---------- Stage 2: runtime ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate /usr/local/bin/migrate

COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]
