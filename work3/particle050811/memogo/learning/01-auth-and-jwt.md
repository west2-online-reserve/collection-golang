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

## 2025-11-05 JWT Token 生成与验证机制

### Q1: AccessToken 和 RefreshToken 有什么区别？代码里如何体现？

**问题背景**：
注册接口返回两个 Token（`biz/handler/memogo/api/memo_go_service.go:57-62`）：
```go
Data: &api.TokenPair{
    AccessToken:      accessToken,
    RefreshToken:     refreshToken,
    AccessExpiresIn:  computeExpiresIn(accessToken),
    RefreshExpiresIn: computeExpiresIn(refreshToken),
}
```

**解答**：两者的区别体现在**有效期**和**使用场景**上。

#### 1. 有效期不同（`pkg/jwt/jwt.go:75-88`）

```go
func GenerateTokenPair(userID uint, username string) (accessToken, refreshToken string, err error) {
    // 访问令牌：15分钟
    accessToken, err = GenerateToken(userID, username, 15*time.Minute)

    // 刷新令牌：7天
    refreshToken, err = GenerateToken(userID, username, 7*24*time.Hour)

    return accessToken, refreshToken, nil
}
```

| Token 类型 | 有效期 | 用途 |
|-----------|--------|------|
| **AccessToken** | 15分钟 | 访问所有受保护的 API（创建待办、查询列表等） |
| **RefreshToken** | 7天 | 仅用于刷新 AccessToken（调用 `/v1/auth/refresh`） |

#### 2. 中间件配置体现（`pkg/middleware/jwt.go:35-36`）

```go
Timeout:     15 * time.Minute,  // Access token 过期时间
MaxRefresh:  7 * 24 * time.Hour, // Refresh token 过期时间
```

#### 3. 使用场景

**AccessToken 使用流程**：
```bash
# 每次 API 调用都需要携带
curl -H "Authorization: Bearer <accessToken>" \
     http://localhost:8080/v1/todos
```

**RefreshToken 使用流程**：
```bash
# 当 AccessToken 过期（15分钟后）时使用
curl -X POST http://localhost:8080/v1/auth/refresh \
     -d '{"refresh_token": "<refreshToken>"}'

# 返回新的 AccessToken
```

#### 4. 双 Token 机制的安全优势

```
用户登录
    ↓
获得 AccessToken (15分钟) + RefreshToken (7天)
    ↓
使用 AccessToken 访问 API
    ↓
AccessToken 过期（15分钟后）
    ↓
使用 RefreshToken 刷新 → 获得新的 AccessToken
    ↓
继续使用新的 AccessToken
    ↓
RefreshToken 也过期（7天后）
    ↓
用户需要重新登录
```

| 安全特性 | AccessToken | RefreshToken |
|---------|------------|--------------|
| **使用频率** | 高（每次 API 请求） | 低（仅刷新时） |
| **被截获风险** | 高 | 低 |
| **泄露损失** | 小（15分钟后失效） | 大（需妥善保管） |
| **存储位置** | 内存 | 安全存储（HttpOnly Cookie） |

**关键**：AccessToken 频繁使用但短期有效，即使泄露影响有限；RefreshToken 使用少但长期有效，降低被截获概率。

---

### Q2: `GenerateTokenPair` 是生成随机 Token 吗？

**问题代码**（`biz/service/auth_service.go:128`）：
```go
accessToken, refreshToken, err := jwtPkg.GenerateTokenPair(user.ID, user.Username)
```

**解答**：不是随机 Token，而是基于 **JWT（JSON Web Token）标准的加密签名令牌**。

#### JWT Token 不是随机的，是可解析的

**生成过程**（`pkg/jwt/jwt.go:37-51`）：

