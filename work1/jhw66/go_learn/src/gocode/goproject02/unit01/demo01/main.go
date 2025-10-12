package main

// os 包的主要功能
// 分类	常用功能	示例函数
// 📁 文件操作	打开、创建、删除文件	os.Open, os.Create, os.Remove
// 📂 目录操作	创建、删除、遍历目录	os.Mkdir, os.RemoveAll, os.ReadDir
// ⚙️ 环境变量	获取或设置系统环境变量	os.Getenv, os.Setenv
// 💀 程序退出	手动退出程序	os.Exit(code)
// 📜 获取信息	获取文件信息、当前工作目录	os.Stat, os.Getwd
import (
	"log"
	"os"
)

// fmt.Println和log.Println的区别
// 特性	          fmt.Println	log.Println
// 目的        	普通输出（给人看）	日志输出（给系统或开发者分析）
// 输出内容	     只打印你写的字符串	自动加上时间戳、文件行号（可选）
// 可定制性     	无	          可以设置输出到文件、格式、前缀等
// 退出行为	       不退出	      不退出（除非用 log.Fatal）
func main() {
	// 创建文件，如果文件已存在会被截断（清空）
	// 	os.Create() 返回的对象类型是：
	// *os.File 也就是指向一个 os.File 结构体的指针。
	//封装了以下函数:
	// 	func (f *File) Write(b []byte) (n int, err error)
	// func (f *File) WriteString(s string) (n int, err error)
	// func (f *File) Read(b []byte) (n int, err error)
	// func (f *File) Close() error
	// 创建文件后，我们通常需要调用 Close 方法来关闭文件，以释放系统资源。
	file, err := os.Create("test.txt")
	if err != nil {
		log.Fatal(err)
		//格式化输出 log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close() // 确保文件关闭

	log.Println("文件创建成功")
}

//log包常见用法：
// log.Println()：打印日志信息；
// log.Fatalf()：打印日志并 调用 os.Exit(1) 结束程序；
// log.Panic()：打印日志并 触发 panic。
