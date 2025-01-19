// Code generated by hertz generator. DO NOT EDIT.

package router

import (

	model "Demo/biz/router/model"
	task "Demo/biz/router/task"
	user "Demo/biz/router/user"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	task.Register(r)

	user.Register(r)

	model.Register(r)

}
