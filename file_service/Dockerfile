FROM golang:1.16.6-alpine3.14
ADD app/ /usr/local/go/src/
WORKDIR /usr/local/go/src/
RUN go clean --modcache
RUN go build -o app cmd/main/app.go

CMD ["/usr/local/go/src/app"]