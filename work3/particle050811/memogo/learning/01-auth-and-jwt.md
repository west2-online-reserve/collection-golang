# 认证与 JWT 学习笔记

本文档记录用户认证、JWT 令牌相关的学习内容。

---

## 2025-10-30 注册功能开发

### Q1: 注册时为什么要检查"用户不存在"的错误？

**问题代码**：
```go
_, err = s.userRepo.GetByUsername(username)
if err == nil {
    return "", "", repository.ErrUserAlreadyExists
}
if !errors.Is(err, repository.ErrUserNotFound) {
    return "", "", err  // ← 为什么这里要报错？
}
```

**解答**：
这是**防御性编程**,区分三种情况：
1. **`err == nil`** - 找到用户 → 用户已存在，返回错误 ❌
2. **`err == ErrUserNotFound`** - 用户不存在 → 正常，继续注册 ✅
3. **`err == 其他错误`** - 数据库故障等系统错误 → 需要报错 ❌

关键点：只有"用户不存在"是预期的正常情况，其他错误（如数据库连接断开）都要抛出。

---

### Q2: Thrift 定义的结构体在哪里？为什么小写变大写？

**Thrift 定义**（`idl/memogo.thrift`）：
```thrift
struct RegisterReq {
  1: string username
  2: string password
}
```

**生成的 Go 代码**（`biz/model/memogo/api/memogo.go:567`）：
```go
type RegisterReq struct {
    Username string `thrift:"username,1" json:"username"`
    Password string `thrift:"password,2" json:"password"`
}
```

**命名规则**：
- **Thrift 字段**：小写 `username`
- **Go 字段**：大写 `Username`（必须 public 才能序列化）
- **JSON 字段**：小写 `username`（通过 tag 控制）

**为什么不能定义私有字段？**
- Thrift 用于跨语言数据交换，所有字段必须可序列化
- 如需私有字段，在 Go 业务模型中手动添加（如 `biz/dal/model/user.go`）

**最佳实践**：
```
API 层（Thrift 生成） ↔️ 业务层（手动定义）
api.RegisterReq      →  model.User (可包含私有字段)
```

---

### Q3: 登录时为什么"用户不存在"和"密码错误"用同一个错误？

**问题代码**：
```go
// 查找用户
user, err := s.userRepo.GetByUsername(username)
if err != nil {
    if errors.Is(err, repository.ErrUserNotFound) {
        return "", "", ErrInvalidCredentials  // ← 返回通用错误
    }
    return "", "", err
}

// 验证密码
if err := hash.VerifyPassword(user.PasswordHash, password); err != nil {
    return "", "", ErrInvalidCredentials  // ← 同样的错误
}
```

**安全原因：防止用户枚举攻击（User Enumeration Attack）**

**攻击场景**：
```bash
# 如果区分错误（不安全）
POST /login {"username":"admin","password":"123"}
→ "密码错误"  # 攻击者知道 admin 存在！

POST /login {"username":"randomuser","password":"123"}
→ "用户不存在"  # 攻击者知道这个用户名没注册

# 如果统一错误（安全）
POST /login {"username":"admin","password":"123"}
→ "用户名或密码错误"

POST /login {"username":"randomuser","password":"123"}
→ "用户名或密码错误"  # 攻击者无法判断用户是否存在
```

**注意事项**：
1. **响应时间**也要一致（bcrypt 验证需要时间，用户不存在时也要执行假验证）
2. **注册接口**必须明确返回"用户名已存在"（但需配合验证码、频率限制）

---

### Q4: 中文网站为什么返回英文错误？业界通用方案是什么？

**业界主流方案对比**：

#### ⭐ 方案1：错误码 + 前端国际化（90% 公司采用）

**后端返回**：
```json
{
  "status": 400,
  "code": "AUTH_INVALID_CREDENTIALS",
  "msg": "Invalid credentials",
  "data": null
}
```

