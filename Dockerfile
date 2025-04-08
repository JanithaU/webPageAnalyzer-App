# Use the official Go image to build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o app ./cmd/web

# Start a new stage from a smaller base image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .

COPY web/ ./web/

EXPOSE 8080

CMD ["./app"]
