package myredis

import (

	"github.com/go-redis/redis"
)

var redisDB *redis.Client
var accountDB *redis.Client

func RedisInit() error {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     `127.0.0.1:6379`,
		Password: "",
		DB:       0,
	})
	accountDB = redis.NewClient(&redis.Options{
		Addr: `127.0.0.1:6379`,
		Password: "",
		DB: 1,
	})
	if _, err := redisDB.Ping().Result(); err != nil {
		return err
	}
	return nil
}

func RedisInsert(key string, data interface{}) error {
	if _, err := redisDB.LPush(key, Struct2Json(data)).Result(); err != nil {
		return err
	}
	return nil
}

func RedisRemove(key string, index int64) error {
	pip := redisDB.TxPipeline()
	if _, err := pip.LSet(key, index, "this value shouldn't exist").Result(); err != nil {
		return err
	}
	if _, err := pip.LRem(key, index, "this value shouldn't exist").Result(); err != nil {
		return err
	}
	if _, err := pip.Exec(); err != nil {
		return err
	}
	return nil
}

func RedisMultRemove(key string, index []int64) error {
	pip := redisDB.TxPipeline()
	for i := range index {
		if _, err := pip.LSet(key, index[i], "this value shouldn't exist").Result(); err != nil {
			return err
		}
	}
	if _, err := pip.Exec(); err != nil {
		return err
	}
	if _, err := pip.LRem(key, 0, "this value shouldn't exist").Result(); err != nil {
		return err
	}
	if _, err := pip.Exec(); err != nil {
		return err
	}
	return nil
}

func RedisRemoveAll(key string) error {
	if _, err := redisDB.Del(key).Result(); err != nil {
		return err
	}
	return nil
}

func RedisGet(key string, index int64) (interface{}, error) {
	item, err := redisDB.LRange(key, index, index).Result()
	if err != nil {
		return nil, err
	}
	var data interface{}
	Json2Struct([]byte(item[0]), &data)
	return data, err
}

func RedisMultGet(key string, index []int64) ([]interface{}, error) {
	items := make([]interface{}, len(index))
	var err error
	for i := range index {
		if items[i], err = RedisGet(key, index[i]); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func RedisGetAll(key string) ([]interface{}, error) {
	items, err := redisDB.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	var data interface{}
	dataSlice := make([]interface{}, 0)
	for i := range items {
		Json2Struct([]byte(items[i]), &data)
		dataSlice = append(dataSlice, data)
	}
	return dataSlice, nil
}

func RedisPop(key string,index int64) (interface{},error){
	pip:=redisDB.TxPipeline()
	item,err:=pip.LRange(key,index,index).Result()
	if err!=nil{
		return nil,err
	}
	if _,err=pip.LRem(key,index,index).Result();err!=nil{
		return nil,err
	}
	var data interface{}
	Json2Struct([]byte(item[0]),&data)
	return data,nil
}

func RedisMultPop(key string,index []int64) ([]interface{},error){
	dataSlice:=make([]interface{},len(index))
	var err error
	for i:=range index{
		if dataSlice[i],err=RedisGet(key,index[i]);err!=nil{
			return nil,err
		}
	}
	if err=RedisMultRemove(key,index);err!=nil{
		return nil,err
	}
	return dataSlice,nil
}

func RedisPopAll(key string) ([]interface{},error){
	dataSlice,err:=RedisGetAll(key)
	if err!=nil{
		return nil,err
	}
	if _,err=redisDB.Del(key).Result();err!=nil{
		return nil,err
	}
	return dataSlice,err
}