```go
func GenerateToken(userID uint, username string, duration time.Duration) (string, error) {
    now := time.Now()

    // 1. 构建 Payload（包含用户信息）
    claims := Claims{
        UserID:   userID,        // 用户ID
        Username: username,      // 用户名
        OrigIat:  now.Unix(),    // 原始签发时间
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(duration)), // 过期时间
            IssuedAt:  jwt.NewNumericDate(now),               // 签发时间
            NotBefore: jwt.NewNumericDate(now),               // 生效时间
        },
    }

    // 2. 使用 HMAC-SHA256 算法创建 JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // 3. 用密钥签名并返回完整的 JWT 字符串
    return token.SignedString(GetJWTSecret())
}
```

#### JWT Token 的结构（三部分用 `.` 分隔）

```
eyJhbGc...JWT9  .  eyJleHAi...MDB9  .  aePaXs5i...uuaw
     ↓                   ↓                   ↓
   Header            Payload            Signature
  (算法信息)         (用户数据)          (防篡改签名)
```

**解码后的内容**：

```json
// Header
{"alg":"HS256","typ":"JWT"}

// Payload（包含用户信息）
{
  "user_id": 2,
  "username": "logintest",
  "exp": 1761837205,      // 过期时间
  "iat": 1761836305,      // 签发时间
  "orig_iat": 1761836305  // 原始签发时间
}

// Signature
// HMAC_SHA256(Header.Payload, secret)
```

#### 随机 Token vs JWT Token

| 特性 | 随机 Token（如 UUID） | JWT Token |
|------|---------------------|-----------|
| **生成方式** | 纯随机字符串 | 包含用户信息 + 加密签名 |
| **可解析性** | 不可解析 | 可直接解析出用户信息 |
| **验证方式** | 必须查数据库/Redis | 验证签名即可（无需查库） |
| **存储需求** | 服务端必须存储 | 服务端无需存储（无状态） |
| **包含信息** | 无意义的随机值 | userID、username、过期时间等 |

**JWT 的优势**：
- ✅ **自包含**：Token 本身包含用户信息，服务器无需查数据库
- ✅ **无状态**：不需要在服务端存储 session
- ✅ **防篡改**：任何修改 Payload 都会导致签名验证失败

**安全机制**：通过 HMAC-SHA256 签名，只有拥有密钥的服务器才能生成和验证有效的 JWT（详见 Q5）。

---

### Q3: JWT RegisteredClaims 的三个时间字段分别是什么意思？

**问题代码**（`pkg/jwt/jwt.go:43-47`）：
```go
RegisteredClaims: jwt.RegisteredClaims{
    ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
    IssuedAt:  jwt.NewNumericDate(now),
    NotBefore: jwt.NewNumericDate(now),
}
```

**解答**：这三个字段定义了 Token 的**生命周期**。

#### 1. ExpiresAt（过期时间）

**含义**：Token 的**失效截止时间**

```go
ExpiresAt: jwt.NewNumericDate(now.Add(duration))
// AccessToken: now + 15分钟
// RefreshToken: now + 7天
```

**验证逻辑**：
```
当前时间：2025-11-05 10:00:00
ExpiresAt：2025-11-05 10:15:00  （15分钟后）

10:14:59 → Token 有效 ✓
10:15:01 → Token 过期 ✗（返回 ErrTokenExpired）
```

#### 2. IssuedAt（签发时间）

**含义**：Token 的**创建时间**

```go
IssuedAt: jwt.NewNumericDate(now)
```

**用途**：
- 记录 Token 何时生成
- 用于审计日志："用户在 10:00 登录并获得 Token"
- 可实现安全策略："拒绝超过 30 天的 Token（即使未过期）"

#### 3. NotBefore（生效时间）

**含义**：Token 的**最早可用时间**

```go
NotBefore: jwt.NewNumericDate(now)  // 立即生效
```

**验证逻辑**：
```
如果 NotBefore = now（当前代码）：
  → Token 立即生效 ✓

如果 NotBefore = now + 1小时（延迟生效）：
  → 1小时内使用会返回 ErrTokenNotValidYet ✗
```

