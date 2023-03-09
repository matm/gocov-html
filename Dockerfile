FROM golang:1.19-alpine3.16

RUN apk update && apk add git
RUN go install github.com/axw/gocov/gocov@latest
RUN git clone https://github.com/matm/gocov-html.git && \
        cd gocov-html && \
        go build github.com/matm/gocov-html/cmd/generator && \
        go generate ./... && \
        go install github.com/matm/gocov-html/cmd/gocov-html
