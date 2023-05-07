FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY ./src/*.go ./src/

# Declare the build argument
ARG VERSION=latest

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o app ./src/

# Run the tests in the container
FROM builder AS tester

RUN go test -v ./...

# Build the final app image
FROM scratch

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]