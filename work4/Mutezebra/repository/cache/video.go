package cache

import "four/types"

func GetVideoInfo(vid uint) *types.VideoInfoResp {
	ok := RedisClient.HExists(VideoInfoKey(vid), "id")
	if !ok.Val() {
		return nil
	}
	result := RedisClient.HGetAll(VideoInfoKey(vid))
	val := result.Val()
	return &types.VideoInfoResp{
		ID:        val["id"],
		CreatedAt: val["created_at"],
		Uid:       val["uid"],
		Title:     val["title"],
		Intro:     val["intro"],
		Tag:       val["tag"],
		Size:      val["size"],
	}
}

func DestroyVideoInfoCache(vid uint) error {
	err := RedisClient.Del(VideoInfoKey(vid)).Err()
	return err
}
