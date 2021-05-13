FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go build -o run ./app

CMD ["/go/src/app/run"]