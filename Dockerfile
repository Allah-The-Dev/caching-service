FROM golang:latest as builder
RUN apk add --no-cache libc6-compat
WORKDIR /go/src/app
COPY . .
RUN go get -v ./...
RUN CGO_ENABLED=1 GOOS=linux go build -o app

FROM alpine:latest
WORKDIR /go/src/app
COPY --from=builder /go/src/app/app .
CMD ["./app"]
EXPOSE 8080
