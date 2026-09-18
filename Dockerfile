FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
