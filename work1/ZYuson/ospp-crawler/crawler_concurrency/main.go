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
	"sync"
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
	MaxCon  = 12
	BaseURL = "https://summer.ospp.ac.cn/2025/api/getProList"
	PDFURL  = "https://summer.ospp.ac.cn/2025/api/publicApplication"
)

var client = http.Client{Timeout: 60 * time.Second}

func main() {
	pageAll := make([][]Project, (FetchPage(1).Total+49)/50)
	ForEach(len(pageAll), func(i int) {
		pageAll[i] = FetchPage(i + 1).Rows
	})
	var project []Project
	for _, n := range pageAll {
		project = append(project, n...)
	}
	os.MkdirAll("PDFs", 0755)
	ForEach(len(project), func(i int) {
		FetchPDF(project[i])
	})
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

func ForEach(n int, task func(i int)) {
	sem := make(chan struct{}, MaxCon)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			task(i)
		}(i)
	}
	wg.Wait()
}