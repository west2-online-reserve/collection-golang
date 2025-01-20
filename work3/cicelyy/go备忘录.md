# HTTP协议的工作原理

1. **客户端请求（Request）**：
   - 客户端（通常是浏览器）发送一个HTTP请求到服务器。这个请求包含了请求方法（如GET、POST、PUT、DELETE等）、请求的资源路径（URL）、HTTP版本和可能的请求头（Headers）和请求体（Body）。
2. **服务器处理请求（Response）**：
   - 服务器接收到请求后，根据请求的类型和路径处理请求。服务器可能需要查询数据库、执行后端逻辑或者直接读取文件系统中的资源。
3. **服务器响应（Response）**：
   - 服务器处理完请求后，会返回一个HTTP响应给客户端。这个响应包含了状态码（如200表示成功，404表示未找到等）、响应头和响应体（Body）。
4. **客户端处理响应（Client Processing）**：
   - 客户端接收到响应后，根据状态码和响应体进行相应的处理。如果是HTML文档，浏览器会解析并渲染页面；如果是图片或CSS文件，浏览器会相应地处理这些资源。
5. **连接关闭（Connection Close）**：
   - 在HTTP/1.0中，默认的连接是短连接，即每次请求/响应后都会关闭TCP连接。而在HTTP/1.1中，默认是长连接（持久连接），即TCP连接可以被多个请求/响应复用，减少了建立和关闭连接的开销。

# Web工作原理

1. **用户输入URL**：
   - 用户在浏览器地址栏输入一个URL。
2. **DNS解析**：
   - 浏览器通过DNS（域名系统）解析URL对应的IP地址。
3. **建立TCP连接**：
   - 浏览器与服务器之间建立TCP连接。
4. **发送HTTP请求**：
   - 浏览器构建HTTP请求并发送给服务器。
5. **服务器处理请求**：
   - 服务器接收请求，并根据请求处理相应的资源。
6. **返回HTTP响应**：
   - 服务器将处理结果以HTTP响应的形式发送回浏览器。
7. **浏览器渲染页面**：
   - 浏览器接收响应，并根据响应内容渲染页面。
8. **关闭TCP连接**（如果是HTTP/1.0）：
   - 请求完成后，TCP连接被关闭。
9. **资源加载**：
   - 页面中的图片、CSS、JavaScript等资源也会通过HTTP协议加载。
10. **用户交互**：
    - 用户可以与页面进行交互，如点击链接、提交表单等，这会触发新的HTTP请求。

# Gin框架

- **发布背景**：Gin是Go语言中最流行的Web框架之一，由Gin团队开发。它的目标是简单、快速、灵活，特别适合构建高性能的Web应用和RESTful API。
- **设计理念**：Gin基于`httprouter`，一个高效的路由器。它通过提供中间件支持、路由组、请求处理等特性，简化了开发者的工作。
- **性能**：Gin在路由匹配速度上较快，支持高并发请求处理。适合中小型Web应用，尤其是在构建RESTful API时。
- **路由**：使用基于Trie树的路由器，支持动态路由、路由分组等特性。
- **中间件**：提供强大的中间件支持，可以方便地为路由添加中间件，如日志、认证、CORS支持等。
- **可扩展性与插件系统**：Gin提供良好的扩展性，允许开发者通过自定义中间件、钩子等方式扩展框架的功能。
- **错误处理与日志**：提供简洁的错误处理机制，支持通过中间件实现日志功能，通常结合第三方库如`logrus`。
- **社区与文档**：Gin拥有庞大的开发者社区和丰富的文档资源，适合快速开发。

# Hertz框架

- **发布背景**：Hertz由字节跳动（ByteDance）开发，旨在提供更高的性能、低延迟和更高效的资源利用。主要面向高并发场景。
- **设计理念**：Hertz是一个高性能的Web框架，提供了高效的路由和请求处理机制。它在性能上进行了优化，尤其是在请求调度、路由匹配和中间件处理上进行了创新设计。
- **性能**：Hertz的设计优化了性能，特别是在高并发、大吞吐量场景下。通过底层优化（如高效的路由匹配算法、异步任务处理等）提升了性能。
- **路由**：Hertz在路由方面进行了更多的优化，采用了一些自定义的高效路由匹配算法，提升了请求匹配速度。
- **中间件**：Hertz的中间件支持也很强大，提供了更多的自定义和灵活性。它的中间件设计上更加高效，支持在请求处理链中进行异步操作。
- **可扩展性与插件系统**：Hertz具备良好的扩展性，支持丰富的插件机制。它的插件系统较为灵活，允许开发者根据需要加载和管理不同的插件。
- **错误处理与日志**：提供了更为完善的错误处理机制，支持在请求处理中更细粒度的错误捕获和管理。日志系统较为高效，能够处理大规模日志的输出。
- **社区与文档**：Hertz作为相对较新的框架，其社区和文档的成熟度较Gin稍逊色，但也得到了不少开发者和企业的关注。