**使用场景**（设置未来时间）：
- 预约系统："这个优惠券明天才能用"
- 定时任务："这个任务令牌晚上 8 点才生效"
- 防止时钟偏差：设置稍微未来的时间

#### 时间关系图

```
时间线：
  ├─────────┼──────────────────────┼─────────→
  NotBefore  IssuedAt              ExpiresAt
  (生效时间) (签发时间)            (过期时间)

本项目中：
  ├─────────────────────────────────┼─────────→
  now (立即生效)                    now + 15分钟/7天
  NotBefore = IssuedAt              ExpiresAt
```

#### 验证实现（`pkg/jwt/jwt.go:54-72`）

```go
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return GetJWTSecret(), nil
    })

    // JWT 库自动验证三个时间字段：
    // 1. 检查 NotBefore：当前时间 >= NotBefore
    // 2. 检查 ExpiresAt：当前时间 < ExpiresAt
    // 3. 检查 IssuedAt：当前时间 >= IssuedAt（可选）

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken  // ExpiresAt 检查失败
        }
        // 注意：jwt.ErrTokenNotValidYet 也会进入这里
        return nil, ErrInvalidToken
    }

    // token.Valid 已经包含了所有时间字段的验证
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, ErrInvalidToken
}
```

**关键**：`token.Valid` 字段由 JWT 库自动设置，已包含对三个时间字段的验证，无需手动检查。

---

### Q4: NotBefore 是必须设置的吗？

**解答**：**不是必须的**，可以省略。

#### JWT 标准中的定义

根据 [RFC 7519](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1)，所有时间字段都是**可选的**：

| 字段 | 简称 | 是否必填 | 说明 |
|------|-----|----------|------|
| exp | ExpiresAt | 可选 | 过期时间 |
| **nbf** | **NotBefore** | **可选** | 生效时间 |
| iat | IssuedAt | 可选 | 签发时间 |

#### 不设置时的行为

```go
// 方式 1：当前的代码（设置为 now）
claims := Claims{
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
        IssuedAt:  jwt.NewNumericDate(now),
        NotBefore: jwt.NewNumericDate(now),  // 可以省略
    },
}

// 方式 2：不设置 NotBefore（完全可行）
claims := Claims{
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
        IssuedAt:  jwt.NewNumericDate(now),
        // NotBefore 为 nil，验证时会跳过检查
    },
}
```

#### JWT 库的验证逻辑

```go
// JWT 库内部实现（简化版）
func (c RegisteredClaims) Valid() error {
    now := time.Now()

    // 如果 NotBefore 为 nil，跳过检查
    if c.NotBefore != nil && now.Before(c.NotBefore.Time) {
        return ErrTokenNotValidYet
    }

    // 如果 ExpiresAt 不为 nil，才检查过期
    if c.ExpiresAt != nil && now.After(c.ExpiresAt.Time) {
        return ErrTokenExpired
    }

    return nil
}
```

**结论**：不设置 `NotBefore` = Token 立即生效（没有限制）

#### 本项目可以简化

由于 `NotBefore` 设置为 `now`（立即生效），与不设置效果相同：

```go
// 简化前（当前代码）
RegisteredClaims: jwt.RegisteredClaims{
    ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
    IssuedAt:  jwt.NewNumericDate(now),
    NotBefore: jwt.NewNumericDate(now),  // ← 可以删除
}

// 简化后（效果相同）
RegisteredClaims: jwt.RegisteredClaims{
    ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
    IssuedAt:  jwt.NewNumericDate(now),
    // NotBefore 省略，Token 立即生效
}
```

#### 什么时候需要设置 NotBefore？

只有需要**延迟生效**时才设置：

