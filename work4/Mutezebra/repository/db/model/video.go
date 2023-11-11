package model

import (
	"four/consts"
	"four/repository/cache"
	"gorm.io/gorm"
	"strconv"
	"time"
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
	View  int
}

func (v *Video) filed() map[string]interface{} {
	return map[string]interface{}{
		"id":         v.ID,
		"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
		"uid":        v.Uid,
		"title":      v.Title,
		"intro":      v.Intro,
		"tag":        v.Tag,
		"size":       v.Size,
		"update":     time.Now().Unix(),
		"last_use":   time.Now().Unix(),
	}
}

// VideoPageSize 每2MB 分成一块
func (v *Video) VideoPageSize(size int64) int64 {
	Pages := size/consts.EveryPageSize + 1
	sizeEveryPage := size / Pages
	return sizeEveryPage + 1
}

func (v *Video) SetVideoInfoCache() error {
	ok := cache.RedisClient.HMSet(cache.VideoInfoKey(v.ID), v.filed())
	if ok.Val() != "OK" {
		return ok.Err()
	}
	return nil
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

func DeleteViewCached(vid uint) {
	cache.RedisClient.Del(cache.VideoViewKey(vid))
}
