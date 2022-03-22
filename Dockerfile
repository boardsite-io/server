FROM golang:1.18-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN go build -o boardsite .

FROM alpine:latest AS deploy
WORKDIR /app
COPY --from=builder /app/boardsite .
EXPOSE 8000
CMD ["./boardsite"]

FROM scratch AS bin
COPY --from=builder /app/boardsite /