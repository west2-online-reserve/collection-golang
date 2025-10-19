package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var gdb *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Comment{})
	if err != nil {
		panic(err)
	}

	gdb = db
}

func GetDB() *gorm.DB {
	if gdb == nil {
		panic("Database has not initialized!")
	}

	return gdb
}
