package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HttpGet(url string)(result string,err error){
		resp,err1:=http.Get(url)
		
		if err1!=nil{
			err = err1
			return 
		}

		defer resp.Body.Close()

		buf:=make([]byte,4096)

		for{
			n, err2:=resp.Body.Read(buf)

			if(n==0){
				break;
			}

			if err2 !=nil&& err2 !=io.EOF{
				err=err2
				return 
			}

			result +=string(buf[:n])

		}
		return
	}


func working(start,end int){
	for i:=start;i<=end;i++{
		url :="https://info22.fzu.edu.cn/lm_list.jsp?totalpage=1025&PAGENUM="+strconv.Itoa(i)+"&urltype=tree.TreeTempUrl&wbtreeid=1460"
		result,err:=HttpGet(url)
		if err!=nil{
			panic(err.Error())
		}
		fmt.Println(string(result))
		
	}
}


func main() {

	start:=time.Now()
	working(1,300)
	time := time.Since(start)
	fmt.Println(time)

}