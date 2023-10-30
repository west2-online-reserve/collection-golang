package model

import (
	"errors"
	"four/consts"
	"four/repository/cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
)

var (
	SearchRecordNotExist = errors.New("search record not exist")
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	NickName       string
	Avatar         string
	Gender         int    `gorm:"default:3"`
	Email          string `gorm:"unique"`
	PasswordDigest string
	Follow         int  `gorm:"unsigned"`
	Fans           int  `gorm:"unsigned"`
	TotpEnable     bool `gorm:"default:false"`
	TotpVerify     bool `gorm:"default:false"`
	TotpUrl        string
	TotpSecret     string
}

func (user *User) SetPassword(password string) (err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), consts.PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(hashPassword)
	return nil
}

func (user *User) CheckPassword(password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err
}

func VideoCount(userName string) int {
	_ = cache.RedisClient.SetNX(cache.VideoCountKey(userName), 0, 0)
	countStr, _ := cache.RedisClient.Get(cache.VideoCountKey(userName)).Result()
	count, _ := strconv.Atoi(countStr)
	return count
}

func AddVideoCount(userName string) {
	cache.RedisClient.Incr(cache.VideoCountKey(userName))
}

func DECVideoCount(userName string) {
	cache.RedisClient.IncrBy(cache.VideoCountKey(userName), -1)
}

func SaveSearchItem(item string, userName string) {
	cache.RedisClient.SAdd(cache.SearchItemKey(userName), item)
}

func GetSearchItem(userName string) (items []string, err error) {
	items, err = cache.RedisClient.SMembers(cache.SearchItemKey(userName)).Result()
	if err != nil {
		return
	}
	if len(items) == 0 { // 说明没有数据
		return nil, SearchRecordNotExist
	}
	return items, nil
}
