# syntax = docker/dockerfile:experimental

FROM golang:1.16.0-alpine AS builder
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -o /out/boardsite .

FROM builder AS deploy
ENV B_PORT=8000
ENV B_CORS_ORIGINS="https://boardsite.io"
EXPOSE 8000
CMD ["/out/boardsite"]

FROM scratch AS bin
COPY --from=builder /out/boardsite /
