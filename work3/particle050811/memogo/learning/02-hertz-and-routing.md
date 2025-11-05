# Hertz 框架与路由学习笔记

本文档记录 Hertz 框架、Thrift、路由相关的学习内容。

---

## 2025-10-31 路由中间件生成规则

### Q: "_deletebyscopeMw() 是不是根据这个事务模块生成的？"

> // 事务模块：删（按 scope 删除：done/todo/all）
> DeleteResp DeleteByScope(1: DeleteByScopeReq req)
>   (api.delete = "/v1/todos")

**结论**：是的，`_deletebyscopeMw()` 由 hz 代码生成器依据 IDL 中的方法 `DeleteByScope` 及其 HTTP 注解 `(api.delete = "/v1/todos")` 生成，用于给该具体路由挂载中间件。

**证据与映射**：
- IDL 方法与注解（idl/memogo.thrift）
  ```thrift
  DeleteResp DeleteByScope(1: DeleteByScopeReq req)
    (api.delete = "/v1/todos")
  ```
- 生成的路由注册（biz/router/memogo/api/memogo.go:22）
  ```go
  _v1.DELETE("/todos", append(_deletebyscopeMw(), api.DeleteByScope)...) // 调用中间件装配函数
  ```
- 生成的中间件函数（biz/router/memogo/api/middleware.go:21）
  ```go
  func _deletebyscopeMw() []app.HandlerFunc {
      // 需要 JWT 认证
      return []app.HandlerFunc{middleware.JWTMiddleware.MiddlewareFunc()}
  }
  ```

**命名规则速记**：
- 路由级：`_` + `方法名小写` + `Mw` → `DeleteByScope` → `_deletebyscopeMw`
- 分组级：路径片段生成 → `/v1` → `_v1Mw()`；`/v1/todos` → `_todosMw()`；重复分组追加序号 → `_todos0Mw()`
- 路径参数：`/{id}` 采用占位编码 → `__7bid_7dMw()`（避免非法标识符）

**实际场景**：
- 需要对"按范围删除"接口加鉴权：在 `_deletebyscopeMw()` 返回 JWT 中间件即可；
- 若整个 `/v1/todos` 组都需鉴权：在 `_todosMw()` 返回 JWT，则其子路由（含 DeleteByScope）自动继承，无需在 `_deletebyscopeMw()` 重复配置。

**延伸阅读**：
- Hertz hz 代码生成与路由绑定: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/
- hz 路由/注解使用: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/router

---

## 2025-11-05 Hertz 路由参数语法：`:id` vs `{id}` 的真相

### Q: 为什么 IDL 中的 `{id}` 导致路由不匹配？

**核心问题**：
Hertz 框架**只支持 `:id` 语法**，不支持 `{id}` 语法。

**问题表现**：
```thrift
// IDL 中使用 {id} 格式
UpdateTodoStatusResp UpdateTodoStatus(1: UpdateTodoStatusReq req)
  (api.patch = "/v1/todos/{id}/status")
```

当 hz 工具生成路由时：
```go
// 生成的路由（错误）
_todos.PATCH("/{id}/status", api.UpdateTodoStatus)
```

**根本原因**：
- Hertz 路由引擎不把 `{id}` 识别为路径参数
- 它会把 `{id}` 当作**普通字符串**处理
- 所以 `/v1/todos/123/status` 无法匹配 `/v1/todos/{id}/status`

**验证方法**：
如果在 `router.go` 中手动添加 `:id` 格式的别名路由，请求立即生效：
```go
// 手动添加别名后才能匹配
todos.PATCH("/:id/status", api.UpdateTodoStatus)
```

这证明 Hertz 只认 `:id` 不认 `{id}`。

---

### 正确的解决方案

#### 方案一：修改 IDL（推荐）

在 IDL 中直接使用 `:id` 语法：
```thrift
// idl/memogo.thrift
UpdateTodoStatusResp UpdateTodoStatus(1: UpdateTodoStatusReq req)
  (api.patch = "/v1/todos/:id/status")

DeleteResp DeleteOne(1: DeleteOneReq req)
  (api.delete = "/v1/todos/:id")
```

重新生成代码：
```bash
hz update -idl idl/memogo.thrift
```

生成的路由会自动使用正确的语法：
```go
// biz/router/memogo/api/memogo.go
_todos.DELETE("/:id", append(_deleteoneMw(), api.DeleteOne)...)
_id.PATCH("/status", append(_updatetodostatusMw(), api.UpdateTodoStatus)...)
```

此时无需在 `router.go` 中添加任何别名路由。

#### 方案二：手动别名（不推荐，临时方案）

如果无法修改 IDL，需要在 `router.go` 中为每个 `{id}` 路由手动添加 `:id` 别名：
```go
// router.go (需要手动维护)
v1 := r.Group("/v1")
todos := v1.Group("/todos", mw.JWTMiddleware.MiddlewareFunc())

// 为每个 {id} 路由添加 :id 别名
todos.PATCH("/:id/status", api.UpdateTodoStatus)
todos.DELETE("/:id", api.DeleteOne)
```

**缺点**：
- 手动维护，容易遗漏
- IDL 与实际路由不一致
- 新增路由时需同步更新两处

---

### 为什么会有这个混淆？

#### 1. 不同框架使用不同语法
| 框架 | 路径参数语法 | 示例 |
|------|-------------|------|
| Hertz | `:param` | `/users/:id` |
| Gin | `:param` | `/users/:id` |
| Echo | `:param` | `/users/:id` |
| OpenAPI/Swagger | `{param}` | `/users/{id}` |
| Spring Boot | `{param}` | `/users/{id}` |

#### 2. Thrift IDL 的默认示例可能使用 `{id}`
CloudWeGo 文档中某些示例可能混用了两种格式，导致误解。

#### 3. hz 工具不会自动转换
hz 工具会**原样**将 IDL 中的路径注解写入生成的代码，不会做 `{id}` → `:id` 的转换。

---

### 最佳实践总结

| 场景 | 建议 |
|------|------|
| **新项目** | IDL 中统一使用 `:id` 语法 |
| **已有项目** | 修改 IDL 并重新生成，删除手动别名 |
| **临时修复** | 在 `router.go` 中添加别名，但标记为技术债 |
| **文档编写** | 明确说明 Hertz 只支持 `:id` |

---

### 相关文件引用

- IDL 定义：`idl/memogo.thrift:178` (UpdateTodoStatus)
- IDL 定义：`idl/memogo.thrift:196` (DeleteOne)
- 生成的路由：`biz/router/memogo/api/memogo.go:26, 34-35`
- 自定义路由：`router.go:13-18`

---

### 实际修复记录（2025-11-05）

**修改前**：
```thrift
// 错误：使用了 {id}
(api.patch = "/v1/todos/{id}/status")
```

**修改后**：
```thrift
// 正确：使用 :id
(api.patch = "/v1/todos/:id/status")
```

**结果**：
- ✅ 路由可以正常匹配 `/v1/todos/123/status`
- ✅ 删除了 `router.go` 中的手动别名
- ✅ 代码更简洁，维护更方便

---

## 延伸阅读

- Hertz hz 代码生成与路由绑定: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/
- hz 路由/注解使用: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/router
- Hertz 路由匹配算法: https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/
