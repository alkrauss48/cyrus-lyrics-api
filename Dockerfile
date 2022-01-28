############################
# STEP 1 build optimized executable binary
############################
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cl-api

############################
# STEP 2 build a small image
############################
FROM alpine

COPY --from=builder /app/cl-api /

EXPOSE 8000

ENTRYPOINT ["/cl-api"]
