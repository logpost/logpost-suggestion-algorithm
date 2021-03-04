FROM golang:1.16-alpine

COPY . /go/src/github.com/logpost/logpost-suggestion-algorithm
WORKDIR /go/src/github.com/logpost/logpost-suggestion-algorithm

RUN go get ./...
RUN go get -u github.com/cosmtrek/air
CMD air -c .air.toml