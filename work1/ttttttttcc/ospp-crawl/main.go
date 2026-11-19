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

// 一条项目的字段
type Project struct {
	ProgramCode string `json:"programCode"`
	Name        string `json:"name"`
	Community   string `json:"community"`
	Difficulty  string `json:"difficulty"`
	Lang        string `json:"lang"`
}

// 接口返回外层结构
type RespData struct {
	Total int       `json:"total"`
	Rows  []Project `json:"rows"`
}

// POST发给服务器的查询参数
type ReqBody struct {
	SupportLanguage      []string `json:"supportLanguage"`
	TechTag              []string `json:"techTag"`
	ProgrammingLanguageTag []string `json:"programmingLanguageTag"`
	ProgramName          string   `json:"programName"`
	Difficulty           []string `json:"difficulty"`
	Lang                 string   `json:"lang"`
	OrgName              []string `json:"orgName"`
	PageNum              string   `json:"pageNum"`
	PageSize             string   `json:"pageSize"`
}

func main() {
	var allProjects []Project
	page := 1
	const size = "50"

	for {
		// 构造请求体
		body := ReqBody{
			SupportLanguage:      []string{},
			TechTag:              []string{},
			ProgrammingLanguageTag: []string{},
			ProgramName:          "",
			Difficulty:           []string{},
			Lang:                 "1",
			OrgName:              []string{},
			PageNum:              fmt.Sprintf("%d", page),
			PageSize:             size,
		}
		jsonBody, _ := json.Marshal(body)

		url := "https://summer.ospp.ac.cn/2025/api/getProList"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("创建请求失败：", err)
			break
		}
		// 必须带上这个头部：告诉服务器我发的是JSON
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Referer", "https://summer.ospp.ac.cn/2025/org/projectlist")

		client := http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("请求失败：", err)
			break
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取失败：", err)
			break
		}

		var res RespData
		err = json.Unmarshal(data, &res)
		if err != nil {
			fmt.Println("❌解析失败，返回内容：")
			fmt.Println(string(data))
			break
		}

		if len(res.Rows) == 0 {
			fmt.Println("✅没有更多项目，爬取结束！")
			break
		}

		allProjects = append(allProjects, res.Rows...)
		fmt.Printf("✅第%d页，拿到%d条，累计：%d\n", page, len(res.Rows), len(allProjects))
		page++
		time.Sleep(700 * time.Millisecond)
	}

	f, err := os.Create("data.json")
	if err != nil {
		fmt.Println("创建文件失败：", err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(allProjects)
	fmt.Printf("\n🎉完成！一共抓取 %d 条项目，保存在data.json\n", len(allProjects))
}
