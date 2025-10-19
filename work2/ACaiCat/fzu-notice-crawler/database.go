package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gdb *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}

	err = db.AutoMigrate(&Notice{})
	if err != nil {
		panic("Failed to create tables: " + err.Error())
	}

	gdb = db
}

func GetDb() *gorm.DB {
	if gdb == nil {
		panic("Database not initialized")
	}

	return gdb
}
