FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY . .
RUN go build -o boardsite .

FROM builder AS deploy
ENV B_PORT=8000
ENV B_CORS_ORIGINS="https://boardsite.io"
EXPOSE 8000
CMD ["/app/boardsite"]

FROM scratch AS bin
COPY --from=builder /app/boardsite /
