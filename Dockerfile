############################
# STEP 1 build executable binary
############################
FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cl-api

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /app/cl-api /cl-api

EXPOSE 8000

# Run the hello binary.
ENTRYPOINT ["/cl-api"]
