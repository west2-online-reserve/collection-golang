# MemoGo 后端文档（与源码一致版）

本项目基于 CloudWeGo Hertz + GORM 构建，当前实现以 SQLite 为默认数据库，采用 Thrift IDL 驱动的代码生成，提供注册/登录/刷新令牌与 Todo 的增删改查、分页、搜索、批量状态更新。

提示：本文档已根据当前源码行为更新，覆盖运行方式、实际路由、鉴权与响应结构等，替换了早期“规划阶段”中关于 MySQL 的描述。

## 技术栈与关键点
- Web 框架：Hertz
- IDL/代码生成：Thrift + `hz`
- 数据库：GORM + SQLite（文件 `memogo.db` 自动创建和迁移）
- 认证：JWT（访问令牌 15 分钟、刷新令牌 7 天）
- 统一响应：`{ status, msg, data }`

## 目录结构（当前）
```
.
├── idl/memogo.thrift           # Thrift IDL 服务定义（路由注解来源）
├── main.go                     # 程序入口，初始化 DB 与 JWT 中间件
├── router.go                   # 自定义与兼容性路由（含 /v1/todos/:id/status 别名）
├── router_gen.go               # 生成的总路由注册（勿改）
├── biz/
│   ├── dal/
│   │   ├── db/init.go         # GORM 初始化（SQLite）与 AutoMigrate
│   │   ├── model/             # User、Todo GORM 模型
│   │   └── repository/        # UserRepository、TodoRepository
│   ├── service/               # AuthService、TodoService（业务规则与分页上限）
│   ├── handler/               # /ping 与 MemoGoService 的实现
│   └── router/                # hz 生成的路由及中间件绑定
├── pkg/
│   ├── hash/                  # bcrypt 封装
│   ├── jwt/                   # Token 生成与解析、密钥读取
│   └── middleware/            # Hertz JWT 中间件封装
└── docs/README.md             # 本文档
```

## 运行与环境
- 直接运行：`go run main.go`
- 二进制：`go build && ./memogo`
- 可选环境变量：
  - `JWT_SECRET`（默认有内置弱密钥，仅用于本地调试，生产务必覆盖）

代码生成相关：
```bash
go install github.com/cloudwego/hertz/cmd/hz@latest
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
hz update -idl idl/memogo.thrift
```
带有“Code generated”注释的文件请勿手改。

## 鉴权与安全
- 除注册/登录/刷新外，所有 Todo 接口需要 `Authorization: Bearer <token>`
- 访问令牌默认 15 分钟、刷新令牌默认 7 天
- 服务端所有 Todo 读写均带 `user_id` 过滤，隔离不同用户数据
- 密码使用 bcrypt 存储（见 `pkg/hash`）

注意：当前登录与刷新接口返回结构由 Hertz JWT 中间件统一输出，`data.refresh_token` 字段与 `access_token` 相同（占位），并包含 `expires_at`；注册接口返回自定义的 token 对，不包含过期时间字段。

## 统一响应结构
```json
{ "status": 200, "msg": "ok", "data": { } }
```
错误示例：`{ "status": 401, "msg": "Unauthorized: ...", "data": null }`

## 路由与示例
基础健康检查：
- `GET /ping` → `{ "message": "pong" }`

认证模块：
- `POST /v1/auth/register`
  - 请求：`{ "username": "u", "password": "p" }`
  - 响应：`{ "status":200, "msg":"Registration successful", "data": { "access_token":"...", "refresh_token":"..." } }`
- `POST /v1/auth/login`
  - 请求：`{ "username": "u", "password": "p" }`
  - 响应：`{ "status":200, "msg":"Login successful", "data": { "access_token":"...", "refresh_token":"...", "expires_at": 1730... } }`
- `POST /v1/auth/refresh`
  - 请求：`{ "refresh_token": "..." }`
  - 响应：同登录结构（由中间件输出）

Todo 模块（均需 Bearer Token）：
- `POST /v1/todos` 新建
  - 请求：`{ "title":"t", "content":"c", "start_time":1730..., "due_time":1730... }`
  - 响应：`{ "status":200, "msg":"ok", "data": { "id":1, "status":0, ... } }`
- `PATCH /v1/todos/{id}/status` 更新单条状态
  - 兼容别名：`PATCH /v1/todos/:id/status`
  - 请求：`{ "status": 1 }`（0=TODO, 1=DONE）
  - 响应：`{ "status":200, "msg":"ok", "data": 1 }`（受影响条数）
- `PATCH /v1/todos/status?from=0&to=1` 批量状态迁移
  - 也支持枚举字符串（如 `from=TODO&to=DONE`），推荐数值
  - 响应：`{ "status":200, "msg":"ok", "data": <affected> }`
- `GET /v1/todos?status=todo|done|all&page=1&page_size=10` 分页列表
  - 响应：`{ "status":200, "msg":"ok", "data": { "items":[], "total": 0 } }`
- `GET /v1/todos/search?q=kw&page=1&page_size=10` 关键词搜索
- `DELETE /v1/todos/{id}` 删除单条
- `DELETE /v1/todos?scope=done|todo|all` 范围删除

分页规范：`page` 从 1 开始；默认 `page_size=10`，最大 50（由服务层限制）。

## 重要实现细节对照
- 数据库使用 SQLite，初始化与迁移见：`biz/dal/db/init.go`
- JWT 配置与登录响应见：`pkg/middleware/jwt.go`
- Token 生成/解析与默认密钥见：`pkg/jwt/jwt.go`
- 路由注册（生成）：`biz/router/memogo/api/memogo.go`
- 自定义别名路由：`router.go`

## 常见问题
- 登录/刷新返回的 `refresh_token` 与 `access_token` 相同：当前中间件使用统一的 `LoginResponse`，此字段为占位输出，不影响受保护接口验证。
- 更新单条状态的请求体枚举：建议使用数值 `0/1`；如使用字符串 `TODO/DONE`，实际行为依绑定器版本可能不同。
- **路由兼容性问题**：由于 Hertz 框架路由优先级限制，单条删除和状态更新接口同时支持 `{id}` 和 `:id` 两种参数格式，确保与各种客户端兼容。

## 调试建议
- 先 `POST /v1/auth/register` 或 `/v1/auth/login` 获取 token
- 所有 Todo 请求带 `Authorization: Bearer <access_token>`
- 通过 `GET /ping` 快速验证服务存活
