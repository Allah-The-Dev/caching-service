FROM golang:latest as builder
WORKDIR /go/src/app
COPY . .
RUN go get -v ./...
RUN GOOS=linux go build -o app

FROM alpine:latest
WORKDIR /go/src/app
COPY --from=builder /go/src/app/app .
CMD ["./app"]
EXPOSE 8080
