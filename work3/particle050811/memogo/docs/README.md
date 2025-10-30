# MemoGo 后端项目文档（简要版）

> 目标：用 Hertz（或 Gin）+ GORM + MySQL 实现一个带用户认证与备忘录（待办）管理的 RESTful API。认证建议使用 `github.com/golang-jwt/jwt/v5`（不要用已知不安全的 `dgrijalva/jwt-go`）。

---

## 一、技术栈与特性

* **Web 框架**：CloudWeGo **Hertz**（也可替换为 Gin）
* **ORM**：GORM v2
* **DB**：MySQL（可选 Redis 加速/限流）
* **认证**：JWT（Access/Refresh 双令牌）
* **接口协议**：RESTful + JSON
* **IDL**：`idl/memogo.thrift`（可用 `hz` 生成骨架与 OpenAPI）
* **统一响应**：`{ status, msg, data }`
* **功能点**：注册、登录、刷新令牌；待办增/查/改/删；分页与搜索；批量状态更新与删除

---

## 二、目录结构（三层架构）

```
.
├─ cmd/
│  └─ server/
│     └─ main.go                 # 程序入口，加载配置/路由/中间件
├─ idl/
│  └─ memogo.thrift              # IDL（通过 hz 生成代码/文档）
├─ internal/
│  ├─ config/                    # 配置加载（env / yaml）
│  ├─ middleware/                # JWT、日志、CORS、Recovery 等
│  ├─ model/                     # GORM 模型（User / Todo）
│  ├─ repo/                      # 数据访问（UserRepo / TodoRepo）
│  ├─ service/                   # 业务逻辑（AuthService / TodoService）
│  ├─ handler/                   # HTTP 适配/参数绑定（若使用 Gin）
│  ├─ pkg/
│  │  ├─ jwtx/                   # JWT 签发/校验/刷新
│  │  ├─ hash/                   # 密码哈希（bcrypt/argon2）
│  │  ├─ resp/                   # 统一返回结构
│  │  └─ pagination/             # 分页工具
│  └─ dal/
│     └─ db.go                   # MySQL/Redis 初始化
├─ scripts/
│  ├─ migrate.sql                # 建表脚本
│  └─ seed.sql                   # 测试数据
├─ docs/
│  ├─ openapi.json               # OpenAPI（可由 hz 或 swag 生成）
│  └─ README.md
├─ configs/config.yaml           # 配置示例
├─ Dockerfile
├─ docker-compose.yml
└─ Makefile
```

---

## 三、数据库表（建议）

**users**

* `id` BIGINT PK, AUTO\_INCREMENT
* `username` VARCHAR(64) UNIQUE NOT NULL
* `password_hash` VARCHAR(255) NOT NULL
* `created_at` BIGINT NOT NULL

**todos**

* `id` BIGINT PK, AUTO\_INCREMENT
* `user_id` BIGINT NOT NULL (FK users.id)
* `title` VARCHAR(128) NOT NULL
* `content` TEXT NOT NULL
* `status` TINYINT NOT NULL DEFAULT 0  (`0=TODO, 1=DONE`)
* `view` INT NOT NULL DEFAULT 0
* `created_at` BIGINT NOT NULL
* `start_time` BIGINT NULL
* `end_time` BIGINT NULL
* `due_time` BIGINT NULL
* 索引建议：`(user_id, status, created_at)`、`(user_id, title)`

> GORM 配置：开启 Prepared Statements；所有 where 条件都带 `user_id` 防止越权。

---

## 四、运行与生成

### 1）环境变量（示例）

```
APP_ADDR=:8080
DB_DSN=user:pass@tcp(127.0.0.1:3306)/memogo?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=replace-with-strong-secret
ACCESS_TOKEN_TTL=900           # 15分钟
REFRESH_TOKEN_TTL=1209600      # 14天
REDIS_ADDR=127.0.0.1:6379
```

### 2）(可选) 使用 `hz` 生成

```bash
go install github.com/cloudwego/hertz/cmd/hz@latest
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
hz new -idl idl/memogo.thrift
# 生成的代码基础上，填充 service/repo/middleware 等实现
```

### 3）启动

```bash
make migrate       # 执行 scripts/migrate.sql
make run           # 或 go run ./cmd/server
# 或: docker compose up -d
```

---

## 五、统一返回结构

```json
{
  "status": 200,          // 参考 HTTP，但允许业务内二次编码
  "msg": "ok",
  "data": { }             // 业务有效载荷（列表/对象/影响条数等）
}
```

**错误示例**

```json
{ "status": 401, "msg": "invalid token", "data": {} }
```

---

## 六、认证与授权

* 注册/登录/刷新令牌无需鉴权
* 其他 **todos** 相关接口需要携带：

  * Header: `Authorization: Bearer <access_token>`
* 使用 `github.com/golang-jwt/jwt/v5` 校验，Access 短期、Refresh 长期
* 所有查询/更新/删除必须 `WHERE user_id = <from JWT claims>`

---

## 七、API 速览与调用示例（cURL）

> 路径以 `/v1` 为前缀；`status` 字段取值：`todo | done | all`；分页 `page` 从 1 开始，`page_size` 建议上限 50。

