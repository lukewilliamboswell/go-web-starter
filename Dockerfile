FROM golang as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app

# Run the tests in the container
FROM builder AS tester

RUN go test -v ./...

# Build the final app image
FROM scratch

COPY --from=builder ./app .

EXPOSE 8080

CMD ["./app"]