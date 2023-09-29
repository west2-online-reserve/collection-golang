package main

import (
	"fmt"
	"os"
)

func main(){
	content:=""
	f,err:=os.OpenFile("ninenine.txt",os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
	if err!=nil{
		fmt.Println(err)
	}
	for i:=1;i<10;i++{
		for j:=1;j<i+1;j++{
			content=fmt.Sprintf("\t%d*%d=%d",j,i,i*j)
			f.Write([]byte(content))
		}
		content="\n"
		f.Write([]byte(content))
	}
	f.Close()
}