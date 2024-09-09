# 使用官方Go镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用轻量级的alpine镜像作为最终镜像
FROM alpine:latest  

WORKDIR /root/

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/main .

# 暴露应用端口（根据您的应用实际使用的端口进行修改）
EXPOSE 8080

# 运行应用
CMD ["./main"]