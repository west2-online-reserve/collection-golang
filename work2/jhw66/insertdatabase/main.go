package Insert

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Item struct {
	Author  string
	Title   string
	Href    string
	Atime   string
	Click   string
	Content string
}

func OpenDB() (*sql.DB, error) {
	// DSN（数据源名称）格式："用户名:密码@tcp(主机:端口)/数据库名
	DSN := "wjhaccount:Wujiahui789@tcp(127.0.0.1:3306)/CRAWL"
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %v", err)
	}
	// 确保连接数据库是否正常
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}
	return db, nil
}
func Insertintotable(db *sql.DB, items []Item) error {
	// 准备一个 INSERT 语句
	stmt, err := db.Prepare(`INSERT INTO articles (author, title, href, publish_time, click, content) 
                              VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("准备插入语句失败: %v", err)
	}
	defer stmt.Close()

	// 批量插入
	for _, item := range items {
		fmt.Println("inserting")
		result, err := stmt.Exec(item.Author, item.Title, item.Href, item.Atime, item.Click, item.Content)
		if err != nil {
			return fmt.Errorf("插入数据失败: %v", err)
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取插入ID失败: %v", err)
		}
		fmt.Printf("插入成功，ID: %d\n", lastInsertID) // 输出插入的ID，确认插入成功
	}

	return nil
}
