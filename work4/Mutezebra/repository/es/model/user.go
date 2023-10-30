package model

import "time"

type User struct {
	Uid             uint   `json:"uid,omitempty"`
	UserName        string `json:"user_name,omitempty"`
	CreateAt        string `json:"create_at,omitempty"`
	CreateTimeStamp int64  `json:"create_timestamp,omitempty"`
	DeleteAt        string `json:"delete_at,omitempty"`
}

func (*User) Index() string {
	return "user_index"
}
func (*User) Mapping() string {
	return `
{
  "mappings": {
    "properties": {
      "uid": {
        "type": "integer"
      },
      "user_name": {
        "type": "text",
        "analyzer": "ik_smart",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 20
          }
        }
      },
      "create_at": {
        "type": "date",
        "null_value": null,
        "format": ["yyyy-MM-dd HH:mm:ss"]
      },
	  "create_timestamp":{
        "type": "integer"
      },
      "delete_at": {
        "type": "date",
        "null_value": null,
        "format": ["yyyy-MM-dd HH:mm:ss"]
      }
    }   
  }
}
`
}

func (u *User) CreateTime() {
	u.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	u.CreateTimeStamp = time.Now().Unix()

}
func (u *User) DeleteAtTime() {
	u.DeleteAt = time.Now().Format("2006-01-02 15:04:05")
}
