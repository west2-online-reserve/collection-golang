package main

import (
	"database/sql"
	"fmt"
	"github.com/bytedance/sonic"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var Db *sql.DB

// 同步并且限制goroutines 数量
var wg = sync.WaitGroup{}
var ch1 = make(chan bool, 60)

// 初始化数据库
func init() {
	var err error
	Db, err = sql.Open("mysql", "root:114514@tcp(127.0.0.1:3306)/fzuA")
	if err != nil {
		panic(err)
	}
}

//在terminal创建表
// CREATE TABLE `bilibili` (
// `id` bigint unsigned AUTO_INCREMENT,
// `uname` varchar(1000) DEFAULT NULL,
// `time` varchar(100) DEFAULT NULL,
// `content` longtext,
// PRIMARY KEY (`id`),
// UNIQUE KEY `id` (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

// B json解析
type B struct {
	Data struct {
		Replies Replies `json:"replies"`
	} `json:"data"`
}
type Replies []struct {
	Rpid_str string `json:"rpid_str"`
	Rcount   int    `json:"rcount,omitempty"`
	Ctime    int    `json:"ctime"`
	Member   struct {
		Uname string `json:"uname"`
	} `json:"member"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
}

// 时间戳转换
func timeEx(Time int) (fmtTime string) {
	t := time.Unix(int64(Time), 0)
	fmtTime = t.Format("2006-01-02 15:04:05")
	return
}

// 懒得写一堆err
func iferr(err error, ques string) {
	if err != nil {
		fmt.Printf("%s err:%v", ques, err)
	}
}

// JSONGET 发送get请求，接收请求并进行处理
func JSONGET(url string) (jsonIn string) {
	resp, err := http.NewRequest(http.MethodGet, url, nil)
	iferr(err, "httpGET")
	resp.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.41")
	resp.Header.Add("authority", "api.bilibili.com")
	resp.Header.Add("Accept", "*/*")
	resp.Header.Add("Referer", "https://www.bilibili.com/video/BV12341117rG/?")
	r, err1 := http.DefaultClient.Do(resp)
	iferr(err1, "doRequest")
	defer func() { _ = r.Body.Close() }()
	buf := make([]byte, 4096)
	for {
		con, err2 := r.Body.Read(buf)
		if err2 != nil && err2 != io.EOF {
		}
		if con == 0 {
			break
		}
		jsonIn += string(buf[:con])
	}
	//fmt.Println(jsonIn)
	return
}

// 使用sonic库，对json进行解析
func jsonMarshal(jsonIn string) {
	t := B{}
	err := sonic.UnmarshalString(jsonIn, &t)
	iferr(err, "unmarshal")
	if len(t.Data.Replies) == 0 {
		//fmt.Println("爬完了")
		return
	}
	//计算子评论页数
	for i := 0; i < len(t.Data.Replies); i++ {
		if t.Data.Replies[i].Rcount == 0 {
			saveToSQL(&t, i)
		} else {
			var page int
			if t.Data.Replies[i].Rcount%10 == 0 {
				page = t.Data.Replies[i].Rcount / 10
			} else {
				page = t.Data.Replies[i].Rcount/10 + 1
			}
			//开启子评论并发
			for j := 1; j <= page; j++ {
				wg.Add(1)
				goSubComment(&t, i, j)
				ch1 <- true
				fmt.Println(runtime.NumGoroutine())
			}
			saveToSQL(&t, i)
		}
	}
	return
}

// 对子评论进行解析
func goSubComment(t *B, i int, j int) {
	s := B{}
	url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + t.Data.Replies[i].Rpid_str + "&ps=10&pn=" + strconv.Itoa(j) + "&web_location=333.788"
	sub := JSONGET(url)
	err := sonic.UnmarshalString(sub, &s)
	iferr(err, "subUnmarshal")
	for k := 0; k < len(s.Data.Replies); k++ {
		saveToSQL(&s, k)
	}
	<-ch1
	wg.Done()
}

// 将处理后的数据存入数据库
func saveToSQL(t *B, i int) {
	fmtTime := timeEx(t.Data.Replies[i].Ctime)
	sql1 := "insert into bilibili (uname,time,content) values (?,?,?)"
	_, err := Db.Exec(sql1, t.Data.Replies[i].Member.Uname, fmtTime, t.Data.Replies[i].Content.Message)
	iferr(err, "saveToSQL err")
}

// 主评论相关
func goMainComment(i int) {
	url := "https://api.bilibili.com/x/v2/reply?&type=1&pn=" + strconv.Itoa(i) + "&oid=420981979&sort=2"
	jsonIn := JSONGET(url)
	jsonMarshal(jsonIn)
	<-ch1
	wg.Done()
}

func main() {
	//pn大于1000时b站返回请求错误
	for i := 1; i < 1000; i++ {
		wg.Add(1)
		ch1 <- true
		go goMainComment(i)
	}
	wg.Wait()
}

//有想过通过再创建一个channel，判断len(t.Data.Replies)=0时返回end来break for循环,结果导致效率太低，不如直接1k页都遍历一遍
