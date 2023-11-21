package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server-redis/cfg"
	"server-redis/datastruct"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// 返回DSN
func DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Username, cfg.Password, cfg.Hostname, cfg.Databasename)
}

// 连接到sql数据库
func ConnectDataBase() error {
	var err error
	db, err = sql.Open("mysql", DSN())
	if err != nil {
		log.Printf("[Error] MySQL: %s", err.Error())
		return err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("[Error] MySQL: %s", err.Error())
		return err
	}
	return nil
}

// 初始化数据表
func MySQLAccountInit() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `todolist_account`(`username` VARCHAR(64) NOT NULL,`password` VARCHAR(64) NOT NULL,PRIMARY KEY ( `username` ))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		log.Printf("[Error] MySQL:%s", err.Error())
		return errors.New("failed to initialize database")
	}
	if _, err = stmt.Exec(); err != nil {
		log.Printf("[Error] MySQL:%s", err.Error())
		return errors.New("failed to initialize database")
	}
	log.Print("[INFO] MySQL: DataTable initialized successfully\n")
	return nil
}

func MySQLAccountSyncPack() ([]datastruct.User, error) {
	userSlice := make([]datastruct.User, 0)
	const SearchCmd = "select * from todolist_account"
	rows, err := db.Query(SearchCmd)
	if err != nil {
		log.Printf("[Error] MySQL(in account searching): %v", err)
		return nil, err
	}

	for rows.Next() {
		var existUser datastruct.User
		rows.Scan(&existUser.Username, &existUser.Password)
		userSlice = append(userSlice, existUser)
	}
	return userSlice, nil
}

func MySQLAccountCreate(username, password string) error {
	const insertCmd = "insert into todolist_account(username,password)values (?,?)"
	if _, err := db.Exec(insertCmd, username, password); err != nil {
		return err
	}

	return nil
}

func MySQLTodolistInit() error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `todolist`(`id` int unsigned AUTO_INCREMENT,`title` varchar(64) not null,`owner` varchar(64) not null,`text` text not null,`deadline` datetime not null,`addtime` datetime not null,`status` boolean default false, primary key(`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;")
	if err != nil {
		log.Printf("[Error] MySQL:%s", err.Error())
		return errors.New("failed to initialize datatable")
	}
	if _, err = stmt.Exec(); err != nil {
		log.Printf("[Error] MySQL:%s", err.Error())
		return errors.New("failed to initialize datatable")
	}
	log.Print("[INFO] MySQL: DataTable initialized successfully\n")
	return nil
}

func MySQLTodolistInsert(what2do datastruct.TodolistBindJSONReceive) (int64, string, error) {
	if len(what2do.Text) == 0 {
		return 0, "", errors.New("text cannot be empty")
	}

	insertData := datastruct.ProcessStruct{
		Title:    what2do.Title,
		Owner:    what2do.Owner,
		Text:     what2do.Text,
		Addtime:  time.Unix(time.Now().Local().Unix(), 0).Format("2006-01-02 15:04:05"),
		Deadline: what2do.Deadline,
		Isdone:   false,
	}

	t1, _ := time.Parse("2006-01-02 15:04:05", insertData.Addtime)
	t2, _ := time.Parse("2006-01-02 15:04:05", insertData.Deadline)
	if t1.After(t2) {
		return 0, "", errors.New("add time is after deadline")
	}

	insertCmd := "insert into todolist(title,owner,text,deadline,addtime)values (?,?,?,?,?)"

	result, err := db.Exec(insertCmd,
		insertData.Title,
		insertData.Owner,
		insertData.Text,
		insertData.Deadline,
		insertData.Addtime)
	if err != nil {
		return 0, "", err
	}
	id, _ := result.LastInsertId()
	return id, insertData.Addtime, nil
}

func MySQLTodoListSyncPack(username string) ([]datastruct.TodolistBindJSONSend, error) {
	dataSlice := make([]datastruct.TodolistBindJSONSend, 0)
	const SearchCmd = "select * from todolist where owner = ?"
	rows, err := db.Query(SearchCmd, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var data datastruct.TodolistBindJSONSend
		rows.Scan(
			&data.Id,
			&data.Title,
			&data.Owner,
			&data.Text,
			&data.Deadline,
			&data.Addtime,
			&data.Status,
		)
		dataSlice = append(dataSlice, data)
	}
	return dataSlice, nil
}

func MySQLTodoListDelete(username string, depend datastruct.TodolistBindRedisCondition) error {
	var err error
	if depend.Method&cfg.DeleteAll == cfg.DeleteAll {
		if _, err = db.Query("delete from todolist where owner=?", username); err != nil {
			return err
		}
	} else {
		deleteCmd := fmt.Sprintf("delete from todolist where owner='%s'", username)
		idSlice := make([]int64, len(depend.Idlist))
		if depend.Method&cfg.DeleteWithId == cfg.DeleteWithId {
			deleteCmd = deleteCmd + " and id=?"
			copy(idSlice, depend.Idlist)
		}

		if depend.Method&cfg.DeleteWithStatus == cfg.DeleteWithStatus {
			deleteCmd = deleteCmd + " and status=?"
		}

		switch depend.Method {
		case cfg.DeleteWithId:
			for index := range idSlice {
				if _, err = db.Query(deleteCmd, idSlice[index]); err != nil {
					return err
				}
			}
		case cfg.DeleteWithStatus:
			if _, err = db.Query(deleteCmd, depend.Status); err != nil {
				return err
			}
		case cfg.DeleteWithId | cfg.DeleteWithStatus:
			for index := range idSlice {
				if _, err = db.Query(deleteCmd, idSlice[index], depend.Status); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func MySQLTodoListModify(username string, depend datastruct.TodolistBindRedisUpdate) error {
	var err error
	if depend.Method&cfg.ModifyAll == cfg.ModifyAll {
		if _, err = db.Query("update todolist set status = ? where owner = ?", depend.NewStatus, username); err != nil {
			return err
		}
	} else {
		idSlice := make([]int64, len(depend.Idlist))
		copy(idSlice, depend.Idlist)
		for index := range idSlice {
			if _, err = db.Query("update todolist set status = ? where owner = ? and id = ?", depend.NewStatus, username, idSlice[index]); err != nil {
				return err
			}
		}
	}
	return nil
}
