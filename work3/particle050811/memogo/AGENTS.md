# AGENTS.md

此文件为 Claude Code 和 Codex 提供在此仓库中工作的指导。

## 基本规范

### 默认用中文写注释和回答用户问题

### 学习笔记整理规则
- **触发条件**：仅当用户明确要求"整理笔记"或"写入笔记"时才执行
- **文件位置**：`LEARNING.md`（系统化学习笔记）
- **笔记内容**：
  - 记录用户提出的问题和对应的解答
  - 包含代码示例和实际场景
  - 提供延伸阅读链接
  - 使用清晰的 Markdown 格式
- **注意事项**：
  - 不要主动整理笔记，必须用户明确要求
  - 整理时要简洁清晰，突出重点
  - 包含日期标记方便后续查阅

## 项目概述

MemoGo 是一个基于 Go 的 RESTful API 服务，用于待办事项（Todo/Memo）管理和用户认证。项目使用 CloudWeGo Hertz 框架，并通过 Thrift IDL 进行代码生成。

## 技术栈（与当前实现一致）

- **Web 框架**：CloudWeGo Hertz
- **IDL/生成**：Apache Thrift + `hz` 代码生成
- **认证**：JWT（访问令牌 15 分钟、刷新令牌 7 天）
- **数据库**：GORM + SQLite（文件 `memogo.db`，已实现 AutoMigrate）
- **架构**：分层（handler / service / repository）

## 开发命令

### 代码生成
```bash
# 更新 Thrift 定义并重新生成代码
hz update -idl idl/memogo.thrift

# 如果未安装 hz 工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# 注意：go.mod 中需要替换 thrift 版本
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
```

### 构建和运行
```bash
# 构建服务
go build

# 运行服务
./memogo

# 直接运行（不构建）
go run main.go
```

### 依赖管理
```bash
# 安装/更新依赖
go mod tidy

# 下载依赖
go mod download
```

## 项目结构（当前）

```
.
├── idl/memogo.thrift            # 服务与路由定义（含 HTTP 注解）
├── main.go                      # 入口：初始化 DB 与 JWT 中间件
├── router.go                    # 自定义/兼容性路由（含 /v1/todos/:id/status 别名）
├── router_gen.go                # 生成的总路由注册（勿改）
├── biz/
│   ├── dal/
│   │   ├── db/init.go          # GORM + SQLite 初始化与迁移
│   │   ├── model/              # User、Todo 模型
│   │   └── repository/         # UserRepository、TodoRepository
│   ├── service/                # AuthService、TodoService
│   ├── handler/                # ping 与业务处理器
│   └── router/                 # 生成的路由与中间件绑定
├── pkg/                         # hash、jwt、middleware 等
└── docs/README.md              # API 文档（本仓库）
```

## 核心架构要点

### 基于 Thrift 的代码生成（已用）
@idl/memogo.thrift
- 服务接口与 HTTP 路由由 `idl/memogo.thrift` 注解定义
- 通过 `hz update -idl idl/memogo.thrift` 生成路由与处理桩
- 带有“Code generated”注释的文件请勿手动编辑
- 实际业务在 `biz/service` 与 `biz/handler` 中实现

### 服务 API 端点（当前生效）

#### 健康检查
- `GET /ping`

#### 认证
- `POST /v1/auth/register`（返回自定义 token 对：仅含 access_token、refresh_token）
- `POST /v1/auth/login`（由中间件返回：access_token、refresh_token=同 access、expires_at）
- `POST /v1/auth/refresh`（由中间件返回：同登录结构）

#### 待办
- `POST /v1/todos` 新建
- `PATCH /v1/todos/{id}/status` 更新单条状态（兼容别名：`/v1/todos/:id/status`）
- `PATCH /v1/todos/status?from=0&to=1` 批量迁移状态（支持 `TODO/DONE` 或 `0/1`）
- `GET /v1/todos?status=todo|done|all&page=1&page_size=10` 分页列表
- `GET /v1/todos/search?q=kw&page=1&page_size=10` 关键词搜索
- `DELETE /v1/todos/{id}` 删除单条
- `DELETE /v1/todos?scope=done|todo|all` 按范围删除

### 认证机制
- 使用 Bearer Token：`Authorization: Bearer <token>`
- 访问令牌 15 分钟、刷新令牌 7 天（见 `pkg/jwt` 与 `pkg/middleware/jwt.go`）
- 除注册/登录/刷新外，所有接口均需认证
- Repository 层所有查询均包含 `user_id` 过滤

### 统一响应格式
```json
{ "status": 200, "msg": "ok", "data": {} }
```

### Hertz 框架特性
- 使用 `server.Default()` 创建默认服务器实例
- 路由注册通过 `register(h)` 函数
- 处理器签名：`func(ctx context.Context, c *app.RequestContext)`
- 中间件在 `biz/router/memogo/api/middleware.go` 中配置

## 开发注意事项

### 代码生成工作流
1. 修改 `idl/memogo.thrift`
2. 执行 `hz update -idl idl/memogo.thrift`
3. 在 `biz/service` 与 `biz/handler` 填充/调整逻辑
4. 生成文件（含“Code generated”）不要手改

### 数据库层（已实现）
- 使用 SQLite 本地文件 `memogo.db`
- `AutoMigrate` 已对 `users` 与 `todos` 表生效
- 所有查询包含 `user_id` 条件防止越权

### 安全实践
- 密码使用 bcrypt 哈希，禁止明文（见 `pkg/hash`）
- JWT 密钥通过 `JWT_SECRET` 配置（默认弱密钥仅用于本地）
- 可扩展：黑名单、限流、CORS 等（当前未内置）

### 分页规范
- `page` 从 1 开始
- `page_size` 默认 10，最大 50（服务层限制）
- 响应包含 `items` 与 `total`

## 测试和调试

- 运行：`go run main.go`
- 健康检查：`GET /ping` → `{ "message": "pong" }`
- 调试示例与 cURL：见 `docs/README.md`
- 建议在 Postman/Apifox 配置：
  - `base_url = http://localhost:8080`
  - `access_token`、`refresh_token`（从注册/登录获取）
