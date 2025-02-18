# Use golang image as a build stage
FROM golang:1.22 as builder

WORKDIR /app

# Set up environment variables for Go build
ENV POSTGRES_DB=aitubet
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=adminpassword

# Copy the Go app source code and SQL initialization file
COPY . .

# Build the Go application

RUN go build -o ./main .

# Use a smaller, production-ready image to run the app
FROM debian:bookworm-slim

WORKDIR /app

# Copy the built Go binary from the previous stage
COPY --from=builder /app/main .

# Install necessary certificates for SSL verification
RUN apt-get update && apt-get install -y ca-certificates

# Expose port for the app
EXPOSE 8080

# Start the application
CMD ["./main"]
