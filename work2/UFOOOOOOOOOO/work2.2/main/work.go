package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strconv"
	// "github.com/antchfx/htmlquery"
	"io"
)

func main() {
	total_page := total_page()
	// detail()
	// int_page_num, _ := strconv.Atoi(page_num)
	var control1 int
	if total_page%20==0{
		control1 = total_page/20
	}else{
		control1 = total_page/20+1
	}
	for i := 1; i <= control1; i++ {
		url := next_url(i)
		message, rpid_str, count := detail(url)
		for j := 0; j <= 17; j++ {
			fmt.Println("main_reply : ", message[j])
			var control int
			if count[j]%10==0{
				control = count[j]/10
			}else{
				control = count[j]/10+1
			}
			for k := 1; k <= control-1; k++{
				replys_url := replys_url(rpid_str[j], k)
				
				reply_message := reply_reply(replys_url)
				for o := 0; o <= 9; o++{
					
					fmt.Println("reply_reply : ", reply_message[o])
					
				}
				// fmt.Println("reply_reply : ", reply_message)
			}
		}

	}
}

func next_url(page int) string {
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&ps=20&pn=%d", page)
	
	return url
}


func replys_url(root string, k int) (string){
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=420981979&type=1&root=%s&ps=10&pn=%d", root, k)
	return url
}


func detail(url string) ([]string, []string, []int){
	fmt.Println(url)
	var client http.Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76")
	resp, _ := client.Do(req) // 发送请求
	body, _ := io.ReadAll(resp.Body)
	r := make(map[string]interface{})
	// fmt.Println([]byte(body))
	json.Unmarshal([]byte(body), &r)

	replies := (r["data"].(map[string]interface{})["replies"].([]interface{}))
	var message []string
	var rpid_str []string
	var count []int
	for i := 0; i <= 17; i++ {
		if num, exist := replies[i].(map[string]interface{}); exist {
			rpid_str = append(rpid_str, num["rpid_str"].(string))
			count_float := num["rcount"].(float64)
			count_int := (int)(count_float)
			count = append(count, count_int)
			if content, exist := num["content"].(map[string]interface{}); exist {
				message = append(message, content["message"].(string))

			}
		}
	}
	return message, rpid_str, count
}


func reply_reply(url string) ([]string){
	fmt.Println(url)
	var client http.Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76")
	resp, _ := client.Do(req) // 发送请求
	body, _ := io.ReadAll(resp.Body)
	r := make(map[string]interface{})
	// fmt.Println([]byte(body))
	json.Unmarshal([]byte(body), &r)

	replies := (r["data"].(map[string]interface{})["replies"].([]interface{}))
	var message []string
	
	for i := 0; i <= 9; i++ {
		if num, exist := replies[i].(map[string]interface{}); exist {
			if content, exist := num["content"].(map[string]interface{}); exist {
				message = append(message, content["message"].(string))
					
			}
		}
	}
	return message
}


func total_page() int {
	var client http.Client
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v2/reply?type=1&oid=420981979&sort=1&pn=1", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.76")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	r := make(map[string]interface{})
	json.Unmarshal([]byte(body), &r)
	pagedata := (r["data"].(map[string]interface{})["page"])
	page, _ := pagedata.(map[string]interface{})
	total_page := page["count"].(float64)
	intpage := (int)(total_page)
	return int(intpage)
}


