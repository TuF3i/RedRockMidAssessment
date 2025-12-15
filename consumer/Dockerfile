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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consumer app.go

# 运行阶段
FROM alpine:latest

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/consumer .

# 运行应用
CMD ["./consumer"]