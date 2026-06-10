# Stage 1: 编译
FROM golang:1.21-alpine AS builder

RUN apk --no-cache add git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o awards-management .

# Stage 2: 运行
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/awards-management .
COPY config.yaml .
RUN mkdir -p uploads/certificates

EXPOSE 8000

CMD ["./awards-management"]
