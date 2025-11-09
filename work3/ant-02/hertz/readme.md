# 1 swagger生成自动文档

![img](D:\desktop\go_west2\3-备忘录\hertz\swagger.png)

# 2 使用三次架构，handler用于接收请求和返回响应结果，service实现核心业务流程，repository负责封装对数据源的操作

# 3 对于一些数据库操作采用事务处理，如分页查询待办事件的数据，查询全部待办的总数量和分页查询待办信息放在同一个事务中完成

# 4 更优秀的json格式

```json
{
  "status": 200,                    // 200 表示正常/成功，500 代表错误。自行了解HTTP状态码。
  "msg": "ok",	                    // 返回信息  
  "data": {	                    // 业务数据。所有的业务信息都应该放到 data 对象上。
    "items": [
      {
        "id": 1,// 待办事项ID
        "title": "更改好了！",        // 主题
        "content": "好耶！",	    // 内容
        "view": 0,	            // 访问次数
        "status": 1,	            // 状态（正在进行/已完成/其他）
        "created_at": 1638257438,   // 创建时间
        "start_time": 1638257437,   // 开始时间
        "end_time": 0	            // 结束时间
      }
    ],
    "pagination": {
    	"page_num": 1,	//页面数
        "page_size": 10,	//页长度
        "total": 1,	//总数量
        "total_pages": 1	//总页数
    }                   // 检索出的匹配全部条目数（不是items的len值）
  }               
}
```

# 5 使用redis存储jwt，对用户进行身份认证。使用redis的list数据结构，通过旁路缓存策略存储待办信息`key`格式为`todo:list:{status}:{uid}`，`status`为待办事务状态，`uid`为用户id。查询先查缓存，命中返回，不命中查询数据库并存入缓存。增加、修改和删除清空缓存，再操作数据库。

# 6 项目结构

## 6.1 项目由go+hertz+redis+mysql设计实现

## 6.2 项目目录为

```sh
hertz.
│  .gitignore
│  .hz
│  go.mod
│  go.sum
│  main.go	// 程序入口，由hz生成
│  readme.md	//本文档
│  router.go	//由hz生成
│  router_gen.go	//由hz生成
│  swagger.png	//文档图片
│
├─biz	//由hz生成
│  ├─handler	//handler层目录
│  │  │  ping.go	//测试http服务是否开启
│  │  │
│  │  ├─todo	//待办事务http接口目录，由
│  │  │      todo_service.go	//待办事务http接口具体实现
│  │  │
│  │  └─user	//用户http接口目录
│  │          user_service.go	//用户http接口具体实现
│  │
│  ├─model	//接口数据类型目录
│  │  ├─api	//hz的api接口数据类型目录
│  │  │      api.pb.go	//hz的api接口数据类型具体代码
│  │  │
│  │  ├─todo	//待办事务的接口数据类型目录
│  │  │      todo.pb.go	//待办事务的接口数据类型具体代码
│  │  │
│  │  └─user	//用户的接口数据类型目录
│  │          user.pb.go	//用户的接口数据类型具体代码
│  │
│  └─router	//路由信息与注册目录
│      │  register.go	//注册代码
│      │
│      ├─todo	//待办事务路由信息目录
│      │      middleware.go	//待办事务路由插入中间件
│      │      todo.go	//待办事务路由具体代码
│      │
│      └─user	//用户路由信息目录
│              middleware.go	//用户路由插入中间件
│              user.go	//用户路由具体代码
│
├─config	//配置信息目录
│      config.go	//提取配置信息代码
│      config.yaml	//配置信息内容
│
├─database	//数据库目录
│      mysql.go	//连接mysql代码
│      redis.go	//连接redis代码
│	
├─docs	//由swagger生成
│      docs.go	//swagger文档源码
│      swagger.json	//swagger文档json数据
│      swagger.yaml	//swagger文档yaml数据
│
├─idl	//proto规则目录，用于hz生成代码
│      api.proto	//hz官方源码
│      todo.proto	//待办事务的规则
│      user.proto	//用户的规则
│
├─pkg	//手写部分
│  ├─middleware	//中间件
│  │      jwt.go	//身份认证中间件
│  │
│  ├─model	//model层目录
│  │      todo.go	//待办事务model
│  │      user.go	//用户model
│  │
│  ├─repository	//repositoy层目录
│  │      todo.go	//待办事务数据库操作
│  │      user.go	//用户数据库操作
│  │
│  └─service	//service层目录
│          todo.go	//待办事务事务处理
│          user.go	//用户事务处理
│
└─util	//工具
        md5.go	//密码加密
```

