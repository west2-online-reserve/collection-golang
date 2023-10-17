package model
type Reply struct{
	Content struct{
		Message string `json:"message"`
	} `json:"content"`
	Rpid int `json:"rpid"`
	Parent int `json:"parent"`
	Member struct{
		Uname string `json:"uname"`
	}`json:"member"`
	Replies []Reply `json:"replies"`
	Like int `json:"like"`
}
type Response struct{
	Code int `json:"code"`
	Data struct{
		Cursor struct{
			Is_end bool `json:"is_end"`
		} `json:"cursor"`
		Replies []Reply `json:"replies"`
	
	} `json:"data"`
}
type StoredReply struct{
	Parent int `gorm:"parent"`
	Uname string `gorm:"uname"`
	Rpid int `gorm:"rpid"`
	Like int `gorm:"like"`
	Message string `gorm:"message"`

}
