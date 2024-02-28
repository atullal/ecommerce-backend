# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the Go Mod and Sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o user-service .

# Make port 50051 available to the world outside this container
EXPOSE 50051

# Run the Go application when the container launches
CMD ["./user-service"]