**前端处理**：
```javascript
const i18n = {
  'zh-CN': { 'AUTH_INVALID_CREDENTIALS': '用户名或密码错误' },
  'en-US': { 'AUTH_INVALID_CREDENTIALS': 'Invalid credentials' }
}
```

**优点**：
- 前端灵活控制显示语言
- 支持多语言切换无需重新请求
- 前端可自定义错误文案
- 减少后端复杂度

**缺点**：
- 需要前后端约定错误码
- 前端需维护翻译文件

---

#### 方案2：后端国际化（根据 Accept-Language）

**客户端请求**：
```http
Accept-Language: zh-CN
```

**后端返回**：
```json
{
  "status": 400,
  "msg": "用户名或密码错误"  // 根据语言返回
}
```

**优点**：
- 前端无需维护翻译
- 适合 SEO（服务端渲染）

**缺点**：
- 后端需维护多语言文件
- 切换语言需重新请求
- 增加后端复杂度

---

**真实案例**：
- **微信开放平台**：`{"errcode": 40001, "errmsg": "invalid credential"}`
- **阿里云 API**：`{"Code": "InvalidParameter", "Message": "..."}`
- **GitHub API**：`{"message": "Bad credentials"}`

**建议**：
- **学习项目**：直接返回中文，简单直接
- **生产项目**：错误码 + 前端国际化（业界标准）

---

## 2025-10-30 JWT 认证与 Hertz 中间件

### Q1: LoginHandler 为何一行能完成登录？它封装了什么？

**问题代码**：
```go
func Login(ctx context.Context, c *app.RequestContext) {
    middleware.JWTMiddleware.LoginHandler(ctx, c)  // ← 只有一行！
}
```

**解答**：`LoginHandler` 不是我们写的函数，是 **Hertz JWT 中间件提供的**。我们做的是配置中间件，告诉它如何处理登录。

**配置的 4 个关键回调函数**（`pkg/middleware/jwt.go`）：

| 函数 | 作用 | 类比 |
|------|------|------|
| **Authenticator** | 验证用户名密码 | 你告诉快递员："这样验证收件人" |
| **PayloadFunc** | 构建 JWT 内容 | 你告诉快递员："包裹里放这些信息" |
| **LoginResponse** | 返回成功响应 | 你告诉快递员："成功时这样回复" |
| **Unauthorized** | 返回错误响应 | 你告诉快递员："失败时这样回复" |

**LoginHandler 的执行流程**：
```
LoginHandler 内部自动执行:
┌────────────────────────────────────┐
│ 1. 调用 Authenticator (我们配置的) │
│    - 解析用户名密码                 │
│    - 验证凭证                       │
│    - 返回用户信息                   │
├────────────────────────────────────┤
│ 2. 调用 PayloadFunc (我们配置的)   │
│    - 将用户信息转为 JWT Claims     │
├────────────────────────────────────┤
│ 3. 生成 JWT Token (中间件自动)     │
│    - 使用 Key 签名                  │
│    - 设置过期时间                   │
├────────────────────────────────────┤
│ 4. 调用 LoginResponse (我们配置的) │
│    - 返回 JSON 响应                 │
└────────────────────────────────────┘
```

**核心思想**：这是**配置驱动**的设计模式，你配置规则，框架执行流程。

---

### Q2: Token 哪里来的？为何暂时使用相同的 token？

**Token 是 Hertz JWT 中间件自动生成的！**

看 `LoginResponse` 的函数签名：
```go
LoginResponse: func(ctx context.Context, c *app.RequestContext,
                     code int, token string, expire time.Time)
                              ↑
                        中间件传给我们的！
```

**标准的双令牌机制应该是**：

| Token 类型 | 有效期 | 用途 |
|-----------|--------|------|
| **access_token** | 短期（15分钟） | API 调用 |
| **refresh_token** | 长期（7天） | 刷新 access_token |

