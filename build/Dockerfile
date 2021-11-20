# Build stage
FROM golang:1.14.2-alpine AS build
WORKDIR /
COPY . .
RUN set -ex; \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o jsonstore ./cmd/jsonstore/main.go

# Final image with no source code
FROM alpine:3.12
WORKDIR /
COPY --from=build /jsonstore .
ENTRYPOINT /jsonstore
