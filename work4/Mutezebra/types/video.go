package types

import "time"

type VideoUploadReq struct {
	Title string `json:"title" form:"title"`
	Intro string `json:"intro" form:"intro"`
	Tag   string `json:"tag" form:"tag"`
}

type VideoShowReq struct {
	VID uint `json:"vid" form:"vid"`
}

type VideoDeleteReq struct {
	Vid uint `json:"vid" form:"vid"`
}

type VideoInfoResp struct {
	ID        uint      `json:"ID,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Uid       uint      `json:"uid,omitempty"`
	Title     string    `json:"title,omitempty"`
	Intro     string    `json:"intro,omitempty"`
	Tag       string    `json:"tag,omitempty"`
	Size      int64     `json:"size,omitempty"`
	Views     int       `json:"views,,omitempty"`
	Url       string    `json:"url,omitempty"`
}
