FROM golang:alpine AS builder

ENV GO111MODOULE=on \
GOPROXY=https://goproxy.cn,direct \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o bluebell .



###################
# 接下来创建一个小镜像
###################
FROM debian:stretch-slim

COPY ./conf /conf
COPY ./wait-for.sh /
COPY ./templates /templates
COPY ./static /static

COPY --from=builder /build/bluebell /

RUN set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends netcat \
    && chmod 755 wait-for.sh

EXPOSE 8081