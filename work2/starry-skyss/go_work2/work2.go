package main

import (
	"database/sql"
	"fmt"
	//"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
	"io"
	_ "io/ioutil"
	"net/http"
	_ "regexp"
	"strconv"
	"strings"
)

const T = 20 //每一页的主评论数
const O = 10 //每页子评论数
const (
	USERNAMES = "root"
	PASSWORDS = "root"
	HOSTS     = "127.0.0.1"
	PORTS     = "3306"
	DBNAMES   = "gowork2"
)

var Db *sql.DB

type Messages struct {
	Message string `json:"message"`
}

/*编码转换函数，但是不符合
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}*/

// 请求头
func request(url string) (req *http.Request) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("req err:", err)
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("authority", "api.bilibili.com")
	req.Header.Set("Referer", "https://www.bilibili.com/video/BV12341117rG/?p=2&spm_id_from=pageDriver&vd_source=b8fc9305d5e9308e757e10418b8886dc")
	//fmt.Println("请求成功")
	return
}

// 内容
func getcontent(url string) (result string, err error) {
	var client http.Client
	//var tmp string
	req := request(url)
	resp, err1 := client.Do(req)
	if err != nil {
		err = err1
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}

		result += string(buf[:n])
	}
	//result = ConvertToString(tmp, "gbk", "utf-8")
	return
}

// 获取rpid以构造子评论url
func getrpid(data string) (list []string) {
	list = make([]string, 0)
	for i := 0; i < T; i++ {
		path := "data.replies." + strconv.Itoa(i) + ".rpid"
		rpid := gjson.Get(data, path)
		list = append(list, rpid.String())
	}
	return list

}

// 取主评论
func spidermain(i int) {
	Url := "https://api.bilibili.com/x/v2/reply/main?type=1&oid=420981979&next=" + strconv.Itoa(i)
	req, err := getcontent(Url)
	if err != nil {
		fmt.Println(err)
	}
	var m Messages
	for i := 0; i < T; i++ {
		path := "data.replies." + strconv.Itoa(i) + ".content.message"
		tmp := gjson.Get(req, path)
		message := tmp.String()
		m.Message = message
		Insertdata(m)
	}
	list := getrpid(req)
	for _, f := range list {
		comments_num(f)
	}
}

// 爬取子评论
func spiderone(url string) {
	data, err := getcontent(url)
	if err != nil {
		fmt.Println(err)
	}
	var m Messages
	for i := 0; i < O; i++ {
		path := "data.replies." + strconv.Itoa(i) + ".content.message"
		tmp := gjson.Get(data, path)
		message := tmp.String()
		m.Message = message
		Insertdata(m)
	}
}

// 求子评论页数
func comments_num(f string) {
	URL := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + f + "&ps=10&pn=1&web_location=333.788"
	data, err := getcontent(URL)
	if err != nil {
		fmt.Println(err)
	}
	paths := "data.page.count"
	tmp := gjson.Get(data, paths)
	num := int(tmp.Int() / O)  //子评论数
	nums := int(tmp.Int() % O) //剩余子评论数
	var flags bool
	if nums > 0 {
		flags = true
	} else {
		flags = false
	}
	for i := 0; i < num; i++ {
		url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + f + "&ps=10&pn=" + strconv.Itoa(i) + "&web_location=333.788"
		spiderone(url)
	}
	if flags {
		url := "https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=" + f + "&ps=10&pn=" + strconv.Itoa(num) + "&web_location=333.788"
		datas, err := getcontent(url)
		if err != nil {
			fmt.Println(err)
		}
		var m Messages
		for i := 0; i < nums; i++ {
			path := "data.replies." + strconv.Itoa(i) + ".content.message"
			tmps := gjson.Get(datas, path)
			message := tmps.String()
			m.Message = message
			Insertdata(m)
		}
	}
}

func InitDb() {
	//root:root@tcp(127.0.0.1:3306)/
	path := strings.Join([]string{USERNAMES, ":", PASSWORDS, "@tcp(", HOSTS, ":", PORTS, ")/", DBNAMES, "?charset=utf8"}, "")
	Db, _ = sql.Open("mysql", path)
	Db.SetConnMaxLifetime(10)
	Db.SetMaxIdleConns(5)
	if err := Db.Ping(); err != nil {
		fmt.Println("open database fail", err)
		return
	}
	fmt.Println("connect success")
}

func Insertdata(message Messages) bool {
	tx, err := Db.Begin()
	if err != nil {
		fmt.Println("begin err", err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO gowork (`message`) VALUES (?)")
	if err != nil {
		fmt.Println("prepare err", err)
		return false
	}
	_, err = stmt.Exec(message.Message)
	if err != nil {
		fmt.Println("exec fail", err)
		return false
	}
	_ = tx.Commit()
	return true
}

func dowork(start, end, pages int, flag bool) {

	for i := start; i < end; i++ {
		spidermain(i)
		fmt.Printf("第%d页爬取成功", i)
	}
	if flag {
		Url := "https://api.bilibili.com/x/v2/reply/main?type=1&oid=420981979&next=" + strconv.Itoa(end)
		req, err := getcontent(Url)
		if err != nil {
			fmt.Println(err)
		}
		var m Messages
		for i := 0; i < pages; i++ {
			path := "data.replies." + strconv.Itoa(i) + ".content.message"
			tmp := gjson.Get(req, path)
			message := tmp.String()
			m.Message = message
			Insertdata(m)
		}
	}
}

// 主评论url
// https://api.bilibili.com/x/v2/reply/main?type=1&oid=420981979&next=1
// https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=5547452670&ps=10&pn=1&web_location=333.788
// https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=104560010768&ps=10&pn=1&web_location=333.788
func main() {
	InitDb()

	URL := "https://api.bilibili.com/x/v2/reply/wbi/main?oid=420981979&type=1&mode=3&pagination_str=%7B%22offset%22:%22%22%7D&plat=1&seek_rpid=&web_location=1315875&w_rid=f9fafcb850df51a2abc9fe5df6c66ccc&wts=1696684049"
	res, err := getcontent(URL)
	if err != nil {
		fmt.Println(err)
	}
	path := "data.cursor.all_count"
	num := gjson.Get(res, path)
	page := num.Int() / T       //评论页数
	pages := int(num.Int() % T) //剩下评论
	var flag bool
	if pages > 0 {
		flag = true
	} else {
		flag = false
	}
	STARTS := 0
	ENDS := int(page)
	fmt.Println(res)
	dowork(STARTS, ENDS, pages, flag)
}
