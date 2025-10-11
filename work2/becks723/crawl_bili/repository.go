package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USER     string = "root"
	PASSWORD string = "123456"
	IP       string = "127.0.0.1"
	PORT     string = "3306"
	DBNAME   string = "bili_comments"
	TBNAME   string = "BV12341117rG"
)

/*
  需要事先在mysql中创建 bili_comments 数据库，再创建两个表：

  user表：
create table user
(
    Uid  bigint       not null comment '用户uid'
        primary key,
    Name varchar(100) null comment '用户名',
    Sex  varchar(20)  null comment '用户性别'
);

  BV12341117rG表：
create table BV12341117rG
(
    Rpid     bigint             primary key comment '评论唯一id',
	Ctime   bigint              null comment '时间戳',
	`Like`  int    default 0    not null comment '点赞数',
	Message varchar(2000)       null comment '评论内容',
	Root    bigint default null null comment '属于哪个主评论。若为主评论，null',
	Uid     bigint              null comment '评论用户uid',
	constraint fk_uid
		foreign key (Uid) references user (Uid)
)
	comment 'BV12341117rG的评论列表';
*/

/* 初始化数据库连接 */
func initDB() (db *sql.DB, err error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		USER, PASSWORD, IP, PORT, DBNAME)
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}
	return
}

/* 保存评论数据 */
func insertComments(db *sql.DB, comments []mainComment) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	// 确保保存失败时回滚数据库
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 准备写入评论
	query := fmt.Sprintf("INSERT INTO %s (`Rpid`, `Ctime`, `Like`, `Message`, `Uid`, `Root`) VALUES (?, ?, ?, ?, ?, ?)", TBNAME)
	stmt, err := db.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	// 准备写入用户信息
	stmtUser, err := db.Prepare("INSERT IGNORE INTO user (`Uid`, `Name`, `Sex`) VALUES (?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtUser.Close()

	for _, mc := range comments {
		fmt.Printf("正在保存评论（rpid=%d）\n", mc.Rpid)
		_, err = stmtUser.Exec(mc.User.Uid, mc.User.Name, mc.User.Sex) // 保存用户信息。必须优先于评论保存，因为user表是主表，外键uid必须先存在
		if err != nil {
			return
		}
		_, err = stmt.Exec(mc.Rpid, mc.Ctime, mc.Like, mc.Message, mc.User.Uid, nil) // 保存主评论
		if err != nil {
			return
		}
		for _, sc := range mc.SubComments {
			fmt.Printf("正在保存子评论（rpid=%d）\n", sc.Rpid)
			_, err = stmtUser.Exec(sc.User.Uid, sc.User.Name, sc.User.Sex) // 保存用户信息
			if err != nil {
				return
			}
			_, err = stmt.Exec(sc.Rpid, sc.Ctime, sc.Like, sc.Message, sc.User.Uid, mc.Rpid) // 保存子评论
			if err != nil {
				return
			}
		}
	}

	_ = tx.Commit()
	return
}
