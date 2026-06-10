# 获奖管理系统 — 安全与存储修复设计

> 状态: 已审批 | 日期: 2026-06-10

## 目标

修复 README 中发现的 23 个问题中可代码直接修复的部分，加上运维基础设施补齐，使项目达到可上线标准。

## 约束

- 图片质量不能损失（不压缩），统一转 PNG 存储
- 不需要外部云服务（OSS 等）
- 单机 Docker Compose 部署
- 无需数据迁移（现有数据为测试数据）
- API 路由和前端接口保持兼容

---

## 一、代码修改

### 1.1 密码安全 — SHA-256 → bcrypt

**涉及文件：**

| 文件 | 改动 |
|------|------|
| `backend/auth/auth.go` | `AuthenticateUser()` 改用 `bcrypt.CompareHashAndPassword` |
| `backend/database/database.go` | `createDefaultAdmin()` 改用 `bcrypt.GenerateFromPassword`（cost=12） |
| `backend/services/admin_service.go` | `ChangePassword()` / `AuthenticateAdmin()` 同步改为 bcrypt |

**兼容性：** 首次启动时默认 admin 用 bcrypt 创建。旧密码哈希自然失效，无需迁移。

### 1.2 证书图片存储 — base64 入库 → PNG 文件存储

**核心变更：** 所有证书图片统一转 PNG 存入 `uploads/certificates/`，DB 仅存 `/static/certificates/<uuid>.png`。

**涉及文件：**

| 文件 | 改动说明 |
|------|----------|
| 新增 `backend/utils/image.go` | `ConvertToPNG(input []byte) ([]byte, error)` — 用 `image.Decode` + `png.Encode`，无损转换 |
| `backend/models/models.go` | `certificate_image` 字段类型：`longtext` → `varchar(500)` |
| `backend/controllers/award_controller.go` | 上传接口：校验 MIME → 限制 10MB → 调用转 PNG → UUID 命名存文件 → 返回路径 |
| `backend/services/person_award_service.go` | 移除 base64 处理逻辑 |
| `backend/services/group_award_service.go` | 移除所有 base64 格式检测/修复代码，精简约 200 行 |

**文件校验规则：**
- 允许 MIME：`image/jpeg`, `image/png`, `image/gif`, `image/bmp`, `image/webp`
- 单文件最大：10MB
- 文件名：`<uuid>.png`（UUID 防冲突）

### 1.3 移除 `group_awards.members` 冗余字段

| 文件 | 改动 |
|------|------|
| `backend/models/models.go` | 删除 `Members string` 字段 |
| `backend/services/group_award_service.go` | 移除 JSON marshal/unmarshal members 的代码，成员信息统一从 `GroupMembers` 关联表读取 |

### 1.4 读操作移除 FOR UPDATE 锁

| 文件 | 位置 | 改动 |
|------|------|------|
| `backend/services/group_award_service.go` | `GetAllGroupAwards()` | 移除 `tx.Begin()` + `FOR UPDATE`，改为普通 `db.Find()` |

### 1.5 文件上传校验（增强）

| 文件 | 位置 | 改动 |
|------|------|------|
| `backend/controllers/award_controller.go` | `UploadCertificateImage` | 校验 Content-Type + 魔数检测 + 大小限制 |
| `backend/controllers/award_controller.go` | `BatchUploadCertificateImages` | 同上 |

### 1.6 请求频率限制

新增 `backend/middleware/rate_limit.go`：
- 公开查询接口：30 req/min per IP
- 认证接口：60 req/min per token（fallback IP）
- 上传接口：10 req/min per token
- 基于令牌桶算法，内存实现

### 1.7 健康检查增强

`routes/routes.go` `/health` 加入 `database.DB.Ping()` 检查。

### 1.8 配置改为环境变量优先

`config/config.go`：
- 数据库密码：优先 `DB_PASSWORD` 环境变量，fallback config.yaml
- JWT Secret：优先 `JWT_SECRET` 环境变量
- 数据库地址：优先 `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER`
- Server Mode：优先 `GIN_MODE`

### 1.9 清理调试日志

`services/group_award_service.go`：移除所有 base64 数据预览日志、`db.Debug()` 调用、冗余的 SQL 打印。仅保留关键操作日志（创建/更新/删除成功）。

---

## 二、运维基础设施

### 文件清单

```
backend/
├── Dockerfile
├── docker-compose.yml
├── nginx.conf
├── .env.example
└── .env                    # gitignore
```

### Dockerfile

```dockerfile
# Stage 1: 编译
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o awards-management .

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
```

### docker-compose.yml 结构

- `mysql` — MySQL 8.0，数据卷 `mysql_data`，健康检查
- `backend` — depends_on mysql (healthy)，环境变量注入，挂载 uploads 卷
- `nginx` — 挂载前端 dist + nginx.conf，端口 80，depends_on backend

### nginx.conf 要点

- 客户端 body 限制 20MB
- `/` → 前端 dist 静态文件
- `/api/` → proxy_pass backend:8000
- `/static/` → proxy_pass backend:8000
- gzip 压缩开启

---

## 三、不变部分

- 前端代码零改动（API 接口完全兼容）
- 数据库表名不变
- API 路由不变
- 监控中间件逻辑不变
- 无外部云服务依赖

---

## 四、上线检查清单

- [ ] `.env` 已填写真实值且加入 `.gitignore`
- [ ] JWT Secret 已更换为随机 256-bit 密钥
- [ ] 默认 admin 密码已修改
- [ ] `config.yaml` server.mode = `release`
- [ ] `uploads/` 目录有持久化卷
- [ ] MySQL 数据卷有备份策略
- [ ] 防火墙仅开放 80/443
- [ ] Docker Compose 已测试启动

---

## 五、变更影响分析

| 维度 | 影响 |
|------|------|
| API 兼容性 | 无破坏性变更（路径和请求格式不变） |
| 前端兼容性 | 完全兼容（证书路径仍是 `/static/certificates/xxx`） |
| 数据库兼容性 | `certificate_image` 列类型变更需 AutoMigrate |
| 已有数据 | 测试数据，不迁移 |
| 性能 | 正向提升（base64 不再加载到内存） |
| 存储 | 数据库体积大幅缩减（base64 移除），文件系统增长可控 |
