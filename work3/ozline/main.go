// Code generated by hertz generator.

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/west2-online-reserve/collection-golang/work3/biz/dal"
)

func Init() {
	dal.Init()
}

func main() {
	Init()

	h := server.Default()

	register(h)
	h.Spin()
}
