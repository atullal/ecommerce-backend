# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the Go Mod and Sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the app
RUN go build -o rest-service .

# Expose port 8080
EXPOSE 8080

# Run the app
CMD ["./rest-service"]
