# 获奖信息管理系统后端 (Go版本)

基于 Go、Gin、Viper、GORM 和 MySQL 实现的获奖信息管理系统后端API。

## 技术栈

- **Go 1.21+** - 编程语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL** - 数据库
- **Viper** - 配置管理
- **JWT** - 身份认证
- **CORS** - 跨域支持

## 功能特性

- 🔐 JWT身份认证和授权
- 👥 管理员账户管理
- 🏆 个人获奖记录管理
- 🏆 团体获奖记录管理
- 🔍 获奖信息查询（支持多条件查询）
- 📤 批量上传获奖记录
- 📝 获奖记录更新和删除
- 📋 团体成员证书管理
- 🌐 CORS跨域支持
- 📊 健康检查接口

## 项目结构

```
backend-go/
├── auth/                 # 认证中间件
├── config/              # 配置管理
├── controllers/         # 控制器层
├── database/           # 数据库连接
├── dto/                # 数据传输对象
├── models/             # 数据模型
├── routes/             # 路由配置
├── services/           # 业务逻辑层
├── config.yaml         # 配置文件
├── go.mod              # Go模块文件
├── main.go             # 主程序入口
├── Makefile            # 构建脚本
└── README.md           # 项目文档
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Git

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd backend-go
   ```

2. **安装依赖**
   ```bash
   go mod download
   go mod tidy
   ```

3. **配置数据库**
   
   创建MySQL数据库：
   ```sql
   CREATE DATABASE awards_management CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

4. **配置文件**
   
   修改 `config.yaml` 中的数据库连接信息：
   ```yaml
   database:
     host: "localhost"
     port: 3306
     username: "your_username"
     password: "your_password"
     dbname: "awards_management"
   ```

5. **运行应用**
   ```bash
   go run main.go
   ```
   
   或使用Makefile：
   ```bash
   make run
   ```

### 使用Docker（可选）

1. **构建镜像**
   ```bash
   make docker-build
   ```

2. **运行容器**
   ```bash
   make docker-run
   ```

## API文档

### 基础信息

- **基础URL**: `http://localhost:8000/api/v1`
- **认证方式**: Bearer Token (JWT)
- **内容类型**: `application/json`

### 主要接口

#### 认证接口

- `POST /auth/login` - 管理员登录
- `GET /auth/me` - 获取当前用户信息
- `POST /auth/create-admin` - 创建新管理员（需要管理员权限）

#### 查询接口（公开）

- `POST /awards/query/person` - 查询个人获奖信息
- `POST /awards/query/group` - 查询团体获奖信息
- `GET /awards/{award_id}` - 获取获奖记录详情

#### 管理接口（需要认证）

- `POST /awards/upload/person` - 上传个人获奖记录
- `POST /awards/upload/group` - 上传团体获奖记录
- `POST /awards/batch-upload/person` - 批量上传个人获奖记录
- `POST /awards/batch-upload/group` - 批量上传团体获奖记录
- `PUT /awards/update/person` - 更新个人获奖记录
- `PUT /awards/update/group` - 更新团体获奖记录
- `DELETE /awards/{award_id}` - 删除获奖记录
- `GET /awards/manage/all` - 获取所有获奖记录
- `GET /awards/group/{award_id}/members` - 获取团体成员信息
- `PUT /awards/group/member/certificate` - 更新成员证书

#### 系统接口

- `GET /` - API信息
- `GET /health` - 健康检查

## 配置说明

### config.yaml

```yaml
server:
  host: "0.0.0.0"          # 服务器地址
  port: 8000               # 服务器端口
  mode: "debug"            # 运行模式: debug, release, test

database:
  host: "localhost"        # 数据库地址
  port: 3306              # 数据库端口
  username: "root"         # 数据库用户名
  password: "password"     # 数据库密码
  dbname: "awards_management" # 数据库名
  charset: "utf8mb4"       # 字符集
  parseTime: true         # 解析时间
  loc: "Local"            # 时区
  maxIdleConns: 10        # 最大空闲连接数
  maxOpenConns: 100       # 最大打开连接数
  connMaxLifetime: 3600   # 连接最大生存时间（秒）

jwt:
  secret: "your-secret-key" # JWT密钥
  expireHours: 8          # 令牌过期时间（小时）

cors:
  allowOrigins: ["*"]     # 允许的源
  allowMethods:           # 允许的方法
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowHeaders: ["*"]     # 允许的头部
  allowCredentials: true  # 允许凭据
```

## 开发工具

### Makefile命令

```bash
make build        # 构建二进制文件
make clean        # 清理构建文件
make test         # 运行测试
make deps         # 下载依赖
make run          # 运行应用
make dev          # 开发模式运行（需要air）
make fmt          # 格式化代码
make vet          # 代码检查
make lint         # 代码规范检查
make help         # 显示帮助信息
```

### 开发模式

安装air进行热重载开发：
```bash
make install-air
make dev
```

## 数据库表结构

### admins - 管理员表
- `id` - 主键
- `username` - 用户名（唯一）
- `password_hash` - 密码哈希
- `role` - 角色
- `created_at` - 创建时间
- `updated_at` - 更新时间

### person_awards - 个人获奖表
- `id` - 主键
- `year` - 年份
- `person_name` - 姓名
- `id_number` - 证件号码
- `category` - 奖项类别
- `award_level` - 奖项等级
- `description` - 描述
- `certificate_image` - 证书图片
- `created_at` - 创建时间
- `updated_at` - 更新时间

### group_awards - 团体获奖表
- `id` - 主键
- `year` - 年份
- `group_name` - 团体名称
- `category` - 奖项类别
- `award_level` - 奖项等级
- `description` - 描述
- `members` - 成员信息（JSON格式）
- `created_at` - 创建时间
- `updated_at` - 更新时间

### group_members - 团体成员表
- `id` - 主键
- `group_award_id` - 团体获奖ID（外键）
- `member_name` - 成员姓名
- `member_id_number` - 成员证件号
- `certificate_image` - 证书图片

## 默认账户

系统启动时会自动创建默认管理员账户：
- **用户名**: admin
- **密码**: admin123

## 部署说明

### 生产环境配置

1. 修改配置文件中的敏感信息
2. 设置运行模式为 `release`
3. 配置生产数据库
4. 更改JWT密钥

### 构建生产版本

```bash
make build
```

### 跨平台编译

```bash
make build-linux    # Linux版本
make build-windows  # Windows版本
make build-all      # 所有平台
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！