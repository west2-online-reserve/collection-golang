package dataprocesser

import (
	"encoding/json"
	"server-redis/cfg"
	"server-redis/datastruct"
	"server-redis/myredis"
	"server-redis/mysql"
	"strings"

	"errors"
)

func TodolistSync(username string) error {
	items, err := mysql.MySQLTodoListSyncPack(username)
	if err != nil {
		return errors.New("sync error in mysql")
	}
	if err = myredis.RedisRemoveAll(username); err != nil {
		return errors.New("sync error in redis")
	}
	for i := range items {
		if err := myredis.RedisInsert(username, items[i]); err != nil {
			return errors.New("sync error in redis")
		}
	}
	myredis.RedisExpire(username)
	return nil
}

func InsertReceiveTodolist(data datastruct.TodolistBindJSONReceive) (int64, int64, error) {
	id, addtime, err := mysql.MySQLTodolistInsert(data)
	if err != nil {
		return -1, -1,err
	}
	exist, err := myredis.RedisExist(data.Owner)
	if err != nil {
		return -1, -1,err
	}
	if exist {
		err = myredis.RedisInsert(data.Owner, datastruct.TodolistBindJSONSend{
			Id:       id,
			Owner:    data.Owner,
			Title:    data.Title,
			Text:     data.Text,
			Deadline: data.Deadline,
			Addtime:  addtime,
			Status:   data.Status,
		})
		if err != nil {
			return -1, -1,err
		}
		myredis.RedisExpire(data.Owner)
	} else {
		if err = TodolistSync(data.Owner); err != nil {
			return -1, -1,err
		}
	}
	total, err := myredis.RedisCount(data.Owner)
	if err != nil {
		return -1, -1,err
	}
	return total, id,nil
}

func GetUserTodolist(username string) (datastruct.TodolistBindJSONSendArray, error) {
	if items, err := myredis.RedisGetAll(username); err != nil {
		return datastruct.TodolistBindJSONSendArray{}, err
	} else {
		if len(items) == 0 {
			TodolistSync(username)
			if items, err = myredis.RedisGetAll(username); err != nil {
				return datastruct.TodolistBindJSONSendArray{}, err
			}
		}
		var result datastruct.TodolistBindJSONSendArray
		result.Todolist.List = make([]datastruct.TodolistBindJSONSend, len(items))
		for i := range items {
			data, _ := json.Marshal(items[i])
			json.Unmarshal(data, &result.Todolist.List[i])
		}
		return result, nil
	}
}

