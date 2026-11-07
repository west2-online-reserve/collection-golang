package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BaseURL   = "https://summer.ospp.ac.cn/2025/api/getProList"
	PageSize  = 50
	userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0"
	OutFile   = "data.json"
)

type listRequest struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}

type listResponse struct {
	Total int    `json:"total"`
	Rows  []row  `json:"rows"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
}

type row struct {
	ProgramCode     string `json:"programCode"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	Difficulty      string `json:"difficulty"`
	SupportLanguage int    `json:"supportLanguage"` // 难度是string但语言为什么要这样
}

type Project struct {
	ProgramCode     string `json:"programCode"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	Difficulty      string `json:"difficulty"`
	SupportLanguage string `json:"supportLanguage"`
}

func NewClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func languageName(code int) string {
	switch code {
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

// FetchPage 返回该页项目列表和总条数和error
func FetchPage(client *http.Client, pageNum int) ([]Project, int, error) {
	payload, err := json.Marshal(listRequest{PageSize: PageSize, PageNum: pageNum})
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", BaseURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("第 %d 页 HTTP %d: %s", pageNum, resp.StatusCode, bytes.TrimSpace(body))
	}

	var listResp listResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, 0, fmt.Errorf("第 %d 页解析 JSON 失败: %w", pageNum, err)
	}

	if listResp.Code != http.StatusOK {
		return nil, 0, fmt.Errorf("第 %d 页接口报错 %d: %s", pageNum, listResp.Code, listResp.Msg)
	}

	var projects []Project
	for _, r := range listResp.Rows {
		projects = append(projects, Project{
			ProgramCode:     r.ProgramCode,
			ProgramName:     r.ProgramName,
			OrgName:         r.OrgName,
			Difficulty:      r.Difficulty,
			SupportLanguage: languageName(r.SupportLanguage),
		})
	}

	return projects, listResp.Total, nil
}

// Save 返回写入条数和error
func Save(projects []Project) (int, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false) // 打开一看全是 \u0026
	if err := enc.Encode(projects); err != nil {
		return 0, err
	}
	if err := os.WriteFile(OutFile, buf.Bytes(), 0o644); err != nil {
		return 0, err
	}
	return len(projects), nil
}
