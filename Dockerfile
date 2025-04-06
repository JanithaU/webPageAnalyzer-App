# Use the official Go image to build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -o app .

# Start a new stage from a smaller base image
FROM alpine:latest

# Install necessary dependencies to run the Go app (like ca-certificates)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

COPY templates/ ./templates/

EXPOSE 8080
CMD ["./app"]
