package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Reply struct {
	Member struct {
		Uname string `json:"uname"`
		Sex   string `json:"sex"`
	} `json:"member"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
}

type Response struct {
	Data struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Count int `json:"count"`
		} `json:"page"`
		Replies []Reply `json:"replies"`
	} `json:"data"`
}

const (
	username = "root"
	password = "zth20041017"
	hostname = "localhost"
	port     = "3306"
	dbname   = "comments"
)

func dsn(dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, hostname, port, dbname)
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func insertComments(db *sql.DB, comment Reply) {
	query := "INSERT INTO comments (uname, sex, message) VALUES (?, ?, ?)"
	_, err := db.Exec(query, comment.Member.Uname, comment.Member.Sex, comment.Content.Message)
	if err != nil {
		log.Printf("插入数据失败: %v\n", err)
	} else {
		log.Printf("成功插入评论: %s\n", comment.Content.Message)
	}
}

func fetchComments(db *sql.DB, page int, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply?oid=420981979&type=1&pn=%d&mode=2", page)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0")
	req.Header.Set("Referer", "https://www.bilibili.com/")
	req.Header.Set("Origin", "https://www.bilibili.com")
	cookies := []*http.Cookie{
		{Name: "buvid3", Value: "0AA097D9-CE06-711B-6546-566F9657BEEE53932infoc"},
		{Name: "b_nut", Value: "1729601953"},
		{Name: "_uuid", Value: "88E97F88-FEA1-108D10-64BA-2F19B5ED5E51054904infoc"},
		{Name: "enable_web_push", Value: "DISABLE"},
		{Name: "home_feed_column", Value: "5"},
		{Name: "browser_resolution", Value: "1528-704"},
		{Name: "buvid4", Value: "390FF2B0-5DEA-C21D-FEF1-2D5CE041F80455779-024102212-OGsZd1CnuZtG2YoUvGR4ulWf9oWinAayNX0scwHjjpbS3Gi351w3o1oiyeGE7xun"},
		{Name: "CURRENT_FNVAL", Value: "4048"},
		{Name: "rpdid", Value: "0zbfVHb7Yk|LyUItId4|4gf|3w1T3etU"},
		{Name: "DedeUserID", Value: "365211151"},
		{Name: "DedeUserID__ckMd5", Value: "5557f4d7d79c2fa1"},
		{Name: "CURRENT_QUALITY", Value: "80"},
		{Name: "header_theme_version", Value: "CLOSE"},
		{Name: "bili_ticket", Value: "eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzExMzA5NTcsImlhdCI6MTczMDg3MTY5NywicGx0IjotMX0.zrQs_tZp2nEr-MiuVbSUHZwsbduBVXB9Rd4qGnx_7Bc"},
		{Name: "bili_ticket_expires", Value: "1731130897"},
		{Name: "SESSDATA", Value: "755d4bb3%2C1746423761%2Cad2cd%2Ab2CjAMFwR-PIMxa_vUIY7oTOE4iD1UkIkytbE9tu4HvMPEexz7g_lcAV-qtDhYaqHEqTwSVmdkenJhTnBaWmpJREpUMkphc3Qtd3F0c2RXb1RGbTg2dGFGa3hGWm5XeEpvbmdINnNqQXhQREc2WWVTX0FwMno0NzNFN09TSlBxMGFpcVhkOGlXblh3IIEC"},
		{Name: "bili_ticket", Value: "eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzExMzA5NTcsImlhdCI6MTczMDg3MTY5NywicGx0IjotMX0.zrQs_tZp2nEr-MiuVbSUHZwsbduBVXB9Rd4qGnx_7Bc"},
		{Name: "bp_t_offset_365211151", Value: "997120402008309760"},
		{Name: "fingerprint", Value: "db0a8bfee79424904a6de4cbcd7ac77a"},
		{Name: "buvid_fp_plain", Value: "undefined"},
		{Name: "buvid_fp", Value: "db0a8bfee79424904a6de4cbcd7ac77a"},
		{Name: "b_lsid", Value: "AF3D4126_1930C473B22"},
		{Name: "sid", Value: "o7mfo1bq"},
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("请求失败,HTTP 状态码: %d,页数：%d\n", resp.StatusCode, page)
		return
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	for _, reply := range response.Data.Replies {
		insertComments(db, reply)
	}
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
	uname VARCHAR(255),
	sex VARCHAR(255),
	message TEXT
	);`)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)
	for i := 1; i <= 455; i++ {
		wg.Add(1)
		go fetchComments(db, i, &wg, semaphore)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
}
