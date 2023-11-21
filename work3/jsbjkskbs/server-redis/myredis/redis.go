package myredis

import (
	"errors"
	"server-redis/cfg"
	"server-redis/mysql"
	"time"

	"github.com/go-redis/redis"
)

var redisDB *redis.Client
var accountDB *redis.Client

func RedisInit() error {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPwd,
		DB:       cfg.TodoListDb,
	})
	accountDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPwd,
		DB:       cfg.AccountDB,
	})
	if _, err := redisDB.Ping().Result(); err != nil {
		panic(err)
	}
	if _, err := accountDB.Ping().Result(); err != nil {
		panic(err)
	}
	return nil
}

func RedisAccountSync() error {
	list, err := mysql.MySQLAccountSyncPack()
	if err != nil {
		return err
	}

	var cursor uint64

	for {
		keys, cursor, err := accountDB.Scan(cursor, "*", 0).Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			if err = accountDB.Del(key).Err(); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}

	for i := range list {
		if err := accountDB.Set(list[i].Username, list[i].Password, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

func RedisExpire(key string) error {
	response, err := redisDB.Expire(key, time.Minute*10).Result()
	if err != nil {
		return err
	}
	if response {
		return nil
	} else {
		return errors.New("set expire failed")
	}
}

func RedisInsert(key string, data interface{}) error {
	if _, err := redisDB.RPush(key, Struct2Json(data)).Result(); err != nil {
		return err
	} else {
		return nil
	}
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

func RedisPop(key string, index int64) (interface{}, error) {
	pip := redisDB.TxPipeline()
	item, err := pip.LRange(key, index, index).Result()
	if err != nil {
		return nil, err
	}
	if _, err = pip.LRem(key, index, index).Result(); err != nil {
		return nil, err
	}
	var data interface{}
	Json2Struct([]byte(item[0]), &data)
	return data, nil
}

func RedisMultPop(key string, index []int64) ([]interface{}, error) {
	dataSlice := make([]interface{}, len(index))
	var err error
	for i := range index {
		if dataSlice[i], err = RedisGet(key, index[i]); err != nil {
			return nil, err
		}
	}
	if err = RedisMultRemove(key, index); err != nil {
		return nil, err
	}
	return dataSlice, nil
}

func RedisPopAll(key string) ([]interface{}, error) {
	dataSlice, err := RedisGetAll(key)
	if err != nil {
		return nil, err
	}
	if _, err = redisDB.Del(key).Result(); err != nil {
		return nil, err
	}
	return dataSlice, err
}

func RedisFindAccount(key string) (bool, error) {
	if _, err := accountDB.Get(key).Result(); err != nil {
		return false, errors.New("account doesn't exist")
	}
	return true, errors.New("account has already existed")
}

func RedisCheckAccount(username, password string) (bool, error) {
	pwd, err := accountDB.Get(username).Result()
	if err != nil {
		return false, err
	}
	if pwd != password {
		return false, errors.New("inavailable account or password")
	}
	return true, nil
}

func RedisCreateAccount(username, password string) error {
	if find, _ := RedisFindAccount(username); find {
		return errors.New("this account exists")
	}
	if err := mysql.MySQLAccountCreate(username, password); err != nil {
		return err
	}
	if _, err := accountDB.Set(username, password, 0).Result(); err != nil {
		return err
	}
	return nil
}

func RedisExist(username string) (bool,error){
	exists,err:=redisDB.Exists(username).Result()
	return exists==1,err
}

func RedisCount(username string) (int64,error){
	return redisDB.LLen(username).Result()
}