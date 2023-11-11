package service

import (
	"context"
	"encoding/json"
	"four/pkg/ctl"
	"four/pkg/e"
	"four/repository/db/dao"
	"four/repository/db/model"
	"four/repository/es/doc"
	esmodel "four/repository/es/model"
	"four/types"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var SearchService *SearchSrv

var SearchOnce sync.Once

type SearchSrv struct {
}

func GetSearchSrv() *SearchSrv {
	SearchOnce.Do(func() {
		SearchService = &SearchSrv{}
	})
	return SearchService
}

// Search 基础搜索
func (s *SearchSrv) Search(ctx context.Context, req *types.SearchReq) (resp interface{}, err error) {
	code := e.SUCCESS

	userInfo, _ := ctl.GetFromContext(ctx)
	if userInfo != nil {
		model.SaveSearchItem(req.Content, userInfo.UserName)
	}

	hit, err := doc.SearchUser(req.Content) // 尝试精准匹配用户名
	u := &esmodel.User{}
	if err != nil && err != doc.RecordNotExist {
		code = e.SearchUserFailed
		return ctl.RespError(code, err), err
	} else if err == doc.RecordNotExist {
		u = nil
	} else { // 如果匹配到精准人名那就绑定该用户的信息
		err = json.Unmarshal(hit.Source, u)
		if err != nil {
			code = e.JsonUnmarshalFailed
			return ctl.RespError(code, err), err
		}
	}

	vids, err := doc.Search(u, req) // 获取匹配到的视频的vid
	if err != nil {
		code = e.SearchFailed
		return ctl.RespError(code, err), err
	}

	videoDao := dao.NewVideoDao(ctx)
	videos, err := videoDao.FindVideos(vids) // 根据vid在数据库搜索数据
	if err != nil {
		code = e.FindVideoFailed
		return ctl.RespError(code, err), err
	}

	var videosInfo []types.VideoInfoResp
	for _, v := range videos {
		info := types.VideoInfoResp{
			ID:        strconv.Itoa(int(v.ID)),
			Title:     v.Title,
			Intro:     v.Intro,
			Views:     strconv.Itoa(v.Views()),
			Url:       v.Url,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		videosInfo = append(videosInfo, info)
	}

	data := types.SearchResp{
		User:   u,
		Videos: videosInfo,
	}

	if u == nil && len(videosInfo) == 0 {
		data.Videos = "暂无搜索结果"
	}

	return ctl.RespSuccessWithData(code, data), nil
}

// Filter 筛选搜索
func (s *SearchSrv) Filter(ctx context.Context, req *types.FilterReq) (resp interface{}, err error) {
	code := e.SUCCESS

	userInfo, _ := ctl.GetFromContext(ctx)
	if userInfo != nil {
		model.SaveSearchItem(req.Content, userInfo.UserName)
	}

	if req.Size == "" {
		req.Size = "10"
	}
	if req.Pages == "" {
		req.Pages = "1"
	}

	if req.ViewStart == "" {
		req.ViewStart = "0"
	}
	if req.ViewEnd == "" {
		req.ViewEnd = strconv.FormatInt(math.MaxInt32, 10)
	}

	// 对时间进行处理
	var timeStart, timeEnd int64
	if req.TimeStart == "" {
		timeStart = 0
	} else {
		str, err := time.Parse("2006-01-02", req.TimeStart)
		if err != nil {

		}
		timeStart = str.Unix()
	}
	if req.TimeEnd == "" {
		timeEnd = time.Now().Unix()
	} else {
		str, err := time.Parse("2006-01-02", req.TimeEnd)
		if err != nil {

		}
		timeEnd = str.Unix()
		if timeEnd > time.Now().Unix() {
			timeEnd = time.Now().Unix()
		}
	}

	// 使用空格分隔字符串
	result := strings.Split(req.Tags, " ")

	// 去除空字符串
	var filteredResult string
	for _, item := range result {
		filteredResult = filteredResult + " " + item
	}

	// 筛选
	vids, err := doc.VideoFilter(timeStart, timeEnd, filteredResult, req)
	if err != nil {
		code = e.VideoFilterFailed
		return ctl.RespError(code, err), err
	}

	if len(vids) == 0 {
		return ctl.RespSuccessWithData(code, "暂无更多数据"), nil
	}

	videoDao := dao.NewVideoDao(ctx)
	videos, err := videoDao.FindVideos(vids) // 根据vid在数据库搜索数据
	if err != nil {
		code = e.FindVideoFailed
		return ctl.RespError(code, err), err
	}

	var videosInfo []types.VideoInfoResp
	for _, v := range videos {
		info := types.VideoInfoResp{
			ID:        strconv.Itoa(int(v.ID)),
			Title:     v.Title,
			Intro:     v.Intro,
			Views:     strconv.Itoa(v.Views()),
			Url:       v.Url,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		videosInfo = append(videosInfo, info)
	}

	data := types.SearchResp{
		Videos: videosInfo,
	}
	return ctl.RespSuccessWithData(code, data), nil
}

// HistorySearchItems 查找历史搜索记录
func (s *SearchSrv) HistorySearchItems(ctx context.Context) (resp interface{}, err error) {
	code := e.SUCCESS
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	items, err := model.GetSearchItem(userInfo.UserName)
	var count int
	if err != nil && err != model.SearchRecordNotExist {
		code = e.GetSearchItemFailed
		return ctl.RespError(code, err), err
	} else if err == model.SearchRecordNotExist {
		items = append(items, "暂无搜索记录")
	}

	if len(items) > 10 { // 只能查看最近十条
		count = 10
	} else {
		count = len(items)
	}

	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = items[i]
	}

	return ctl.RespSuccessWithData(code, result), err
}