总结来说，如果你的项目主要关注性能并且需要处理大量的并发请求，Hertz可能是更好的选择；如果你更注重开发速度和社区支持，Gin可能会更合适。

## 概述

三高：高扩展性、高易用性、高性能





# RESTful API

是一种软件架构风格，用于设计网络应用程序。它基于HTTP协议，使用HTTP方法来执行资源的操作。以下是RESTful API的一些核心规范：

1. **使用HTTP方法**：
   - **GET**：用于获取资源。
   - **POST**：用于创建新资源。
   - **PUT**：用于更新现有资源。
   - **DELETE**：用于删除资源。
   - **PATCH**：用于对资源进行部分修改。
2. **无状态**：
   - 每个请求包含所有必要的信息，服务器不需要保存会话信息。这意味着对于相同的请求，服务器总是返回相同的响应。
3. **统一接口**：
   - 系统组件之间的交互应该是统一的，即通过一组标准的接口进行。
4. **资源导向**：
   - API应该以资源为中心，资源通常对应于现实世界中的对象或概念，如用户、订单等。
5. **URI（统一资源标识符）**：
   - 每个资源都有一个唯一的URI，用于识别和定位资源。
   - URI应该是自描述的，清晰地表达资源的语义。
6. **资源的表述**：
   - 客户端请求的资源可以通过多种格式（如JSON、XML）来表述，通常在请求的`Accept`头中指定。
7. **超媒体驱动**（HATEOAS，Hypermedia as the Engine of Application State）：
   - 服务器响应应该包含超媒体链接，这些链接指向其他资源，客户端可以通过这些链接来发现API的其他部分。
8. **分层系统**：
   - 客户端和服务器之间的通信应该是分层的，客户端不应该了解服务器后面的数据存储细节。
9. **缓存**：
   - 响应应该被标记为可缓存或不可缓存，以利用HTTP缓存机制提高效率。
10. **代码重用**：
    - 通过无状态操作和统一接口，可以在不同的上下文中重用代码。
11. **客户端-服务器分离**：
    - 将用户界面关注点与数据存储关注点分离，提高跨多种平台的界面可移植性。
12. **按需代码**（可选）：
    - 服务器可以根据请求发送代码给客户端执行，例如JavaScript。

遵循这些规范，可以设计出易于理解和使用的API，同时也能够提高API的可维护性和可扩展性。

