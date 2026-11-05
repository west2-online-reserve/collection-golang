package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProListResp struct {
	Total	int       	`json:"total"`
	Rows	[]Project 	`json:"rows"`
}

type Project struct {
	ID        	string	`json:"programCode"`
	Title     	string 	`json:"programName"`
	Community	string 	`json:"orgName"`
	Level     	string 	`json:"difficulty"`
	Language  	int    	`json:"supportLanguage"`
	ProID     	int    	`json:"proId"`
}

const (
	BaseURL	= "https://summer.ospp.ac.cn/2025/api/getProList"
	PDFURL	= "https://summer.ospp.ac.cn/2025/api/publicApplication"
)

var client = http.Client{Timeout: 60 * time.Second}

func main() {
	var project []Project
	pageNum := (FetchPage(1).Total + 49) / 50
	for page := 1; page <= pageNum; page++ {
		project = append(project, FetchPage(page).Rows...)
	}
	os.MkdirAll("PDFs", 0755)
	for _, p := range project {
		FetchPDF(p)
	}
	Save(project)
}

func DoPost(URL, body string) ([]byte, error) {
	resp, err := client.Post(URL, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func FetchPage(page int) ProListResp {
	data, err := DoPost(BaseURL, fmt.Sprintf(`{"lang":"zh","pageNum":%d,"pageSize":50}`, page))
	if err != nil {
		return ProListResp{}
	}
	var result ProListResp
	json.Unmarshal(data, &result)
	return result
}

func FetchPDF(project Project) {
	if file, err := os.Stat(filepath.Join("PDFs", project.ID+".pdf")); err == nil && file.Size() > 0 {
		return
	}
	pdf, err := DoPost(PDFURL, fmt.Sprintf(`{"proId":%d}`, project.ProID))
	if err != nil || !bytes.HasPrefix(pdf, []byte("%PDF")) {
		return
	}
	os.WriteFile(filepath.Join("PDFs", project.ID+".pdf"), pdf, 0644)
}