**问题**：中间件只生成一个 token（根据 `Timeout: 15 * time.Minute`），所以暂时返回了两次。

**解决方案**：

```go
// 方案 A：单令牌模式（推荐简单场景）
LoginResponse: func(..., token string, expire time.Time) {
    c.JSON(200, utils.H{
        "access_token": token,  // 只返回一个
        "expires_at": expire.Unix(),
    })
}

// 方案 B：真双令牌（手动生成 refresh_token）
LoginResponse: func(..., token string, expire time.Time) {
    claims := hertzJWT.ExtractClaims(ctx, c)
    userID := uint(claims["user_id"].(float64))

    // 使用自己的 jwt 包生成长期 token
    _, refreshToken, _ := jwt.GenerateTokenPair(userID, username)

    c.JSON(200, utils.H{
        "access_token": token,        // 15分钟
        "refresh_token": refreshToken, // 7天
    })
}
```

---

### Q3: 为何中间件不默认支持双 token？

**4 个核心原因**：

1. **JWT 标准没有定义双令牌**
   - RFC 7519 只定义了如何创建/验证 token
   - 双令牌是**安全最佳实践**，不是标准

2. **双令牌机制有多种实现方式**
   - 不同应用需求完全不同（过期时间、存储方式、撤销机制）
   - 强制一种实现会限制灵活性

3. **Refresh Token 通常需要持久化**
   - Access Token: 短期、无状态、不需要存储 ✓
   - Refresh Token: 长期、需要可撤销、必须存储 ✗
   - 中间件无法假设你用什么数据库！

4. **框架设计哲学：提供机制，不强制策略**
   - ✅ 提供生成/验证/刷新 token 的能力
   - ❌ 不强制如何使用这些能力

**Hertz JWT 的变相双令牌方案**：

```go
Timeout:    15 * time.Minute,     // Token 15分钟后过期
MaxRefresh: 7 * 24 * time.Hour,   // 但在 7 天内可刷新
```

**工作原理**：
```
生成的 token 包含两个时间戳:
{
  "exp": 1234567890,      // 过期时间（15分钟后）
  "orig_iat": 1234567000  // 原始签发时间
}

刷新流程:
1. Token 15分钟后过期
2. 但在 7 天内，可以用这个"过期"的 token 换新 token
3. 超过 7 天，必须重新登录
```

**结论**：对于待办应用，Hertz 的方案已经够用！

---

### Q4: Hertz token 是存到内存里的吗？服务端重启会退出登录吗？

**关键结论：JWT Token 不存储在服务端！服务端重启不影响登录！**

**JWT 的无状态特性**：

```
登录流程:
┌─────────┐                    ┌─────────┐
│  客户端  │ ─(用户名密码)────→ │ 服务端   │
│         │                    │         │
│         │ ←──(JWT token)──── │ 生成token│
│ 存储到   │                    │ 不保存！ │
│ 本地     │                    └─────────┘
└─────────┘

后续请求:
┌─────────┐                    ┌─────────┐
│  客户端  │ ─(请求+token)────→ │ 服务端   │
│         │                    │         │
│ 从本地   │                    │ 1.验证签名│
│ 取token  │                    │ 2.检查过期│
│         │                    │ 3.提取信息│
│         │ ←───(响应)──────── │ 不查数据库│
└─────────┘                    └─────────┘
```

**对比不同存储方式**：

| 存储方式 | Token 存哪里？ | 服务端重启影响？ |
|---------|--------------|----------------|
| **JWT** | 客户端（浏览器/App） | ❌ 不影响 |
| **Session** | 服务端内存/Redis | ✅ 影响 |

**为什么服务端不需要存储？**

JWT token 包含了所有需要的信息：
```json
{
  "user_id": 2,
  "username": "logintest",
  "exp": 1761837205,  // 过期时间
  "iat": 1761836305   // 签发时间
}
```

