FROM golang:alpine AS build

WORKDIR /go/src/solrapi

COPY . .

RUN go build -o /go/bin/solrapi main.go

EXPOSE 8080

CMD /go/bin/solrapi
