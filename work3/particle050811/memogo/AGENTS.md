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

## 技术栈

- **Web 框架**: CloudWeGo Hertz（高性能 HTTP 框架）
- **IDL**: Apache Thrift 用于服务定义和代码生成
- **认证**: JWT 双令牌机制（Access Token + Refresh Token）
- **数据库**: 设计使用 MySQL + GORM（待实现）
- **架构**: 清晰的分层架构（handler/service/repository）

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

## 项目结构

```
.
├── idl/memogo.thrift           # Thrift IDL 服务定义文件
├── main.go                     # 应用入口
├── router.go                   # 自定义路由注册
├── router_gen.go               # 生成的路由注册（勿手动编辑）
├── build.sh                    # 构建脚本
├── script/bootstrap.sh         # 运行时启动脚本
├── biz/
│   ├── handler/                # HTTP 处理器
│   │   ├── ping.go            # 健康检查端点
│   │   └── memogo/api/        # 生成的服务处理器
│   ├── router/                # 路由注册
│   │   ├── register.go        # 生成的路由注册
│   │   └── memogo/api/        # API 路由和中间件
│   └── model/                 # 数据模型（生成）
└── docs/README.md             # 详细的 API 文档
```

## 核心架构要点

### 基于 Thrift 的代码生成
@idl/memogo.thrift
- 服务接口和 HTTP 路由在 `idl/memogo.thrift` 中定义
- `hz` 工具生成 handler、model 和路由代码
- **重要**: 带有 "Code generated" 注释的文件不应手动编辑
- 业务逻辑实现在 handler 文件中完成

### 服务 API 端点

#### 认证模块
- `POST /v1/auth/register` - 用户注册
- `POST /v1/auth/login` - 用户登录
- `POST /v1/auth/refresh` - 刷新令牌

#### 待办事项模块
- `POST /v1/todos` - 创建待办
- `PATCH /v1/todos/{id}/status` - 更新单条状态
- `PATCH /v1/todos/status` - 批量更新状态（通过 from/to 参数）
- `GET /v1/todos` - 列表查询（维持状态过滤和分页）
- `GET /v1/todos/search` - 关键词搜索（支持分页）
- `DELETE /v1/todos/{id}` - 删除单条
- `DELETE /v1/todos` - 按范围批量删除（通过 scope 参数）

### 认证机制
- 使用 Bearer Token 认证：`Authorization: Bearer <token>`
- 实现访问令牌（短期）和刷新令牌（长期）双令牌模式
- 除认证端点外，所有待办操作都需要认证
- 所有操作必须带 `user_id` 过滤，防止越权访问

### 统一响应格式
```json
{
  "status": <HTTP状态码>,
  "msg": "<消息>",
  "data": <响应数据>
}
```

### Hertz 框架特性
- 使用 `server.Default()` 创建默认服务器实例
- 路由注册通过 `register(h)` 函数
- 处理器签名：`func(ctx context.Context, c *app.RequestContext)`
- 中间件在 `biz/router/memogo/api/middleware.go` 中配置

## 开发注意事项

### 代码生成工作流
1. 修改 `idl/memogo.thrift` 定义
2. 运行 `hz update -idl idl/memogo.thrift`
3. 实现生成的 handler 桩代码中的业务逻辑
4. 不要修改任何标记为"Code generated"的文件

### 数据库层（待实现）
- 需要实现 GORM 模型（User、Todo）
- 需要实现 Repository 层进行数据访问
- 所有查询必须包含 `user_id` 条件防止越权
- 使用 Prepared Statements 防止 SQL 注入

### 安全实践
- 密码使用 bcrypt/argon2 哈希存储，禁止明文
- JWT 密钥必须通过环境变量配置，不要硬编码
- 实现令牌黑名单机制（可用 Redis）
- 对登录和搜索接口实施限流
- 配置 CORS 仅允许可信域名

### 分页规范
- `page` 从 1 开始
- `page_size` 建议默认 10，最大限制 50
- 返回数据包含 `items`（当前页数据）和 `total`（符合条件的总数）

## 测试和调试

详细的 API 调用示例和 cURL 命令请参考 `docs/README.md`。

推荐使用 Postman 或 Apifox 进行接口测试，配置环境变量：
- `base_url`: http://localhost:8080
- `access_token`: 从登录/注册获取
- `refresh_token`: 从登录/注册获取
