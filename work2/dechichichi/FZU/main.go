package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"
)

// 正则表达式全局变量
var (
	writer = regexp.MustCompile(`target=_blank class="lm_a" style="float:left;">【((.*?))】<\/a>`)
	title  = regexp.MustCompile(`target=_blank title="((.*?))" style="">`)
	text   = regexp.MustCompile(`<a href="((.*?))" target=_blank title=`)
	time   = regexp.MustCompile(`<span class="fr">((.*?))</span>`)
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result := ""
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		result += string(buf[:n])
	}
	return result, nil
}

func HttpGetSingle(i int, page chan<- int, f *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1015&PAGENUM=%d&wbtreeid=1460", i)
	result, err := HttpGet(url)
	if err != nil {
		fmt.Printf("爬取失败: %s\n", err.Error())
		return
	}

	var ans string
	ans1 := writer.FindAllStringSubmatch(result, -1)
	ans2 := title.FindAllStringSubmatch(result, -1)
	ans3 := text.FindAllStringSubmatch(result, -1)
	ans4 := time.FindAllStringSubmatch(result, -1)
	for i := 0; i < len(ans1) && i < len(ans2) && i < len(ans3) && i < len(ans4); i++ {
		ans += fmt.Sprintf("%s  %s  %s  %s\n", ans1[i][1], ans2[i][1], ans3[i][1], ans4[i][1])
	}

	_, err = f.WriteString(ans)
	if err != nil {
		fmt.Printf("写入文件失败: %s\n", err.Error())
		return
	}

	page <- i
}

func work(start int, end int) {
	page := make(chan int, end-start+1)
	f, err := os.Create("data.txt")
	if err != nil {
		fmt.Printf("创建文件失败: %s\n", err.Error())
		return
	}
	defer f.Close()

	var wg sync.WaitGroup
	for i := start; i <= end; i++ {
		wg.Add(1)
		go HttpGetSingle(i, page, f, &wg)
	}

	go func() {
		wg.Wait()
		close(page)
	}()

	for i := start; i <= end; i++ {
		fmt.Printf("第%d页爬取成功\n", <-page)
	}
}

func main() {
	start := 1
	end := 4
	work(start, end)
}