[字节开源WEB框架Hertz太香啦！安装Hertz命令行工具 请确保您的Go版本在1.15及以上版本，笔者用的版本是1. - 掘金](https://juejin.cn/post/7124337913352945672)

在 RESTful API 中，使用的主要是以下五种HTTP方法：

1. GET，表示读取服务器上的资源
2. POST，表示在服务器上创建资源
3. PUT,表示更新或者替换服务器上的资源
4. DELETE，表示删除服务器上的资源
5. PATCH，表示更新/修改资源的一部分

两个`GET`方法的示例，第一个表示获取所有用户的信息；第二个表示获取`id`为`123`用户的信息

```path
HTTP GET https://www.flysnow.org/users
HTTP GET https://www.flysnow.org/users/123
```



创建一个用户，会通过`POST`给服务器提供创建这个用户所需的全部信息。注意这里`users`是个复数

```path
HTTP POST https://www.flysnow.org/users
```



表示要更新/替换`id`为`123`的这个用户，在更新的时候，会通过`PUT`提供更新这个用户需要的全部用户信息。这里`PUT`和`POST`不太一样的是 ，从URL看，`PUT`操作的是单个资源，比如这里`id`为`123`的这个用户。

```path
HTTP PUT https://www.flysnow.org/users/123
```



删除`id`为`123`的这个用户。

```path
HTTP DELETE https://www.flysnow.org/users/123
```



`PATCH`也更新资源，它和`PUT`不一样的是，它只能更新这个资源的部分信息，而不是全部(这种也叫替换)，是部分更新。

```path
HTTP PATCH https://www.flysnow.org/users/123
```



`Gin`提供了`Any`方法，可以一次性注册以上这些`HTTP Method`方法。

```push
// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
    group.handle("GET", relativePath, handlers)
    group.handle("POST", relativePath, handlers)
    group.handle("PUT", relativePath, handlers)
    group.handle("PATCH", relativePath, handlers)
    group.handle("HEAD", relativePath, handlers)
    group.handle("OPTIONS", relativePath, handlers)
    group.handle("DELETE", relativePath, handlers)
    group.handle("CONNECT", relativePath, handlers)
    group.handle("TRACE", relativePath, handlers)
    return group.returnObj()
}
```



如果你只想注册其中某两个、或者三个方法，`Gin`就没有这样的便捷方法了，不过`Gin`为我们提供了通用的`Handle`方法，我们可以包装一下使用。

```push
func Handle(r *gin.Engine, httpMethods []string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
    var routes gin.IRoutes
    for _, httpMethod := range httpMethods {
        routes = r.Handle(httpMethod, relativePath, handlers...)
    }
    return routes
}
```

```push
Handle(r, []string{"GET", "POST"}, "/", func(c *gin.Context) {
   //同时注册GET、POST请求方法
})
```

虽然这种方式比较便利，但是并不太推荐，因为他破坏了Resultful 规范中HTTP Method的约束。



# 项目结构

```text
ZZZZPROJECTS3/
└── memo/
    ├── api/
    │   ├── main.go      // API 路由入口
    │   ├── tasks.go    // 任务相关的 API 处理函数
    │   └── user.go     // 用户相关的 API 处理函数
    ├── conf/
    │   └── jwt.go      // JWT 配置和中间件
    ├── middleware/
    │   └── jwt.go      // JWT 中间件实现
    ├── model/
    │   ├── init.go     // 数据库初始化
    │   ├── migrate.go  // 数据库迁移脚本
    │   ├── task.go    // 任务数据模型
    │   └── user.go    // 用户数据模型
    ├── pkg/
    │   └── utils/
    │       └── utils.go  // 工具函数
    ├── routes/
    │   └── routes.go  // 路由配置
    ├── serializer/
    │   ├── common.go  // 通用序列化函数
    │   ├── task.go    // 任务序列化函数
    │   └── user.go    // 用户序列化函数
    ├── service/
    │   ├── task.go    // 任务业务逻辑
    │   └── user.go    // 用户业务逻辑
    ├── go.mod        // Go Modules 配置文件
    ├── go.sum        // 依赖版本锁定文件
    └── main.go        // 应用入口
```

## api

使用了 Gin Web 框架来处理 HTTP 请求，`service` 包中的服务来执行业务逻辑，并使用 `logrus` 库来记录错误日志。每个函数都遵循相同的模式：解析请求数据，执行业务逻辑，然后返回结果或错误。

## conf

定义了应用程序的配置管理。通过将配置信息存储在外部文件中，可以方便地进行修改和维护，而无需更改代码。使用 Go 语言的 `ini` 包来解析配置文件，使得配置加载过程简单明了。

## model

定义结构体和函数，实现了与数据库的交互。使用 GORM 库来简化数据库操作，并使用 `bcrypt` 库来安全地存储和验证密码。通过自动迁移功能，确保了数据库结构与代码中定义的模型结构保持一致。

## pkg

提供了生成和解析 JWT 的功能。通过 `GenerateToken` 函数，在用户登录成功后生成一个 JWT，然后将其发送给客户端。客户端在随后的请求中使用这个 JWT 来证明其身份。服务器通过 `ParseToken` 函数来验证 JWT 的有效性，从而确认请求者的身份和权限。 

## routes

设置 Gin Web 框架的路由和中间件，以处理不同的 HTTP 请求。通过定义路由组和中间件，可以组织代码并实现身份验证等功能。用户操作（如注册和登录）不需要身份验证，而其他操作（如任务管理）需要通过 JWT 中间件进行身份验证。

## serializer

定义数据传输对象（DTO），这些对象用于在 API 层和客户端之间传输数据。通过将数据库模型序列化为这些 DTO，可以更容易地控制发送给客户端的数据格式和内容。

## service

封装了应用的业务逻辑，包括用户注册、登录、任务的创建、展示、列出、更新、搜索和删除。



