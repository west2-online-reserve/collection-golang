# MemoGo 学习笔记

本文件是学习笔记的导航索引。详细笔记按主题分类存放在 `learning/` 文件夹中。

---

## 📂 笔记结构

```
learning/
├── 01-auth-and-jwt.md      # 认证与 JWT 相关
├── 02-hertz-and-routing.md # Hertz 框架与路由
└── 03-redis-cache.md        # Redis 缓存相关
```

---

## 📚 笔记目录

### 01. 认证与 JWT [`learning/01-auth-and-jwt.md`]

**主题内容**：
- 用户注册与登录
- 防御性编程与错误处理
- 用户枚举攻击防护
- Thrift 结构体生成规则
- JWT 认证机制
- Hertz JWT 中间件配置
- 双令牌机制
- JWT 防篡改原理（HMAC-SHA256）
- JWT vs Cookie 对比

**关键问题**：
- Q: 为什么登录时"用户不存在"和"密码错误"用同一个错误？
- Q: LoginHandler 为何一行能完成登录？
- Q: JWT Token 是存到内存里的吗？
- Q: JWT 是如何防篡改的？

---

### 02. Hertz 框架与路由 [`learning/02-hertz-and-routing.md`]

**主题内容**：
- Hertz 中间件生成规则
- Thrift IDL 与路由映射
- 路由优先级与兼容性问题
- 参数格式处理（`:id` vs `{id}`）

**关键问题**：
- Q: `_deletebyscopeMw()` 是如何生成的？
- Q: 为什么需要为 `:id` 形式的路由添加特殊处理？

---

### 03. Redis 缓存 [`learning/03-redis-cache.md`]

**主题内容**：
- Cache-Aside 模式实现
- TodoRepository 缓存策略
- UserRepository 缓存策略
- 缓存键命名规范
- TTL 过期时间设置
- Redis vs Go Map 对比
- 多实例共享缓存
- 自动降级机制

**关键问题**：
- Q: 如何实现"先尝试 Redis"？
- Q: Redis 相当于 map 吗？
- Q: 为什么写操作要删除所有用户缓存？
- Q: 什么时候用 Map？什么时候用 Redis？

---

## 🎯 快速查找

### 按主题查找

| 主题 | 文件 |
|-----|------|
| 用户认证 | `learning/01-auth-and-jwt.md` |
| JWT 令牌 | `learning/01-auth-and-jwt.md` |
| 安全实践 | `learning/01-auth-and-jwt.md` |
| 路由配置 | `learning/02-hertz-and-routing.md` |
| 中间件 | `learning/02-hertz-and-routing.md` |
| 缓存实现 | `learning/03-redis-cache.md` |
| 性能优化 | `learning/03-redis-cache.md` |

### 按日期查找

| 日期 | 主题 | 文件 |
|-----|------|------|
| 2025-10-30 | 注册功能开发 | `learning/01-auth-and-jwt.md` |
| 2025-10-30 | JWT 认证与中间件 | `learning/01-auth-and-jwt.md` |
| 2025-10-31 | 路由中间件生成规则 | `learning/02-hertz-and-routing.md` |
| 2025-10-31 | JWT vs Cookie | `learning/01-auth-and-jwt.md` |
| 2025-11-01 | Hertz 路由兼容性 | `learning/02-hertz-and-routing.md` |
| 2025-11-03 | Redis 缓存实现 | `learning/03-redis-cache.md` |

---

## 🔖 常用链接

### 官方文档
- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Redis 官方文档](https://redis.io/docs/)
- [JWT 规范 RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519)
- [Go 官方博客](https://go.dev/blog/)

### 安全相关
- [OWASP Web 安全测试指南](https://owasp.org/www-project-web-security-testing-guide/)
- [OWASP JWT 安全实践](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)
- [OWASP CSRF 防护](https://owasp.org/www-community/attacks/csrf)

### 架构模式
- [Cache-Aside Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside)
- [RESTful API 设计规范](https://restfulapi.net/)

---

## 📝 笔记规范

1. **文件命名**：使用数字前缀 + 主题，如 `01-auth-and-jwt.md`
2. **内容组织**：按日期和问题组织，每个问题包含解答和代码示例
3. **代码引用**：包含文件路径和行号，如 `biz/service/auth_service.go:92`
4. **图表说明**：使用 ASCII 图表和表格说明流程
5. **延伸阅读**：每个主题末尾提供相关链接

---

*本笔记持续更新中...*
*最后更新：2025-11-03*
