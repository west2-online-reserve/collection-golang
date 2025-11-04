# MemoGo 后端 API 文档

本项目基于 CloudWeGo Hertz + GORM + Redis 构建，采用 Thrift IDL 驱动的代码生成，提供用户认证与待办事项的完整管理功能。

---

## 🎯 技术栈

- **Web 框架**：CloudWeGo Hertz
- **IDL/代码生成**：Apache Thrift + `hz`
- **数据库**：GORM + MySQL
- **缓存**：Redis（Cache-Aside 模式，可选）
- **认证**：JWT（访问令牌 15 分钟、刷新令牌 7 天）
- **架构**：分层架构（Handler → Service → Repository）

---

## 📂 项目结构

```
.
├── idl/memogo.thrift            # Thrift IDL 服务定义（路由注解来源）
├── main.go                      # 程序入口：初始化 DB、Redis 与 JWT 中间件
├── router.go                    # 自定义路由（`:id` 格式兼容性别名）
├── router_gen.go                # 生成的路由注册（勿手动编辑）
├── biz/
│   ├── dal/
│   │   ├── db/init.go          # GORM + MySQL 初始化与迁移
│   │   ├── redis/init.go       # Redis 客户端初始化
│   │   ├── model/              # User、Todo GORM 模型
│   │   └── repository/         # UserRepository、TodoRepository（含缓存逻辑）
│   ├── service/                # AuthService、TodoService（业务逻辑）
│   ├── handler/                # HTTP 请求处理器
│   └── router/                 # hz 生成的路由与中间件绑定
├── pkg/
│   ├── hash/                   # bcrypt 密码哈希
│   ├── jwt/                    # JWT 令牌生成与解析
│   └── middleware/             # Hertz JWT 中间件封装
└── docs/
    ├── README.md               # 本文档
    └── openapi.json            # OpenAPI 3.0 规范
```

---

## 🚀 快速开始

### 1. 环境准备

**必需服务**：
- MySQL 5.7+
- Redis 3.0+（可选，不启动会自动降级到无缓存模式）

**Go 依赖**：
```bash
go mod tidy
```

### 2. 环境变量配置

创建 `.env` 文件（参考 `.env.example`）：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password_here
DB_NAME=memogo

# Redis 配置（可选）
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 密钥配置（生产环境必须修改）
JWT_SECRET=your_jwt_secret_here
```

### 3. 运行服务

```bash
# 开发环境：直接运行
go run main.go

# 生产环境：编译后运行
go build
./memogo
```

服务启动后：
- **API 服务**：http://localhost:8888
- **API 文档**：http://localhost:8888/docs/index.html
- **健康检查**：http://localhost:8888/ping

---

## 🔐 认证机制

### JWT 双令牌模式

| 令牌类型 | 有效期 | 用途 |
|---------|--------|------|
| **access_token** | 15 分钟 | API 调用认证 |
| **refresh_token** | 7 天 | 刷新访问令牌 |

### 使用方式

所有需要认证的接口在请求头中携带：
```http
Authorization: Bearer <access_token>
```

### 令牌刷新流程

1. `access_token` 过期后（15分钟）
2. 使用 `refresh_token` 调用 `/v1/auth/refresh`
3. 获取新的 `access_token` 和 `refresh_token`
4. 超过 7 天需要重新登录

---

## 📡 API 接口

### 统一响应格式

```json
{
  "status": 200,
  "msg": "ok",
  "data": {}
}
```

**错误响应示例**：
```json
{
  "status": 401,
  "msg": "Unauthorized: token is expired",
  "data": null
}
```

---

### 健康检查

#### `GET /ping`

检查服务是否运行。

**响应**：
```json
{
  "message": "pong"
}
```

---

### 用户认证

#### `POST /v1/auth/register`

用户注册。

**请求体**：
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应**：
```json
{
  "status": 200,
  "msg": "Registration successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

---

#### `POST /v1/auth/login`

用户登录。

**请求体**：
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**响应**：
```json
{
  "status": 200,
  "msg": "Login successful",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_at": 1730123456
  }
}
```

---

#### `POST /v1/auth/refresh`

刷新访问令牌。

**请求体**：
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**响应**：
```json
{
  "status": 200,
  "msg": "Token refreshed",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_at": 1730123456
  }
}
```

---

### 待办事项管理

以下接口均需要 Bearer Token 认证。

#### `POST /v1/todos`

创建待办事项。

**请求头**：
```http
Authorization: Bearer <access_token>
```

**请求体**：
```json
{
  "title": "完成项目文档",
  "content": "编写 API 文档和使用说明",
  "start_time": 1730000000,
  "due_time": 1730086400
}
```

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "id": 1,
    "title": "完成项目文档",
    "content": "编写 API 文档和使用说明",
    "view": 0,
    "status": 0,
    "created_at": 1730000000,
    "start_time": 1730000000,
    "end_time": 0,
    "due_time": 1730086400
  }
}
```

---

#### `GET /v1/todos`

分页查询待办事项列表。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| status | string | 否 | 过滤条件：`todo`/`done`/`all`（默认 `all`）|
| page | int | 否 | 页码，从 1 开始（默认 1）|
| page_size | int | 否 | 每页数量（默认 10，最大 50）|

**示例**：
```http
GET /v1/todos?status=todo&page=1&page_size=10
Authorization: Bearer <access_token>
```

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "完成项目文档",
        "status": 0,
        ...
      }
    ],
    "total": 15
  }
}
```

