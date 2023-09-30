package dataprocesser

import (
	"encoding/json"
	"server-redis/datastruct"
	"server-redis/myredis"
	"time"
)

func InsertReceiveTodolist(data datastruct.TodolistBindJSONReceive) (int64,error) {
	type ProcessStruct struct {
		Title    string `json:"title"`
		Owner    string `json:"owner"`
		Text     string `json:"text"`
		Addtime  string `json:"addtime"`
		Deadline string `json:"deadline"`
		Isdone   bool   `json:"isdone"`
	}

	return myredis.RedisInsert(data.Owner,ProcessStruct{
		Title: data.Title,
		Owner: data.Owner,
		Text: data.Text,
		Addtime: time.Unix(time.Now().Local().Unix(), 0).Format("2006-01-02 15:04:05"),
		Deadline: data.Deadline,
		Isdone: false,
	})

}

func GetUserTodolist(username string) ([]datastruct.TodolistBindJSONSend,error){
	if items,err:= myredis.RedisGetAll(username);err!=nil{
		return nil,err
	}else{
		slice:=make([]datastruct.TodolistBindJSONSend,len(items))
		for i:=range items{
			data,_:=json.Marshal(items[i])
			json.Unmarshal(data,&slice[i])
			slice[i].Id=int64(i)-1
		}
		return slice,nil
	}
}