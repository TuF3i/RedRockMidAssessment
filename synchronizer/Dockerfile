# 构建阶段
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 设置GOPROXY
ENV GOPROXY=https://goproxy.cn,direct

# 安装必要的系统依赖
RUN apk add --no-cache git tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
# CGO_ENABLED=0 表示禁用CGO，静态编译
# -ldflags "-w -s" 去除调试信息和符号表，减小二进制文件大小
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o synchronizer app.go

# 第二阶段：运行阶段
FROM alpine:latest

# 创建非root用户
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/synchronizer .

# 设置静态文件权限
RUN chown -R appuser:appuser /root/

# 切换到非root用户
USER appuser

# 运行应用
CMD ["./synchronizer"]