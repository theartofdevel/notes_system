FROM golang:1.16.6-alpine3.14
ADD api_service/app/ /usr/local/go/src/
WORKDIR /usr/local/go/src/
RUN go clean --modcache
RUN go build -o api_service cmd/main/app.go

CMD ["/usr/local/go/src/api_service"]