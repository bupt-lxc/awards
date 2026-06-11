# 获奖信息管理系统 (Awards Management System)

## 项目概述
基于 Go + Vue 3 的获奖证书信息管理系统，支持个人/团体获奖记录的上传、查询、证书管理。

## 技术栈
- **后端**: Go 1.26.4, Gin v1.9.1, GORM v1.25.5, MySQL 8.0
- **前端**: Vue 3.5.x, Vite 6.x, Ant Design Vue 4.2.x
- **部署**: Docker Compose (MySQL 8.0 + Go 后端 + Nginx)

## 运行环境 (Windows 11 本地)
| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.26.4 | `C:\Program Files\Go\bin\go.exe`，需手动加入 PATH |
| Node.js | v24.15.0 | npm 11.12.1 |
| Docker | 29.4.2 | MySQL 容器 `awards-mysql`，端口 3306 |
| 后端二进制 | `backend/awards-management.exe` | 已编译，可直接运行 |

## 数据库
- 地址: `127.0.0.1:3306`
- 数据库: `awards_management`
- 用户: `awards_user` / `awards_pass123`
- 4 张表: `admins`, `person_awards`, `group_awards`, `group_members`
- Docker 容器: `awards-mysql` (docker-compose 启动)

## 默认管理员
- 用户名: `admin`
- 密码: `admin123`

## 启动命令
```bash
# 确保 Docker MySQL 运行中: awards-mysql 容器
# 启动后端
cd backend && ./awards-management.exe

# 启动前端 (需在 backend 目录外开新终端)
cd frontend && npx vite --host 0.0.0.0
```

## 服务地址
- 后端 API: http://localhost:8000
- 前端公共查询: http://localhost:3000
- 后台管理: http://localhost:3000/admin.html
- 证书验证: http://localhost:3000/verify.html

## 关键修复 (2026-06-11)

### 1. Vite 代理配置 (`frontend/vite.config.js`)
前端开发服务器添加了 proxy 规则，`/api` 和 `/static` 转发到 `localhost:8000`，否则开发模式下前后端通不了。

### 2. 公共查询分流转发 (`frontend/src/components/PublicQueryView.vue`)
原来不管选"个人查询"还是"团队成员查询"都调 `queryPersonAwards`，改为根据 `query_type` 分流调 `queryPersonAwards` 或 `queryGroupAwards`。

### 3. 查询条件全部改为可选
所有必填校验去掉，年份、姓名、证件号、团队名都改为可选，可以只用部分条件查询。

### 4. 团队查询放宽限制 (`backend/controllers/award_controller.go` + `backend/services/group_award_service.go`)
原来要求姓名和证件号必须配对使用，改为每个字段可独立使用。

### 5. 环境配置适配 (`backend/config.yaml`)
数据库 host 从 `1Panel-mysql-lAXR` 改为 `127.0.0.1`，username/password 改为 awards_user/awards_pass123，JWT secret 填入开发用值，mode 改为 debug。

## 证书图片
- 存储路径: `backend/uploads/certificates/`
- 已有图片: `111.jpg`, `123456.jpg`, `1234567.jpg`, `234567.jpg`
- 源路径: `E:\BUPT_PS\P_0\uploads\uploads\certificates\`
- URL 路由: `/static/certificates/<filename>`

## 当前数据状态
- admins: 1 条 (admin)
- group_awards: 3 条 (团队1, 想吃汉堡, 想吃汉堡包)
- group_members: 3 条
- person_awards: 0 条

## 注意事项
- 当前 shell 需手动 `export PATH="$PATH:/c/Program Files/Go/bin"` 才能用 go 命令（Go 安装时 shell 未加载新 PATH）
- GitHub (`https://github.com/Lonlybutfull/awards`) 通过 WebFetch/curl 均无法访问，可能被企业网络拦截
- 项目本身无种子数据脚本，数据通过后台管理页面录入
- `package-lock.json` 变动为 npm 版本差异，不影响功能
