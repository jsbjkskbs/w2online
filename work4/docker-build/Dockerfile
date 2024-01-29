FROM golang:1.21-alpine AS builder

COPY biz/ /go/src/biz/
COPY pkg/constants /go/src/pkg/constants
COPY go.mod go.sum *.go /go/src/

ENV GO111MODULE=on \
    GOOS=linux \
    FOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/
RUN go mod tidy
RUN go build -o main

FROM jrottenberg/ffmpeg:4.3-alpine

RUN mkdir "/app"
COPY --from=builder /go/src/main /app/main

RUN chmod +x /app/main

EXPOSE 12345
ENTRYPOINT [ "/app/main" ]