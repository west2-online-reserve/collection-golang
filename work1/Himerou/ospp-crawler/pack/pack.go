package pack

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Rows struct {
	ProgramCode     string `json:"programCode"`
	Difficulty      string `json:"difficulty"`
	ProgramName     string `json:"programName"`
	OrgName         string `json:"orgName"`
	SupportLanguage int    `json:"supportLanguage"`
	ProId           int    `json:"proId"`
	//对应项目的pdf文件ID 使用
	//https://summer.ospp.ac.cn/2025/api/publicApplication
}

type Response struct {
	Total int    `json:"total"`
	Rows  []Rows `json:"rows"`
}

type Request struct {
	PageSize int `json:"pageSize"` // 每一页含有的数据数量
	PageNum  int `json:"pageNum"`
}
type RequestPdf struct {
	ProId int `json:"proId"`
}

func BuildBody(pageNum int) ([]byte, error) {
	reqBody := Request{}
	reqBody.PageNum = pageNum
	reqBody.PageSize = 50
	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	return reqJson, nil
}

func BuildRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	return req, nil
}

func Fetch(pageNum int, url string, client *http.Client) (*http.Response, error) {
	reqBody, err := BuildBody(pageNum)
	if err != nil {
		return nil, fmt.Errorf("请求体构建错误: %v", err)
	}

	req, err := BuildRequest("POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("请求构建错误: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("响应错误: %v", err)
	}
	return resp, nil

}
func DownloadPDF(url string, client *http.Client, rows Rows, dir string) error {
	reqBody := RequestPdf{ProId: rows.ProId}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("请求体编码错误: %v", err)
	}
	req, err := BuildRequest("POST", url, reqBodyJson)

	if err != nil {
		return fmt.Errorf("请求体构建错误: %v", err)
	}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("请求错误: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("proId %d 状态码异常: %v", rows.ProId, resp.StatusCode)
	}

	filename := fmt.Sprintf("ProID%d.pdf", rows.ProId)
	path := filepath.Join(dir, filename)
	// 已存在就跳过,重跑时只补缺失的
	_, err = os.Stat(path)
	if err == nil {
		return nil
	}

	temp := path + "_tmp"
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("数据读取失败: %w", err)
	}
	if resp.ContentLength > 0 && int64(len(data)) != resp.ContentLength {
		_ = os.Remove(temp)
		return fmt.Errorf("ProID %d下载不完整%d/%d", rows.ProId, len(data), resp.ContentLength)
	}
	err = os.WriteFile(temp, data, 0644)
	if err != nil {
		_ = os.Remove(temp)
		return fmt.Errorf("文件写入失败: %w", err)
	}
	return os.Rename(temp, path)
}
func SavetoJson(data Response) {
	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("文件创建失败", err)
		return
	}
	defer file.Close()
	dataJson, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.Write(dataJson)
	if err != nil {
		fmt.Println("写入错误", err)
		return
	}
	fmt.Println("JSON保存结束")
}
func SavetoDatabase(data Response, databaseName string) {
	db, err := sql.Open("sqlite", databaseName)
	if err != nil {
		fmt.Println("数据库打开失败", err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS project(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		programCode TEXT,
		difficulty TEXT,
		programName TEXT,
		orgName TEXT,
		supportLanguage INTEGER,
		proId INTEGER)
		`)
	if err != nil {
		fmt.Println("表建立失败", err)
		return
	}
	for _, row := range data.Rows {
		_, err = db.Exec(`
		INSERT INTO project(
				programCode,
		        difficulty,
		        programName,
		        orgName,
		        supportLanguage,
		        proId
		)
		VALUES(?,?,?,?,?,?)
		`, row.ProgramCode,
			row.Difficulty,
			row.ProgramName,
			row.OrgName,
			row.SupportLanguage,
			row.ProId)
		if err != nil {
			fmt.Println("写入错误", err)
			return
		}
	}
	fmt.Println("SQLite 保存结束")
}
