package main

import (
	"lyt112/BilibiliComments"
	"lyt112/FZUSpider"
	"lyt112/config"
)

func main() {
	config.InitDataBase()
	FZUSpider.GetData()
	BilibiliComments.GetComments()
}
