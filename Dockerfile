# Use the official Golang base image
FROM golang:1.23 AS builder

# Install Python3
RUN apt-get update && apt-get install -y python && ln -s /usr/bin/python /usr/bin/python

# Set the working directory for the Go app
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -tags netgo -ldflags '-s -w' -o app

# Specify the runtime environment
FROM golang:1.20

# Install Python3 again in the runtime container
RUN apt-get update && apt-get install -y python3 && ln -s /usr/bin/python3 /usr/bin/python

# Copy the built Go binary from the builder stage
COPY --from=builder /app/app /app/

# Set the working directory for runtime
WORKDIR /app

# Command to run the app when the container starts
CMD ["./app"]
