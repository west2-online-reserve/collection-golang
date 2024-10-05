package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// 定义数据结构以解析JSON
type ReplyData struct {
	Data struct {
		Replies []struct {
			Rpid string `json:"rpid"`
		} `json:"replies"`
		Page struct {
			Num int `json:"num"`
		} `json:"page"`
	} `json:"data"`
}

// 正则表达式全局变量
var (
	writer   = regexp.MustCompile(`target=_blank class="lm_a" style="float:left;">【((.*?))】<\/a>`)
	title    = regexp.MustCompile(`target=_blank title="((.*?))" style="">`)
	text     = regexp.MustCompile(`<a href="((.*?))" target=_blank title=`)
	time     = regexp.MustCompile(`<span class="fr">((.*?))</span>`)
	maintext = regexp.MustCompile(`<META Name="description" Content=((.*?))/>`)
)

func fetch(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func getSecondId(data string) ([]string, int) {
	var replyData ReplyData
	err := json.Unmarshal([]byte(data), &replyData)
	if err != nil {
		return nil, 0
	}

	ids := make([]string, len(replyData.Data.Replies))
	for i, reply := range replyData.Data.Replies {
		ids[i] = reply.Rpid
	}
	return ids, replyData.Data.Page.Num
}

func checkAndWriteToFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func main() {
	for i := 1; i <= 200; i++ {
		fmt.Println(i)
		data, err := fetch("https://api.bilibili.com/x/v2/reply?&type=1&oid=420981979&pn=" + strconv.Itoa(i) + "&sort=1")
		if err != nil {
			fmt.Println(err)
			continue
		}
		ids, pn := getSecondId(data)
		for _, id := range ids {
			subData, err := fetch("https://api.bilibili.com/x/v2/reply/reply?oid=420981979&pn=1&ps=10&root=" + id + "&type=1")
			if err != nil {
				fmt.Println(err)
				continue
			}
			path := "第 " + strconv.Itoa(pn) + "页" + id + ".txt"
			checkAndWriteToFile(path, subData)
		}
	}
}
