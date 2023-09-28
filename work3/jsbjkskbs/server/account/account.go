package account

import (
	"server/mysql"
	"server/datastruct"
)

func InsertUserTodoList(username string,what2do datastruct.TodolistBindJSONReceive) (int64,error){
	var id int64
	var err error

	if id,err=mysql.MySQLTodolistInsert(username,what2do);err!=nil{
		return 0,err
	}

	return id,nil
}