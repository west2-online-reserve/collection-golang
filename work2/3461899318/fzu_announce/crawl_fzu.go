package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql" //这里要加下划线
	"golang.org/x/net/publicsuffix"
)

// 配置信息
const (
	ssoLoginURL   = "https://sso.fzu.edu.cn/login"                                                                                                                                  // 正确的SSO登录接口
	targetService = "https://info22.fzu.edu.cn/ssjg_0509.jsp?wbtreeid=1460&order=0&zzjg=0&pageSize001=100&currentnum=69&_vt=&_vn=&_vk=&_vm=&_va=&_vs=&_vc=&_vb=&_ve=&searchScope=1" // 登录后要访问的目标服务
	username      = "102400432"                                                                                                                                                     // 用户名（学号）
	password      = "zhangqi20060212"
	public_head   = "https://info22.fzu.edu.cn/" // 密码
)

var (
	data = make([]announceInfo, 0)
	db   *sql.DB
)

type announceInfo struct {
	date     string
	author   string
	title    string
	text     string
	clickNum string
}

// 起始网址：https://info22.fzu.edu.cn/ssjg_0509.jsp?wbtreeid=1460&order=0&zzjg=0&pageSize001=100&currentnum=94&_vt=&_vn=&_vk=&_vm=&_va=&_vs=&_vc=&_vb=&_ve=&searchScope=1
// 最后一页：https://info22.fzu.edu.cn/ssjg_0509.jsp?wbtreeid=1460&order=0&zzjg=0&pageSize001=100&currentnum=69&_vt=&_vn=&_vk=&_vm=&_va=&_vs=&_vc=&_vb=&_ve=&searchScope=1
// 初始化带Cookie的客户端（保持会话）
func createClient() *http.Client {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		fmt.Printf("创建Cookie容器失败: %v\n", err)
		return nil
	}
	// 允许自动跟随重定向（登录成功后会跳转到targetService）
	return &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("重定向到: %s\n", req.URL) // 调试用，查看重定向路径
			return nil
		},
	}
	//https: //info22.fzu.edu.cn/ssjg_0509.jsp?wbtreeid=1460
	//https://info22.fzu.edu.cn/content.jsp?urltype=news.NewsContentUrl&wbtreeid=1355&wbnewsid=39427
}

// 登录SSO平台
func login(client *http.Client) bool {
	// 构建登录参数（使用CAS协议的service参数）
	formData := url.Values{
		"username": {username},      // 用户名（学号）
		"password": {password},      // 密码（修正顺序）
		"service":  {targetService}, // 目标服务地址（登录后跳转）
	}

	// 创建POST请求，发送到SSO登录接口
	req, err := http.NewRequest("POST", ssoLoginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Printf("创建登录请求失败: %v\n", err)
		return false
	}

	// 设置请求头（模拟浏览器）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 表单提交必须
	req.Header.Set("Referer", ssoLoginURL)                              // 模拟从SSO登录页发起请求
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// 发送登录请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("登录请求失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// 验证登录结果（200=登录页跳转完成，302=正在重定向，均可能表示成功）
	fmt.Printf("登录响应状态码: %d\n", resp.StatusCode)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound {
		return true
	}
	return false
}

// 跳转到列表页并获取列表页所有通知页面的url（登录后）
func spider(client *http.Client, url string) (string, bool) {
	// 用已登录的客户端访问目标服务（自动携带Cookie）
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("访问目标页面失败: %v\n", err)
		return "", true
	}
	defer resp.Body.Close()

	fmt.Printf("目标页面响应状态码: %d\n", resp.StatusCode)
	//读取页面内容（验证是否登录成功）
	content := getRECALLinfo(resp)
	//检查是否成功读取HTML
	//fmt.Println(content)
	//通知页url的正则
	reg1 := regexp.MustCompile("<a href=\"(.*?)\" target=\"_blank\">")
	////(?s:(.*?))这个正则匹配换行符
	link_url1 := reg1.FindAllStringSubmatch(content, -1)
	//确实是100个url
	//	fmt.Printf("%d", len(link_url))
	num := len(link_url1) //101

	for i := 0; i < num-1; i++ {
		announceUrl := public_head + link_url1[i][1]
		webDate := get_info(client, announceUrl)
		data = append(data, webDate)
		//fmt.Println(announceUrl)
	}
	//选择器匹配
	//document, err := goquery.NewDocumentFromReader(tem.Body)
	//if err != nil {
	//	fmt.Println("解析失败")
	//}
	//s := document.Find("body > div.wa1200w > div.conth > table > tbody > tr > td > a:nth-child(3)")
	//nextUrl, check := s.Attr("href")
	//if check {
	//	fmt.Println(check)
	//}
	//fmt.Println(nextUrl)
	//return nextUrl
	check := false
	s := strings.Split(data[len(data)-1].date, "-")
	if s[0] == "2019" {
		check = true
	}
	//下一页url的正则(用于自动翻页)
	reg2 := regexp.MustCompile("<a href=\"(.*?)\">下一页</a>")
	link_url2 := reg2.FindAllStringSubmatch(content, -1)
	///ssjg_0509.jsp?wbtreeid=1460&order=0&zzjg=0&pageSize001=100&currentnum=95&_vt=&_vn=&_vk=&_vm=&_va=&_vs=&_vc=&_vb=&_ve=&searchScope=1
	//检查使用
	//fmt.Println("link_url2:", link_url2[0][1])
	//<a href="content.jsp?urltype=news.NewsContentUrl&wbtreeid=1294&wbnewsid=5025" target="_blank">
	//https://info22.fzu.edu.cn/content.jsp?urltype=news.NewsContentUrl&wbtreeid=1305&wbnewsid=23604 =
	//(?s:(.*?))\
	return public_head + link_url2[0][1], check
}

