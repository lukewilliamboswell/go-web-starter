FROM golang as app-builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app

EXPOSE 8080

CMD ["./app"]