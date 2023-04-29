FROM golang:1.20 as golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build
ADD . /build

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o ./lark-vkm ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash
COPY --from=golang /build/lark-vkm /app
COPY --from=golang /build/role_list.yaml /app

EXPOSE 9000
ENTRYPOINT ["/app/lark-vkm", "server"]
