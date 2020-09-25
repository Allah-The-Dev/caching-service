FROM golang:1.14.2-alpine as builder
RUN apk add alpine-sdk
WORKDIR /go/app
COPY . .
RUN get -v ./...
RUN GOOS=linux GOARCH=amd64 go build -o caching-service -tags musl

FROM alpine:latest as runner
WORKDIR /root/
COPY --from=builder /go/app/caching-service .
ENTRYPOINT /root/caching-service
EXPOSE 8080