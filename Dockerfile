FROM golang:latest as builder
RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu
WORKDIR /go/src/app
COPY . .
RUN go get -v ./...
RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -o app .

FROM alpine:latest
WORKDIR /go/src/app
COPY --from=builder /go/src/app/app .
CMD ["./app"]
EXPOSE 8080
