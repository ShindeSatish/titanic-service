# Use an official Go runtime as a parent image
FROM golang:1.21-alpine

# Set the working directory to /app
WORKDIR /app

# Enable CGO
ENV CGO_ENABLED=1

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy the current directory contents into the container at /app
COPY . .

# Install any dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
