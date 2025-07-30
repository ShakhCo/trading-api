# syntax=docker/dockerfile:1
FROM golang:1.24rc1-alpine

WORKDIR /app

# Install git and bash
RUN apk add --no-cache git bash

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o main .

# Run the executable
CMD ["./main"]
