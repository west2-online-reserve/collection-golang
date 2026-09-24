package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Save(projects []Project) {
	db, err := gorm.Open(sqlite.Open("ospp.db"), &gorm.Config{})
	if err != nil {
		return
	}
	if err := db.AutoMigrate(&Project{}); err != nil {
		return
	}
	db.Clauses(clause.OnConflict{
		Columns:	[]clause.Column{{Name: "id"}},
		UpdateAll:	true,
	}).Create(&projects)
}