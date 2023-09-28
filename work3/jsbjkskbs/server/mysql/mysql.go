package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/cfg"
	"server/datastruct"
	encrypt2 "server/encrypt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// MySQL信息
	Username     = "root"           //username 		like 	`root`
	Password     = "357920"         //password 		like 	`123456`
	Hostname     = "127.0.0.1:3306" //hostname 		like 	`127.0.0.1:3306`
	Databasename = "test"           //databasename 	like 	`databasename`
)

var db *sql.DB

// 返回DSN
func DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", Username, Password, Hostname, Databasename)
}

// 连接到sql数据库,并返回数据库指针
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

func MySQLAccountSearch(username string) ([]datastruct.User, error) {
	userSlice := make([]datastruct.User, 0)
	const SearchCmd = "select * from todolist_account where username=?"
	rows, err := db.Query(SearchCmd, username)
	if err != nil {
		log.Printf("[Error] MySQL(in account searching): %v", err)
		return nil, err
	}

	for rows.Next() {
		var existUser datastruct.User
		rows.Scan(&existUser.UserName, &existUser.Password)
		userSlice = append(userSlice, existUser)
	}

	return userSlice, nil
}

func MySQLAccountCheck(username, shaPassword string) ([]datastruct.User, error) {
	userSlice := make([]datastruct.User, 0)
	const CheckCmd = "select * from todolist_account where username=? and password=?"
	rows, err := db.Query(CheckCmd, username, shaPassword)
	if err != nil {
		log.Printf("[Error] MySQL(in account checking): %v", err)
		return userSlice, err
	}

	for rows.Next() {
		var existUser datastruct.User
		rows.Scan(&existUser.UserName, &existUser.Password)
		userSlice = append(userSlice, existUser)
	}

	return userSlice, nil
}

func MySQLAccountCreate(username, password string) error {
	const insertCmd = "insert into todolist_account(username,password)values (?,?)"
	if _, err := db.Exec(insertCmd, username, encrypt2.SHA256(password)); err != nil {
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

func MySQLTodolistInsert(username string, what2do datastruct.TodolistBindJSONReceive) (int64, error) {
	if len(what2do.Text) == 0 {
		return 0, errors.New("text cannot be empty")
	}

	insertData := datastruct.TodolistBindMysqlInsert{
		JsonData: what2do,
		Addtime:  time.Unix(time.Now().Local().Unix(), 0).Format("2006-01-02 15:04:05"),
	}

	t1, _ := time.Parse("2006-01-02 15:04:05", insertData.Addtime)
	t2, _ := time.Parse("2006-01-02 15:04:05", insertData.JsonData.Deadline)
	if t1.After(t2) {
		return 0, errors.New("add time is after deadline")
	}

	insertCmd := "insert into todolist(title,owner,text,deadline,addtime)values (?,?,?,?,?)"

	result, err := db.Exec(insertCmd,
		insertData.JsonData.Title,
		insertData.JsonData.Owner,
		insertData.JsonData.Text,
		insertData.JsonData.Deadline,
		insertData.Addtime)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func MySQLTodoListSearch(depend datastruct.TodolistBindMysqlSearch) ([]datastruct.TodolistBindJSONSend, error) {
	satisSlice := make([]datastruct.TodolistBindJSONSend, 0)
	var rows *sql.Rows
	var err error
	if depend.SearchMethod&cfg.SearchAll == cfg.SearchAll {
		rows, err = db.Query("select * from todolist where owner=?", depend.Username)
	} else {

		searchCmd := fmt.Sprintf("select * from todolist where owner='%s'", depend.Username)
		if depend.SearchMethod&cfg.SearchWithStatus == cfg.SearchWithStatus {
			searchCmd = searchCmd + " and status=?"
		}

		if depend.SearchMethod&cfg.SearchWithKeyWord == cfg.SearchWithKeyWord {
			searchCmd = searchCmd + fmt.Sprintf(" and (text like '%%%s%%' or title like '%%%s%%')", depend.Key, depend.Key)
		}
		switch depend.SearchMethod {
		case cfg.SearchWithKeyWord:
			rows, err = db.Query(searchCmd)
		case cfg.SearchWithStatus:
			rows, err = db.Query(searchCmd, depend.Status)
		case cfg.SearchWithKeyWord | cfg.SearchWithStatus:
			rows, err = db.Query(searchCmd, depend.Status)
		}

	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item datastruct.TodolistBindJSONSend
		rows.Scan(&item.Id, &item.Title, &item.Owner, &item.Text, &item.Deadline, &item.Addtime, &item.Status)
		satisSlice = append(satisSlice, item)
	}

	return satisSlice, nil
}

func MySQLTodoListDelete(depend datastruct.TodolistBindMysqlDelete) error{
	var err error
	if depend.DeleteMethod&cfg.DeleteAll==cfg.DeleteAll{
		if _,err=db.Query("delete from todolist where owner=?", depend.Username);err!=nil{
			return err
		}
	}else{
		deleteCmd:=fmt.Sprintf("delete from todolist where owner='%s'",depend.Username)
		idSlice:=make([]int64,len(depend.Idlist))
		if depend.DeleteMethod&cfg.DeleteWithId==cfg.DeleteWithId{
			deleteCmd=deleteCmd+" and id=?"
			copy(idSlice,depend.Idlist)
		}

		if depend.DeleteMethod&cfg.DeleteWithStatus==cfg.DeleteWithStatus{
			deleteCmd=deleteCmd+" and status=?"
		}

		switch depend.DeleteMethod{
		case cfg.DeleteWithId:
			for index := range idSlice {			if _,err=db.Query(deleteCmd,idSlice[index]);err!=nil{
					return err
				}
			}
		case cfg.DeleteWithStatus:
			if _,err=db.Query(deleteCmd,depend.Status);err!=nil{
				return err
			}
		case cfg.DeleteWithId|cfg.DeleteWithStatus:
			for index :=range idSlice{
				if _,err=db.Query(deleteCmd,idSlice[index],depend.Status);err!=nil{
					return err
				}
			}
		}
	}
	return nil
}

func MySQLTodoListModify(depend datastruct.TodolistBindMysqlModify) error{
	var err error
	if depend.ModifyMethod&cfg.ModifyAll==cfg.ModifyAll{
		if _,err=db.Query("update todolist set status = ? where owner = ?",depend.Status,depend.Username);err!=nil{
			return err
		}
	}else{
		idSlice:=make([]int64,len(depend.Idlist))
		copy(idSlice,depend.Idlist)
		for index :=range idSlice{
			if _,err=db.Query("update todolist set status = ? where owner = ? and id = ?",depend.Status,depend.Username,idSlice[index]);err!=nil{
				return err
			}
		}
	}
	return nil
}