服务端验证时只需要：
1. 验证签名（token 没被篡改）
2. 检查过期时间（exp < now）
3. 提取用户信息（user_id）

**不需要查数据库！不需要查 Redis！**

---

### Q5: JWT 是如何防篡改的？是用 RSA 加密吗？

**关键概念：JWT 不是加密，是签名！**

**JWT 的三部分结构**（用 `.` 分隔）：

```
eyJhbGc...       .  eyJ1c2Vy...    .  aePaXs5iQ...
   ↓                    ↓                  ↓
 Header            Payload             Signature
（头部）            （载荷）             （签名）
```

**1. Header 和 Payload - Base64 编码（不加密）**

```json
// Header
{"alg":"HS256","typ":"JWT"}

// Payload
{"user_id":2,"username":"logintest","exp":1761837205}
```

⚠️ **任何人都能解码！不要在 token 里放敏感信息！**

**2. Signature - 这是防篡改的关键**

我们使用的是 **HS256**（对称加密）：

```javascript
signature = HMAC-SHA256(
  base64(header) + "." + base64(payload),
  secret_key  // 服务端保密的密钥
)
```

**两种签名算法对比**：

| 算法 | 类型 | 密钥 | 适用场景 |
|------|------|------|---------|
| **HS256** | 对称加密 | 一个密钥 | 单体应用 ✅ |
| **RS256** | 非对称加密（RSA） | 公钥+私钥 | 微服务架构 |

**HMAC-SHA256 算法原理**：

```
输入: Header.Payload (明文)
密钥: memogo-default-secret...
         ↓
    HMAC-SHA256 算法
         ↓
    32字节的签名

特性:
✓ 单向：无法从签名推导出密钥
✓ 确定性：相同输入+密钥 = 相同签名
✓ 雪崩效应：输入改1位，签名完全不同
```

**防篡改验证流程**：

```
客户端请求:
Header: Authorization: Bearer eyJhbGc...

服务端验证:
┌────────────────────────────────────┐
│ 1. 分割 Token                      │
│    header, payload, signature      │
├────────────────────────────────────┤
│ 2. 重新计算签名                    │
│    new_sig = HMAC(header.payload,  │
│                   secret_key)      │
├────────────────────────────────────┤
│ 3. 对比签名                        │
│    if new_sig != signature:        │
│        return "Invalid Token" ✗    │
├────────────────────────────────────┤
│ 4. 检查过期时间                    │
│    if exp < now():                 │
│        return "Token Expired" ✗    │
├────────────────────────────────────┤
│ 5. 提取用户信息 ✓                  │
└────────────────────────────────────┘
```

**为什么无法篡改？**

| 攻击方式 | 能成功吗？ | 原因 |
|---------|-----------|------|
| 修改 user_id | ❌ | 签名不匹配 |
| 修改过期时间 | ❌ | 签名不匹配 |
| 伪造签名 | ❌ | 不知道密钥，暴力破解需要几十亿年 |
| 重放旧 token | ✅ | **需要额外防护**（黑名单） |
| 中间人攻击 | ✅ | **必须使用 HTTPS** |

**安全最佳实践**：

```go
// 1. 密钥必须保密
func GetJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")  // ✓ 从环境变量
    if secret == "" {
        // ✗ 生产环境必须改！
        secret = "memogo-default-secret-change-in-production"
    }
    return []byte(secret)
}

// 2. Payload 不放敏感信息
// ✗ 错误：{"password": "123456"}
// ✓ 正确：{"user_id": 2, "username": "alice"}

// 3. 必须使用 HTTPS
// 防止中间人截获 token

// 4. 设置合理的过期时间
Timeout: 15 * time.Minute  // 不要设置太长
```

**类比总结**：

JWT 签名就像快递的防伪封条：
- 包裹内容（Payload）：任何人都能看
- 防伪封条（Signature）：只有邮局（服务端）能验证真伪
- 如果有人打开包裹修改内容，封条会破损（签名不匹配），邮局会拒收

