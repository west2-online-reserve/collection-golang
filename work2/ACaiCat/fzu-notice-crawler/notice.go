package main

import (
	"time"
)

type Notice struct {
	Index       int `gorm:"primaryKey;index;autoincrement"`
	ReleaseTime time.Time
	Author      string
	Title       string
	Body        string
	URL         string
}
