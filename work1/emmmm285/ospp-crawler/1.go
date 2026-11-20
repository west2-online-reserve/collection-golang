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

type ReqPayload struct {
	Lang                   string        `json:"lang"`
	PageNum                string        `json:"pageNum"`
	PageSize               string        `json:"pageSize"`
	ProgramName            string        `json:"programName"`
	Difficulty             []interface{} `json:"difficulty"`
	ProgrammingLanguageTag []interface{} `json:"programmingLanguageTag"`
	SupportLanguage        []interface{} `json:"supportLanguage"`
	TechTag                []interface{} `json:"techTag"`
}
type Resp struct {
	Total int   `json:"total"`
	Rows  []Row `json:"rows"`
}
type Row struct {
	ProgramCode     string `json:"programCode"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	Difficulty      string `json:"difficulty"`
	SupportLanguage int    `json:"supportLanguage"`
}
type OutItem struct {
	ProjectCode string `json:"项目编号"`
	ProjectName string `json:"项目名称"`
	OrgName     string `json:"社区名称"`
	Difficulty  string `json:"项目难度"`
	LangText    string `json:"支持语言"`
}

func getLangText(n int) string {
	switch n {
	case 0:
		return "中文&English"
	case 1:
		return "中文"
	case 2:
		return "English"
	default:
		return "未知"
	}
}
func main() {
	var allData []OutItem
	page := 1
	pageSize := 50
	client := &http.Client{Timeout: 30 * time.Second}
	for {
		payload := ReqPayload{
			Lang:                   "zh",
			PageNum:                fmt.Sprintf("%d", page),
			PageSize:               fmt.Sprintf("%d", pageSize),
			ProgramName:            "",
			Difficulty:             []interface{}{},
			ProgrammingLanguageTag: []interface{}{},
			SupportLanguage:        []interface{}{},
			TechTag:                []interface{}{},
		}
		jsonBody, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("序列化请求失败：%v\n", err)
			return
		}
		req, err := http.NewRequest("POST", "https://summer.ospp.ac.cn/2025/api/getProList", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("创建请求失败：%v\n", err)
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Referer", "https://summer.ospp.ac.cn/2025/org/projectlist")
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("请求第%d页失败，重试\n", page)
			time.Sleep(2 * time.Second)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("读取响应失败：%v\n", err)
			return
		}

		var res Resp
		err = json.Unmarshal(body, &res)
		if err != nil {
			fmt.Printf("解析json失败，第%d页：%v\n返回内容：%s\n", page, err, string(body))
			return
		}

		for _, r := range res.Rows {
			allData = append(allData, OutItem{
				ProjectCode: r.ProgramCode,
				ProjectName: r.ProgramName,
				OrgName:     r.OrgName,
				Difficulty:  r.Difficulty,
				LangText:    getLangText(r.SupportLanguage),
			})
		}
		total := res.Total
		totalPage := (total + pageSize - 1) / pageSize
		fmt.Printf("✅ 爬取第 %d / %d 页，本页 %d 条\n", page, totalPage, len(res.Rows))

		if page >= totalPage {
			break
		}
		page++
		time.Sleep(1 * time.Second)
	}
	f, err := os.OpenFile("data.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("创建文件失败：%v\n", err)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(allData)
	if err != nil {
		fmt.Printf("写入json失败：%v\n", err)
		return
	}
	fmt.Printf("\n🎉 爬取完成！一共 %d 条项目，保存到 data.json\n", len(allData))
}
