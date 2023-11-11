package types

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
	ID        string `json:"ID,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Uid       string `json:"uid,omitempty"`
	Title     string `json:"title,omitempty"`
	Intro     string `json:"intro,omitempty"`
	Tag       string `json:"tag,omitempty"`
	Size      string `json:"size,omitempty"`
	Views     string `json:"views,,omitempty"`
	Url       string `json:"url,omitempty"`
}
