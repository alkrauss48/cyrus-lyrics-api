# Build a static executable with Go
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/cyrus-lyrics-api .

# Copy the executable to a small image
FROM alpine:3.24

COPY --from=builder /out/cyrus-lyrics-api /cyrus-lyrics-api

# Run as an unprivileged user (nobody)
USER 65534:65534

EXPOSE 8000

ENTRYPOINT ["/cyrus-lyrics-api"]
