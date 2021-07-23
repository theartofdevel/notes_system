FROM golang:1.16.6-alpine3.14
ADD . /usr/local/go/src/
WORKDIR /usr/local/go/src/
RUN go clean --modcache
RUN go build -o api_service/build/api_service api_service/app/cmd/main/app.go