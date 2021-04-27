FROM golang:1.15-buster as builder

ENV GO111MODULE "on"

WORKDIR /usr/local/go/src/echo-http
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
  -o server \
  -ldflags "-X main.BuildDatetime=$(date --iso-8601=seconds)" \
  ./cmd/server.go

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /usr/local/go/src/echo-http/server /app/
RUN apk add curl jq --no-cache
EXPOSE 8000
ENTRYPOINT ["/app/server"]
LABEL maintainer="flsixtyfour@gmail.com"
LABEL org.label-schema.vcs-url="https://github.com/fl64/http-echo"
LABEL org.label-schema.docker.cmd="docker run --rm -p 8000:8000 fl64/echo-http:latest"
