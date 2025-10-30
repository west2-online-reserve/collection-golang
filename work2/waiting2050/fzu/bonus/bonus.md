### 使用并发爬取，同时给出加速比（加速比：相较于普通爬取，快了多少倍）
![普通爬取](https://github.com/waiting2050/go_learning/blob/main/erlun/task/fzu/chuanxing.png?raw=true)
![并发爬取](https://github.com/waiting2050/go_learning/blob/main/erlun/task/fzu/bingfa.png?raw=true)
==加速比==：16.81

### 搜集每个通知的访问人数
由于点击量数据不在初始 HTML 中，而是需要通过接口动态获取，且接口依赖网页脚本中的隐藏参数，因此先提取参数，再模拟浏览器发起请求，最后解析结果，获取点击量。

### 将爬取的数据存入数据库，原生SQL或ORM映射都可以
![示例图片](https://github.com/waiting2050/go_learning/blob/main/erlun/task/fzu/shili.png?raw=true)