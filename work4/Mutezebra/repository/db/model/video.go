package model

import (
	"four/consts"
	"four/repository/cache"
	"gorm.io/gorm"
	"strconv"
)

type Video struct {
	gorm.Model
	Uid   uint
	Title string
	Intro string
	Tag   string
	Star  int // 收藏数
	Pages int // 用来记录存视频用了几页
	Size  int64
	Url   string
}

// VideoPageSize 每2MB 分成一块
func (v *Video) VideoPageSize(size int64) int64 {
	Pages := size/consts.EveryPageSize + 1
	sizeEveryPage := size / Pages
	return sizeEveryPage + 1
}

func (v *Video) Views() int {
	_ = cache.RedisClient.SetNX(cache.VideoViewKey(v.ID), 0, 0)
	countStr, _ := cache.RedisClient.Get(cache.VideoViewKey(v.ID)).Result()
	count, _ := strconv.Atoi(countStr)
	return count
}

func (v *Video) AddView() {
	cache.RedisClient.Incr(cache.VideoViewKey(v.ID))
	cache.RedisClient.ZIncrBy(cache.VideoRankKey, 1, cache.VideoViewKey(v.ID))
}
