package FZUSpider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"lyt112/FZUSpider/model"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func GetClick(clickId int) int {
	clicktype := "wbnews"
	owner := "1768654345"
	clickid := strconv.Itoa(clickId)
	url := fmt.Sprintf("https://info22.fzu.edu.cn/system/resource/code/news/click/dynclicks.jsp?clickid=%s&owner=%s&clicktype=%s", clickid, owner, clicktype)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("发生错误：", err.Error())
		return 0
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("无法获取数据:", err)
		return 0
	}
	//fmt.Println("点击量: ", string(responseBody))
	clickNum, _ := strconv.Atoi(string(responseBody))
	return clickNum
}
func GetTotalPageNumber() int {
	url := "https://info22.fzu.edu.cn/lm_list.jsp?totalpage=957&PAGENUM=1&wbtreeid=1460"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	//解析网页内容
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 选择元素
	parentElement := document.Find("div.list.fl")
	re := regexp.MustCompile(`共(\d+)条`)
	text := parentElement.Text()
	match := re.FindStringSubmatch(text)
	if len(match) == 2 {
		pageNumber := match[1]
		count, _ := strconv.Atoi(pageNumber)
		return count/20 + 1
	} else {
		fmt.Println("未找到页数")
		return 0
	}
}

func getData(url string, wg *sync.WaitGroup) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}
	defer func() {
		wg.Done()
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
	}()
	//解析网页内容
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Printf("%v", err.Error())
		return
	}
	parentElement := document.Find("div.list.fl")
	parentElement.Find("p").Each(func(index int, element *goquery.Selection) {
		title := element.Find("a[title]").AttrOr("title", "")
		writer := element.Find("a.lm_a").Text()
		link := "https://info22.fzu.edu.cn/" + element.Find("a[title]").AttrOr("href", "")
		date := element.Find("span.fr").Text()
		if IsLate(date) || IsEarly(date) {
			return
		}
		response, err = http.Get(link)
		if err != nil {
			log.Fatal(err)
		}

		document, err = goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		re := regexp.MustCompile(`wbnewsid=(\d+)`)
		match := re.FindStringSubmatch(link)
		clickNum := 0
		if len(match) == 2 {
			wbnewsid := match[1]
			clickid, _ := strconv.Atoi(wbnewsid)
			clickNum = GetClick(clickid)
		} else {
			fmt.Println("未找到 wbnewsid 参数")
		}
		parentElement = document.Find("div.conth")
		texts := ""
		parentElement.Find("p").Each(func(i int, element *goquery.Selection) {
			text := element.Text()
			texts += fmt.Sprintf("%s\n", text)
		})
		parentElement.Find("#vsb_content .v_news_content").Each(func(i int, s *goquery.Selection) {
			// 获取正文部分的文本内容
			text := s.Text()
			texts += fmt.Sprintf("%s\n", text)
		})
		article := model.Article{
			Time:     date,
			Title:    title,
			Content:  texts,
			Author:   writer,
			ClickNum: clickNum,
		}
		err = model.InsertArticle(&article, model.DB)
		if err == nil {
			fmt.Println("插入成功!")
		} else {
			panic(err.Error())
		}
	})
}
func GetData() {
	var wg sync.WaitGroup
	totalNumber := GetTotalPageNumber()
	endPage := GetEndPageNumber(totalNumber)
	startPage := GetStartPageNumber(totalNumber)
	for pageNum := endPage; pageNum <= startPage; pageNum++ {
		url := "https://info22.fzu.edu.cn/lm_list.jsp" + fmt.Sprintf("?totalpage=%d&PAGENUM=%d", totalNumber, pageNum) + "&wbtreeid=1460"
		wg.Add(1)
		go getData(url, &wg)
	}
	wg.Wait()
}

// SearchEndPage 二分查找要爬取的日期
func SearchEndPage(wg *sync.WaitGroup, l, r int, resultChan chan<- int, totalPage int, mutex *sync.Mutex) {
	defer wg.Done()
	if l < r {
		mid := (l + r) / 2
		if !IsEndPage(mid, totalPage) {
			wg.Add(1)
			go SearchEndPage(wg, l, mid, resultChan, totalPage, mutex)
		} else {
			wg.Add(1)
			go SearchEndPage(wg, mid+1, r, resultChan, totalPage, mutex)
		}
	} else {
		mutex.Lock()
		resultChan <- l
		mutex.Unlock()
	}
}

