# 获奖信息管理系统

用于录入、管理和公开查询个人/团队获奖记录，支持证书图片、证书编号、团队成员证书、防伪二维码、批量录入、名称联想和运行监控。

线上查询站点：<https://hkicfest.com>

> 本仓库只保存程序源码和脱敏配置示例。生产服务器地址、账号、密钥、数据库备份及日常运维命令记录在单独的交接文档中，不进入 Git。

## 功能

- 个人获奖记录：姓名、证件号码、证书编号、年份、组别、奖项、证书图片
- 团队获奖记录：团队名称、团队证书编号、团队证书、成员及成员个人证书
- 公开查询：按个人或团队查询，团队名称支持模糊联想
- 管理后台：新增、编辑、删除、批量导入、图片上传、数据导出与系统监控
- 防伪验证：为个人证书和团队成员证书生成带二维码的查看/下载版本
- 图片存储：支持腾讯云 COS；未配置或不可用时回退至本地目录
- 安全措施：JWT 鉴权、登录/查询限流、上传格式校验、10 MB 请求限制

## 技术结构

```text
浏览器
  └─ Nginx（静态页面、TLS、反向代理）
       ├─ Vue 3 + Vite + Ant Design Vue
       └─ Go + Gin API
            ├─ MySQL 8
            └─ 腾讯云 COS / 本地 uploads 回退
```

主要目录：

```text
backend/                 Go API、数据库迁移、Docker 配置
  controllers/           HTTP 控制器
  services/              业务逻辑
  models/                GORM 模型
  storage/               COS 与本地存储适配
frontend/                Vue 前端
  src/                   页面与组件源码
  dist/                  生产构建产物
docs/                    设计/修复记录
```

## 本地开发

### 依赖

- Go 1.21+
- Node.js 20+
- MySQL 8，或 Docker + Docker Compose

### 使用 Docker Compose

```bash
cd backend
cp .env.example .env
# 编辑 .env，至少设置数据库密码、JWT_SECRET 和 ADMIN_DEFAULT_PASSWORD
docker compose up -d --build
```

首次启动时会自动迁移数据库并创建管理员。`ADMIN_DEFAULT_PASSWORD` 只在管理员尚不存在时使用；请在首次登录后立即修改密码。

构建前端：

```bash
cd frontend
npm ci
npm run build
```

### 分别运行

先创建 MySQL 数据库并修改 `backend/config.yaml`，也可以使用环境变量覆盖配置：

```bash
cd backend
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=awards_user
export DB_PASSWORD='replace-me'
export DB_NAME=awards_management
export JWT_SECRET='replace-with-a-long-random-value'
export ADMIN_DEFAULT_PASSWORD='replace-me-on-first-start'
go run .
```

另开终端启动前端：

```bash
cd frontend
npm ci
npm run dev
```

Vite 默认监听 `http://localhost:3000`，API 默认监听 `http://localhost:8000`。

## 配置

| 变量 | 用途 | 必填 |
| --- | --- | --- |
| `DB_HOST` / `DB_PORT` | MySQL 地址和端口 | 是 |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | MySQL 账号和数据库 | 是 |
| `JWT_SECRET` | 管理后台令牌签名 | 是 |
| `ADMIN_DEFAULT_PASSWORD` | 首次创建管理员的初始密码 | 首次启动必填 |
| `GIN_MODE` | `debug` 或 `release` | 否 |
| `COS_SECRET_ID` / `COS_SECRET_KEY` | 腾讯云访问凭据 | 使用 COS 时必填 |
| `COS_REGION` / `COS_BUCKET` | COS 地域和存储桶 | 使用 COS 时必填 |
| `COS_PUBLIC_URL` | 证书图片公开前缀 | 否，默认 `/static/certificates` |

不要提交 `.env`、云密钥、SSH 私钥、证书原图、业务 Excel 或数据库备份。

## 常用接口

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| `GET` | `/health` | 应用与数据库健康检查 | 否 |
| `POST` | `/api/v1/auth/login` | 管理员登录 | 否 |
| `GET` | `/api/v1/awards/suggestions` | 团队名称联想 | 否 |
| `POST` | `/api/v1/awards/query/person` | 查询个人获奖记录 | 否 |
| `POST` | `/api/v1/awards/query/group` | 查询团队获奖记录 | 否 |
| `POST` | `/api/v1/awards/upload/certificate` | 上传单张证书图片 | Bearer Token |
| `POST` | `/api/v1/awards/upload/person` | 新增个人记录 | Bearer Token |
| `POST` | `/api/v1/awards/upload/group` | 新增团队记录 | Bearer Token |
| `PUT` | `/api/v1/awards/update/person` | 更新个人记录 | Bearer Token |
| `PUT` | `/api/v1/awards/update/group` | 更新团队记录 | Bearer Token |
| `DELETE` | `/api/v1/awards/:award_id` | 删除记录 | Bearer Token |

完整路由以 `backend/routes/routes.go` 为准。

## 自检

提交前至少执行：

```bash
cd backend
gofmt -w .
go test ./...
go vet ./...

cd ../frontend
npm ci
npm run build
```

部署后检查：

```bash
curl -fsS https://hkicfest.com/health
```

并人工验证个人查询、团队查询、后台登录、单张证书上传和记录编辑。

## 部署说明

通用 Docker 部署说明见 [`DEPLOYMENT.md`](./DEPLOYMENT.md)。生产环境的主机、容器、挂载目录、备份位置和维护流程只记录在单独交接文档中，避免将基础设施细节和敏感信息放入公开仓库。

## 数据与安全约定

- `certificate_number` 是系统展示和对外交付使用的证书编号；证件号码使用 `id_number`
- 新增记录前先备份数据库，完成后通过公开查询页逐条抽查
- 证书图片先走上传接口，记录中只保存接口返回的路径
- 禁止直接把生产 `.env`、`config.yaml`、`SecretKey.csv`、数据库 SQL 备份或证书原图提交到 Git
- 发现密钥曾被复制到非受控位置时，应立即轮换，而不是只删除文件
