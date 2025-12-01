## Task1:爬取爬取福大通知、文件系统

### bonus1：使用并发爬取，同时给出加速比（最大并发数为10没有多开，福大的网站有点脆弱OWO）

#### 单线程：33m1.6689406s  多线程：3m57.4190083s

#### 加速比约为：9.3

### bounus2：搜集每个通知的访问人数

#### 动态数据，拿到url的参数用客户端去GET就行

### bonuse3 ：将爬取的数据存入数据库，原生SQL或ORM映射都可以（√）

## Task2：爬取Bilibili视频评论

### bonus1：给出Bilibili爬虫检测阈值

#### 500ms的请求延迟，爬了会被ban了，测了10ms的延迟秒ban

### bonus2：给出爬取的流程图，使用mermaid或者excalidraw

![mermaid-diagram (1)](C:\Users\ZOE\Desktop\新建文件夹 (2)\work2\mermaid-diagram (1).png)

### bonus3：给出接口返回的json中每个参数所代表的意义

#### 写在程序注释（其实API文档都有(●'◡'●)）

