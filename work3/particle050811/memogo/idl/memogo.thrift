namespace go memogo.api

typedef i64 Timestamp  // unix seconds

// ---------- 基础模型 ----------
enum TodoStatus {
  TODO = 0,
  DONE = 1
}

struct Todo {
  1: i64        id
  2: string     title
  3: string     content
  4: i32        view
  5: TodoStatus status
  6: Timestamp  created_at
  7: Timestamp  start_time
  8: Timestamp  end_time
  9: Timestamp  due_time
}

// ---------- 认证与用户 ----------
struct RegisterReq {
  1: string username
  2: string password
}
struct LoginReq {
  1: string username
  2: string password
}
struct RefreshReq {
  1: string refresh_token
}

struct TokenPair {
  1: string access_token
  2: string refresh_token
  3: i64    access_expires_in       // seconds
  4: i64    refresh_expires_in      // seconds
}

// 统一认证返回，贴合示例的 status/msg/data 结构
struct AuthResp {
  1: i32       status
  2: string    msg
  3: TokenPair data
}

// ---------- 待办 - 创建 ----------
struct CreateTodoReq {
  1: optional string    authorization (api.header = "Authorization") // Bearer <token>
  2: string             title
  3: string             content
  4: optional Timestamp start_time
  5: optional Timestamp due_time
}

struct CreateTodoResp {
  1: i32    status
  2: string msg
  3: Todo   data
}

// ---------- 待办 - 更新状态（单条 / 批量） ----------
struct UpdateTodoStatusReq {
  1: optional string authorization (api.header = "Authorization")
  2: i64             id          (api.path   = "id")
  3: TodoStatus      status      // 目标状态：TODO 或 DONE
}
struct UpdateTodoStatusResp {
  1: i32    status
  2: string msg
  3: i32    data      // 受影响条数（单条通常为1）
}

// 将所有满足 from_status 的事项批量改为 to_status
struct UpdateAllStatusReq {
  1: optional string authorization (api.header = "Authorization")
  2: TodoStatus      from_status   (api.query = "from")
  3: TodoStatus      to_status     (api.query = "to")
}
struct UpdateAllStatusResp {
  1: i32    status
  2: string msg
  3: i32    data       // 受影响条数
}

// ---------- 待办 - 查询与搜索（分页） ----------
struct ListTodosReq {
  1: optional string    authorization (api.header = "Authorization")
  2: optional string    status        (api.query = "status") // "todo" | "done" | "all"
  3: i32                page          (api.query = "page")
  4: i32                page_size     (api.query = "page_size")
}
struct ItemsTodoData {
  1: list<Todo> items
  2: i64        total
}
struct ListTodosResp {
  1: i32           status
  2: string        msg
  3: ItemsTodoData data
}

struct SearchTodosReq {
  1: optional string authorization (api.header = "Authorization")
  2: string          q             (api.query = "q")         // 关键词
  3: i32             page          (api.query = "page")
  4: i32             page_size     (api.query = "page_size")
}
struct SearchTodosResp {
  1: i32           status
  2: string        msg
  3: ItemsTodoData data
}

// ---------- 待办 - 游标分页（高效遍历，O(n) 复杂度） ----------
struct ListTodosCursorReq {
  1: optional string authorization (api.header = "Authorization")
  2: optional string status        (api.query = "status")  // "todo" | "done" | "all"
  3: i64             cursor         (api.query = "cursor")  // 上一页最后一条的 ID，首次传 0
  4: i32             limit          (api.query = "limit")   // 每页数量，默认 10，最大 100
}

struct CursorTodoData {
  1: list<Todo> items
  2: i64        next_cursor  // 下一页的游标，0 表示无下一页
  3: bool       has_more     // 是否还有更多数据
}

struct ListTodosCursorResp {
  1: i32            status
  2: string         msg
  3: CursorTodoData data
}

struct SearchTodosCursorReq {
  1: optional string authorization (api.header = "Authorization")
  2: string          q             (api.query = "q")      // 关键词
  3: i64             cursor         (api.query = "cursor")
  4: i32             limit          (api.query = "limit")
}

struct SearchTodosCursorResp {
  1: i32            status
  2: string         msg
  3: CursorTodoData data
}

// ---------- 待办 - 删除 ----------
struct DeleteOneReq {
  1: optional string authorization (api.header = "Authorization")
  2: i64             id            (api.path   = "id")
}
struct DeleteByScopeReq {
  1: optional string authorization (api.header = "Authorization")
  2: string          scope         (api.query = "scope") // "done" | "todo" | "all"
}
struct DeleteResp {
  1: i32    status
  2: string msg
  3: i32    data     // 受影响条数
}

// ========== Service 定义（HTTP 映射） ==========
service MemoGoService {

  // 用户模块：注册 / 登录 / 刷新令牌
  AuthResp Register(1: RegisterReq req)     (api.post = "/v1/auth/register")
  AuthResp Login(1: LoginReq req)           (api.post = "/v1/auth/login")
  AuthResp RefreshToken(1: RefreshReq req)  (api.post = "/v1/auth/refresh")

  // 事务模块：增
  CreateTodoResp CreateTodo(1: CreateTodoReq req) (api.post = "/v1/todos")

  // 事务模块：改（单条改状态）
  UpdateTodoStatusResp UpdateTodoStatus(1: UpdateTodoStatusReq req) (api.patch = "/v1/todos/:id/status")

  // 事务模块：改（全部/批量改状态；用 from/to 表示"将所有 X 改为 Y"）
  UpdateAllStatusResp UpdateAllStatus(1: UpdateAllStatusReq req) (api.patch = "/v1/todos/status")

  // 事务模块：查（分页 + 状态过滤）
  ListTodosResp ListTodos(1: ListTodosReq req) (api.get = "/v1/todos")

  // 事务模块：查（分页 + 关键词搜索）
  SearchTodosResp SearchTodos(1: SearchTodosReq req) (api.get = "/v1/todos/search")

  // 事务模块：查（游标分页，高效遍历）
  ListTodosCursorResp ListTodosCursor(1: ListTodosCursorReq req) (api.get = "/v1/todos/cursor")

  // 事务模块：查（关键词 + 游标分页）
  SearchTodosCursorResp SearchTodosCursor(1: SearchTodosCursorReq req) (api.get = "/v1/todos/search/cursor")

  // 事务模块：删（单条）
  DeleteResp DeleteOne(1: DeleteOneReq req) (api.delete = "/v1/todos/:id")

  // 事务模块：删（按 scope 删除：done/todo/all）
  DeleteResp DeleteByScope(1: DeleteByScopeReq req) (api.delete = "/v1/todos")
}
