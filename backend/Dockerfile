# Build Stage
FROM golang:alpine AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV GIN_MODE=release

# Copy go.mod and go.sum first for efficient caching
COPY go.mod go.sum ./

# go mod tidy
RUN go mod tidy

COPY . ./

# Install necessary dependencies
RUN apk --no-cache add gcc g++ make

# Copy the .env file into the image
COPY .env ./

# Copy the MySQL Deployment Script
COPY service-desk.mysql ./

# App Config
COPY config/app.conf ./

# Build the binary
RUN CGO_ENABLED=0 go build -o /backend .

# Final Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV GO111MODULE=on
ENV GIN_MODE=release

WORKDIR /backend

COPY --from=builder /backend .

# Create non-root user
RUN addgroup -g 1001 docker && \
    adduser -S -u 1001 -G docker docker

RUN chown -R docker:docker /backend && \
    chmod 755 /backend

USER docker:docker

EXPOSE 4306
EXPOSE 3306
EXPOSE 7788

ENTRYPOINT ["./backend"]
