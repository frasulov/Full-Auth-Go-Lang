FROM golang:1.16.7 AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 4001
CMD go run main.go --start-service