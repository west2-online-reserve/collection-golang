package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"four/config"
	"four/pkg/ctl"
	"four/pkg/e"
	"four/pkg/log"
	"four/pkg/myutils"
	"four/repository/cache"
	"four/repository/db/dao"
	"four/repository/db/model"
	"four/repository/es/doc"
	esmodel "four/repository/es/model"
	"four/types"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
	"time"
)

type VideoSrv struct {
}

var VideoSrvOnce sync.Once

var VideoService *VideoSrv

func GetVideoSrv() *VideoSrv {
	VideoSrvOnce.Do(func() {
		VideoService = &VideoSrv{}
	})
	return VideoService
}

func (s *VideoSrv) Upload(ctx context.Context, req *types.VideoUploadReq, videoHeader *multipart.FileHeader) (resp interface{}, err error) { //TODO styles
	code := e.SUCCESS
	videoFile, err := videoHeader.Open()
	if err != nil {
		code = e.OpenVideoHeaderFailed
		return ctl.RespError(code, err), err
	}

	err = myutils.IsValidVideoSize(videoHeader.Size) // 判断文件大小是否超限
	if err != nil {
		code = e.InvalidVideo
		return ctl.RespError(code, err), err
	}

	videoDao := dao.NewVideoDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}

	// 开一个goroutine以获得新的视频的id
	v := &model.Video{}
	ch := make(chan *model.Video)
	chErr := make(chan error)
	go videoDao.Create(ch, chErr)
	v = <-ch // 拿到新的视频的id

	ext := myutils.GetFileSuffix(videoHeader.Filename)
	// 保存视频到本地
	if config.Config.System.LocalMode == "local" {
		dst := config.Config.Local.DefaultVideoPath + userInfo.UserName // 保存路径
		err = config.DirExistAndCreate(dst)
		if err != nil {
			code = e.ERROR
			return ctl.RespError(code, err), err
		}
		size := v.VideoPageSize(videoHeader.Size) // 每一个分块的大小
		for i := 0; ; i++ {
			buf := make([]byte, size)
			n, err := videoFile.Read(buf)
			if err != nil && err != io.EOF {
				code = e.ReadVideoFileFailed
				return ctl.RespError(code, err), err
			}
			if n == 0 {
				v.Pages = i // 视频被分片成的页数
				break
			}
			url := fmt.Sprintf(".%s/%d_%d%s", dst, v.ID, i, ext)

			file, err := os.OpenFile(url, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil && err != io.EOF {
				code = e.OpenFileFailed
				return ctl.RespError(code, err), err
			}
			_, err = file.Write(buf)
			if err != nil {
				code = e.VideoWriteToFileFailed
				return ctl.RespError(code, err), err
			}
			err = file.Close()
			if err != nil {
				code = e.CloseFileFailed
				return ctl.RespError(code, err), err
			}
			v.Url = dst
		}
	} else {
		dst := userInfo.UserName + "/" + strconv.Itoa(int(v.ID)) + ext
		buf, err := io.ReadAll(videoFile)
		if err != nil {
			code = e.ReadVideoFileFailed
			return ctl.RespError(code, err), err
		}
		_, err = myutils.UploadVideo(ctx, dst, buf)
		if err != nil {
			code = e.OSSUploadVideoFailed
			return ctl.RespError(code, err), err
		}
		v.Url = dst
	}
	// 赋值，等下通过chan来更新
	v.Uid = userInfo.ID
	v.Tag = req.Tag
	v.Size = videoHeader.Size
	v.Title = req.Title
	v.Intro = req.Intro

	// 向管道中发送值,在goroutine中更新信息
	ch <- v
	err = <-chErr
	if err != nil {
		code = e.CreateVideoFailed
		return ctl.RespError(code, err), err
	}

	model.AddVideoCount(userInfo.UserName) // 增加用户的发布视频数量

	// 向elasticsearch中添加数据
	video := esmodel.Video{
		Vid:   v.ID,
		Uid:   v.Uid,
		Title: v.Title,
		Intro: v.Intro,
		Tag:   v.Tag,
		Star:  v.Star,
	}
	video.CreateTime()
	err = doc.DocCreate(&video)
	if err != nil {
		code = e.CreateDocFailed
		log.LogrusObj.Errorln(err, e.GetMsg(code))
	}

	// 返回resp
	data := &types.VideoInfoResp{
		ID:        strconv.Itoa(int(v.ID)),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Uid:       strconv.Itoa(int(v.Uid)),
		Size:      strconv.FormatInt(v.Size, 10),
		Url:       v.Url,
	}
	return ctl.RespSuccessWithData(code, data), err
}

