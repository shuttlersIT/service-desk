# Build Stage
FROM golang:alpine AS builder

WORKDIR /Users/itplatforms/Dev/Projects/service-desk/backend

ENV GO111MODULE=off
ENV GIN_MODE=release

# Copy go.mod and go.sum first for efficient caching
COPY go.mod go.sum ./
# RUN go mod tidy

COPY . ./

# Install necessary dependencies
RUN apk --no-cache add gcc g++ make

# Build the binary
RUN GOOS=linux go build -o /backend .

# Final Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

ENV GO111MODULE=on
ENV GIN_MODE=release

WORKDIR /backend

COPY --from=builder /backend .

RUN go mod tidy

# Create non-root user
RUN addgroup -g 1001 docker && \
    adduser -S -u 1001 -G docker docker

RUN chown -R docker:docker /backend && \
    chmod 755 /backend

USER docker:docker

EXPOSE 7788

ENTRYPOINT ["./backend"]