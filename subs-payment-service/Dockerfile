FROM golang:1.22.4 AS build

WORKDIR /app

COPY . .

# Change to /app
WORKDIR /app

RUN go build -o main .

CMD ["/app/main"]