**JWT 安全性依赖于：密钥保密 + HTTPS + 合理的过期时间**

---

## 2025-10-31 JWT vs Cookie 认证选择

### Q: 为什么项目里使用 JWT（签名，Authorization 头携带），而很多网站用 Cookie 区分用户？

**结论**：本项目是面向多端的无状态 REST API，更适合使用 JWT；传统主要面向浏览器的站点更适合 Cookie + 服务端会话。

**为什么本项目用 JWT**
- 无状态、易水平扩展：不依赖服务端会话存储或粘滞会话。
- 多端/跨域友好：移动端、Postman/Apifox、SPA 都可直接用 `Authorization: Bearer <token>`。
- 降低 CSRF 风险：令牌不随浏览器自动携带，默认不受第三方站点跨站请求影响（仍需防 XSS）。
- 微服务友好：下游服务可独立校验签名，无需回源查会话。

**为什么很多网站用 Cookie**
- 浏览器原生支持：自动携带，配合 `HttpOnly/SameSite/Secure` 易控管。
- 撤销/风控强：服务端集中失效会话即可（踢人、权限变更即时生效）。
- SEO/SSR/后台系统：以浏览器为主的产品形态更贴合会话模型。

**对比要点**
- 存储位置：JWT 在客户端（Header/Storage）；Cookie 会话在服务端（内存/Redis）+ 客户端保存会话 ID。
- 安全关注：JWT 关注泄露与刷新机制；Cookie 关注 CSRF（配合 SameSite/CSRF Token）。
- 撤销与权限变更：JWT 需黑名单或短期+刷新；Cookie 服务器集中失效即可。
- 跨域：JWT（Header）更直接；Cookie 需处理 `CORS` 与 `SameSite=None; Secure`。

**你提到"用 username 区分用户"**
- 建议使用不可变的唯一标识 `user_id` 作为主身份声明；`username` 仅用于展示或冗余。
- 切记不要把敏感信息放入 JWT；全程使用 HTTPS 传输。

**代码/请求示例**
```bash
# 使用 JWT（推荐本项目）
curl -H "Authorization: Bearer <ACCESS_TOKEN>" \
     "http://localhost:8080/v1/todos?page=1&page_size=10"

# 使用 Cookie 会话（典型网站）
curl -H "Cookie: sid=abcdef123456; Path=/; HttpOnly; Secure" \
     "http://example.com/dashboard"
```

**实践建议**
- 当前项目继续用 JWT：统一从 `Authorization` 头读取并在中间件校验；令牌建议短期+刷新，配合黑名单（如 Redis）。
- 主要面向浏览器时，可采用"JWT 装进 HttpOnly Cookie"的混合方案，降低 XSS 窃取风险，同时保持无状态验证。
- 如需切换到 Cookie 会话，我可以协助改造：新增会话存储、CSRF 防护、SameSite 策略与登录/登出流程。

---

## 延伸阅读

- OWASP 用户枚举攻击防护：https://owasp.org/www-project-web-security-testing-guide/
- Go 错误处理最佳实践：https://go.dev/blog/error-handling-and-go
- RESTful API 设计规范：https://restfulapi.net/
- JWT 官方规范 RFC 7519: https://datatracker.ietf.org/doc/html/rfc7519
- JWT.io 在线调试工具: https://jwt.io/
- Hertz JWT 中间件文档: https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/middleware/jwt/
- OWASP JWT 安全最佳实践: https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html
- HMAC-SHA256 算法详解: https://en.wikipedia.org/wiki/HMAC
- JWT vs Session 深度对比：https://auth0.com/blog/session-vs-token-based-authentication/
- OWASP CSRF 防护：https://owasp.org/www-community/attacks/csrf
- CORS 与 Cookie 的 SameSite：https://web.dev/articles/samesite-cookies-explained