```go
// 场景 1：明天才能用的优惠券
NotBefore: jwt.NewNumericDate(now.Add(24 * time.Hour))

// 场景 2：只在晚上 8-10 点有效的 Token
start := time.Date(2025, 11, 5, 20, 0, 0, 0, time.Local)
NotBefore: jwt.NewNumericDate(start)
ExpiresAt: jwt.NewNumericDate(start.Add(2 * time.Hour))
```

**建议**：
- ✅ 如果不需要延迟生效，可以删除 `NotBefore` 这行代码
- ✅ 保留也可以，代码更明确（"立即生效"的显式声明）
- ✅ 只保留 `ExpiresAt` 和 `IssuedAt` 即可满足大部分场景

---

### Q5: HS256 签名算法是如何实现 JWT 签名的？

**问题代码**（`pkg/jwt/jwt.go:49-51`）：
```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(GetJWTSecret())
```

**解答**：`HS256 = HMAC-SHA256`，通过**对称加密**实现防篡改签名。

#### HS256 算法说明

```
HS256 = HMAC-SHA256
```

- **HMAC**：Hash-based Message Authentication Code（基于哈希的消息认证码）
- **SHA256**：使用 SHA-256 哈希算法
- **对称加密**：签名和验证使用**相同的密钥**

#### JWT 签名的完整流程

```go
// 步骤 1：构建 Header（固定格式）
header := {
    "alg": "HS256",     // 算法
    "typ": "JWT"        // 类型
}
headerBase64 := base64UrlEncode(json.Marshal(header))
// 结果：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9

// 步骤 2：构建 Payload（用户数据）
payload := {
    "user_id": 123,
    "username": "testuser",
    "exp": 1730801700,
    "iat": 1730800800
}
payloadBase64 := base64UrlEncode(json.Marshal(payload))
// 结果：eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoidGVzdHVzZXIi...

// 步骤 3：拼接待签名字符串
message := headerBase64 + "." + payloadBase64
// 结果：eyJhbGci...JWT9.eyJ1c2Vy...MDB9

// 步骤 4：使用 HMAC-SHA256 生成签名
secret := GetJWTSecret()  // "memogo-default-secret-change-in-production"
signature := HMAC_SHA256(message, secret)
signatureBase64 := base64UrlEncode(signature)
// 结果：aePaXs5iQoxZaqYVnHYdksHUaD5Ofwxy09019m3uuaw

// 步骤 5：拼接最终 Token
jwt := headerBase64 + "." + payloadBase64 + "." + signatureBase64
```

#### HMAC-SHA256 算法原理

```go
// HMAC-SHA256 的伪代码实现
func HMAC_SHA256(message string, secret []byte) []byte {
    // 1. 密钥预处理（标准化为 64 字节）
    if len(secret) > 64 {
        secret = SHA256(secret)  // 太长就 hash
    }
    if len(secret) < 64 {
        secret = padWithZeros(secret, 64)  // 太短就填充
    }

    // 2. 创建两个密钥衍生值
    opad := XOR(secret, 0x5c5c5c5c...)  // 外部填充（0x5c 重复 64 次）
    ipad := XOR(secret, 0x36363636...)  // 内部填充（0x36 重复 64 次）

    // 3. 两次哈希（双层安全）
    innerHash := SHA256(ipad + message)        // 内层哈希
    finalHash := SHA256(opad + innerHash)      // 外层哈希

    return finalHash  // 32 字节的签名
}
```

**关键特性**：
- ✅ **单向性**：无法从签名反推密钥
- ✅ **确定性**：相同输入 + 密钥 = 相同签名
- ✅ **雪崩效应**：输入改 1 位，签名完全不同

#### 签名验证流程（`pkg/jwt/jwt.go:54-72`）

