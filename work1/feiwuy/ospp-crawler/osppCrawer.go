package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"os"
	"encoding/json"
)

const pageSize = 50

func main() {
	url := "https://summer.ospp.ac.cn/2025/api/getProList"
	pageNum := 1
	var dataSlice []Info

	for {
		rows, adjust := RequestList(url, pageNum)
		dataSlice = append(dataSlice, rows...) // 解包符: 把切片"拆开"成独立元素, 将[]Info里的数据一个个填入(append接受元素而不是切片)
		if !adjust {
			break
		}
		pageNum++
	}

	jsonBytes, err := json.MarshalIndent(dataSlice, "", "  ") // 自动换行并缩进
	// jsonBytes, err := json.Marshal(dataSlice)这么写的话数据都在一行
	if err != nil {
		fmt.Println("Error occured in MarshalIndent:", err)
		return
	}

	_ = os.WriteFile("data.json", jsonBytes, 0644)
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
	Code int `json:"code"`
	Msg string `json:"message"`
	Rows []Info `json:"rows"` // 需要的数据
}

func RequestList(url string, pageNum int) ([]Info, bool) {
	// 定义请求
	jsonStr := fmt.Sprintf(`{"pageNum":%d,"pageSize":%d}`, pageNum, pageSize)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error occured in RequestList:", err)
		return nil, false
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error occured in get response:", err)
		return nil, false
	}
	defer resp.Body.Close()

	// 获取返回的数据([]byte格式)
	data, _ := io.ReadAll(resp.Body)

	// 将获取的JSON数据反序列化为[code, message, rows(数组)]
	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("Error occured in Unmarshal:", err)
		return nil, false
	}

	// 若没有数据直接返回
	if len(res.Rows) == 0 {
		return nil, false
	}

	// 判断是否继续
	if (pageNum * pageSize >= res.Total) {
		return nil, false
	}
	return res.Rows, true
}
