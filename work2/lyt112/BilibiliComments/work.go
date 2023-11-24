package BilibiliComments

import (
	"encoding/json"
	"fmt"
	"io"
	"lyt112/config"
	"net/http"
	"strconv"
)

func getCommentData(url string) map[string]interface{} {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(response.Body)
	var comment map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&comment); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}
	return comment
}
func GetReplies(root string, oid string, commentCount int) int {
	replyCount := commentCount
	startNumber := "1"
	replyUrl := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=%s&type=1&root=%s&pn=%s", oid, root, startNumber)
	comment := getCommentData(replyUrl)
	if comment == nil {
		return 0
	}
	data := comment["data"].(map[string]interface{})
	page := data["page"].(map[string]interface{})
	counts := page["count"].(float64)
	size := page["size"].(float64)
	count := int(counts)
	endNumber := count + commentCount
	totalNumber := count/int(size) + 1
	if count == 0 {
		return count
	}
	for pageNumbers := 1; replyCount <= endNumber && pageNumbers <= totalNumber; pageNumbers++ {
		pageNumber := strconv.Itoa(pageNumbers)
		replyUrl = fmt.Sprintf("https://api.bilibili.com/x/v2/reply/reply?oid=%s&type=1&root=%s&pn=%s", oid, root, pageNumber)
		comments := getCommentData(replyUrl)
		if comments == nil {
			break
		}
		data, ok := comments["data"].(map[string]interface{})
		if !ok {
			break
		}
		replies, ok := data["replies"].([]interface{})
		if !ok {
			break
		}
		for _, reply := range replies {
			replyData, ok := reply.(map[string]interface{})
			if !ok {
				continue
			}
			content, ok := replyData["content"].(map[string]interface{})
			if !ok {
				continue
			}
			message, ok := content["message"].(string)
			if !ok {
				continue
			}
			fmt.Println("Reply:", replyCount, message)
			replyCount++
		}
	}
	return count
}

func GetComments() {
	bvNumber := config.GetBvNumber()
	startNumber := "1"
	apiUrl := fmt.Sprintf("https://api.bilibili.com/x/v2/reply?jsonp=jsonp&pn=%s&type=1&oid=%s&sort=2", startNumber, bvNumber)
	comment := getCommentData(apiUrl)
	if comment == nil {
		return
	}
	data := comment["data"].(map[string]interface{})
	page := data["page"].(map[string]interface{})
	counts := page["count"].(float64)
	size := page["size"].(float64)
	count := int(counts)
	totalNumber := count/int(size) + 1
	endNumber := count
	count = 1
	for pageNumbers := 1; count <= endNumber && pageNumbers <= totalNumber; pageNumbers++ {
		pageNumber := strconv.Itoa(pageNumbers)
		apiUrl = fmt.Sprintf("https://api.bilibili.com/x/v2/reply?jsonp=jsonp&pn=%s&type=1&oid=%s", pageNumber, bvNumber)
		comments := getCommentData(apiUrl)
		if comments == nil {
			return
		}
		data, ok := comments["data"].(map[string]interface{})
		if !ok {
			break
		}
		replies, ok := data["replies"].([]interface{})
		if !ok {
			break
		}
		for _, reply := range replies {
			replyData, ok := reply.(map[string]interface{})
			if !ok {
				continue
			}
			rpid := replyData["rpid"].(float64)
			oids := replyData["oid"].(float64)
			root := strconv.FormatFloat(rpid, 'f', -1, 64)
			oid := strconv.FormatFloat(oids, 'f', -1, 64)
			commentReplies := replyData["replies"]
			content, ok := replyData["content"].(map[string]interface{})
			if !ok {
				continue
			}
			message, ok := content["message"].(string)
			if !ok {
				continue
			}
			fmt.Println("comment:", count, message)
			count++
			if commentReplies != nil {
				count += GetReplies(root, oid, count)
			}
		}
	}
}
