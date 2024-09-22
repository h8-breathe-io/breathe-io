FROM golang:1.22.4 AS build

WORKDIR /email-notif-service

COPY . .

RUN go build -o main .

CMD ["/email-notif-service/main"]