// 在通知页面提取信息
func get_info(client *http.Client, url string) announceInfo {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("响应失败")
	}
	//content := getRECALLinfo(resp)
	//用于检查是否成功爬取网址源码
	//fmt.Println(content)
	info := announceInfo{}
	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	info.text = document.Find("#vsb_content > div").Text()
	info.title = document.Find("body > div.wa1200w > div.conth > form > div.conth1").Text()
	//fmt.Printf("标题：%s\n", info.title)
	//fmt.Printf("正文%s\n", info.text)
	elseData := document.Find("body > div.wa1200w > div.conth > form > div.conthsj").Text()
	//fmt.Println(elseData)
	datePrefix := "日期： "
	infoPrefix := "信息来源： "
	clickPrefix := "点击数:"
	// 找到日期前缀的位置
	dateStart := strings.Index(elseData, datePrefix)
	if dateStart == -1 {
		fmt.Println("未找到日期前缀")
		return info
	}
	// 日期内容从日期前缀之后开始，到信息来源前缀之前结束
	dateEnd := strings.Index(elseData, infoPrefix)
	if dateEnd == -1 {
		fmt.Println("未找到信息来源前缀")
		return info
	}
	dateStr := elseData[dateStart+len(datePrefix) : dateEnd]
	// 去除前后空白字符（包括空格、特殊空格等）
	info.date = strings.TrimSpace(dateStr)

	// 提取信息来源
	infoStart := strings.Index(elseData, infoPrefix)
	if infoStart == -1 {
		fmt.Println("未找到信息来源前缀")
		return info
	}
	// 信息来源内容从信息来源前缀之后开始，到点击数前缀之前结束
	infoEnd := strings.Index(elseData, clickPrefix)
	if infoEnd == -1 {
		fmt.Println("未找到点击数前缀")
		return info
	}
	infoStr := elseData[infoStart+len(infoPrefix) : infoEnd]
	info.author = strings.TrimSpace(infoStr)
	//提取动态数据属性的特征值
	staFeature := "1768654345, "
	endFeature := ")"
	staIndex := strings.Index(elseData, staFeature)
	if staIndex == -1 {
		fmt.Printf("未找到特征前缀\n")
		return info
	}
	endIndex := strings.Index(elseData, endFeature)
	if endIndex == -1 {
		fmt.Printf("未找到特征后缀")
		return info
	}
	feature := elseData[staIndex+len(staFeature) : endIndex]
	//fmt.Printf("feature: %s\n", feature)
	// 输出结果(测试使用)
	//fmt.Printf("日期：%s\n", info.date)
	//fmt.Printf("信息来源：%s\n", info.author)
	staUrl := "https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid="
	endUrl := "&owner=1768654345&clicktype=wbnews"
	click, err := client.Get(staUrl + feature + endUrl)
	if err != nil {
		fmt.Println("获取点击率失败")
	}
	defer click.Body.Close()
	//fmt.Println(staUrl + feature + endUrl)
	doc, err := goquery.NewDocumentFromReader(click.Body)
	if err != nil {
		fmt.Println("响应体解析失败")
	}
	info.clickNum = doc.Find("body").Text()
	//fmt.Println(info.clickNum)
	fmt.Println("一则通知爬取成功")
	return info
}

// 获得html文本
func getRECALLinfo(response *http.Response) string {
	buf := make([]byte, 4096)
	buf, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("信息读取失败")
	}
	return string(buf)
}

// 连接数据库
func loginMysql() {
	dsn := "root:zhangqi20060212@tcp(127.0.0.1:3306)/fzu"
	db, _ = sql.Open("mysql", dsn)
	err := db.Ping()
	if err != nil {
		fmt.Println("没得ping通数据库")
	}
	fmt.Println("连接数据库成功(●'◡'●)")
}

// 插入数据
func insertTable(info *announceInfo) {
	sql := "insert into announcement(title,publishTime,mainContent,author,clickNum) values(?,?,?,?,?)"
	click, err := strconv.ParseInt(info.clickNum, 10, 64)
	//提前对
	result, err := db.Exec(sql, info.title, info.date, info.text, info.author, click)
	if err != nil {
		fmt.Println("数据插入失败！！！")
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Println("获取数据失败")
	}
}

// 爬取一个列表页的所有
func main() {
	sta := time.Now()
	// 1. 创建客户端
	client := createClient()
	if client == nil {
		return
	}

	// 2. 执行登录
	loginSuccess := login(client)
	fmt.Printf("登录结果: %v\n", loginSuccess)
	if !loginSuccess {
		return
	}
	curUrl := targetService
	// 3. 访问目标服务页面
	for {
		nextUrl, check := spider(client, curUrl)
		curUrl = nextUrl
		if check {
			break
		}
	}
	//fmt.Println(data)
	fmt.Println("爬取完成")
	gap := time.Since(sta)
	fmt.Printf("单线程版时间：%v\n", gap) //33m1.6689406s
	num := len(data)
	loginMysql()
	for i := 0; i < num; i++ {
		insertTable(&data[i])
	}
	fmt.Println("数据传输完成")
}
