package myredis

import "encoding/json"

func Struct2Json(data interface{}) ([]byte){
	d,_:=json.Marshal(data)
	return d
}

func Json2Struct(jsonBytes []byte,data interface{}) error{
	return json.Unmarshal(jsonBytes,&data)
}