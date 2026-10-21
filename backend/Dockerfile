# Stage 1: 编译
FROM golang:1.21-alpine AS builder

RUN apk --no-cache add git

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
# 先复制所有源码（go.sum 可能缺新依赖条目，需 go mod tidy）
COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o awards-management .

# Stage 2: 运行
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata curl && \
    rm -rf /var/cache/apk/*
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/awards-management .
COPY config.yaml .
RUN mkdir -p uploads/certificates

EXPOSE 8000

CMD ["./awards-management"]
