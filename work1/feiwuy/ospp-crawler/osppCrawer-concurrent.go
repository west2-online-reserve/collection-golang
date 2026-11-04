package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"os"
	"encoding/json"
	"sync"
)

const (
	pageSize = 50 // 每次读取项目数
	MAXCON = 8 // 最大线程数
)

func main() {
	url := "https://summer.ospp.ac.cn/2025/api/getProList"

	totalNum := GetTotal(url) // 获取项目数
	pages := (totalNum + pageSize - 1) / pageSize // 获取总页数

	dataSlice := make([][]Info, pages) // 一个长度为pages的类型为[]Info的切片

	var wait sync.WaitGroup // 等待
	maxCon := make(chan struct{}, MAXCON) // 使用通道设置并发上限为8, 当不传入数据时一般使用struct{}

	for i := 0; i < pages; i++ {
		wait.Add(1) // 一个进程启动
		maxCon <- struct{}{} // 占位
		go func(i int) {
			defer wait.Done() // 结束后排队次数减1
			defer func() { <-maxCon}() // 占位少1, defer需要调用函数, 定义在参数保证只调用一次
			rows, err := RequestList(url, i+1) // 获取一页的内容, 编号从1开始, 下标从0开始
			if err != nil {
				fmt.Println("Error occured in RequestList:", err)
				return
			}
			dataSlice[i] = rows // 将rows的值赋值给对应下标的切片, 为什么要用切片组是为了防止多线程同时对切片操作时导致错误(内存覆盖, 重叠...)
		}(i)
	}
	wait.Wait() // 等待进程完成

	var dataSlices []Info // 总的数据切片
	for _, rows := range dataSlice {
		dataSlices = append(dataSlices, rows...) // 分散的切片传入
	}

	jsonBytes, err := json.MarshalIndent(dataSlices, "", "  ") // 自动换行并缩进
	// jsonBytes, err := json.Marshal(dataSlice)这么写的话数据都在一行
	if err != nil {
		fmt.Println("Error occured in MarshalIndent:", err)
		return
	}

	_ = os.WriteFile("data-con.json", jsonBytes, 0644) // 写入文件
	fmt.Println("写入完成!")
}

// rows中的具体数据
type Info struct {
	ProgramCode string `json:"programCode"`
	Name string `json:"programName"`
	OrgId string `json:"orgId"`
	Difficulty string `json:"difficulty"`
	SupportLanguage int `json:"supportLanguage"`
}

// 返回的JSON数据分为code, message, rows(我们需要的存数据的数组)
type Response struct {
	Total int `json:"total"` // 总项目数
	// Code int `json:"code"`
	// Msg string `json:"message"`
	Rows []Info `json:"rows"` // 需要的数据
}

// 获取总的项目数, 也可以对RequestList改造, 不过我不想重写就这样吧
func GetTotal(url string) int {
	// 定义请求
	jsonStr := fmt.Sprintf(`{"pageNum":%d,"pageSize":%d}`, 1, 1)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error occured in RequestList:", err)
		return 0
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error occured in get response:", err)
		return 0
	}
	defer resp.Body.Close()

	// 获取返回的数据([]byte格式)
	data, _ := io.ReadAll(resp.Body)

	// 将获取的JSON数据反序列化为total:
	var res struct{Total int `json:"total"`}
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("Error occured in Unmarshal:", err)
		return 0
	}

	return res.Total
}

func RequestList(url string, pageNum int) ([]Info, error) {
	// 定义请求
	jsonStr := fmt.Sprintf(`{"pageNum":%d,"pageSize":%d}`, pageNum, pageSize)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error occured in RequestList:", err)
		return nil, err
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error occured in get response:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 获取返回的数据([]byte格式)
	data, _ := io.ReadAll(resp.Body)

	// 将获取的JSON数据反序列化为[code, message, rows(数组)]
	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("Error occured in Unmarshal:", err)
		return nil, err
	}

	// 若没有数据直接返回
	if len(res.Rows) == 0 {
		return nil, err
	}

	return res.Rows, nil
}
