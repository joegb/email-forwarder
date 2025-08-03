# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY . .

# 下载依赖
RUN go mod download

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o email-forwarder ./cmd/email-forwarder/

# 最终运行阶段
FROM alpine:3.18

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/email-forwarder .

# 复制初始化脚本
COPY init.sql .

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./email-forwarder"]