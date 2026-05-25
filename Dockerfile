FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o signal-server ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/signal-server /signal-server

EXPOSE 8080

ENTRYPOINT ["/signal-server"]