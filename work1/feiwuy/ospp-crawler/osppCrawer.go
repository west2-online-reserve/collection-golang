package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"os"
	"encoding/json"
)

func main() {
	// url := "https://summer.ospp.ac.cn/2025/api/getProDetail"
	// Request(url)
	url2 := "https://summer.ospp.ac.cn/2025/api/getProList"
	RequestList(url2)
}

type Info struct {
	ProgramCode string `json:"programCode"`
	Name string `json:"programName"`
	OrgId string `json:"orgId"`
	Difficulty string `json:"difficulty"`
	SupportLanguage int `json:"supportLanguage"`
}

// 返回的JSON数据分为code, message, rows(我们需要的存数据的数组)
type Response struct {
	Code int `json:"code"`
	Msg string `json:"message"`
	Rows []Info `json:"rows"`
}

func RequestList(url string) {
	// 定义请求
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"pageNum":1,"pageSize":50}`)))
	if err != nil {
		fmt.Println("Error occured in RequestList:", err)
		return
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error occured in get response:", err)
		return
	}
	defer resp.Body.Close()

	// 获取返回的数据([]byte格式)
	data, _ := io.ReadAll(resp.Body)

	// 将获取的JSON数据反序列化为[code, message, rows(数组)]
	var res Response
	_ = json.Unmarshal(data, &res)
	// fmt.Println(res)

	file, _ := os.Create("data.json")
	defer file.Close()

	// 从数组中获取数据
	for _, item := range res.Rows {
		byteJSON, _ := json.Marshal(item)	
		file.WriteString(string(byteJSON))
	}
}
