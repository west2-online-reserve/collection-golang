package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var defaultHeaders map[string]string
var client = &http.Client{}

func init() {
	file, _ := os.Open("bilibili_headers.json")
	defer file.Close()
	json.NewDecoder(file).Decode(&defaultHeaders)

	//videoURL := "https://www.bilibili.com/video/BV12341117rG"
	// data := fetch(videoURL)
	// os.WriteFile("video.html", data, 0644)

}

func fetch(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", defaultHeaders["user-agent"])
	req.Header.Set("Cookie", defaultHeaders["cookie"])

	resp, err := client.Do(req)
	if err != nil {
		log.Println("请求失败:", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

// 1) 如果只有 BV，先换成 aid
func getAidByBvid(bvid string) int64 {
	u := "https://api.bilibili.com/x/web-interface/view?bvid=" + url.QueryEscape(bvid)
	b := fetch(u)
	var r struct {
		Code int
		Data struct {
			Aid int64 `json:"aid"`
		}
	}
	if err := json.Unmarshal(b, &r); err != nil {
		log.Println(err)
	}
	return r.Data.Aid
}

// 2) 获取 wbi 两个 key（从 nav 里提取 .../wbi/*.png 的哈希）
func getWbiKeys() (imgKey, subKey string) {
	b := fetch("https://api.bilibili.com/x/web-interface/nav")
	var r struct {
		Data struct {
			WbiImg struct {
				ImgUrl string `json:"img_url"`
				SubUrl string `json:"sub_url"`
			} `json:"wbi_img"`
		} `json:"data"`
	}

	if err := json.Unmarshal(b, &r); err != nil {
		log.Println(err)
	}
	re := regexp.MustCompile(`/([^/]+)\.png`)
	//os.WriteFile("b.html", b, 0644)
	//log.Println(r.Data)
	//log.Println(re.FindStringSubmatch(r.Data.WbiImg.ImgUrl))
	//log.Println(re.FindStringSubmatch(r.Data.WbiImg.SubUrl))
	imgKey = re.FindStringSubmatch(r.Data.WbiImg.ImgUrl)[1]
	subKey = re.FindStringSubmatch(r.Data.WbiImg.SubUrl)[1]
	return
}

// 3) 按官方逆向规则生成 mixinKey（固定 64 位索引表）
func mixinKey(imgKey, subKey string) string {
	// 索引表来自社区文档（WBI 签名）简化版
	idx := []int{46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 22, 56, 4, 6, 54, 11, 34, 51, 25, 21, 62, 20, 57, 30, 44, 52, 59, 36, 63}
	s := (imgKey + subKey)
	var b strings.Builder
	for _, i := range idx {
		if i < len(s) {
			b.WriteByte(s[i])
		}
	}
	return b.String()[:32]
}

// 4) 生成 w_rid（对“按 key 排序后的 query + &wts=xxx”进行 md5，再拼 mixinKey）
func signWbi(q map[string]string, mix string) (wts string, wRid string) {
	wts = fmt.Sprintf("%d", time.Now().Unix())
	// 按 key 排序并过滤值中某些特殊字符（与前端一致，保守实现）
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		v := q[k]
		v = strings.ReplaceAll(v, "%20", "+")
		parts = append(parts, k+"="+v)
	}
	base := strings.Join(parts, "&") + "&wts=" + wts
	h := md5.Sum([]byte(base + mix))
	wRid = hex.EncodeToString(h[:])
	return
}
