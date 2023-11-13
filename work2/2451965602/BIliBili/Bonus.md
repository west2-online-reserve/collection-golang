## 给出Bilibili爬虫检测阈值（请求频率高于这个阈值将会被ban）

## 给出爬取的流程图，使用mermaid

graph LR
    A[获取该视频总评论数] -->B(计算总页码,获取页面范围) 
    B --> |迭代页码| C(获取每一页所有主评论的rpid)
    C -->|迭代页码rpid| D[rpid与api结合生成单层评论url]
    D -->|获取该层总评论数| E(计算该层总页码,获取页面范围)
    E -->|迭代页码| F(获取该层所有评论的内容)
    F -->|创建 第 + rpid + 楼 的txt文件| G(写入评论)
## 给出接口返回的json中每个参数所代表的意义
评论api  https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/comment/list.md