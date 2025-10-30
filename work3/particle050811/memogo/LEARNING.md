# MemoGo 学习笔记

本文件记录开发过程中遇到的问题和解答，供后续学习参考。

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
这是**防御性编程**，区分三种情况：
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

## 延伸阅读

- OWASP 用户枚举攻击防护：https://owasp.org/www-project-web-security-testing-guide/
- Go 错误处理最佳实践：https://go.dev/blog/error-handling-and-go
- RESTful API 设计规范：https://restfulapi.net/

---

*本笔记持续更新中...*
