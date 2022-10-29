FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o ./antibrut ./cmd/antibrut

FROM scratch

WORKDIR /app

COPY --from=builder /app .

ENTRYPOINT ["./antibrut"]

EXPOSE 8080