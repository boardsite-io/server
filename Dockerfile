# syntax = docker/dockerfile:experimental

FROM golang:1.15.5-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -o /out/boardsite ./cmd/boardsite

FROM builder AS unit-test
ENV B_REDIS_HOST=localhost
ENV B_REDIS_PORT=63000
RUN --mount=type=cache,target=/root/.cache/go-build \
go test -v ./test

FROM builder AS debug
ENV B_PORT=8000
ENV B_REDIS_HOST=b-redis-debug
ENV B_REDIS_PORT=6379
CMD ["/out/boardsite"]

FROM scratch AS bin
COPY --from=builder /out/boardsite /