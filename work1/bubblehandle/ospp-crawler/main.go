// 第 7 关：把抓到的 565 条整理好，存成 data.json
//
// 前面六关你已经把"网上把数据捞回来"这件事做完了。
// 这一关只干最后一步：把筐里的东西**写到硬盘上**，变成文件。
//
// 为什么作业要求存成文件？因为程序一关，内存里的 all 就没了。
// 存成 data.json，这 565 条就一直躺在硬盘上，可以交作业、可以用别的工具打开看。
//
// 运行（在 go-crawler-101 文件夹里执行）：
//
//	go run ./homework/step6
//
// 跑完你会看到 go-crawler-101 文件夹里多出一个 data.json。
// 注意：文件会落在"你敲命令时站的那个文件夹"里，不是 main.go 旁边。
//
//	（在 go-crawler-101 里执行，data.json 就出现在 go-crawler-101 里）
//
// 本关新词：
//
//	os            operating system 的缩写，中文"操作系统"。
//	              这个包里管文件和文件夹的事，比如写文件、读文件、看目录。
//	WriteFile     写文件。os.WriteFile(文件名, 内容, 权限)。
//	byte          中文"字节"，读作"拜特"。文件的内容在电脑里都是一串字节。
//	[]byte        一串字节。json.Marshal 吐出来的就是这个东西。
//	indent        中文"缩进"。就是每行开头空几个格。
//	MarshalIndent 带缩进的 Marshal。内容和 Marshal 一样，只是排版好看。
//	0644          文件权限，八进制写法（第一个 0 表示这是八进制数）。
//	              意思大致是"我能读能写，别人只能读"。写 0644 就行，先别深究。
//	data.json     文件名。json 是给大家看的文字文件，.json 是它的后缀。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ask struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

// 注意：这个结构体正好就是作业要求的 5 个字段！
//
//	programCode      项目编号
//	programName      项目名称
//	orgName          社区名称
//	difficulty       项目难度
//	supportLanguage  支持语言
//
// 所以抓完之后不需要再转换，直接把 all 存成 JSON 就行。
type project struct {
	ProgramCode     string `json:"programCode"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	Difficulty      string `json:"difficulty"`
	SupportLanguage int    `json:"supportLanguage"`
}

type answer struct {
	Total int       `json:"total"`
	Rows  []project `json:"rows"`
}

func main() {
	client := &http.Client{Timeout: 10 * time.Second}

	var all []project

	pageSize := 50
	page := 1

	for {
		text, err := json.Marshal(ask{PageNum: page, PageSize: pageSize})
		if err != nil {
			fmt.Println("造纸条失败:", err)
			return
		}

		resp, err := client.Post("https://summer.ospp.ac.cn/2025/api/getProList",
			"application/json", bytes.NewReader(text))
		if err != nil {
			fmt.Println("请求失败:", err)
			return
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("读内容失败:", err)
			return
		}

		var out answer
		err = json.Unmarshal(body, &out)
		if err != nil { // 先查错（上次说的顺序，这里改对了）
			fmt.Println("解析失败:", err)
			return
		}

		all = append(all, out.Rows...) // 再装筐

		if len(all) >= out.Total { // 最后才判断停不停
			break
		}

		page++
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("抓完了：总共", len(all), "条")

	// ============================================================
	//            下面就是这一关的新东西，只有这几行
	// ============================================================

	// ---------- TODO(你) ①：把整筐数据变成 JSON 文字 ----------
	//   以前你用 json.Marshal 转过一个 ask（一张纸条）。
	//   这次要转的是 all（一整个筐，565 条），用法完全一样。
	//
	//   唯一区别：这次用 json.MarshalIndent，多两个参数，专门为了排版。
	//     json.MarshalIndent(要转的东西, 每行开头加什么, 每层缩进用什么)
	//       - 第 2 个参数写 ""      = 每行开头什么都不加
	//       - 第 3 个参数写 "  "    = 每往里一层，缩进两个空格
	//
	//   为什么不用 Marshal？565 条挤成一行十几万字符，
	//   编辑器打开会卡，人也没法检查。Indent 之后一行一条，好读。
	//
	//   形状：jsonText, err := json.MarshalIndent(___, "", "  ")
	//   写好把下面两行占位符删掉。
	jsonText, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		fmt.Println("转换出错了 瞄")
	}

	// ---------- TODO(你) ②：把 JSON 文字写进文件 ----------
	//   os.WriteFile(文件名, 内容, 权限)  三个参数，顺序别错：
	//     文件名  -> "data.json"（就写这个名字，不加路径也行）
	//     内容    -> 上面那个 jsonText
	//     权限    -> 0644
	//
	//   注意：内容必须是 []byte 这种"一串字节"。
	//   正好 json.MarshalIndent 吐出来的就是 []byte，可以直接放进来。
	//
	//   它也会返回一个错误，记得检查。
	//   形状：err = os.WriteFile(___, ___, ___)
	//
	//   写完删掉下面这行占位符。
	err = os.WriteFile("data.json", jsonText, 0644)
	if err != nil {
	}

	fmt.Println("已写入 data.json，条数:", len(all))

	// ---- 跑完请自己检查三件事 ----
	// 1. go-crawler-101 文件夹里有没有多出 data.json？多大？（应该有几百 KB）
	// 2. 用 GoLand 打开它，数一数大括号。最外层是 [ ... ] 还是 { ... }？
	//    为什么是这种括号？（提示：all 的类型是 []project，方括号代表"一筐/一串"）
	// 3. 把文件拉到最后，看看最后一条是谁。再跟你之前打印的 all[564] 对一下。
}
