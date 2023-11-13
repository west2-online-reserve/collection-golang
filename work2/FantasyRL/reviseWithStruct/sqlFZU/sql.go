package sqlFZU

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reptile/model"
)

var Db *sql.DB

// CREATE TABLE `FZUannounce` (
// `id` bigint unsigned AUTO_INCREMENT,
// `title` varchar(1000) DEFAULT NULL,
// `time` varchar(1000) DEFAULT NULL,
// `origin` varchar(1000) DEFAULT NULL,
// `click` int DEFAULT NULL,
// `content` longtext,
// PRIMARY KEY (`id`),
// UNIQUE KEY `id` (`id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
// 存至mysql
func SaveToSQL(ann model.Announce) {
	sql := "insert into FZUannounce (title, time, origin, click, content) values (?, ?, ?, ?, ?)"
	res, err := Db.Exec(sql, ann.Title, ann.Time, ann.Origin, ann.Click, ann.Content)
	if err != nil {
		fmt.Println("save err:", err)
	}
	annId, _ := res.LastInsertId()
	ann.Id = int(annId)
}
