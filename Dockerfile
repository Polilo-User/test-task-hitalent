FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./migrations ./migrations

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=0 /app/server ./
COPY --from=0 /app/migrations ./migrations

EXPOSE 8080

CMD ["./server"]