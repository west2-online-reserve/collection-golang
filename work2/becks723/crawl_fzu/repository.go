package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"
	PASSWORD = "123456"
	IP       = "127.0.0.1"
	PORT     = "3306"
	DBNAME   = "fzu"
	TBNAME   = "notice"
)

/* 初始化数据库。建库建表 */
func initDB() (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4",
		USERNAME, PASSWORD, IP, PORT)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}

	// ping
	if err = db.Ping(); err != nil {
		return
	}

	// 创建数据库
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", DBNAME)
	_, err = db.Exec(query)
	if err != nil {
		return
	}

	// 重新连接至刚刚创建的数据库。sql的 USE DBNAME 语句有坑，不要用。
	_ = db.Close()
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		USERNAME, PASSWORD, IP, PORT, DBNAME)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}

	// 建表
	query = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
    	Title varchar(128) COMMENT '标题',
    	Date  varchar(10)  COMMENT '日期',
    	Author varchar(20) COMMENT '作者',
    	Clicks int NOT NULL DEFAULT 0 COMMENT '点击数',
    	Body MEDIUMTEXT COMMENT '正文'
) COMMENT '通知文件'`, TBNAME) // 经统计，正文最大长度为71352，超过TEXT的65535byte，因此使选用mediumtext（2^24-1byte）
	_, err = db.Exec(query)
	if err != nil {
		return
	}

	return
}

/* 插入文章数据 */
func insertArticles(db *sql.DB, articles []article) (err error) {
	// 新建事务
	tx, err := db.Begin()
	if err != nil {
		return
	}

	// 出错回滚
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := db.Prepare("INSERT INTO " + TBNAME + "(Title, Date, Author, Clicks, Body) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()

	// 写入数据
	for _, a := range articles {
		_, err = stmt.Exec(a.Title, a.Date, a.Author, a.Clicks, a.Body)
		if err != nil {
			return
		}
	}

	return
}
