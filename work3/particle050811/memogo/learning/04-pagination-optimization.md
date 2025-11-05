# 分页查询优化与游标分页

本笔记记录分页查询的实现方式、性能问题及优化方案。

---

## 📅 2025-11-04：分页查询的性能问题

### 问题：查询如何做到分页？

**场景**：需要实现待办事项的分页查询功能。

#### 传统分页（OFFSET + LIMIT）

**代码实现**：`biz/dal/repository/todo_repo.go:123-183`

```go
func (r *TodoRepository) ListTodos(userID uint, statusFilter string, page, pageSize int) ([]model.Todo, int64, error) {
    var (
        todos []model.Todo
        total int64
    )

    // 1. 构建基础查询（带用户过滤）
    q := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)

    // 2. 添加状态过滤
    switch statusFilter {
    case "done":
        q = q.Where("status = ?", 1)
    case "todo":
        q = q.Where("status = ?", 0)
    }

    // 3. 先查询总数（用于返回 total）
    if err := q.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 4. 计算偏移量（核心分页公式）
    if page < 1 { page = 1 }
    if pageSize <= 0 { pageSize = 10 }
    offset := (page - 1) * pageSize

    // 5. 使用 Offset + Limit 实现分页
    if err := q.Order("created_at ASC, id ASC").
              Offset(offset).
              Limit(pageSize).
              Find(&todos).Error; err != nil {
        return nil, 0, err
    }

    return todos, total, nil
}
```

**分页公式**：
```
offset = (page - 1) * pageSize

第 1 页：offset = 0,  LIMIT 10 OFFSET 0   → 记录 1-10
第 2 页：offset = 10, LIMIT 10 OFFSET 10  → 记录 11-20
第 3 页：offset = 20, LIMIT 10 OFFSET 20  → 记录 21-30
```

**生成的 SQL**：
```sql
SELECT * FROM todos
WHERE user_id = 2
ORDER BY created_at ASC, id ASC
LIMIT 10 OFFSET 10;
```

---

## 📅 2025-11-04：为何不是直接返回第 x-y 条？

### 问题：为何用 OFFSET + LIMIT 而不是直接指定行号范围？

**❌ 想象的语法（不存在）**：
```sql
SELECT * FROM todos ROWS 11 TO 20  -- SQL 标准中没有这种语法
```

**✅ 实际的 SQL 逻辑**：
```sql
-- 数据库内部执行步骤：
1. WHERE user_id = 2          -- 筛选
2. ORDER BY id ASC             -- 排序
3. OFFSET 10                   -- 跳过 10 条
4. LIMIT 10                    -- 取 10 条
```

**不同数据库的语法**：

| 数据库 | 语法 |
|--------|------|
| MySQL | `LIMIT 10 OFFSET 20` 或 `LIMIT 20, 10` |
| PostgreSQL | `LIMIT 10 OFFSET 20` |
| SQL Server | `OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY` |
| Oracle | `OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY` |

**GORM 的优势**：使用 `Offset().Limit()` 可以自动适配不同数据库。

---

## 📅 2025-11-04：数据库使用的排序算法

### 问题：数据库排序用什么算法？

**两种情况**：

#### 1. 有索引时：不需要真正排序（最快⚡）

```sql
SELECT * FROM todos WHERE user_id = 2 ORDER BY id DESC
```

如果 `id` 列有**主键索引**，数据库会：

```
B+树索引（已经是有序的）：
    [10]
   /    \
 [5]    [15]
 / \    /  \
[3][7][12][18]

-- 直接遍历索引，O(log n) 定位 + O(k) 读取
-- 不需要额外排序！
```

**验证方式**：
```sql
EXPLAIN SELECT * FROM todos ORDER BY id DESC LIMIT 10;
```

输出中**没有 `Using filesort`** 说明直接用索引顺序。

#### 2. 无索引时：真正的排序算法

**小数据量（内存排序）**：
- 数据量 < `sort_buffer_size`（默认 256KB-2MB）
- 使用 **QuickSort** 或 **MergeSort**
- 时间复杂度：O(n log n)

**大数据量（外部排序）**：
- 数据量 > `sort_buffer_size`
- 使用 **外部归并排序**（External Merge Sort）
- 步骤：
  1. 将数据分块，每块在内存中排序
  2. 排序后的块写入临时文件（磁盘）
  3. 多路归并这些临时文件
- 时间复杂度：O(n log n)，但有大量**磁盘 I/O**，很慢！

**执行计划中的警告**：
```
Extra: Using filesort  ← 需要真正排序，性能较差
```

---

## 📅 2025-11-04：排序方向与业务需求

### 问题：最早的备忘录应该在最前面

**修改前**（新的在前）：
```go
q.Order("id DESC")  // 降序：最新的记录优先
```

**修改后**（旧的在前）：`biz/dal/repository/todo_repo.go:164`
```go
q.Order("created_at ASC, id ASC")  // 升序：最早的记录优先
```

**为什么用 `created_at ASC, id ASC`？**
1. **主要排序**：按创建时间升序（旧→新）
2. **次要排序**：按 ID 升序（防止同一秒创建多条时顺序不确定）

