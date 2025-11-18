package util

import (
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

func GetHtmlNodeValByKey(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func ParseShowDynClicks(html string) (clicktype string, owner, clickid int64, ok bool) {
	// 正则表达式匹配 _showDynClicks("...", number, number)
	re := regexp.MustCompile(`_showDynClicks\(\s*["']([^"']+)["']\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`)

	matches := re.FindStringSubmatch(html)
	if len(matches) != 4 {
		return "", 0, 0, false
	}

	clicktype = matches[1]
	owner = mustAtoi64(matches[2])
	clickid = mustAtoi64(matches[3])
	return clicktype, owner, clickid, true
}

func mustAtoi64(s string) int64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return v
	}
	panic("invalid number: " + s)
}
