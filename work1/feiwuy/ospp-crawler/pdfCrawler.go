package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"bytes"
	"encoding/json"
	// "regexp"
)

const (
	pageSize = 50
	listUrl = "https://summer.ospp.ac.cn/2025/api/getProList" // 获取code和名称获取pdf
	pdfUrl = "https://summer.ospp.ac.cn/2025/api/publicApplication" // 从publicApplication的api里获取pdf
)

func main() {
	pageNum := 1
	var dataSlice []Info
	fmt.Println("正在爬取数据")
	for {
		rows, adjust := RequestList(pageNum)
		dataSlice = append(dataSlice, rows...) // 解包符: 把切片"拆开"成独立元素, 将[]Info里的数据一个个填入(append接受元素而不是切片)
		if !adjust {
			break
		}
		pageNum++
	}

	jsonBytes, err := json.MarshalIndent(dataSlice, "", "  ") // 自动换行并缩进
	// jsonBytes, err := json.Marshal(dataSlice)这么写的话数据都在一行
	if err != nil {
		fmt.Println("Error occured in MarshalIndent:", err)
		return
	}

	_ = os.WriteFile("pdfdata.json", jsonBytes, 0644)
	fmt.Println("爬取完成!")

	fmt.Println("正在爬取pdf")
	_ = os.Mkdir("pdfs", 0755) // 新建文件夹存放pdf

	// 开始最基础的 for 循环遍历
	for i, item := range dataSlice {
		fmt.Printf("[%d/%d] 正在下载: %s\n", i+1, len(dataSlice), item.Name)
		
		DownloadPDF(item)

	}

	fmt.Println("爬取完成!")
}

// rows中的具体数据
type Info struct {
	ProId int `json:"proId"`
	ProgramCode string `json:"programCode"`
	Name string `json:"programName"`
}

// 返回的JSON数据分为code, message, rows(我们需要的存数据的数组)
type Response struct {
	Total int `json:"total"` // 总项目数
	Rows []Info `json:"rows"` // 需要的数据
}

func RequestList(pageNum int) ([]Info, bool) {
	// 定义请求
	jsonStr := fmt.Sprintf(`{"pageNum":%d,"pageSize":%d}`, pageNum, pageSize)
	request, err := http.NewRequest("POST", listUrl, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error occured in RequestList:", err)
		return nil, false
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Error occured in get response:", err)
		return nil, false
	}
	defer resp.Body.Close()

	// 获取返回的数据([]byte格式)
	data, _ := io.ReadAll(resp.Body)

	// 将获取的JSON数据反序列化为[code, message, rows(数组)]
	var res Response
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println("Error occured in Unmarshal:", err)
		return nil, false
	}

	// 若没有数据直接返回
	if len(res.Rows) == 0 {
		return nil, false
	}

	// 判断是否继续
	if (pageNum * pageSize >= res.Total) {
		return nil, false
	}
	return res.Rows, true
}

func DownloadPDF(info Info) {
	jsonStr := fmt.Sprintf(`{"proId":%d}`, info.ProId) // 构造请求参数
	req, err := http.NewRequest("POST", pdfUrl, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Println("Error occured in download pdf:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")

	resp, err := http.DefaultClient.Do(req) // 发送请求
	if err != nil {
		fmt.Println("Error occured in response:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体数据给 bodyData
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occured in reading body:", err)
		return
	}

	// 处理文件名中的非法字符（如冒号、斜杠等）
	// reg := regexp.MustCompile(`[\\/:*?"<>|]`)
	// cleanName := reg.ReplaceAllString(info.Name, "_")
	
	// 设置文件名
	// 把编号放在前面，防止重名
	filename := fmt.Sprintf("pdfs/%s_%s.pdf", info.ProgramCode, cleanName)

	// 写入文件
	err = os.WriteFile(filename, bodyData, 0644) // 将二进制流的内容写入filename(.pdf)中
	if err != nil {
		fmt.Printf("Error occured in write:", err)
		return
	}
}
