FROM golang:1.8

WORKDIR /go/src/github.com/richardpanda/quick-poll/server
COPY . .

RUN go-wrapper download
RUN go-wrapper install
