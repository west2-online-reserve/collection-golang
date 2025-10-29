package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"west2/model"
	"west2/util"

	"golang.org/x/net/html"
)

func GetHtmlNode(url string) *html.Node {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("无法连接网址： %s, 错误：%v", url, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("连接失败，状态码: %d", resp.StatusCode)
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("无法解析网站: %v", err)
		return nil
	}

	return doc
}

func ParseNode(n *html.Node, target *model.Node, ch chan *html.Node) {
	var wg sync.WaitGroup

	var addKeyVal func(*html.Node, *sync.WaitGroup)
	addKeyVal = func(n *html.Node, wg *sync.WaitGroup) {
		defer wg.Done()

		if n == nil {
			return
		}

		if n.Type == target.Type && n.Data == target.Data && util.GetHtmlNodeValByKey(n, "class") == target.ClassName {
			ch <- n
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			wg.Add(1)
			go addKeyVal(c, wg)
		}
	}

	wg.Add(1)
	go addKeyVal(n, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func ParseNodeAndDeal(n *html.Node, target *model.Node, deal func (n *html.Node)) {
	if n == nil {
		return
	}

	if n.Type == target.Type && n.Data == target.Data && util.GetHtmlNodeValByKey(n, "class") == target.ClassName {
		deal(n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ParseNodeAndDeal(c, target, deal)
	}
}

func GetFZUClickCount(clicktype string, owner, clickid int64) (string, error) {
	url := fmt.Sprintf(
		os.Getenv("CLICK_BASE_URL")+"?clickid=%d&owner=%d&clicktype=%s",
		clickid, owner, clicktype,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	clickCount := strings.TrimSpace(string(body))
	return clickCount, nil
}
