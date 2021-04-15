FROM golang:1.16-buster as builder

ENV GO111MODULE "on"

WORKDIR /usr/local/go/src/echo-http
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -o server \
  -ldflags "-X main.BuildDatetime=$(date --iso-8601=seconds)" \
  ./cmd/server.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /usr/local/go/src/echo-http/server /app/
EXPOSE 8000
ENTRYPOINT ["/app/server"]
