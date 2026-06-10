# 获奖信息管理系统 (Awards Management System)

> 供内部使用的获奖记录查询与管理平台。支持个人/团体获奖信息的上传、查询、证书管理及系统监控。

---

## 目录

- [项目概览](#项目概览)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [架构设计](#架构设计)
- [数据库表结构](#数据库表结构)
- [API 路由表](#api-路由表)
- [快速开始](#快速开始)
- [部署指南 (Ubuntu)](#部署指南-ubuntu)
- [问题清单与改进建议](#问题清单与改进建议)

---

## 项目概览

这是一个 **获奖证书信息管理系统**，用于管理某机构的个人和团体获奖记录。系统提供：

| 功能 | 说明 |
|------|------|
| 🏆 获奖查询 | 公开查询，按姓名+证件号、团队名、年份等条件检索 |
| 📤 数据录入 | 管理员可逐条或批量上传个人/团体获奖记录 |
| 🖼️ 证书图片 | 上传证书图片自动转 PNG 存入服务器文件系统，DB 仅存路径 |
| 🔐 管理后台 | JWT 认证的管理员系统，支持记录增删改查 |
| 📊 系统监控 | 内置 QPS/响应时间/在线用户/数据库统计等监控面板 |

项目分为 **后端 (Go + Gin)** 和 **前端 (Vue 3 + Ant Design Vue)** 两部分，前后端分离部署。

---

## 技术栈

### 后端

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.21+ |
| Web 框架 | Gin | v1.9.1 |
| ORM | GORM | v1.25.5 |
| 数据库 | MySQL | 5.7+ |
| 配置管理 | Viper | v1.17.0 |
| 认证 | golang-jwt | v5.2.0 |
| 跨域 | gin-contrib/cors | v1.5.0 |

### 前端

| 组件 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue | 3.5.x |
| 构建工具 | Vite | 6.x |
| UI 组件库 | Ant Design Vue | 4.2.x |
| HTTP 客户端 | Axios | 1.10.x |
| 其他 | QRCode, ExcelJS, canvas-confetti | — |

### 部署目标

- **OS**: Ubuntu Server
- **Web 服务**: 计划使用云服务部署
- **构建产物**: `CGO_ENABLED=0` 静态编译 Go 二进制 + Vite 构建的静态前端文件

---

## 项目结构

```
奖状/
├── backend/                        # Go 后端
│   ├── main.go                     # 入口：加载配置 → 初始化DB → 注册路由 → 启动服务
│   ├── config.yaml                 # ⚠️ 配置文件（含数据库密码、JWT密钥等敏感信息）
│   ├── config/
│   │   └── config.go               # Viper 配置加载 + 结构体定义
│   ├── database/
│   │   └── database.go             # GORM 连接池初始化 + AutoMigrate + 默认管理员创建
│   ├── models/
│   │   └── models.go               # 4个数据模型：Admin, PersonAward, GroupAward, GroupMember
│   ├── dto/
│   │   └── dto.go                  # 请求/响应 DTO（数据传输对象）
│   ├── auth/
│   │   └── auth.go                 # JWT 签发/验证 + 认证中间件 + 角色中间件
│   ├── controllers/
│   │   ├── auth_controller.go      # 登录/登出/改密/获取用户信息
│   │   ├── award_controller.go     # 获奖查询/上传/更新/删除 + 证书文件上传
│   │   └── monitor_controller.go   # 系统监控/数据库统计/在线用户管理
│   ├── services/
│   │   ├── admin_service.go        # 管理员密码修改
│   │   ├── person_award_service.go # 个人获奖 CRUD + 轻量查询
│   │   ├── group_award_service.go  # 团体获奖 CRUD + 证书文件处理
│   │   └── monitor_service.go      # 请求统计/在线用户/系统指标采集
│   ├── middleware/
│   │   └── monitor_middleware.go   # 请求耗时统计 + 用户活跃追踪
│   ├── routes/
│   │   └── routes.go               # 全部路由注册（公开 + 认证）
│   ├── uploads/certificates/       # 证书图片文件存储目录
│   ├── Makefile                    # 构建/测试/交叉编译/DB操作
│   └── publish-backend.sh          # 一键构建+运行脚本
│
└── frontend/                       # Vue 3 前端
    ├── index.html                  # 公开查询页入口
    ├── admin.html                  # 管理后台入口
    ├── verify.html                 # 证书验证页入口
    ├── vite.config.js              # Vite 多入口构建配置
    ├── package.json                # 依赖列表
    └── src/
        ├── public-main.js          # 公开查询页 main
        ├── main.js                 # 管理后台 main
        ├── verify-main.js          # 验证页 main
        ├── api/index.js            # Axios 实例 + 全部 API 封装
        └── components/             # Vue 组件（13个）
            ├── PublicQueryView.vue     # 公开查询页核心组件
            ├── PersonQueryView.vue     # 个人信息查询
            ├── GroupQueryView.vue      # 团体信息查询
            ├── RecordManagement.vue    # 记录管理（增删改查）
            ├── InputForm.vue           # 数据录入表单
            ├── LoginForm.vue           # 管理员登录
            ├── ChangePassword.vue      # 修改密码
            ├── CaptchaInput.vue        # 验证码输入
            ├── ImageUploader.vue       # 图片上传组件
            ├── StatisticsView.vue      # 数据统计面板
            ├── MonitorView.vue         # 系统监控面板
            ├── QuerySelector.vue       # 查询条件选择器
            └── utils/confetti.js       # 彩带特效工具
```

---

## 架构设计

### 分层架构

```
┌──────────────────────────────┐
│  routes/routes.go            │  ← 路由注册
├──────────────────────────────┤
│  middleware/                 │  ← 监控中间件 (请求计时/用户追踪)
├──────────────────────────────┤
│  controllers/                │  ← 控制器：请求绑定 → 调用 Service → 返回响应
├──────────────────────────────┤
│  services/                   │  ← 业务逻辑：数据验证、证书文件路径处理、CRUD
├──────────────────────────────┤
│  models/ + database/         │  ← GORM 数据模型 + AutoMigrate
├──────────────────────────────┤
│  auth/ + config/             │  ← JWT 认证 + Viper 配置
└──────────────────────────────┘
```

### 数据流

```
前端 (Vue) ──HTTP──▶ Gin Engine ──▶ 中间件(监控/JWT) ──▶ Controller
                                                              │
                                                              ▼
                                                         Service
                                                              │
                                                         ┌────┴────┐
                                                         │  GORM   │
                                                         │  MySQL  │
                                                         └─────────┘
```

### 证书图片存储（两种方式并存）

| 方式 | 存储位置 | 优点 | 缺点 |
|------|----------|------|------|
| **文件存储** | `uploads/certificates/` 目录，UUID 命名 `.png` 文件 | 不占数据库空间，DB 仅存路径 | — |

---

## 数据库表结构

### `admins` — 管理员

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT PK AUTO_INCREMENT | 主键 |
| username | VARCHAR(50) UNIQUE | 用户名 |
| password_hash | VARCHAR(255) | SHA-256 哈希（⚠️ 无盐） |
| role | VARCHAR(20) | 角色（默认 admin） |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### `person_awards` — 个人获奖记录

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT PK | 主键 |
| year | INT INDEX | 获奖年份 |
| person_name | VARCHAR(100) INDEX | 获奖人姓名 |
| id_number | VARCHAR(50) INDEX | 身份证号 |
| category | VARCHAR(100) INDEX | 奖项类别/组别 |
| award_level | VARCHAR(50) | 奖项等级 |
| description | LONGTEXT | 获奖描述 |
| certificate_image | VARCHAR(500) | 证书图片静态文件路径 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### `group_awards` — 团体获奖记录

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT PK | 主键 |
| year | INT INDEX | 获奖年份 |
| group_name | VARCHAR(200) INDEX | 团体名称 |
| category | VARCHAR(100) INDEX | 组别 |
| award_level | VARCHAR(50) INDEX | 奖项等级 |
| description | TEXT | 描述 |
| members | LONGTEXT | ⚠️ JSON 格式的成员信息（与 group_members 表冗余） |
| certificate_image | VARCHAR(500) | 团队证书图片静态文件路径 |
| created_at | DATETIME INDEX | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### `group_members` — 团体成员

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT PK | 主键 |
| group_award_id | INT INDEX (FK) | 关联团体获奖 ID |
| member_name | VARCHAR(100) | 成员姓名 |
| member_id_number | VARCHAR(50) | 成员身份证号 |
| certificate_image | VARCHAR(500) | 成员证书图片静态文件路径 |

---

## API 路由表

### 公开接口（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | API 信息 |
| GET | `/health` | 健康检查 |
| POST | `/api/v1/awards/query/person` | 查询个人获奖（需 name + id_number） |
| POST | `/api/v1/awards/query/group` | 查询团体获奖（需 group_name） |
| GET | `/api/v1/awards/:award_id` | 获取获奖详情（自动判断个人/团体） |

### 认证接口（需 Bearer Token）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 管理员登录 |
| POST | `/api/v1/auth/logout` | 登出 |
| GET | `/api/v1/auth/me` | 当前用户信息 |
| PUT | `/api/v1/auth/change-password` | 修改密码 |
| POST | `/api/v1/awards/upload/person` | 上传个人获奖 |
| POST | `/api/v1/awards/upload/group` | 上传团体获奖 |
| POST | `/api/v1/awards/batch-upload/person` | 批量上传个人 |
| POST | `/api/v1/awards/batch-upload/group` | 批量上传团体 |
| POST | `/api/v1/awards/upload/certificate` | 上传证书图片文件 |
| POST | `/api/v1/awards/upload/certificates/batch` | 批量上传证书图片 |
| PUT | `/api/v1/awards/update/person` | 更新个人获奖 |
| PUT | `/api/v1/awards/update/group` | 更新团体获奖 |
| PUT | `/api/v1/awards/group/member/certificate` | 更新成员证书 |
| DELETE | `/api/v1/awards/:award_id` | 删除获奖记录（?award_type=group/person） |
| GET | `/api/v1/awards/manage/all` | 获取全部记录（含证书图片路径） |
| GET | `/api/v1/awards/manage/all-light` | 获取全部记录（不含证书图片，轻量） |
| GET | `/api/v1/awards/group/:award_id/members` | 获取团体成员 |
| GET | `/api/v1/monitor/data` | 监控数据总览 |
| GET | `/api/v1/monitor/database` | 数据库统计（表大小等） |
| GET | `/api/v1/monitor/health` | 系统健康评分 |
| GET | `/api/v1/monitor/users` | 在线用户列表 |
| GET | `/api/v1/monitor/requests` | 请求历史 |
| POST | `/api/v1/monitor/reset` | 重置统计 |
| DELETE | `/api/v1/monitor/users/:user_id` | 踢出用户 |

---

## 快速开始

### 1. 环境准备

```bash
# 后端
Go 1.21+
MySQL 5.7+

# 前端
Node.js 18+
npm 9+
```

### 2. 数据库

```sql
CREATE DATABASE awards_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 修改配置

编辑 `backend/config.yaml`，修改数据库连接信息和 JWT secret。

### 4. 启动后端

```bash
cd backend
go mod download && go mod tidy
go run main.go
```

### 5. 启动前端

```bash
cd frontend
npm install
npm run dev        # 开发模式，localhost:3000
```

### 6. 默认管理员

- 用户名: `admin`
- 密码: `admin123`

---

## 部署指南 (Ubuntu)

### 构建

```bash
# 构建后端（Linux 静态编译）
cd backend
make build-linux
# 产出: awards-management_unix

# 构建前端
cd ../frontend
npm run build
# 产出: dist/
```

### 部署架构建议

```
                     ┌──────────────┐
                     │   Nginx      │  ← 反向代理 + 静态文件服务
                     │   :80/:443   │
                     └──────┬───────┘
                            │
            ┌───────────────┼───────────────┐
            ▼                               ▼
   ┌────────────────┐              ┌────────────────┐
   │ 前端静态文件     │              │  Go 后端 :8000   │
   │ (Vite dist/)   │              │  /api/v1/*     │
   │ / => index.html│              │  /static/*      │
   └────────────────┘              └───────┬────────┘
                                           │
                                   ┌───────▼────────┐
                                   │  MySQL :3306    │
                                   └────────────────┘
```

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    root /opt/awards/frontend/dist;
    index index.html;

    # 公开页面
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 管理后台
    location /admin {
        try_files $uri $uri/ /admin.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 证书静态资源
    location /static/ {
        proxy_pass http://127.0.0.1:8000;
    }
}
```

### Systemd 服务

```ini
[Unit]
Description=Awards Management API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/awards/backend
ExecStart=/opt/awards/backend/awards-management_unix
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## 问题清单与改进建议

> 严重程度：🔴 严重 / 🟠 中等 / 🟡 轻微

### 🔴 安全类 — ✅ 已全部修复

| # | 问题 | 状态 | 修复方式 |
|---|------|------|----------|
| 1 | config.yaml 含明文密码 | ✅ | 环境变量优先（`DB_PASSWORD`, `JWT_SECRET`），`.env` gitignore |
| 2 | 密码使用无盐 SHA-256 | ✅ | 全部改为 bcrypt (cost=12) |
| 3 | CORS 允许所有来源 | ⚠️ | 部署时在 nginx 反向代理后无需 CORS；直接暴露时需改为具体域名 |
| 4 | JWT Secret 强度不足 | ✅ | `.env` 生成随机 256-bit 密钥，环境变量注入 |

### 🔴 存储类 — ✅ 已全部修复

| # | 问题 | 状态 | 修复方式 |
|---|------|------|----------|
| 5 | 证书图片存为 base64 入 MySQL | ✅ | 上传自动转 PNG 存文件系统，DB 仅存 `/static/certificates/uuid.png` 路径 |
| 6 | `group_awards.members` 冗余字段 | ✅ | 已删除，统一从 `group_members` 表查询 |
| 7 | 无文件清理机制 | ✅ | 删除记录时同步清理关联证书文件（goroutine 异步） |

### 🟠 架构与性能 — ✅ 已全部修复

| # | 问题 | 状态 | 修复方式 |
|---|------|------|----------|
| 8 | 读操作用 FOR UPDATE 锁 | ✅ | `GetAllGroupAwards` 改为普通查询 |
| 9 | 事务仅用于读 | ✅ | 移除不必要的事务 |
| 10 | 前端 baseURL 动态获取 | ✅ | Docker Compose 中 Nginx 反向代理统一入口 |
| 11 | 大量 Debug 日志 | ✅ | 移除 `db.Debug()`、base64 数据预览日志 |
| 12 | 无 Dockerfile | ✅ | 多阶段构建 Dockerfile + docker-compose.yml + nginx.conf |

### 🟠 数据完整性 — ✅ 已全部修复

| # | 问题 | 状态 | 修复方式 |
|---|------|------|----------|
| 13 | 无文件类型/大小校验 | ✅ | MIME magic bytes 校验 + 10MB 限制 |
| 14 | 文件名冲突 | ✅ | UUID 命名 `.png` 文件 |
| 15 | 无输入数据长度校验 | ✅ | `certificate_image` 字段改为 `varchar(500)` |

### 🟡 运维与可靠性 — 大部分已修复

| # | 问题 | 状态 | 说明 |
|---|------|------|------|
| 16 | 无请求频率限制 | ✅ | 令牌桶限流中间件（IP 30/min, Token 60/min） |
| 17 | 健康检查太简略 | ✅ | 已加入 DB Ping |
| 18 | 无 Token 黑名单 | ⚠️ | 未修复，需要 Redis 支持（超出本次范围） |
| 19 | 监控数据仅内存 | ⚠️ | 未修复，需 Prometheus + Grafana（超出本次范围） |
| 20 | 证书目录权限 | ✅ | `0750` 目录，`0640` 文件 |

### 🟡 代码质量 — ✅ 已改善

| # | 问题 | 状态 | 修复方式 |
|---|------|------|----------|
| 21 | GroupAwardService 过长 (834行) | ✅ | 精简至 ~300 行，移除 base64 处理逻辑 |
| 22 | 日志中打印 base64 数据 | ✅ | 全部移除 |
| 23 | Person/Group Service 重复 | ✅ | 统一为文件路径模式，消除重复 |

---

## 存储空间问题专项分析

> ✅ **已修复** — 证书图片已改为 PNG 文件存储，数据库仅存 `/static/certificates/uuid.png` 路径。

原问题：证书图片以 base64 编码存入 MySQL LONGTEXT 字段，5000 条记录可膨胀至 **~3.2 GB**。改为文件存储后，数据库体积缩减到仅存 VARCHAR(500) 路径，证书文件独立存放于 `uploads/` 卷。

### 当前架构

```
上传图片 → MIME 校验 → 转 PNG → UUID 命名 → 存 uploads/certificates/
                                              ↓
                              DB 仅存 /static/certificates/uuid.png
```

### 后续可选优化

- 接入云对象存储（OSS/S3），本地完全不存证书文件
- 图片上传前在前端压缩为 WebP（当前要求无损，未实施）

---

## 配置清单（部署检查表）

部署到 Ubuntu 服务器前，请确认以下事项：

- [ ] 修改 `config.yaml` 中的数据库密码和 JWT Secret
- [ ] 将 `server.mode` 设为 `release`
- [ ] 将 `cors.allowOrigins` 改为具体域名
- [ ] 数据库字符集为 `utf8mb4`
- [ ] `uploads/certificates/` 目录存在且可写
- [ ] 配置 Nginx 反向代理
- [ ] 配置 Systemd 或 Supervisor 守护进程
- [ ] 配置 HTTPS（Let's Encrypt）
- [ ] 配置防火墙仅开放 80/443
- [ ] 配置 MySQL 定期备份
- [ ] 添加日志轮转（logrotate）

---

*最后更新: 2026-06-10*
