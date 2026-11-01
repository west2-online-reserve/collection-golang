package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
)

func main() {
	url := "https://summer.ospp.ac.cn/2025/api/getProDetail"
	Request(url)
	// jsonRequest()
}

func Request(url string) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"programId": "250900034"}`)))
	if err != nil {
		fmt.Println("Error occured in request:", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:159.0) Gecko/20100101 Firefox/159.0")
	request.Header.Set("Referer", "https://summer.ospp.ac.cn/2025/org/prodetail/250900034?lang=zh&list=pro")
	request.Header.Set("Origin", "https://summer.ospp.ac.cn")
	// 注意：Cookie 可能有有效期，如果失效需要重新从浏览器复制
	request.Header.Set("Cookie", "tgt=1792498790.235.7697.274501|2125b4d7ab6829ae59e8509cd0133f66; tgt=1792498790.796.7697.985932|180c5bd6c6a9187530d255592766b85d")

	resp, err := http.DefaultClient.Do(request)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error occured in get the response:", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error occured in process data:", err)
		return
	}

	fmt.Println("HTTP StatusCode:", resp.StatusCode)
	fmt.Println("Return:", string(data))

}
