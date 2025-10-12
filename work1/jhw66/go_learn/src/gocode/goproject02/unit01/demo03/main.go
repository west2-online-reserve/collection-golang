// Go 语言也提供了多种写入文件的方式，包括逐行写入、一次性写入等。
// 我们可以使用 os 包来创建和写入文件。
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

// Go 1.16 把 ioutil 里的很多函数迁移到了 os 或 io 包里
// 包	主要职责
// os	文件和操作系统交互（打开文件、读写、路径）
// io	数据流操作（Reader、Writer 接口）
// bufio	提高读写效率（加缓冲层）
func main() {
	err := os.Remove("../test.txt")
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}
	fmt.Println("File deleted successfully!")
	time.Sleep(time.Second * 2)

	file, err := os.OpenFile("../test.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	// 	| 模式组合          | 含义                       | 常用场景                  |
	// | ------------- | ------------------------ | --------------------- |
	// | `os.O_RDONLY` | 只读                       | 只读取文件内容               |
	// | `os.O_WRONLY` | 只写                       | 只写文件内容（文件必须已存在，否则报错）  |
	// | `os.O_RDWR`   | 读写                       | 同时读写文件                |
	// | `os.O_APPEND` | 追加写入                     | 写操作自动追加到文件末尾          |
	// | `os.O_CREATE` | 不存在则创建                   | 配合写模式使用               |
	// | `os.O_TRUNC`  | 清空文件                     | 打开文件时清空原内容（覆盖写）       |
	// | `os.O_EXCL`   | 配合 `O_CREATE` 使用，文件存在则报错 | 用于防止覆盖已有文件            |
	// | `os.O_SYNC`   | 同步写入（立即写入磁盘，不走缓存）        | 写入要求非常严格的场合，如日志或数据库文件 |
	// 	0644
	// 文件权限（Linux/macOS 常见，但 Windows 也接受）。
	// 数字	权限含义
	// 6 = 4+2	拥有者可读可写
	// 4	组用户只读
	// 4	其他用户只读

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File create successfully!")
	defer file.Close()

	//文件输入

	// 	os.File.WriteString：直接写入磁盘，适合小文件。
	// bufio.Writer：先写进内存缓冲，最后一次性写入，适合大量写入。
	// bufio.Writer会更快

	// 方式1：直接写入字符串
	file.WriteString("直接写入字符串\n")

	// 方式2：写入字节切片
	data := []byte("写入字节切片\n")
	file.Write(data)

	// 方式3：使用fmt.Fprintf格式化写入
	fmt.Fprintf(file, "格式化写入: %d\n", 123)

	//方法4：使用读写器
	//os.File 是直接读写文件。
	//bufio 是带缓冲的读写层，包在 os.File 外面，用来提高效率。
	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, "Hello, World!")
	writer.WriteString("World,Hello!\n")
	writer.Write([]byte("World!Hello!\n"))
	writer.Flush()

	var content []byte
	//方法5：使用os.WriteFile一次性写入
	content = []byte("Hello!World!\n")
	err = os.WriteFile("../test.txt", content, 0644)
	// 	os.WriteFile 会重新创建或覆盖目标文件内容。
	// 所以这行代码：
	// err = os.WriteFile("../test.txt", content, 0644)
	// 会把你前面写入的所有内容清空！😅
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully!")

	//文件读取
	// 	bufio.Reader
	// 它一次性从文件中多读一大块（比如 4KB 或 8KB）到内存缓冲区中，
	// 之后每次调用 ReadString()，只是从内存里取数据。
	// 会比os.Read()快

	file.Seek(0, 0)
	//Go 的文件读写是共用一个文件指针,不分读写指针,所以这里要位移一下
	//方法一：使用读取器
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n') //读取一行(直到 '\n'),包括/n
	if err != nil {
		fmt.Println("读取结束或出错:", err)
	}
	fmt.Println("读取到：", line)

	//在 Go 里，如果你想手动控制读取位置，要用：
	// file.Seek(offset, whence)
	// 参数说明：
	// offset：要移动的字节数（可以为负）
	// whence：基准点
	// 0 → 从文件开头算（io.SeekStart）
	// 1 → 从当前位置算（io.SeekCurrent）
	// 2 → 从文件末尾算（io.SeekEnd）

	file.Seek(0, 0)

	//方法二：扫描器（最常用于按行读）
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //自动去掉/n
		line := scanner.Text()
		fmt.Println("读取到：", line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("出错：", err)
	}

	file.Seek(0, 0)

	//方法三：使用os.ReadFile一次性读取
	content, err = os.ReadFile("../test.txt") //包括/n
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println("读取到：", string(content))

	file.Seek(0, 0)

	//方法四：使用os.Read
	buf := make([]byte, 64)  // 只读 64 字节
	n, err := file.Read(buf) //包括\n
	if err != nil {
		fmt.Println("读取出错：", err)
	}
	fmt.Println("读取到：", string(buf[:n]))

}