func SearchUserTodoList(username string, condition datastruct.TodolistBindRedisCondition, searchPage int) (datastruct.SendingJSONData, bool, error) {
	if searchPage <= 0 {
		return datastruct.SendingJSONData{}, false, errors.New("search page cannot be lower than zero")
	}
	var end int
	totalData, err := GetUserTodolist(username)
	if err != nil {
		return datastruct.SendingJSONData{}, false, err
	}
	if len(totalData.Todolist.List) == 0 {
		if err = TodolistSync(username); err != nil {
			return datastruct.SendingJSONData{}, false, err
		}
		if totalData, err = GetUserTodolist(username); err != nil {
			return datastruct.SendingJSONData{}, false, err
		}
	} else {
		myredis.RedisExpire(username)
	}
	result := make([]datastruct.TodolistBindJSONSend, 0)
	switch condition.Method {
	case cfg.SearchAll:
		if (searchPage-1)*cfg.ItemsCountInPage > len(totalData.Todolist.List) {
			return datastruct.SendingJSONData{}, true, nil
		} else if searchPage*cfg.ItemsCountInPage > len(totalData.Todolist.List) {
			end = len(totalData.Todolist.List)
		} else {
			end = searchPage * cfg.ItemsCountInPage
		}
		return datastruct.SendingJSONData{
			Items: totalData.Todolist.List[(searchPage-1)*cfg.ItemsCountInPage : end],
			Count: end - (searchPage-1)*cfg.ItemsCountInPage,
			Total: len(totalData.Todolist.List),
		}, searchPage*cfg.ItemsCountInPage >= len(totalData.Todolist.List), nil
	case cfg.SearchWithStatus:
		for i := range totalData.Todolist.List {
			if totalData.Todolist.List[i].Status == condition.Status {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		if (searchPage-1)*cfg.ItemsCountInPage > len(result) {
			return datastruct.SendingJSONData{}, true, nil
		} else if searchPage*cfg.ItemsCountInPage > len(result) {
			end = len(result)
		} else {
			end = searchPage * cfg.ItemsCountInPage
		}
		return datastruct.SendingJSONData{
			Items: result[(searchPage-1)*cfg.ItemsCountInPage : end],
			Count: end - (searchPage-1)*cfg.ItemsCountInPage,
			Total: len(result),
		}, searchPage*cfg.ItemsCountInPage >= len(result), nil
	case cfg.SearchWithKeyWord:
		for i := range totalData.Todolist.List {
			if strings.Contains(totalData.Todolist.List[i].Title, condition.Keyword) || strings.Contains(totalData.Todolist.List[i].Text, condition.Keyword) {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		if (searchPage-1)*cfg.ItemsCountInPage > len(result) {
			return datastruct.SendingJSONData{
				Items: nil,
				Count: 0,
				Total: 0}, true, nil
		} else if searchPage*cfg.ItemsCountInPage > len(result) {
			end = len(result)
		} else {
			end = searchPage * cfg.ItemsCountInPage
		}
		return datastruct.SendingJSONData{
			Items: result[(searchPage-1)*cfg.ItemsCountInPage : end],
			Count: end - (searchPage-1)*cfg.ItemsCountInPage,
			Total: len(result),
		}, searchPage*cfg.ItemsCountInPage >= len(result), nil
	case cfg.SearchWithKeyWord | cfg.SearchWithStatus:
		for i := range totalData.Todolist.List {
			if totalData.Todolist.List[i].Status == condition.Status && (strings.Contains(totalData.Todolist.List[i].Title, condition.Keyword) || strings.Contains(totalData.Todolist.List[i].Text, condition.Keyword)) {
				result = append(result, totalData.Todolist.List[i])
			}
		}
		if (searchPage-1)*cfg.ItemsCountInPage > len(result) {
			return datastruct.SendingJSONData{
				Items: nil,
				Count: 0,
				Total: 0}, true, nil
		} else if searchPage*cfg.ItemsCountInPage > len(result) {
			end = len(result)
		} else {
			end = searchPage * cfg.ItemsCountInPage
		}
		return datastruct.SendingJSONData{
			Items: result[(searchPage-1)*cfg.ItemsCountInPage : end],
			Count: end - (searchPage-1)*cfg.ItemsCountInPage,
			Total: len(result),
		}, searchPage*cfg.ItemsCountInPage >= len(result), nil
	}
	return datastruct.SendingJSONData{}, false, errors.New("wrong method")
}

func DeleteUserTodoList(username string, condition datastruct.TodolistBindRedisCondition) error {
	if (condition.Method&cfg.DeleteAll) != cfg.DeleteAll &&
		(condition.Method&cfg.DeleteWithId) != cfg.DeleteWithId &&
		(condition.Method&cfg.DeleteWithStatus) != cfg.DeleteWithStatus {
		return errors.New("wrong method")
	}
	if err := mysql.MySQLTodoListDelete(username, condition); err != nil {
		return err
	}
	return TodolistSync(username)
}

func UpdateTodoList(username string, msg datastruct.TodolistBindRedisUpdate) error {
	if (msg.Method&cfg.ModifyAll) != cfg.ModifyAll &&
		(msg.Method&cfg.ModifyWithId) != cfg.ModifyWithId {
		return errors.New("wrong method")
	}
	mysql.MySQLTodoListModify(username, msg)
	return TodolistSync(username)
}
