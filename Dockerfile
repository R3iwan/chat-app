# Use the official Golang image as the base
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application
RUN go build -o main cmd/user-service/main.go

# Expose the necessary port
EXPOSE 8080

# Run the application
CMD ["./main"]
