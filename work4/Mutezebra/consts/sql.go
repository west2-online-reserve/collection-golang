package consts

import (
	"strconv"
)

const (
	NewFansTable1 = "CREATE TABLE IF NOT EXISTS `fans"
	NewFansTable2 = "` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `deleted_at` datetime DEFAULT NULL,`uid` int(10) unsigned NOT NULL,`follower_id` int(10) unsigned NOT NULL,PRIMARY KEY (`id`), KEY `idx_deleted_at` (`deleted_at`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	NewCommentTable1  = "CREATE TABLE IF NOT EXISTS `comment"
	NewCommentTable2  = "` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `deleted_at` datetime DEFAULT NULL,`video_id` int(10) unsigned NOT NULL,`root` int(10) unsigned NOT NULL,`reply_id` int(10) unsigned NOT NULL,`uid` int(10) unsigned NOT NULL,`reply_uid` int(10) unsigned NOT NULL,`content` text,PRIMARY KEY (`id`), KEY `idx_deleted_at` (`deleted_at`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	InsertNewComment1 = "INSERT INTO `"
	InsertNewComment2 = "` (created_at,updated_at,deleted_at,video_id,reply_id,uid,content,reply_uid,root) VALUES (NOW(),NOW(),null" // ?,?,?,?,?,?)"

)

func CreateNewFansTable(index uint) string {
	return NewFansTable1 + strconv.Itoa(int(index)) + NewFansTable2
}

func CreateNewCommentTable(index uint) string {
	return NewCommentTable1 + strconv.Itoa(int(index)) + NewCommentTable2
}
