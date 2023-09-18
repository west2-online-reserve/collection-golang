package cache

import (
	"fmt"
	"strconv"
)

const RankKey = "rank"

func TaskViewKey(id uint) string {
	return fmt.Sprintf("view:task:%s", strconv.Itoa(int(id)))
}
