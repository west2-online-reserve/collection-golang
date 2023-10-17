package main

import (
	"bilibili/dao"
	"bilibili/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

)

func main(){
	flag:=true;//标记isEnd
	dao.InitMysql()
	now:=time.Now()
	next:=1
	for flag {
		url:=fmt.Sprintf("https://api.bilibili.com/x/v2/reply/main?oid=420981979&type=1&mode=3&next=%d",next)
		cookie:="buvid3=E2D6E0D3-FA7B-399D-565E-F23583D744C824174infoc; b_nut=1694333124; CURRENT_FNVAL=4048; rpdid=|(umul|Y)Rkk0J'uYmRJJJukY; _uuid=FE417731-FA24-489D-5786-18C81057D9C7C62701infoc; buvid4=BAD09737-8625-68E9-860E-87827040B51B62596-023091020-623AK9X1GgdYv9LcM1groA%3D%3D; buvid_fp=9de7d317750ea5443b7022b3f114d732; DedeUserID=457351159; DedeUserID__ckMd5=3dca149814affb63; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MjgwNjcsImlhdCI6MTY5NzM2ODgwNywicGx0IjotMX0.RHRJaHdwTj7eJQqMwG6l7pgkH8iyy2keBV4Ez1Q66x8; bili_ticket_expires=1697628007; SESSDATA=45671587%2C1713077965%2C77a6a%2Aa2CjATPFRI0TYMJ12rSevbtJDq2RPLkMhD8l7PYlTSLyuPNZJllck_qgJqPXjQJIFplE0SVnNDOWQ5WlN4Z0dfQ0VTcm4zQzI5LXQ2a1I4NFlKLVlxS0xNeUo1R1RLdWc0ZWNqYVNtcmVCRnNNU29sdVdWX3JaT2lkUWVyUUU0cXZmZ1NGY1JKdDR3IIEC; bili_jct=19388e8d014a2aa4277bf3fc8c1e8a49; sid=6zkilw7x; b_lsid=89E9D782_18B3D3B769F; PVID=6"
		client:=&http.Client{}
		req,err:=http.NewRequest("GET",url,nil)
		if err!=nil{
			fmt.Println(err)
		}
		req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
		req.Header.Set("Cookie",cookie)
		res,err:=client.Do(req)
		if err!=nil{
			fmt.Println(err)
		}
		body,err:=io.ReadAll(res.Body)
		if err!=nil{
			fmt.Println(err)
		}
		var back model.Response
		json.Unmarshal(body,&back)
		flag=!back.Data.Cursor.Is_end  //检测是否到了末尾
		for _,data :=range back.Data.Replies{
			var Reply model.StoredReply
			Reply.Message=data.Content.Message
			Reply.Parent=data.Parent
			Reply.Rpid=data.Rpid
			Reply.Uname=data.Member.Uname
			Reply.Like=data.Like
			dao.DB.Create(Reply)
			if len(data.Replies)>0 {//如果有人回复主评论，就继续存储
				for _,r:=range data.Replies{
					var Reply model.StoredReply
					Reply.Message=r.Content.Message
					Reply.Parent=r.Parent
					Reply.Rpid=r.Rpid
					Reply.Uname=r.Member.Uname
					Reply.Like=r.Like
					dao.DB.Create(Reply)
				}
			}

		}
		next++
	}
	duration:=time.Since(now)
	fmt.Println(duration)
}