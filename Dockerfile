FROM golang:alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download
RUN go mod verify

# Copy the source code
COPY . .

# Declare the build argument
ARG VERSION=latest

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o app ./src/

# Run the tests in the container
FROM builder AS tester

# Run the tests
RUN go test -v ./...

# Build the final app image
# FROM scratch <-- This is the smallest image possible but cannot connect a shell for debugging
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port for the app
EXPOSE 8080

# Run the app binary
CMD ["./app"]