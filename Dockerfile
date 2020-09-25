# first stage
FROM golang:1.14.2-stretch as build_base
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

# second stage
FROM build_base AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo -o /bin/app -ldflags "-w -s -X" ./cmd/**/main.go

# final stage
FROM debian:10.3-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
COPY --from=builder /bin/app /bin/app
EXPOSE 8080
ENTRYPOINT ["/bin/app"]