# syntax=docker/dockerfile:1

FROM golang:1.21.3-bookworm

# Set destination for COPY
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y netcat-openbsd

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /proxysql-clustering

# Run
CMD ["/proxysql-clustering"]
