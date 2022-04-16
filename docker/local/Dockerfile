FROM golang:1.16-alpine

WORKDIR /app

COPY . .
RUN go mod download

ENTRYPOINT ["go", "run", "."]
