# FROM golang:1.19.0-alpine AS builder

FROM golang:latest AS builder

# ENV GOLANG_VERSION 1.21.7

# Support CGO and SSL (if needed)
RUN apk --no-cache add gcc g++ make
RUN apk add git

WORKDIR /backend  # Adjust to match your Go module path

# ENV GO111MODULE=on
ENV GIN_MODE=release

# Copy go.mod and go.sum first for efficient caching
COPY go.mod go.sum ./
RUN go mod tidy
#RUN go mod download
#RUN go mod vendor
#RUN go mod verify

COPY . ./
RUN ls -la ./*

# Build the binary (adjust path if needed)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /backend .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

#ENV GO111MODULE=on

# Adjust WORKDIR to match build output
WORKDIR /service-desk

ENV GIN_MODE=release

COPY --from=builder /service-desk .

# Create non-root user
RUN addgroup -g 1001 docker && \
    adduser -S -u 1001 -G docker docker

RUN chown -R docker:docker /service-desk && \
    chmod 755 /service-desk

USER docker:docker

EXPOSE 7788

ENTRYPOINT ["./service-desk"]