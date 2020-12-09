FROM golang:alpine

ENV GO111MODULE=auto
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build
COPY . .

RUN go build -o app .

# 声明服务端口
EXPOSE 8080

# 启动容器时运行的命令
CMD ["sleep 1m"]
CMD ["/build/app"]