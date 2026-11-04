package main

import (
	"Crawler/pack"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

const (
	apiurl string = "https://summer.ospp.ac.cn/2025/api/getProList"
	pdfurl string = "https://summer.ospp.ac.cn/2025/api/publicApplication"
)

var ToDownloadpdf bool = true

func main() {
	client := &http.Client{}
	pages := map[int][]pack.Rows{} //确保顺序保持不变
	sem := make(chan struct{}, 10) //手动限制最多允许10个goroutine同时运行 加速比大致在10
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup
	var mu sync.Mutex  //避免并发map写入
	var once sync.Once //Ai 指导使用once 避免多个goroutine 同时发送信号导致堵塞
loop:
	for pageNum := 1; ; pageNum++ {
		select {
		case <-done:
			break loop
		case sem <- struct{}{}:
			wg.Add(1)
			go func(pageNum int) {
				defer wg.Done()
				defer func() { <-sem }()
				resp, err := pack.Fetch(pageNum, apiurl, client)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != 200 {
					return
				}
				resp_body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("数据读取错误", err)
					return
				}
				resp_data := pack.Response{}
				err = json.Unmarshal(resp_body, &resp_data)
				if err != nil {
					fmt.Println(err)
					return
				}
				if len(resp_data.Rows) == 0 {
					once.Do(func() {
						close(done)
					})
					return
				}
				mu.Lock()
				pages[pageNum] = resp_data.Rows
				mu.Unlock()
				fmt.Printf("第%d页完成\n", pageNum)
			}(pageNum)
		}
	}
	wg.Wait()
	fmt.Println("爬取结束")
	data := pack.Response{}
	for page := 1; ; page++ {
		rows, ok := pages[page]
		if !ok {
			break
		}
		data.Total += len(rows)
		data.Rows = append(data.Rows, rows...)
	}
	//sqlite 保存
	pack.SavetoDatabase(data, "data")

	//json 格式保存
	pack.Savetojson(data)

	//下载pdf申请书
	if ToDownloadpdf {
		fmt.Println("开始下载申请书PDF")
		dir := "PDF"
		os.MkdirAll(dir, 0755)
		file, err := os.ReadFile("data.json")
		if err != nil {
			fmt.Println("读取json失败", err)
			return
		}
		err = json.Unmarshal(file, &data)
		if err != nil {
			fmt.Println("解码失败", err)
			return
		}
		sem = make(chan struct{}, 3)
		var wg sync.WaitGroup
		for _, row := range data.Rows {
			sem <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() { <-sem }()
				defer wg.Done()
				err := pack.Downloadpdf(pdfurl, client, row, dir)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("编号%d申请书下载完成\n", row.ProId)
			}()
		}
		wg.Wait()
		fmt.Println("PDF申请书下载完成")
	}

}
