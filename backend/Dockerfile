FROM golang:1.26-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY ./cmd ./cmd
COPY ./config ./config
COPY ./pkg ./pkg
COPY ./internal ./internal
COPY ./migrations ./migrations
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download
RUN go mod vendor
RUN go build -a -o start_point cmd/main.go

FROM alpine:3.19

# Install wget for healthcheck and ffmpeg for media metadata extraction
RUN apk add --no-cache wget ffmpeg

WORKDIR /app

COPY --from=builder /build/start_point .

COPY --from=builder /build/migrations ./backend/migrations
COPY --from=builder /build/config ./backend/config

EXPOSE 17777 18888

ENTRYPOINT ["./start_point"]