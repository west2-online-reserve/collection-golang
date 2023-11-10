package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func Fetch1(idx int, page chan int) {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=965&PAGENUM=" + strconv.Itoa(idx) + "&wbtreeid=1460"
	resp, _ := http.Get(url)
	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	News := make([]string, 0)
	Cont := make([]string, 0)
	dom.Find(".clearfloat").Each(func(i int, s *goquery.Selection) {
		Newsifmt := s.Text()
		Content, err1 := s.Find("a").Eq(1).Attr("href")
		if err1 != true {
			fmt.Println(err1)
			return
		}
		News = append(News, Newsifmt)
		Cont = append(Cont, Content)
	})
	SaveFile(idx, News, Cont)
	page <- idx
}
func SaveFile(idx int, News, Cont []string) {
	file, err2 := os.Create("D:/work/project/hellogo/collection-golang/work2/aichicqqqq" + "第" + strconv.Itoa(idx) + "页.txt")
	if err2 != nil {
		fmt.Println(err2)
	}
	defer file.Close()
	for i := 0; i < len(News); i++ {
		file.WriteString(News[i] + Cont[i])
		file.WriteString("--------------------------------------------------------------------")
	}
}
func working(start, end int) {
	page := make(chan int)
	for i := start; i <= end; i++ {
		go Fetch1(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Print("%d", <-page)
	}
}
func main() {
	var start, end int
	fmt.Scan(&start)
	fmt.Scan(&end)
	working(start, end)

}