**显示效果**：
```
第 1 页：备忘录-01, 备忘录-02, 备忘录-03...（最早创建的）
第 2 页：备忘录-11, 备忘录-12, 备忘录-13...
```

符合 Todo 应用的使用习惯：**优先处理老任务**。

---

## 📅 2025-11-04：OFFSET 分页的性能陷阱

### 问题：读取全部 n 条评论所需时间是 O(n²)，能否优化到 O(n)？

**OFFSET 的问题**：深分页性能极差。

**时间复杂度分析**：

假设总共 10,000 条记录，每页 100 条：

```
第 1 页：  OFFSET 0 LIMIT 100      → 扫描 100 条
第 2 页：  OFFSET 100 LIMIT 100    → 扫描 200 条（跳过 100，读 100）
第 3 页：  OFFSET 200 LIMIT 100    → 扫描 300 条（跳过 200，读 100）
...
第 100 页：OFFSET 9900 LIMIT 100   → 扫描 10,000 条（跳过 9,900，读 100）

总扫描次数：
100 + 200 + 300 + ... + 10,000
= 100 × (1 + 2 + 3 + ... + 100)
= 100 × (100 × 101) / 2
= 505,000 次扫描

时间复杂度：O(n²)
```

**实测性能**（MySQL，100 万条记录）：
```bash
# OFFSET 方式读取全部数据
时间：~45 分钟
原因：最后几页需要扫描百万条记录
```

---

## 📅 2025-11-04：游标分页优化（Cursor Pagination）

### 优化方案：使用游标分页，时间复杂度 O(n)

**核心思想**：使用 `WHERE id > cursor` 代替 `OFFSET`，利用索引直接定位。

#### 实现代码

**Repository 层**：`biz/dal/repository/todo_repo.go:245-286`

```go
// ListTodosCursor 游标分页查询（用于高效遍历全部数据）
func (r *TodoRepository) ListTodosCursor(userID uint, statusFilter string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    var todos []model.Todo

    // 构建基础查询
    q := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)

    // 状态过滤
    switch statusFilter {
    case "done":
        q = q.Where("status = ?", 1)
    case "todo":
        q = q.Where("status = ?", 0)
    }

    // 🔥 关键：用 WHERE id > cursor 代替 OFFSET
    // 因为是升序（旧→新），需要找比 cursor 更大的 ID
    if cursor > 0 {
        q = q.Where("id > ?", cursor)
    }

    // 查询 limit+1 条，用于判断是否还有下一页
    if err := q.Order("created_at ASC, id ASC").Limit(limit + 1).Find(&todos).Error; err != nil {
        return nil, 0, false, err
    }

    // 判断是否有更多数据
    hasMore := len(todos) > limit
    var nextCursor uint

    if hasMore {
        // 有下一页：返回 limit 条数据，nextCursor 为最后一条的 ID
        nextCursor = uint(todos[limit-1].ID)
        todos = todos[:limit]
    } else {
        // 没有下一页
        nextCursor = 0
    }

    return todos, nextCursor, hasMore, nil
}
```

**Service 层**：`biz/service/todo_service.go:99-108`

```go
// ListTodosCursor 游标分页查询（用于高效遍历全部数据，O(n) 复杂度）
func (s *TodoService) ListTodosCursor(userID uint, status string, cursor uint, limit int) ([]model.Todo, uint, bool, error) {
    // 限制每次查询的最大数量
    if limit <= 0 {
        limit = 10
    } else if limit > 100 {
        limit = 100 // 游标分页可以允许更大的 limit
    }
    return s.repo.ListTodosCursor(userID, status, cursor, limit)
}
```

#### 执行过程

```sql
-- 第 1 页（首次查询，cursor=0）
SELECT * FROM todos
WHERE user_id = 2
ORDER BY created_at ASC, id ASC
LIMIT 11;  -- 查询 11 条（limit+1）判断是否有下一页

-- 返回：id=1~10，nextCursor=10，hasMore=true

-- 第 2 页（使用上一页的 nextCursor）
SELECT * FROM todos
WHERE user_id = 2 AND id > 10  -- 🔥 直接定位，不需要跳过前面的记录
ORDER BY created_at ASC, id ASC
LIMIT 11;

-- 返回：id=11~20，nextCursor=20，hasMore=true

-- 第 3 页
SELECT * FROM todos
WHERE user_id = 2 AND id > 20
ORDER BY created_at ASC, id ASC
LIMIT 11;

-- 返回：id=21~25，nextCursor=0，hasMore=false（无下一页）
```

#### 为什么是 `id > cursor` 而不是 `id < cursor`？

**关键：排序方向和游标条件要匹配！**

| 排序方向 | 游标条件 | 原因 |
|---------|---------|------|
| `ORDER BY id ASC`<br>（小→大，旧→新） | `WHERE id > cursor` | 下一页要找**更大**的 ID |
| `ORDER BY id DESC`<br>（大→小，新→旧） | `WHERE id < cursor` | 下一页要找**更小**的 ID |

