# Use the official Golang image to create a build artifact.
FROM golang:1.22 as builder

# Create a working directory.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod ./

# Download dependencies.
RUN go mod download

# Copy the source code into th container.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/client

# Use a small image to run the server
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
