### DCL-介绍
DCL英文全称是Data Control Language(数据控制语言)，用来管理数据库用户、控制数据库的访问权限。

***

### DCL-管理用户
1. **查询用户**
```sql
USE mysql;
SELECT * FROM user;
```
2. **创建用户**
```sql
CREATE USER '用户名'@'主机名' IDENTIFIED BY '密码';
```
3. **修改用户密码**
```sql
ALTER USER '用户名'@'主机名' IDENTIFIED WITH mysql_native_password BY '新密码';
```
4. **删除用户**
```sql
DROP USER '用户名'@'主机名';
```

#### 注意
- 主机名可以使用 `%` 通配。
- 这类 SQL 开发人员操作的比较少，主要是 DBA（Database Administrator，数据库管理员）使用。

***

### DCL-权限控制

MySQL 中定义了很多种权限，但常用的有以下几种：

| 权限                | 说明                 |
| ------------------- | --------------------|
| ALL, ALL PRIVILEGES | 所有权限             |
| SELECT              | 查询数据             |
| INSERT              | 插入数据             |
| UPDATE              | 修改数据             |
| DELETE              | 删除数据             |
| ALTER               | 修改表               |
| DROP                | 删除数据库/表/视图    |
| CREATE              | 创建数据库/表        |

***

### DCL-权限控制
1. **查询权限**
```sql
SHOW GRANTS FOR '用户名'@'主机名';
```
2. **授予权限**
```sql
GRANT 权限列表 ON 数据库名.表名 TO '用户名'@'主机名';
```
3. **撤销权限**
```sql
REVOKE 权限列表 ON 数据库名.表名 FROM '用户名'@'主机名';
```

#### 注意
- 多个权限之间，使用逗号分隔。
- 授权时，数据库名和表名可以使用 `*` 进行通配，代表所有。

***

