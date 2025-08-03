#!/bin/bash
# build-with-docker.sh

# 清理旧文件
rm -f email-forwarder Dockerfile

# 使用 Docker 容器编译 Go 应用
echo "使用 Docker 容器编译 Go 应用..."
docker run --rm -v "$PWD":/app -w /app/src golang:1.21-alpine \
    sh -c "go mod download && CGO_ENABLED=0 GOOS=linux go build -o ../email-forwarder ./cmd/"

# 检查编译结果
if [ ! -f "email-forwarder" ]; then
    echo "编译失败: 未生成二进制文件"
    exit 1
fi

# 创建 Dockerfile
cat > Dockerfile <<'EOL'
FROM alpine:3.18
WORKDIR /app
COPY email-forwarder .
COPY init.sql .
EXPOSE 8080
CMD ["./email-forwarder"]
EOL

# 创建 init.sql (如果不存在)
if [ ! -f "init.sql" ]; then
    echo "CREATE DATABASE IF NOT EXISTS email_forwarder;" > init.sql
fi

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t email-forwarder .

# 清理临时文件
rm -f email-forwarder Dockerfile

echo ""
echo "构建成功！"
echo "运行容器: docker run -d -p 8080:8080 email-forwarder"