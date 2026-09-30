FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git & build dependencies
RUN apk add --no-cache git gcc musl-dev

# Install golang-migrate & air
RUN go install github.com/air-verse/air@latest && \
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Production stage
FROM alpine:3.19 AS runner

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server .
COPY --from=builder /app/.env.example ./.env

EXPOSE 8080

CMD ["./server"]
