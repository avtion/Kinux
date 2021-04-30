# 编译环境
FROM golang:1.16.3-alpine3.13 AS builder
# 安装GCC
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add gcc g++
# go proxy镜像源
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /kinux
COPY . .
RUN go mod download
RUN go build -o kinux .

# 生产环境
FROM alpine:3.13 AS prod
WORKDIR /kinux
COPY --from=builder /kinux/kinux kinux
EXPOSE 80
CMD ["./kinux"]