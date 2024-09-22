FROM golang:1.22.4 AS build

WORKDIR /subs-payment-service

COPY . .

RUN go build -o main .

CMD ["/subs-payment-service/main"]