```go
func ParseToken(tokenString string) (*Claims, error) {
    // JWT 库内部验证逻辑：

    // 1. 分割 Token
    parts := strings.Split(tokenString, ".")
    // parts[0] = Header (base64)
    // parts[1] = Payload (base64)
    // parts[2] = Signature (base64)

    // 2. 重新计算签名
    message := parts[0] + "." + parts[1]
    expectedSignature := HMAC_SHA256(message, GetJWTSecret())
    actualSignature := base64Decode(parts[2])

    // 3. 对比签名
    if expectedSignature != actualSignature {
        return nil, ErrInvalidToken  // 签名不匹配，Token 被篡改！
    }

    // 4. 解析 Payload
    claims := json.Unmarshal(base64Decode(parts[1]))

    // 5. 验证时间字段
    if claims.ExpiresAt < time.Now().Unix() {
        return nil, ErrExpiredToken
    }

    return claims, nil
}
```

#### 防篡改原理

**场景：攻击者尝试修改 Token**

```
原始 Token：
  Header.Payload.Signature_A

攻击者修改 Payload：
  Header.Payload_Modified.Signature_A

服务器验证：
  重新计算 = HMAC_SHA256(Header.Payload_Modified, secret)
  计算结果 ≠ Signature_A
  → 验证失败 ❌ 拒绝请求
```

| 攻击方式 | 能成功吗？ | 原因 |
|---------|-----------|------|
| 修改 user_id | ❌ | 签名不匹配 |
| 修改过期时间 | ❌ | 签名不匹配 |
| 伪造签名 | ❌ | 不知道密钥，暴力破解需要几十亿年 |
| 重放旧 token | ✅ 需要额外防护 | 黑名单机制（Redis） |
| 中间人攻击 | ✅ 必须使用 HTTPS | 传输层加密 |

#### 签名算法对比

| 算法 | 类型 | 密钥 | 特点 | 适用场景 |
|------|------|------|------|---------|
| **HS256** | 对称 | 单个密钥（签名=验证） | 简单、快速 | 单体服务 ✅ |
| **RS256** | 非对称（RSA） | 公钥+私钥（私钥签名，公钥验证） | 安全、可分发公钥 | 微服务 |
| **ES256** | 非对称（椭圆曲线） | 椭圆曲线密钥对 | 更短的密钥、更高安全性 | 高安全场景 |

**本项目使用 HS256（对称算法）**，适合单体应用。

#### 密钥保护（关键）

```go
func GetJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        // ⚠️ 默认密钥仅用于开发环境
        secret = "memogo-default-secret-change-in-production"
    }
    return []byte(secret)
}
```

**安全最佳实践**：
1. ✅ **生产环境必须修改密钥**（至少 32 字节随机字符串）
2. ✅ **密钥不能泄露**（不写入代码，使用环境变量）
3. ✅ **Payload 不放敏感信息**（任何人都能 Base64 解码）
4. ✅ **必须使用 HTTPS**（防止中间人截获 Token）
5. ✅ **设置合理的过期时间**（15 分钟，不要设置太长）

#### 图解签名过程

```
┌─────────────────────────────────────────────┐
│              生成 JWT Token                  │
└─────────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
    Header       Payload       Secret
   (算法信息)    (用户数据)    (密钥)
        │            │            │
        ▼            ▼            │
  Base64Encode  Base64Encode     │
        │            │            │
        └────────┬───┘            │
                 ▼                │
          "Header.Payload"        │
                 │                │
                 └────────┬───────┘
                          ▼
                    HMAC-SHA256
                          │
                          ▼
                   Base64Encode
                          │
                          ▼
                      Signature
                          │
                          ▼
         "Header.Payload.Signature"
```

**类比总结**：

JWT 签名就像快递的防伪封条：
- **包裹内容（Payload）**：任何人都能看（Base64 解码）
- **防伪封条（Signature）**：只有邮局（服务端）能验证真伪
- **如果有人修改内容**：封条会破损（签名不匹配），邮局拒收

**JWT 安全性依赖于：密钥保密 + HTTPS + 合理的过期时间 + 不在 Payload 存敏感信息**

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
