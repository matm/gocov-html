FROM golang:1.19-alpine3.16

RUN apk update && apk add git
RUN go install github.com/axw/gocov/gocov@latest
RUN go install github.com/matm/gocov-html/cmd/gocov-html@latest
