############################
# STEP 1 build executable binary
############################
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /cl-api

# Go seems to need an env file, or else it won't let you access any env vars
RUN touch /env.empty

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /cl-api /cl-api
COPY --from=builder /env.empty /.env

EXPOSE 8000

# Run the hello binary.
ENTRYPOINT ["/cl-api"]
