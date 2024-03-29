FROM golang:1.19-buster as builder

ENV GO111MODULE "on"

ARG BUILD_VER

WORKDIR /usr/local/go/src/echo-http
COPY . .
RUN go mod download
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build \
  -v \
  -ldflags "-w -s -X 'main.BuildDatetime=$(date --iso-8601=seconds)' -X 'main.BuildVer=${BUILD_VER}'" \
  -o server \
  ./cmd/main.go

FROM alpine:3.17.0
WORKDIR /app
COPY --from=builder /usr/local/go/src/echo-http/server /app/
# test cert/key
COPY tls.crt /app/
COPY tls.key /app/
RUN apk add curl jq iproute2 bind-tools --no-cache
EXPOSE 8000
ENTRYPOINT ["/app/server"]
LABEL maintainer="flsixtyfour@gmail.com"
LABEL org.label-schema.vcs-url="https://github.com/fl64/echo"
LABEL org.label-schema.docker.cmd="docker run --rm -p 8000:8000 -p 8443:8443 -p 1234:1234 fl64/echo:latest"