func (s *VideoSrv) ShowVideoContent(ctx context.Context, req *types.VideoShowReq, resp chan []byte) {
	videoDao := dao.NewVideoDao(ctx)
	video, err := videoDao.FindVideoByVideoID(req.VID)
	if err != nil {
		resp <- nil
		return
	}

	u, err := dao.NewUserDao(ctx).FindUserByUserID(video.Uid)
	if err != nil {
		resp <- nil
		return
	}

	url := video.Url
	if config.Config.System.LocalMode == "local" {
		for i := 0; i < video.Pages; i++ {
			file, err := os.Open(fmt.Sprintf(".%s/%d_%d.mp4", url, video.ID, i))
			if err != nil {
				log.LogrusObj.Errorln(err)
				break
			}
			buf, err := io.ReadAll(file)
			if err != nil {
				log.LogrusObj.Errorln(err)
				break
			}
			resp <- buf
			err = file.Close()
			if err != nil {
				log.LogrusObj.Errorln(err)
				break
			}
		}
		resp <- nil
		close(resp)
		return
	} else {
		name := u.UserName + "/" + strconv.Itoa(int(video.ID)) + ".mp4"
		err = myutils.DownloadVideo(name, resp)
		if err != nil {
			log.LogrusObj.Errorln(err)
		}
		close(resp)
	}
}

func (s *VideoSrv) Comment(ctx context.Context, req *types.VideoCommentReq) (resp interface{}, err error) {
	code := e.SUCCESS
	videoDao := dao.NewVideoDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	comment := &model.Comment{}
	comment.Uid = userInfo.ID
	comment.VideoID = req.VideoID
	comment.Content = req.Content
	err = videoDao.CreateComment(comment)
	if err != nil {
		code = e.VideoCommentCreateFailed
		return ctl.RespError(code, err), err
	}
	return ctl.RespSuccess(code), nil
}

func (s *VideoSrv) Reply(ctx context.Context, req *types.VideoCommentReq) (resp interface{}, err error) {
	code := e.SUCCESS
	if req.ReplyID == 0 || req.ReplyUid == 0 {
		code = e.ReplyRecordNotExist
		err = errors.New("reply failed")
		return ctl.RespError(code, err), err
	}
	videoDao := dao.NewVideoDao(ctx)
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	c := model.Comment{}
	c.ReplyID = req.ReplyID
	c.ReplyUid = req.ReplyUid
	c.VideoID = req.VideoID
	c.Content = req.Content
	c.Uid = userInfo.ID
	root, err := videoDao.FindCommentRoot(&c)
	if err != nil {
		code = e.FindCommentRootFailed
		return ctl.RespError(code, err), err
	}
	if root == 0 {
		c.Root = c.ReplyID
	} else {
		c.Root = root
	}
	err = videoDao.CreateComment(&c)
	return ctl.RespSuccess(code), err
}

func (s *VideoSrv) Show(ctx context.Context, req *types.VideoShowReq) (resp interface{}, err error) {
	code := e.SUCCESS
	videoDao := dao.NewVideoDao(ctx)

	// 先查看缓存中有没有video的info
	videoResp := cache.GetVideoInfo(req.VID)
	var data *types.VideoInfoResp
	if videoResp != nil {
		data = videoResp
	} else {
		video, err := videoDao.FindVideoByVideoID(req.VID)
		if err != nil {
			code = e.FindVideoFailed
			return ctl.RespError(code, err), err
		}
		err = video.SetVideoInfoCache()
		if err != nil {
			code = e.CachedVideoFailed
			return ctl.RespError(code, err), err
		}
		data = &types.VideoInfoResp{
			ID:        strconv.Itoa(int(video.ID)),
			CreatedAt: video.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: video.UpdatedAt.Format("2006-01-02 15:04:05"),
			Uid:       strconv.Itoa(int(video.Uid)),
			Title:     video.Title,
			Intro:     video.Intro,
			Tag:       video.Tag,
		}
	}

	video := &model.Video{}
	video.ID = req.VID
	views := video.Views()
	data.Views = strconv.Itoa(views)

	// 将数据同步到es
	v := &esmodel.Video{}
	hit, id, err := doc.DocSearch(v, "vid", video.ID)
	if err != nil {
		code = e.SearchDocFailed
		return ctl.RespError(code, err), err
	}
	_ = json.Unmarshal(hit.Source, v)
	err = doc.DocUpdate(v, "view", views, id)
	if err != nil {
		code = e.UpdateDocFailed
		return ctl.RespError(code, err), err
	}

	video.AddView()
	return ctl.RespSuccessWithData(code, data), err
}

func (s *VideoSrv) Delete(ctx context.Context, req *types.VideoDeleteReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userInfo, err := ctl.GetFromContext(ctx)
	if err != nil {
		code = e.GetUserInfoFailed
		return ctl.RespError(code, err), err
	}
	videoDao := dao.NewVideoDao(ctx)
	v, err := videoDao.FindVideoByVideoID(req.Vid)
	if err != nil {
		code = e.FindVideoFailed
		return ctl.RespError(code, err), err
	}
	if v.Uid != userInfo.ID && userInfo.UserName != "admin" { // 验证
		code = e.ERROR
		err = errors.New("no operation permission")
		return ctl.RespError(code, err), err
	}

	err = videoDao.Delete(req.Vid)
	if err != nil {
		code = e.DeleteVideoFailed
		return ctl.RespError(code, err), err
	}
	model.DECVideoCount(userInfo.UserName)
	model.DeleteViewCached(req.Vid)
	return ctl.RespSuccess(code), nil
}
