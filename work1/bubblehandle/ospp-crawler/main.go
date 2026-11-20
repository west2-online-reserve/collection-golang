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

	
	jsonText, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		fmt.Println("转换出错了 瞄")
	}

	
	err = os.WriteFile("data.json", jsonText, 0644)
	if err != nil {
	}

	fmt.Println("已写入 data.json，条数:", len(all))

	
}
