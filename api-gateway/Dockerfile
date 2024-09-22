FROM golang:1.22.4 AS build

WORKDIR /api-gateway

COPY . .

RUN go build -o main .

CMD ["/api-gateway/main"]
