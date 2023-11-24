package main

import(
	"fmt"
	"time"
	"net/http"
	"strconv"
	"github.com/antchfx/htmlquery"
)


func total_page_num() (string){
	var client http.Client
	req, err := http.NewRequest("GET", "https://info22.fzu.edu.cn/lm_list.jsp?wbtreeid=1460", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.61")
	resp, err:= client.Do(req)  // 发送请求
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc, _ := htmlquery.Parse(resp.Body)
	target, err := htmlquery.Query(doc, "/html/body/div[6]/div/div[2]/div[2]/div/span[2]/span[9]/a")
	page_num := htmlquery.InnerText(target)
	return page_num
}


func make_url (page_num string, page int) (string){
	url := fmt.Sprintf("https://info22.fzu.edu.cn/lm_list.jsp?totalpage=%s&PAGENUM=%d&urltype=tree.TreeTempUrl&wbtreeid=1460", page_num, page)
	return url
}


func get_date(url string) (string, string){
	var client http.Client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.61")
	resp, err:= client.Do(req)  // 发送请求
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc, _ := htmlquery.Parse(resp.Body)
	target, err := htmlquery.Query(doc, "/html/body/div[6]/div/div[2]/div[2]/ul/li/p/span")
	date := htmlquery.InnerText(target)
	href, err := htmlquery.Query(doc, "/html/body/div[6]/div/div[2]/div[2]/ul/li/p/a[2][@href]")
	href_str := htmlquery.SelectAttr(href, "href")
	return date, href_str
}


func select_date_url(date_str string, href_str string) (string, string){
	// fmt.Println(date_str)
	date_init := "2006-01-02"
	parse, err := time.Parse(date_init, date_str)
	if err != nil {
		fmt.Println("err")
	}
	start_date := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end_date := time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC)
	if parse.After(start_date)&&parse.Before(end_date) {
		base_url :="https://info22.fzu.edu.cn/"
		detail_url := fmt.Sprintf("%s%s", base_url, href_str)
		return date_str, detail_url
	}else{
		return "", ""
	}
	
}


func detail_analysis(detail_url string) (string, string, string){
	var client http.Client
	req, err := http.NewRequest("GET", detail_url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/118.0.2088.61")
	resp, err:= client.Do(req)  // 发送请求
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	doc, _ := htmlquery.Parse(resp.Body)
	target1, err := htmlquery.Query(doc, "/html/body/div[5]/div[2]/form/div[1]")
	title := htmlquery.InnerText(target1)
	target2, err := htmlquery.Query(doc, "/html/body/div[5]/div[2]/form/div[3]")
	text := htmlquery.InnerText(target2)
	target3, err := htmlquery.Query(doc, "/html/body/div[5]/div[1]/div/div/a[4]")
	author := htmlquery.InnerText(target3)
	return title, text, author
}


func main(){
	total_page_num := total_page_num()
	int_page_num, _ := strconv.Atoi(total_page_num)
	for i := 1; i < int_page_num; i++{
		url := make_url(total_page_num, i)
		date, href_str :=get_date(url)
		fmt.Println(date)
		date_str, selected_url := select_date_url(date, href_str)
		if selected_url != "" {
			title, text, author := detail_analysis(selected_url)
			fmt.Println(date_str)
			fmt.Println(author)
			fmt.Println(title)
			fmt.Println(text)
		}
		
	}
}