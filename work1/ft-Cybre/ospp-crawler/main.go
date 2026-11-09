package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	url := "https://summer.ospp.ac.cn/2025/api/getProList"

	file, _ := os.Create("data.json")
	defer file.Close()
	file.WriteString("[\n")

	for page := 1; page <= 12; page++ {

		payloadStr := fmt.Sprintf(`{"supportLanguage":[],"techTag":[],"programmingLanguageTag":[],"programName":"","difficulty":[],"pageNum":%d,"pageSize":50,"lang":"zh","orgName":[]}`, page)
		payload := strings.NewReader(payloadStr)

		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "Mozilla/5.0")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, _ := client.Do(req)

		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if page < 12 {
				file.Write(body)
				file.WriteString(",\n")
			} else {
				file.Write(body)
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	file.WriteString("\n]")
}
