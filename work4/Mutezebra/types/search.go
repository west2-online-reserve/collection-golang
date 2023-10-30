package types

type SearchReq struct {
	Content string `json:"content" form:"content"`
	Size    int    `json:"size" form:"size"`
	Pages   int    `json:"pages" form:"pages"` // 页数,默认一页十条记录
}

type FilterReq struct {
	Content   string `json:"content" form:"content,required"`
	Tags      string `json:"tags" form:"tags"`
	Size      string `json:"size" form:"size"`
	Pages     string `json:"pages" form:"pages"`
	TimeStart string `json:"time_start" form:"time_start"`
	TimeEnd   string `json:"time_end" form:"time_end"`
	ViewStart string `json:"view_start" form:"view_start"`
	ViewEnd   string `json:"view_end" form:"view_end"`
}

type SearchResp struct {
	User   interface{} `json:"user,omitempty"`
	Videos interface{} `json:"videos,omitempty"`
}