### 1）用户模块

#### 注册

`POST /v1/auth/register`

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"fanone","password":"P@ssw0rd"}'
```

**响应**

```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "access_expires_in": 900,
    "refresh_expires_in": 1209600
  }
}
```

#### 登录

`POST /v1/auth/login`（同注册参数/返回）

#### 刷新令牌

`POST /v1/auth/refresh`

```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh>"}'
```

---

### 2）待办模块

#### 新增待办

`POST /v1/todos`

```bash
curl -X POST http://localhost:8080/v1/todos \
  -H "Authorization: Bearer <access>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "学习 Hertz",
    "content": "完成 CRUD 与分页",
    "start_time": 1730100000,
    "due_time": 1730300000
  }'
```

**响应**

```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "id": 1,
    "title": "学习 Hertz",
    "content": "完成 CRUD 与分页",
    "view": 0,
    "status": 0,
    "created_at": 1730080000,
    "start_time": 1730100000,
    "end_time": 0,
    "due_time": 1730300000
  }
}
```

#### 将**单条**设置为 **DONE/TODO**

`PATCH /v1/todos/{id}/status`

```bash
curl -X PATCH http://localhost:8080/v1/todos/1/status \
  -H "Authorization: Bearer <access>" \
  -H "Content-Type: application/json" \
  -d '{"status":"DONE"}'
```

**响应**

```json
{ "status": 200, "msg": "ok", "data": 1 }
```

#### 将**所有**从 X 改为 Y（批量）

`PATCH /v1/todos/status?from=TODO&to=DONE`

```bash
curl -X PATCH "http://localhost:8080/v1/todos/status?from=TODO&to=DONE" \
  -H "Authorization: Bearer <access>"
```

**响应**

```json
{ "status": 200, "msg": "ok", "data": 12 }
```

#### 查询（分页 + 状态过滤）

`GET /v1/todos?status=todo|done|all&page=1&page_size=10`

```bash
curl "http://localhost:8080/v1/todos?status=all&page=1&page_size=10" \
  -H "Authorization: Bearer <access>"
```

**响应（示例）**

```json
{
  "status": 200,
  "msg": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "学习 Hertz",
        "content": "完成 CRUD 与分页",
        "view": 3,
        "status": 1,
        "created_at": 1730080000,
        "start_time": 1730100000,
        "end_time": 1730123456,
        "due_time": 1730300000
      }
    ],
    "total": 35
  }
}
```

#### 搜索（分页 + 关键词）

`GET /v1/todos/search?q=关键词&page=1&page_size=10`

```bash
curl "http://localhost:8080/v1/todos/search?q=Hertz&page=1&page_size=10" \
  -H "Authorization: Bearer <access>"
```

#### 删除单条

`DELETE /v1/todos/{id}`

```bash
curl -X DELETE http://localhost:8080/v1/todos/1 \
  -H "Authorization: Bearer <access>"
```

**响应**

```json
{ "status": 200, "msg": "ok", "data": 1 }
```

#### 删除（按范围）

`DELETE /v1/todos?scope=done|todo|all`

```bash
curl -X DELETE "http://localhost:8080/v1/todos?scope=done" \
  -H "Authorization: Bearer <access>"
```

**响应**

```json
{ "status": 200, "msg": "ok", "data": 7 }
```

---

## 八、分页与查询规范

* `page` 从 1 开始；`page_size` 建议 10，最大 50（服务端做上限）
* 返回字段：

  * `data.items`：当前页数据数组
  * `data.total`：**符合条件的总条目数**（不是 `items.length`）

---

## 九、状态码与错误

* 正常请求：**HTTP 200** + `status: 200`
* 客户端错误：400（参数错误）、401（未认证/令牌失效）、403（无权限）
* 服务器错误：500（内部错误）
* `msg` 放可读文案；详细错误仅记录在服务端日志

---

## 十、安全与最佳实践

* **JWT**：短期 Access + 长期 Refresh；支持黑名单（可用 Redis）
* **密码**：使用 bcrypt/argon2 存储；禁止明文
* **越权防护**：所有 Todo 操作均以 `user_id` 作为查询准入条件
* **SQL**：GORM + Prepared Statements；禁止手写拼接
* **限流**：登录与搜索可加入 IP/User 限流（Redis + 令牌桶）
* **CORS**：仅放行可信前端域名

---

## 十一、接口文档与调试

* **Apifox / Postman**：配置环境变量（`base_url`、`access_token`）与集合测试
* **自动文档**：

  * 使用 **hz** 根据 `memogo.thrift` 生成路由及 OpenAPI，再导入 Swagger/Apifox
  * 或使用 `swag` 注解生成 `docs/openapi.json`

---

## 十二、快速联调清单

1. `POST /v1/auth/register` → 保存 Access/Refresh
2. `POST /v1/auth/login`（可复用注册用户）→ 覆盖 Access/Refresh
3. 带 `Authorization: Bearer <access>` 调用 `/v1/todos` 系列
4. Access 过期后：`POST /v1/auth/refresh` 换新 Access/Refresh
5. 验证多用户隔离：不同账号只能看到/修改自己的 todos