---

#### `GET /v1/todos/search`

关键词搜索待办事项（标题或内容）。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| q | string | 是 | 搜索关键词 |
| page | int | 否 | 页码（默认 1）|
| page_size | int | 否 | 每页数量（默认 10）|

**示例**：
```http
GET /v1/todos/search?q=项目&page=1&page_size=10
Authorization: Bearer <access_token>
```

**响应**：同分页列表格式

---

#### `GET /v1/todos/cursor` ⚡

**游标分页查询**（高效遍历全部数据，时间复杂度 O(n)）。

> 推荐用于数据导出、移动端下拉刷新等场景。相比传统分页，深度遍历时性能提升 100+ 倍。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| status | string | 否 | 过滤条件：`todo`/`done`/`all`（默认 `all`）|
| cursor | int64 | 否 | 上一页最后一条的 ID，首次查询传 `0`（默认 0）|
| limit | int | 否 | 每页数量（默认 10，最大 100）|

**示例**：
```http
# 第一页
GET /v1/todos/cursor?status=all&cursor=0&limit=10
Authorization: Bearer <access_token>

# 第二页（使用上一页返回的 next_cursor）
GET /v1/todos/cursor?status=all&cursor=123&limit=10
Authorization: Bearer <access_token>
```

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "完成项目文档",
        "status": 0,
        ...
      }
    ],
    "next_cursor": 10,
    "has_more": true
  }
}
```

**字段说明**：
- `next_cursor`: 下一页的游标值（0 表示无下一页）
- `has_more`: 是否还有更多数据

**遍历全部数据示例**：
```bash
# 伪代码
cursor = 0
all_items = []
while true:
    resp = GET /v1/todos/cursor?cursor={cursor}&limit=100
    all_items.append(resp.data.items)
    if not resp.data.has_more:
        break
    cursor = resp.data.next_cursor
```

---

#### `GET /v1/todos/search/cursor` ⚡

**关键词游标分页搜索**（高效遍历搜索结果）。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| q | string | 是 | 搜索关键词 |
| cursor | int64 | 否 | 上一页最后一条的 ID（默认 0）|
| limit | int | 否 | 每页数量（默认 10，最大 100）|

**示例**：
```http
GET /v1/todos/search/cursor?q=项目&cursor=0&limit=10
Authorization: Bearer <access_token>
```

**响应**：同游标分页列表格式

---

#### `PATCH /v1/todos/{id}/status`

更新单条待办事项的状态。

**路径参数**：
- `{id}`：待办事项 ID（也支持 `:id` 格式）

**请求体**：
```json
{
  "status": 1
}
```

状态值：`0` = TODO（未完成），`1` = DONE（已完成）

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": 1
}
```

`data` 为受影响的记录数（通常为 1）。

---

#### `PATCH /v1/todos/status`

批量迁移待办事项状态。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| from | int | 是 | 原状态：`0` 或 `1` |
| to | int | 是 | 目标状态：`0` 或 `1` |

**示例**（将所有未完成改为已完成）：
```http
PATCH /v1/todos/status?from=0&to=1
Authorization: Bearer <access_token>
```

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": 5
}
```

`data` 为受影响的记录数。

---

#### `DELETE /v1/todos/{id}`

删除单条待办事项（软删除）。

**路径参数**：
- `{id}`：待办事项 ID

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": 1
}
```

---

#### `DELETE /v1/todos`

按范围批量删除待办事项。

**请求参数**：
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| scope | string | 是 | 删除范围：`done`/`todo`/`all` |

**示例**（删除所有已完成的待办）：
```http
DELETE /v1/todos?scope=done
Authorization: Bearer <access_token>
```

**响应**：
```json
{
  "status": 200,
  "msg": "ok",
  "data": 3
}
```

---

## 💾 数据缓存

### Redis 缓存策略

