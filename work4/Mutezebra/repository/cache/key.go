package cache

import (
	"fmt"
)

const VideoRankKey = "videoRank"

func VideoViewKey(vid uint) string {
	return fmt.Sprintf("view:video:%d", vid)
}

func VideoInfoKey(vid uint) string {
	return fmt.Sprintf("info:video:%d", vid)
}

func SearchItemKey(userName string) string {
	return fmt.Sprintf("search:item:%s", userName)
}

func VideoCountKey(userName string) string {
	return fmt.Sprintf("count:video:%s", userName)
}
