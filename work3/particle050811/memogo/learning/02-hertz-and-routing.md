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

## 2025-11-01 Hertz 路由兼容性问题与解决方案

### Q: 为什么需要为 `:id` 形式的路由添加特殊处理？Hertz 和 Thrift 没有兼容吗？

**问题背景**：
在 `router.go` 中看到需要额外注册 `:id` 格式的路由：
```go
// 状态更新路由别名
todos.PATCH("/:id/status", api.UpdateTodoStatus)

// 单条删除路由别名
todos.DELETE("/:id", api.DeleteOne)
```

**问题分析**：

#### 1. Thrift 生成的路由格式
从 `biz/router/memogo/api/memogo.go` 可以看到 Thrift 生成器使用的是 `{id}` 格式：
```go
_todos.DELETE("/{id}", append(_deleteoneMw(), api.DeleteOne)...)
__7bid_7d.PATCH("/status", append(_updatetodostatusMw(), api.UpdateTodoStatus)...)
```

#### 2. Hertz 路由匹配机制
Hertz 框架的路由匹配存在**优先级问题**：
- **精确匹配** 优先于 **参数化路由**
- 当同时存在 `/v1/todos` 和 `/v1/todos/{id}` 时，`DELETE /v1/todos/123` 可能匹配到错误的路由

#### 3. 参数格式兼容性
- Thrift 生成器使用 `{id}` 格式
- 但某些客户端、环境或历史代码可能使用 `:id` 格式
- 添加 `:id` 格式的别名确保向后兼容

**具体冲突场景**：

```go
// 生成的路由（存在优先级冲突）
_v1.DELETE("/todos", append(_deletebyscopeMw(), api.DeleteByScope)...)     // 精确匹配
_todos.DELETE("/{id}", append(_deleteoneMw(), api.DeleteOne)...)           // 参数化路由

// 请求：DELETE /v1/todos/123
// 可能匹配到：DELETE /v1/todos (按范围删除) 而不是 DELETE /v1/todos/{id} (单条删除)
```

**解决方案**：
通过 `router.go` 中的自定义路由注册兼容两种格式：
```go
// 兼容性路由别名
// 说明：由于 Hertz 框架路由优先级问题，需要额外注册 ":id" 格式的路由
// 确保与 Thrift 生成的 "{id}" 格式路由同时可用
v1 := r.Group("/v1")
todos := v1.Group("/todos", mw.JWTMiddleware.MiddlewareFunc())

// 状态更新路由别名
todos.PATCH("/:id/status", api.UpdateTodoStatus)

// 单条删除路由别名
todos.DELETE("/:id", api.DeleteOne)
```

**为什么 Hertz 和 Thrift 没有完全兼容？**

1. **框架设计差异**：
   - Thrift 专注于接口定义和代码生成
   - Hertz 专注于 HTTP 路由和中间件处理
   - 两者在路由匹配算法上存在差异

2. **参数格式标准化**：
   - 不同 Web 框架使用不同的参数格式（`:id`、`{id}`、`<id>`）
   - Thrift 选择了 `{id}` 格式，但需要确保与各种客户端兼容

3. **路由优先级处理**：
   - Hertz 的路由匹配算法在处理精确匹配和参数化路由时不够智能
   - 需要手动处理潜在的冲突

**最佳实践**：

1. **预防性兼容**：对于关键路由，同时注册两种格式
2. **统一客户端格式**：在文档中明确推荐使用 `{id}` 格式
3. **测试覆盖**：确保两种格式都能正常工作
4. **监控日志**：记录实际使用的参数格式，逐步淘汰不推荐的格式

**结论**：
虽然理想情况下 Thrift 和 Hertz 应该完全兼容，但实际使用中确实需要这些兼容性处理来确保系统的稳定性和向后兼容性。这些特殊路由是解决框架限制的实用方案，体现了防御性编程的思想。

---

## 延伸阅读

- Hertz hz 代码生成与路由绑定: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/
- hz 路由/注解使用: https://www.cloudwego.io/docs/hertz/tutorials/tool/hz/router
- Hertz 路由匹配算法: https://www.cloudwego.io/docs/hertz/tutorials/basic-feature/route/