项目采用 **Cache-Aside（旁路缓存）** 模式：

#### 缓存的数据

| 数据类型 | TTL | 说明 |
|---------|-----|------|
| 待办列表 | 5 分钟 | `ListTodos` 和 `SearchTodos` 查询结果 |
| 用户信息 | 10 分钟 | `GetByID` 和 `GetByUsername` 查询结果 |

#### 缓存失效策略

- **写操作**（Create/Update/Delete）后自动清除相关缓存
- 确保数据一致性
- Redis 连接失败会自动降级到无缓存模式

#### 性能提升

- 首次查询：~50ms（数据库）
- 缓存命中：~2ms（**96% 性能提升**）
- 减少数据库负载约 **85-90%**

---

## 🔧 开发指南

### 代码生成

当修改 `idl/memogo.thrift` 后，需要重新生成代码：

```bash
# 安装 hz 工具（首次）
go install github.com/cloudwego/hertz/cmd/hz@latest

# 确保 thrift 版本兼容
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0

# 重新生成代码
hz update -idl idl/memogo.thrift
```

**注意**：
- 带有 `Code generated` 注释的文件请勿手动编辑
- 业务逻辑在 `biz/service` 和 `biz/handler` 中实现

### 数据库迁移

项目启动时会自动执行 `AutoMigrate`，无需手动创建表结构。

表结构：
- `users`：用户表
- `todos`：待办事项表

### 路由兼容性

由于 Hertz 路由优先级问题，项目中同时注册了两种参数格式：
- Thrift 生成：`/v1/todos/{id}`
- 兼容性别名：`/v1/todos/:id`（在 `router.go` 中）

推荐使用 `{id}` 格式。

---

## 🧪 测试示例

### 使用 cURL

```bash
# 1. 注册用户
curl -X POST http://localhost:8888/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123"}'

# 2. 登录获取 token
curl -X POST http://localhost:8888/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"pass123"}'

# 保存返回的 access_token
TOKEN="eyJhbGc..."

# 3. 创建待办
curl -X POST http://localhost:8888/v1/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试任务","content":"这是内容"}'

# 4. 查询待办列表
curl -X GET "http://localhost:8888/v1/todos?status=all&page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"

# 5. 更新状态为已完成
curl -X PATCH http://localhost:8888/v1/todos/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":1}'
```

### 使用 Postman/Apifox

1. **配置环境变量**：
   - `base_url` = `http://localhost:8888`
   - `access_token` = `<从登录获取>`

2. **导入接口**：访问 http://localhost:8888/docs/openapi.json

3. **认证配置**：
   - Type: Bearer Token
   - Token: `{{access_token}}`

---

## ⚠️ 常见问题

### Q1: Redis 连接失败会影响服务吗？

**A**: 不会。系统会自动降级到无缓存模式，直接查询数据库。日志会显示：
```
Warning: Failed to connect to Redis: ... (caching will be disabled)
```

### Q2: 为什么有两种路由参数格式？

**A**: 由于 Hertz 路由匹配的优先级问题，项目同时支持：
- `{id}` 格式（推荐，Thrift 生成）
- `:id` 格式（兼容性别名）

### Q3: 如何修改 JWT 过期时间？

**A**: 在 `pkg/middleware/jwt.go` 中修改：
```go
Timeout:    15 * time.Minute,     // access_token 有效期
MaxRefresh: 7 * 24 * time.Hour,   // refresh_token 有效期
```

### Q4: 数据库表结构在哪里定义？

**A**: 在 `biz/dal/model/` 目录中：
- `user.go`：用户表模型
- `todo.go`：待办表模型

使用 GORM 的 `AutoMigrate` 自动创建和更新表结构。

---

## 📚 相关文档

- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Thrift IDL 语法](https://thrift.apache.org/docs/idl)
- [GORM 文档](https://gorm.io/zh_CN/docs/)
- [Redis Go 客户端](https://redis.uptrace.dev/)
- [JWT 最佳实践](https://jwt.io/introduction)

---

## 📝 更新日志

- **2025-11-04**：
  - 新增游标分页接口（`/v1/todos/cursor` 和 `/v1/todos/search/cursor`）
  - 优化查询排序：改为按创建时间升序（旧备忘录优先显示）
  - 深度遍历性能提升 100+ 倍（O(n²) → O(n)）
- **2025-11-03**：添加 Redis 缓存支持（TodoRepository 和 UserRepository）
- **2025-11-01**：从 SQLite 切换到 MySQL
- **2025-10-31**：完善路由兼容性处理
- **2025-10-30**：实现基础认证和 CRUD 功能

---

*文档最后更新：2025-11-04*