// SearchStartPage 二分查找要爬取的日期
func SearchStartPage(wg *sync.WaitGroup, l, r int, resultChan chan<- int, totalPage int, mutex *sync.Mutex) {
	defer wg.Done()
	if l < r {
		mid := (l + r) / 2
		if IsStartPage(mid, totalPage) {
			wg.Add(1)
			go SearchStartPage(wg, l, mid, resultChan, totalPage, mutex)
		} else {
			wg.Add(1)
			go SearchStartPage(wg, mid+1, r, resultChan, totalPage, mutex)
		}
	} else {
		mutex.Lock()
		resultChan <- l
		mutex.Unlock()
	}
}

// GetEndPageNumber 获取起始页数
func GetEndPageNumber(totalPage int) int {
	var (
		l          = 1
		r          = totalPage
		resultChan = make(chan int, 1)
		wg         sync.WaitGroup
		mutex      sync.Mutex
	)

	wg.Add(1)
	go SearchEndPage(&wg, l, r, resultChan, totalPage, &mutex)
	wg.Wait()
	close(resultChan)
	endPage := <-resultChan
	return endPage
}
func GetStartPageNumber(totalPage int) int {
	var (
		l          = 1
		r          = totalPage
		resultChan = make(chan int, 1)
		wg         sync.WaitGroup
		mutex      sync.Mutex
	)

	wg.Add(1)
	go SearchStartPage(&wg, l, r, resultChan, totalPage, &mutex)
	wg.Wait()
	close(resultChan)
	endPage := <-resultChan
	return endPage
}

func compareEarlyTime(startTime, checkTime, format string) (bool, error) {
	start, err := time.Parse(format, startTime)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return false, err
	}
	check, err := time.Parse(format, checkTime)
	if err != nil {
		return false, err
	}
	if check.After(start) {
		return true, nil
	}
	return false, nil
}
func compareEndTime(endTime, checkTime, format string) (bool, error) {
	end, err := time.Parse(format, endTime)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return false, err
	}
	check, err := time.Parse(format, checkTime)
	if err != nil {
		return false, err
	}
	if check.Before(end) {
		return true, nil
	}
	return false, nil
}
func IsEarly(time string) bool {
	startTime := "2020-01-01"
	//endTime := "2021-09-01"
	format := "2006-01-02"
	ok, err := compareEarlyTime(startTime, time, format)
	if err != nil {
		return false
	}
	return !ok
}
func IsLate(time string) bool {
	endTime := "2021-09-01"
	format := "2006-01-02"
	ok, err := compareEndTime(endTime, time, format)
	if err != nil {
		return false
	}
	return !ok
}
func IsStartPage(mid, totalPage int) bool {
	url := "https://info22.fzu.edu.cn/lm_list.jsp" + fmt.Sprintf("?totalpage=%d&PAGENUM=%d", totalPage, mid) + "&wbtreeid=1460"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	// 解析网页内容
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 选择元素
	parentElement := document.Find("div.list.fl")
	var isStartPage bool
	parentElement.Find("p").Each(func(index int, element *goquery.Selection) {
		date := element.Find("span.fr").Text()
		ok := IsEarly(date)
		if ok {
			isStartPage = true
			return
		}
	})
	return isStartPage
}
func IsEndPage(mid, totalPage int) bool {
	url := "https://info22.fzu.edu.cn/lm_list.jsp" + fmt.Sprintf("?totalpage=%d&PAGENUM=%d", totalPage, mid) + "&wbtreeid=1460"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	// 解析网页内容
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	parentElement := document.Find("div.list.fl")
	var isEndPage bool
	parentElement.Find("p").Each(func(index int, element *goquery.Selection) {
		date := element.Find("span.fr").Text()
		ok := IsLate(date)
		if ok {
			isEndPage = true
			return
		}
	})
	return isEndPage
}

// MeasureTime 测量时间
func MeasureTime(f func()) time.Duration {
	startTime := time.Now()

	f()

	endTime := time.Now()

	return endTime.Sub(startTime)
}
