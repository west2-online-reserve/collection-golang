package model

import "time"

type Video struct {
	Vid             uint   `json:"vid,omitempty"`
	Uid             uint   `json:"uid,omitempty"`
	Title           string `json:"title,omitempty"`
	Intro           string `json:"intro,omitempty"`
	Tag             string `json:"tag,omitempty"`
	Star            int    `json:"star,omitempty"`
	View            int    `json:"view" form:"view"`
	CreateAt        string `json:"createAt,omitempty"`
	CreateTimeStamp int64  `json:"create_timestamp"`
	DeleteAt        string `json:"deleteAt,omitempty"`
}

func (*Video) Index() string {
	return "video_index"
}
func (*Video) Mapping() string {
	return `
{
  "mappings": {
    "properties": {
      "uid": {
        "type": "integer"
      },
      "vid": {
        "type": "integer"
      },
      "title": {
        "type": "text",
        "analyzer": "ik_smart",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 70
          }
        }
      },
      "tag": {
        "type": "text",
        "analyzer": "ik_smart",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 70
          }
        }
      },
      "intro": {
        "type": "text",
        "analyzer": "ik_smart"
      },
      "star": {
        "type": "integer"
      },
	  "view": {
        "type": "integer"
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

func (v *Video) CreateTime() {
	v.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	v.CreateTimeStamp = time.Now().Unix()
}

func (v *Video) DeleteAtTime() {
	v.DeleteAt = time.Now().Format("2006-01-02 15:04:05")
}
