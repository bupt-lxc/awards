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

### 6. 前端依赖补全 (`frontend/package.json`)
项目缺少 `vue-router` 依赖（`PublicQueryView.vue` 中 `import { useRouter } from 'vue-router'` 会失败），已通过 `npm install vue-router --legacy-peer-deps` 安装。版本冲突原因：vue-router@5.x 要求 vite@^7/^8，项目用 vite@6.x。

### 7. 个人查询后端彻底改为可选 (`backend/dto/dto.go` + `backend/services/person_award_service.go`)
`PersonAwardQuery` 的 `person_name` / `id_number` 残留 `binding:"required"` 标签，且 service 层 `QueryPersonAwards` 硬编码 `WHERE person_name = ? AND id_number = ?` 强制两个条件同时存在。修复：
- DTO 去掉 `required` 并补上 `omitempty`
- Service 改为各自独立的可选 `if` 条件

### 8. 证书验证页团体查询修复 (`frontend/src/VerifyApp.vue`)
`getUrlParams()` 参数校验对所有 query type 都要求 `person_name` + `id_number`，但团体验证链接（由 `GroupQueryView` 生成）只传 `type/group_name/year/category/award_id`，直接报"缺少必要参数"。
- 拆分为个人/团体两套校验规则
- 个人：必须 `person_name, id_number, year, category`
- 团体：必须 `group_name, year, category`

### 9. 验证页匹配逻辑修复 (`frontend/src/VerifyApp.vue`)
`targetAward` 初始化为 `null`，只有传了 `award_id` 且精确匹配时才赋值；没传 `award_id` 时 API 返回数据也被丢弃报"未找到匹配的获奖记录"。改为默认取 API 返回的第一条结果。

### 10. 公共查询结果无 award_type 标记 (`frontend/src/components/PublicQueryView.vue` + `frontend/src/components/GroupQueryView.vue`)
API 返回数据不含 `award_type` 字段，但模板用 `record.award_type === 'group'` 判断渲染逻辑。团队查询时 `award_type` 为 `undefined`，走到个人分支，但团队数据无 `person_name`，结果卡片内容为空。已为两个组件的查询结果注入 `award_type` 标记。

### 11. GroupQueryView 参赛编码取消必填 (`frontend/src/components/GroupQueryView.vue`)
团队查询表单中"参赛编码"（id_number）原来设了 `required: true`，不填无法通过表单校验。已改为可选。

### 12. PersonQueryView 团队查询两处 bug (`frontend/src/components/PersonQueryView.vue`)
这是首页 (`/`) 上渲染的组件，支持个人/团队切换查询。之前团队查询有两个 bug：
- **参数错误**：`queryParams.person_name = queryForm.group_name` 把团队名称塞到了 `person_name` 字段，后端团队查询接口收不到 `group_name`
- **API 调用错误**：始终调用 `queryPersonAwards`，未根据 `query_type` 分流到 `queryGroupAwards`
修复：团队查询时正确传 `group_name` 字段并调用 `queryGroupAwards`。

### 13. PersonQueryView 团队结果展示修复 (`frontend/src/components/PersonQueryView.vue`)
模板卡片中标题、姓名/团队名称行、证书编码行、下载文件名均硬编码了 `award.person_name`。团队查询返回的数据是 `group_name`，导致卡片不显示团队名。改为根据 `query_type` 动态取值：团体时取 `award.group_name`。

### 14. GroupQueryView AutoComplete 输入丢失修复 (`frontend/src/components/GroupQueryView.vue`)
`<a-auto-complete>` 的 `v-model:value` 仅在用户从下拉框点击选中时才更新，直接输入后点查询则 `queryForm.group_name` 为空。修复：在 `handleTeamSearch` 中实时同步 `queryForm.group_name = value`，并将用户输入值加入选项数组防止被丢弃。

## 前端组件与页面关系
| 文件 | 用途 | 路由 |
|------|------|------|
| `src/PublicApp.vue` | 公共查询主页，路由分发 | `/` (index.html) |
| `src/App.vue` | 管理后台主页 | `/admin.html` |
| `src/components/PersonQueryView.vue` | 首页查询组件，支持个人/团队切换 | 通过 PublicApp 加载 |
| `src/components/GroupQueryView.vue` | 团队获奖查询（独立路由版） | `/group-query` |
| `src/components/PublicQueryView.vue` | 公共查询（另一版本） | - |
| `src/VerifyApp.vue` | 证书验证页 | `/verify.html` |

## 团队查询"参赛编码"语义
在团队查询上下文中，`id_number` / "参赛编码"实际上是**团队成员（指挥）的身份证号/编码**，后端通过子查询匹配 `group_members.member_id_number`。个人查询中则是获奖者本人的证件号。

## 证书图片
- 存储路径: `backend/uploads/certificates/`
- URL 路由: `/static/certificates/<filename>`

## 注意事项
- 当前 shell 需手动 `export PATH="$PATH:/c/Program Files/Go/bin"` 才能用 go 命令（Go 安装时 shell 未加载新 PATH）
- GitHub (`https://github.com/Lonlybutfull/awards`) 通过 WebFetch/curl 均无法访问，可能被企业网络拦截
- 项目本身无种子数据脚本，数据通过后台管理页面录入
- `package-lock.json` 变动为 npm 版本差异，不影响功能
- 每次启动需检查端口占用：旧进程可能导致端口冲突（8000 后端 / 3000 前端），需手动 kill 后再启
