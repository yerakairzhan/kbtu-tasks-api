FROM golang:1.24.9-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -h /app appuser
COPY --from=builder /out/app /app/app
COPY database /app/database
COPY docs /app/docs
COPY .env.example /app/.env.example
USER appuser
EXPOSE 8080
CMD ["/app/app"]
