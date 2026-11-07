package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	BaseURL        = "https://summer.ospp.ac.cn/2025/api/getProList"
	ApplicationURL = "https://summer.ospp.ac.cn/2025/api/publicApplication"
	PageSize       = 50
	userAgent      = "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0"
	OutFile        = "data.json"
	ApplicationDir = "applications"
)

type listRequest struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}

type applicationRequest struct {
	ProId int `json:"proId"`
}

type listResponse struct {
	Total int    `json:"total"`
	Rows  []row  `json:"rows"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
}

type row struct {
	ProId           int    `json:"proId"`
	ProgramCode     string `json:"programCode"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	Difficulty      string `json:"difficulty"`
	SupportLanguage int    `json:"supportLanguage"` // 难度是string但语言为什么要这样
}

type Project struct {
	ProId           int    `json:"proId"`
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
		return nil, 0, fmt.Errorf("第 %d 页	HTTP %d: %s", pageNum, listResp.Code, listResp.Msg)
	}

	var projects []Project
	for _, r := range listResp.Rows {
		projects = append(projects, Project{
			ProId:           r.ProId,
			ProgramCode:     r.ProgramCode,
			ProgramName:     r.ProgramName,
			OrgName:         r.OrgName,
			Difficulty:      r.Difficulty,
			SupportLanguage: languageName(r.SupportLanguage),
		})
	}

	return projects, listResp.Total, nil
}

// Save 保存爬取到的数据, 返回写入条数和error
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

// FetchApplication 返回申请书内容
func FetchApplication(client *http.Client, proId int) ([]byte, error) {
	payload, err := json.Marshal(applicationRequest{ProId: proId})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ApplicationURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proId %d HTTP %d: %s", proId, resp.StatusCode, bytes.TrimSpace(data))
	}

	// 没有申请书
	if len(data) == 0 {
		return nil, nil
	}

	return data, nil
}

// SaveApplication 保存申请书, 返回文件路径
func SaveApplication(p Project, data []byte) (string, error) {
	if err := os.MkdirAll(ApplicationDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(ApplicationDir, fmt.Sprintf("%s.pdf", p.ProgramCode)) //	因为项目名称里会有"/"之类的
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
