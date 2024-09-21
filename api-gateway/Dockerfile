FROM golang:1.22.4 AS build

WORKDIR /gateway

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["/gateway/main"]
