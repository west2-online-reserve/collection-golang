package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver   = "mysql"
	dbUser     = "root"
	dbPassword = "123456"
	dbName     = "bilibili"
)

// 实现get函数 爬取指定页面所有数据
// 参数：URL 返回：页面内容
func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		return "", err1
	}
	//关闭响应体，防止资源泄露
	defer resp.Body.Close()
	//存储二进制数据
	buf := make([]byte, 4096)
	//循环爬取整页数据
	for {
		//返回值：读取的字节数、错误
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil {
			if err2 != io.EOF {
				return "", err2
			}
			break
		}
		//字符串累加
		result += string(buf[:n])
	}
	return result, nil
}

// fetchArticleURLs 爬取网页中大概页面范围内的所有文章的URL并暂时存储起来
// 返回：范围内文章URL的数组，和一个错误信息
func fetchArticleURLs() ([]string, error) {
	var articleURLs []string
	for i := 224; i <= 341; i++ {
		url := fmt.Sprintf("https://info22-fzu-edu-cn.fzu.edu.cn/lm_list.jsp?totalpage=1023&PAGENUM=%d&wbtreeid=1460&urltype=tree.TreeTempUrl&#34", i)

		pageContent, err := HttpGet(url)
		if err != nil {
			fmt.Printf("Error fetching page %d: %s\n", i, err)
			continue
		}
		// 使用正则表达式提取文章URL
		// 注意：正则表达式需要根据实际页面结构调整
		ret := regexp.MustCompile(`<a href="([^"]+)" target=_blank title="`)
		matches := ret.FindAllStringSubmatch(pageContent, -1)
		for _, match := range matches {
			if len(match) > 1 {
				articleURLs = append(articleURLs, match[1])
			}
		}
	}
	return articleURLs, nil
}

// 提取文章中的时间，检验是否在时间范围内
// 参数：页面内容、时间戳范围
func isPageDateInRange(pageContent, datePattern string, startTimestamp, endTimestamp int64) (bool, error) {
	// 编译时间匹配的正则表达式
	re := regexp.MustCompile(datePattern)
	matches := re.FindStringSubmatch(pageContent)

	if len(matches) == 0 {
		return false, nil // 没有找到时间信息
	}

	// 将找到的时间字符串转换为 time.Time 对象
	dateStr := matches[0]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false, err // 时间格式转换失败
	}

	// 将时间转换为时间戳
	pageTimestamp := date.Unix()

	// 检查时间戳是否在范围内
	return startTimestamp <= pageTimestamp && pageTimestamp <= endTimestamp, nil
}

// 爬取时间范围内页面中的标题、作者、发布日期以及正文
// 参数：文章页面的URL 返回：提取的内容
func working(articleURL string) (map[string]string, error) {
	//定义时间范围
	startDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC)
	startTimestamp, endTimestamp := startDate.Unix(), endDate.Unix()

	//先调用HttpGet，获取文章的所有内容
	pageContent, err := HttpGet(articleURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching page content: %v", err)
	}

	//调用isPageDateInRange 检查是否在时间范围内
	inRange, err := isPageDateInRange(pageContent, "2006-01-02", startTimestamp, endTimestamp)
	if err != nil {
		return nil, fmt.Errorf("error checking page date: %v", err)
	}
	if !inRange {
		return nil, nil // 日期不在范围内，跳过
	}
	// 提取文章的标题、作者、发布日期和正文
	titleRegex := regexp.MustCompile(`<title>([^<]+)</title>`)
	authorRegex := regexp.MustCompile(`<h3 class="fl">([^<]+)</h3>`)
	dateRegex := regexp.MustCompile(`<div class="conthsj" >日期：([^<]+) &nbsp;`)
	contentRegex := regexp.MustCompile(`<META Name="description" Content="(.*?)" />`)

	title := titleRegex.FindStringSubmatch(pageContent)
	author := authorRegex.FindStringSubmatch(pageContent)
	date := dateRegex.FindStringSubmatch(pageContent)
	content := contentRegex.FindStringSubmatch(pageContent)

	extractedData := map[string]string{
		"Title":   title[1],
		"Author":  author[1],
		"Date":    date[1],
		"Content": content[1],
	}

	return extractedData, nil
}

// 将爬取到的数据存入数据库
func storeDataInDatabase(data map[string]string) {
	// 连接数据库
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(192.168.141.27:3306)/%s", dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表
	const createTableQuery = `CREATE TABLE IF NOT EXISTS FZUnews (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255),
		date DATE,
		writer VARCHAR(255),
		content LONGTEXT
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	// 插入数据
	insertQuery := "INSERT INTO FZUnews (title, date, writer, content) VALUES (?, ?, ?, ?)"

	_, err = db.Exec(insertQuery, data["Title"], data["Date"], data["Author"], data["Content"])
	if err != nil {
		log.Fatal(err)
	}

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data stored in database:", data)
}

func main() {
	//fzuUrl:="https://info22-fzu-edu-cn.fzu.edu.cn/lm_list.jsp?urltype=tree.TreeTempUrl&wbtreeid=1460&#34"
	//调用fetchArticleURLs，进行循环
	articleURLs, err := fetchArticleURLs()
	if err != nil {
		fmt.Printf("Error fetching article URLs: %s\n", err)
		return
	}

	for _, url := range articleURLs {
		extractedData, err := working(url)
		if err != nil {
			fmt.Printf("Error processing article URL %s: %s\n", url, err)
			continue
		}
		if extractedData != nil {
			storeDataInDatabase(extractedData)
		}
	}
}
