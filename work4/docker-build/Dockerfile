FROM golang:1.21-alpine AS builder

COPY biz/ /go/src/biz/
COPY pkg/constants /go/src/pkg/constants
COPY pkg/errmsg /go/src/pkg/errmsg
COPY pkg/qiniuyun_oss /go/src/pkg/qiniuyun_oss
COPY pkg/utils /go/src/pkg/utils
COPY go.mod go.sum *.go /go/src/

ENV GO111MODULE=on \
    GOOS=linux \
    FOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/
RUN go mod tidy
RUN go build -o main

FROM jrottenberg/ffmpeg:4.4-alpine

RUN mkdir "/app"
COPY --from=builder /go/src/main /app/main
COPY config.yaml .

RUN chmod +x /app/main

ENTRYPOINT [ "/app/main" ]