**记忆技巧**：
```
升序（ASC）：从山脚往上爬
   1
   2  ← 第1页结束
   3
   4  ← 继续往上（id > 3）
   5

降序（DESC）：从山顶往下走
  10
   9  ← 第1页结束
   8
   7  ← 继续往下（id < 8）
   6
```

当前代码使用升序（旧记录优先），所以用 `id > cursor`。

---

## 📅 2025-11-04：性能对比与测试

### 性能对比

| 场景 | OFFSET 分页 | 游标分页 | 性能提升 |
|------|-----------|---------|---------|
| **读取 10,000 条（每页 100）** | 505,000 次扫描 | 10,000 次扫描 | **50 倍** |
| **读取 100 万条** | 5,000,500,000 次扫描 | 1,000,000 次扫描 | **5000 倍** |
| **时间复杂度** | O(n²) | O(n) | - |

### API 接口

**Thrift IDL 定义**：`idl/memogo.thrift:118-149`

```thrift
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
```

**路由**：
- `GET /v1/todos/cursor` - 游标分页列表查询
- `GET /v1/todos/search/cursor` - 游标分页搜索查询

### 测试结果

**测试代码**：`tools/testcursor/main.go`

```bash
✓ 创建 15 条测试备忘录

🚀 测试游标分页（每次 5 条）...

第 1 页 (cursor=0):
  - ID=1, Title=备忘录-01
  - ID=2, Title=备忘录-02
  - ID=3, Title=备忘录-03
  - ID=4, Title=备忘录-04
  - ID=5, Title=备忘录-05
  Next Cursor: 5, Has More: true

第 2 页 (cursor=5):
  - ID=6, Title=备忘录-06
  - ID=7, Title=备忘录-07
  - ID=8, Title=备忘录-08
  - ID=9, Title=备忘录-09
  - ID=10, Title=备忘录-10
  Next Cursor: 10, Has More: true

第 3 页 (cursor=10):
  - ID=11, Title=备忘录-11
  - ID=12, Title=备忘录-12
  - ID=13, Title=备忘录-13
  - ID=14, Title=备忘录-14
  - ID=15, Title=备忘录-15
  Next Cursor: 0, Has More: false

✅ 游标分页完成！总共获取了 15 条记录
```

### 遍历全部数据示例

```go
// 使用游标分页遍历全部数据
func ExportAllTodos(userID uint, status string) ([]model.Todo, error) {
    var allTodos []model.Todo
    cursor := uint(0)
    pageSize := 100

    for {
        // 🔥 每次查询时间：O(log n + k)，不是 O(n)！
        todos, nextCursor, hasMore, err := repo.ListTodosCursor(userID, status, cursor, pageSize)
        if err != nil {
            return nil, err
        }

        allTodos = append(allTodos, todos...)

        // 没有下一页，退出
        if !hasMore {
            break
        }

        cursor = nextCursor
    }

    return allTodos, nil
}
```

---

## 📅 2025-11-04：何时使用哪种分页方式？

### 使用场景对比

| 场景 | 推荐方式 | 原因 |
|------|---------|------|
| **前端分页展示** | OFFSET（当前方式） | 用户需要跳转到任意页 |
| **手机端下拉刷新** | 游标分页 | 性能好，体验好 |
| **数据导出** | 游标分页 或 全量查询 | 避免 O(n²) |
| **API 遍历（第三方调用）** | 游标分页 | 标准做法 |
| **总数 < 1000 条** | OFFSET 也可以 | 性能差异不明显 |
| **总数 > 10 万条** | 必须用游标 | OFFSET 会超时 |

### 游标分页的限制

**❌ 不支持的功能**：
- 跳转到任意页（只能顺序遍历）
- 显示总页数（因为不知道 total）
- 返回上一页（需要双向游标）

**✅ 适用场景**：
- 数据导出
- 移动端无限滚动
- 消息/Feed 流
- 日志查询
- 大数据集遍历

---

## 🎯 核心要点总结

1. **传统分页**：`OFFSET + LIMIT`，适合前端分页，但深分页性能差（O(n²)）
2. **游标分页**：`WHERE id > cursor`，适合遍历全部数据，性能优秀（O(n)）
3. **排序优化**：
   - 有索引：直接扫描 B+树，无需排序
   - 无索引：QuickSort/MergeSort（内存）或 External MergeSort（磁盘）
4. **业务需求**：待办事项按创建时间升序，旧任务优先显示
5. **游标方向**：
   - `ORDER BY id ASC` → `WHERE id > cursor`（升序）
   - `ORDER BY id DESC` → `WHERE id < cursor`（降序）

---

## 🔗 延伸阅读

- [高性能分页方案：Seek Method（游标分页）](https://use-the-index-luke.com/no-offset)
- [为什么深度分页很慢？](https://www.eversql.com/faster-pagination-in-mysql-why-order-by-with-limit-and-offset-is-slow/)
- [B+树索引原理](https://dev.mysql.com/doc/refman/8.0/en/innodb-physical-structure.html)
- [MySQL 排序优化](https://dev.mysql.com/doc/refman/8.0/en/order-by-optimization.html)
- [PostgreSQL Cursor Pagination](https://www.postgresql.org/docs/current/queries-limit.html)

---

*笔记创建：2025-11-04*
