package main

//协程最大开到300左右
import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Reply struct {
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
	Rpid   int `json:"rpid"`
	Parent int `json:"parent"`
	Member struct {
		Uname string `json:"uname"`
	} `json:"member"`
	Replies []Reply `json:"replies"`
	Like    int     `json:"like"`
}
type Response struct {
	Code int `json:"code"`
	Data struct {
		Cursor struct {
			Is_end bool `json:"is_end"`
		} `json:"cursor"`
		Replies []Reply `json:"replies"`
	} `json:"data"`
}
type Info struct {
	Uname   string
	Message string
}

func main() {
	start := time.Now()
	init_bili()
	c := colly.NewCollector()
	rule := &colly.LimitRule{
		RandomDelay: time.Second,
		Parallelism: 300, //协程数量
	}
	_ = c.Limit(rule)

	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("authority", "api.bilibili.com") //设置请求头
		req.Headers.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Headers.Set("cookie", "buvid3=F2CED409-B34C-500F-9EDC-F8A1F7B65B9559593infoc; buvid4=5F4C75EA-8F59-C96C-AA37-A668DE2A3BB852219-022071714-tKxIwyBtfcfVXhPDwYtfaw%3D%3D; _uuid=48DD8188-10B9E-10DC1-11AE-97884D9B751C61059infoc; buvid_fp=59ec17c905482edaa86f0e8f3aecb351; rpdid=|(J~|Ru|l|||0J'uY)Yk)kRYJ; b_nut=100; header_theme_version=CLOSE; enable_web_push=DISABLE; CURRENT_BLACKGAP=0; CURRENT_FNVAL=4048; SESSDATA=033af731%2C1714638191%2Cbd2e5%2Ab1CjCT8yqxGxLQeCqFYGF8lr3zm6AJvYeUyMDOb2trgawiGOknOGq4gIOFVBkqcpYxeCQSVmlreUR6cnI5RkdXMElEX2lXYkhvQ3VRVl9JemVva3NUQ0NKOVBNd2NnaGdIS0FYNHE2NGZyWVB5bDNucUMwSGhzUkRaTzVYV3JIUjNPeXJTX3I2V3FRIIEC; bili_jct=747dbd6c710520ad85c7510011bb8430; DedeUserID=33455959; DedeUserID__ckMd5=ba47923f1cc50b21; home_feed_column=5; browser_resolution=1659-838; hit-dyn-v2=1; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTkzNDk0MDksImlhdCI6MTY5OTA5MDE0OSwicGx0IjotMX0.oBREz37XsAM6k6SBICVLFS8D6Z6-bKFxbXxyFlVj7u8; bili_ticket_expires=1699349349; bp_video_offset_33455959=860072690016321575; sid=81hlzd19; b_lsid=510D69641_18BA427B7E4; PVID=1")
	})
	container := Response{}
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response received", r.StatusCode) //打印状态码 成功访问为200
		err := json.Unmarshal(r.Body, &container)      //r.Body为得到的json数据 []byte类型 进行反序列化 字节序列转化成对象
		if err != nil {
			fmt.Println("error", err)
			log.Fatal(err)
		}
	})
	for idx := 0; idx < 500; idx++ { //评论最多到500页
		url := "https://api.bilibili.com/x/v2/reply?&type=1&pn=" + strconv.Itoa(idx) + "&oid=420981979&sort=2"
		c.Visit(url)
		for _, data := range container.Data.Replies {
			var res Info
			res.Message = data.Content.Message
			res.Uname = data.Member.Uname
			fmt.Printf("data.Member.Uname: %v\n", data.Member.Uname)
			fmt.Printf("data.Content.Message: %v\n", data.Content.Message)
			DB.Create(&res)
			if len(data.Replies) > 0 { //评论的评论
				for _, r := range data.Replies {
					var re Info
					re.Message = r.Content.Message
					re.Uname = r.Member.Uname
					DB.Create(&re)
				}
			}
		}
	}
	tm := time.Since(start)
	fmt.Printf("tm: %v\n", tm)
}

var DB *gorm.DB

func init_bili() { //连接到数据库
	dsn := "root:1914@tcp(127.0.0.1:3306)/mybiligorm?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("open db failed,err:", err)
		return
	}
	DB = d
	err = DB.AutoMigrate(&Info{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
