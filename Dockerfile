FROM golang:1.14.2-alpine as builder
RUN apk add alpine-sdk
WORKDIR /go/src/app
COPY . .
RUN go get -v ./...
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest
WORKDIR /go/src/app
COPY --from=builder /go/src/app/app .
CMD ["./app"]
EXPOSE 8080
