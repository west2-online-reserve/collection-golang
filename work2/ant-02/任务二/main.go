package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

type Comment struct {
	UserName string `json:"uname"`
	Sex      string `json:"sex"`
	Message  string `json:"message"`
	Time     string `json:"time"`
}

func getHeader() map[string]string {
	cookie, err := os.ReadFile("cookie.txt")
	if err != nil {
		return nil
	}
	return map[string]string{
		"Cookie":     string(cookie),
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0",
	}
}

func getOid(client *http.Client, bv string) string {
	u := fmt.Sprintf("https://www.bilibili.com/video/%s/?p=14&spm_id_from=pageDriver&vd_source=cd6ee6b033cd2da64359bad72619ca8a", bv)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal(err)
	}

	params := getHeader()
	for k, v := range params {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	htmlContent := string(body)

	// 使用正则表达式提取 OID
	// 尝试多种模式，提高匹配成功率
	patterns := []string{
		fmt.Sprintf(`"aid":(\d+),.*?"bvid":"%s"`, bv), // 模式1: 结合bvid匹配
		`"aid"\s*:\s*(\d+)`,                           // 模式2: 标准格式
		`"aid":\s*(\d+)`,                              // 模式3: 简化格式
		`window\.__INITIAL_STATE__.*?"aid":(\d+)`,     // 模式4: 从初始化状态中提取
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(htmlContent)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func md5Hash(code string) string {
	hash := md5.Sum([]byte(code))
	return hex.EncodeToString(hash[:])
}

func URLQuote(s string, safe string) string {
	// 先进行标准 URL 编码
	encoded := url.QueryEscape(s)

	// 将安全字符的编码还原
	for _, char := range safe {
		encodedChar := url.QueryEscape(string(char))
		encoded = strings.ReplaceAll(encoded, encodedChar, string(char))
	}

	return encoded
}

func parseUrl(u string, client *http.Client) map[string]interface{} {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal(err)
	}

	params := getHeader()
	for k, v := range params {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("请求失败:", err)
	}
	defer resp.Body.Close()

	var j map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		log.Fatal("解析失败:", err)
	}

	return j
}

func getBody(oid, pageID string, client *http.Client) map[string]interface{} {
	var u string
	if pageID == "" {
		pagination_str := `{"offset":""}`
		code := fmt.Sprintf("mode=2&oid=%s&pagination_str=%s&plat=1&seek_rpid=&type=1&web_location=1315875&wts=%d", oid, url.QueryEscape(pagination_str), time.Now().Unix()) + "ea1db124af3c7062474693fa704f4ff8"
		w_rid := md5Hash(code)
		u = fmt.Sprintf("https://api.bilibili.com/x/v2/reply/wbi/main?oid=%s&type=1&mode=2&pagination_str=%s&plat=1&seek_rpid=&web_location=1315875&w_rid=%s&wts=%d", oid, URLQuote(pagination_str, ":"), w_rid, time.Now().Unix())

	} else {
		pagination_str := fmt.Sprintf(`{"offset":"%s"}`, pageID)
		code := fmt.Sprintf("mode=2&oid=%s&pagination_str=%s&plat=1&type=1&web_location=1315875&wts=%d", oid, url.QueryEscape(pagination_str), time.Now().Unix()) + "ea1db124af3c7062474693fa704f4ff8"
		w_rid := md5Hash(code)
		u = fmt.Sprintf("https://api.bilibili.com/x/v2/reply/wbi/main?oid=%s&type=1&mode=2&pagination_str=%s&plat=1&web_location=1315875&w_rid=%s&wts=%d", oid, URLQuote(pagination_str, ":"), w_rid, time.Now().Unix())

	}
	return parseUrl(u, client)
}

func extractFirstNumber(s string) (int, error) {
	// 编译正则表达式
	re := regexp.MustCompile(`\d+`)

	// 查找所有匹配
	matches := re.FindAllString(s, -1)

	// 检查是否找到匹配
	if len(matches) == 0 {
		return 0, fmt.Errorf("未找到数字")
	}

	// 转换第一个匹配为整数
	return strconv.Atoi(matches[0])
}

func parseComments(reply map[string]interface{}) {
	comment := &Comment{}
	if member, ok := reply["member"]; ok {
		if uname, ok := member.(map[string]interface{})["uname"]; ok {
			comment.UserName = uname.(string)
		}
		if sex, ok := member.(map[string]interface{})["sex"]; ok {
			comment.Sex = sex.(string)
		}
	}
	if content, ok := reply["content"]; ok {
		if message, ok := content.(map[string]interface{})["message"]; ok {
			comment.Message = message.(string)
		}
	}
	if ctime, ok := reply["ctime"]; ok {
		comment.Time = time.Unix(int64(ctime.(float64)), 0).Format("2006-01-02 15:04:05")
	}
	addComment(comment)
}

func getComments(oid string, replies []interface{}, client *http.Client) {
	for _, reply := range replies {
		parseComments(reply.(map[string]interface{}))

		if n_replies, ok := reply.(map[string]interface{})["replies"]; ok {
			if n_replies != nil {
				for _, n_reply := range n_replies.([]interface{}) {
					parseComments(n_reply.(map[string]interface{}))
				}
			}
		}

		if rpid, ok := reply.(map[string]interface{})["rpid_str"]; ok {
			if sub_reply_entry_text, ok := reply.(map[string]interface{})["reply_control"].(map[string]interface{})["sub_reply_entry_text"]; ok {
				rereply, err := extractFirstNumber(sub_reply_entry_text.(string))
				if err != nil {
					for i := 1; i <= rereply; i++ {
						u := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=%s&type=1&root=%s&ps=10&pn=%d&web_location=333.788", oid, rpid.(string), i)
						if data, ok := parseUrl(u, client)["data"]; ok {
							if replies, ok := data.(map[string]interface{})["replies"]; ok {
								for _, reply := range replies.([]interface{}) {
									parseComments(reply.(map[string]interface{}))
								}
							}
						}
					}
				}
			}
		}
	}
}

func initDB() {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
		return
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&Comment{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
		return
	}
}

func addComment(comment *Comment) {
	if err := db.Create(comment).Error; err != nil {
		log.Fatalf("failed to add ni: %v\nerr: %v", comment, err)
		return
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initDB()

	bv := "BV12341117rG"

	client := &http.Client{}

	oid := getOid(client, bv)

	pageID := ""

	for pageID != "0" {
		if data, ok :=  getBody(oid, pageID, client)["data"]; ok {
			if cursor, ok := data.(map[string]interface{})["cursor"]; ok {
				if pagination_reply, ok := cursor.(map[string]interface{})["pagination_reply"]; ok {
					if next_offset, ok := pagination_reply.(map[string]interface{})["next_offset"]; ok {
						pageID = next_offset.(string)
					}
				}
			}

			if replies, ok := data.(map[string]interface{})["replies"]; ok {
				if replies != nil {
					getComments(oid, replies.([]interface{}), client)
				}
			}
		}
	}

}
