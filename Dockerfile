FROM golang:1.19-alpine AS builder

WORKDIR /app

RUN apk update && apk add build-base

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN touch .env # fallback

ENV CGO_ENABLED=1

RUN go build -tags=sqlite_unlock_notify -o ./bin/antibrut ./cmd/antibrut

FROM alpine

WORKDIR /app

COPY --from=builder /app/bin/antibrut .
COPY --from=builder /app/*.env .
COPY --from=builder /app/data ./data

ENTRYPOINT ["./antibrut", "run"]

EXPOSE 9090