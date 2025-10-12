package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer func() {
		conn.Close()
		fmt.Println("conn关闭")
	}()

	for {
		//创建一个切片，准备将读入的数据放入切片
		buf := make([]byte, 1024)
		//从conn连接中读取数据
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		fmt.Println(string(buf[0:n]))
	}
}
func main() {
	fmt.Println("服务器端启动")
	//进行监听，需要指定服务器TCP协议，服务器端的IP+PORT
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("监听失败,err:", err)
		return
	}

	//监听成功以后要等待客户端连接
	//循环等待
	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println("客户端等待失败,err2:", err2)
		} else {
			fmt.Printf("等待连接成功,con=%v,接收到的客户端信息:%v\n", conn, conn.RemoteAddr().String())
		}

		//准备一个协程，协程处理客户端服务请求：
		go process(conn) //不同客户端请求，连接conn不一样

	}

}

//运行时要先启动服务器端
