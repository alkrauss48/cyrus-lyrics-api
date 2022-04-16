# Build the executable with Go
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ARG GOOS=linux
ARG GOOS=linux
ARG CGO_ENABLED=0

COPY . .

RUN set -x \
    && go generate \
    && go build -ldflags="-w -s"

## Copy the executable to a small image
FROM alpine

COPY --from=builder /app/cyrus-lyrics-api /

ENTRYPOINT ["/cyrus-lyrics-api"]
