FROM golang:1.16.5-alpine3.14 AS builder

WORKDIR /go/src/app

ADD app/ /usr/local/go/src/
WORKDIR /usr/local/go/src/

RUN go build -mod=readonly -o app cmd/main/app.go

FROM alpine:3.14

COPY --from=builder /usr/local/go/src/app /

CMD ["/app"]