package dataprocesser

import (
	"encoding/json"
	"server-redis/cfg"
	"server-redis/datastruct"
	"server-redis/myredis"
	"strings"
	"time"

	"errors"
)

func InsertReceiveTodolist(data datastruct.TodolistBindJSONReceive) (int64, error) {
	type ProcessStruct struct {
		Title    string `json:"title"`
		Owner    string `json:"owner"`
		Text     string `json:"text"`
		Addtime  string `json:"addtime"`
		Deadline string `json:"deadline"`
		Isdone   bool   `json:"isdone"`
	}

	return myredis.RedisInsert(data.Owner, ProcessStruct{
		Title:    data.Title,
		Owner:    data.Owner,
		Text:     data.Text,
		Addtime:  time.Unix(time.Now().Local().Unix(), 0).Format("2006-01-02 15:04:05"),
		Deadline: data.Deadline,
		Isdone:   false,
	})

}

func GetUserTodolist(username string) (datastruct.TodolistBindJSONSendArray, error) {
	if items, err := myredis.RedisGetAll(username); err != nil {
		return datastruct.TodolistBindJSONSendArray{}, err
	} else {
		var result datastruct.TodolistBindJSONSendArray
		result.Todolist.List = make([]datastruct.TodolistBindJSONSend, len(items))
		for i := range items {
			data, _ := json.Marshal(items[i])
			json.Unmarshal(data, &result.Todolist.List[i])
			result.Todolist.List[i].Id = int64(i) +1
		}
		return result, nil
	}
}

func SearchUserTodoList(username string,condition datastruct.TodolistBindRedisCondition) (datastruct.SendingJSONData, error) {
	totalData, err := GetUserTodolist(username)
	if err != nil {
		return datastruct.SendingJSONData{}, err
	}
	result := make([]datastruct.TodolistBindJSONSend, 0)
	switch condition.Method {
	case cfg.SearchAll:
		return datastruct.SendingJSONData{
			Items: totalData.Todolist.List,
			Count: len(totalData.Todolist.List),
			Total: len(totalData.Todolist.List),
		}, nil
	case cfg.SearchWithStatus:
		for i := range totalData.Todolist.List {
			if totalData.Todolist.List[i].Status == condition.Status {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		return datastruct.SendingJSONData{
			Items: result,
			Count: len(result),
			Total: len(totalData.Todolist.List),
		}, nil
	case cfg.SearchWithKeyWord:
		for i := range totalData.Todolist.List {
			if strings.Contains(totalData.Todolist.List[i].Title, condition.Keyword) || strings.Contains(totalData.Todolist.List[i].Text, condition.Keyword) {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		return datastruct.SendingJSONData{
			Items: result,
			Count: len(result),
			Total: len(totalData.Todolist.List),
		}, nil
	case cfg.SearchWithKeyWord | cfg.SearchWithStatus:
		for i := range totalData.Todolist.List {
			if totalData.Todolist.List[i].Status == condition.Status && (strings.Contains(totalData.Todolist.List[i].Title, condition.Keyword) || strings.Contains(totalData.Todolist.List[i].Text, condition.Keyword)) {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		return datastruct.SendingJSONData{
			Items: result,
			Count: len(result),
			Total: len(totalData.Todolist.List),
		}, nil
	}
	return datastruct.SendingJSONData{}, errors.New("wrong method")
}

func DeleteUserTodoList(username string,condition datastruct.TodolistBindRedisCondition) (int,error){
	totalData, err := GetUserTodolist(username)
	if err != nil {
		return 0, err
	}
	if len(totalData.Todolist.List)==0{
		return 0 ,errors.New("no data exist")
	}
	index:=make([]int64,0)
	switch condition.Method {
	case cfg.DeleteAll:
		if err:=myredis.RedisRemoveAll(username);err!=nil{
			return 0,err
		}
		return len(totalData.Todolist.List), nil
	case cfg.DeleteWithStatus:
		for i := range totalData.Todolist.List {
			if totalData.Todolist.List[i].Status == condition.Status {
				index=append(index, int64(i)-1)
			}
		}
		if err:=myredis.RedisMultRemove(username,index);err!=nil{
			return 0,err
		}
		return len(index), nil
	case cfg.DeleteWithId:
		for i:=range condition.Idlist{
			condition.Idlist[i]--
		}
		if err:=myredis.RedisMultRemove(username,condition.Idlist);err!=nil{
			return 0,err
		}
		return len(condition.Idlist), nil
	case cfg.DeleteWithId| cfg.DeleteWithStatus:
		for i := range condition.Idlist {
			if totalData.Todolist.List[i].Status == condition.Status {
				index=append(index, int64(i)-1)
			}
		}
		if err:=myredis.RedisMultRemove(username,index);err!=nil{
			return 0,err
		}
		return len(index), nil
	}
	return 0, errors.New("wrong method")
}

func UpdateTodoList(username string,msg datastruct.TodolistBindRedisUpdate) (int,error){
	totalData, err := GetUserTodolist(username)
	if err != nil {
		return 0, err
	}
	if len(totalData.Todolist.List)==0{
		return 0 ,errors.New("no data exist")
	}
	switch msg.Method {
	case cfg.ModifyAll:
		if items,err:=myredis.RedisPopAll(username);err!=nil{
			return 0,err
		}else{
			for i:=range items{
				items[i].(map[string]interface{})["isdone"]=msg.NewStatus	
				myredis.RedisInsert(username,items[i])
			}
			return len(items),nil
		}
	case cfg.ModifyWithId:
		if len(msg.Idlist)==0{
			return 0,errors.New("idlist empty")
		}
		for i:=range msg.Idlist{
			msg.Idlist[i]--
		}
		if items,err:=myredis.RedisMultPop(username,msg.Idlist);err!=nil{
			return 0,err
		}else{
			if len(items)==0{
				return 0,errors.New("no data exists")
			}
			for i:=range items{
				items[i].(map[string]interface{})["isdone"]=msg.NewStatus	
				myredis.RedisInsert(username,items[i])
			}
			return len(items), nil
		}
	}
	return 0, errors.New("wrong method")
}