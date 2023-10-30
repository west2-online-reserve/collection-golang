package api

import (
	"context"
	"four/config"
	"four/pkg/log"
	"four/service"
	"four/types"
	"github.com/cloudwego/hertz/pkg/app"
	resp2 "github.com/cloudwego/hertz/pkg/protocol/http1/resp"
	"net/http"
)

func VideoUploadHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.VideoUploadReq
		err := c.Bind(&req)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		video, err := c.FormFile("video")
		if err == nil {
			l := service.GetVideoSrv()
			resp, err := l.Upload(ctx, &req, video)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func VideoWatchContentHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Content-Type", "video/mp4")
		c.Header("Transfer-Encoding", "chunked")
		var req types.VideoShowReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetVideoSrv()
			if config.Config.System.LocalMode == "local" {
				resp := make(chan []byte, 2)
				go l.ShowVideoContent(ctx, &req, resp)
				c.Response.HijackWriter(resp2.NewChunkedBodyWriter(&c.Response, c.GetWriter()))
				for i := 0; ; {
					buf := <-resp
					i++
					if buf != nil {
						_, _ = c.Write(buf)
						if i%5 == 0 {
							_ = c.Flush()
						}
					} else {
						_ = c.Flush()
						break
					}
				}
				return
			} else {
				resp := make(chan []byte, 2)
				go l.ShowVideoContent(ctx, &req, resp)
				for {
					buf := <-resp
					if buf == nil {
						break
					}
					_, err = c.Write(buf)
					if err != nil {
						log.LogrusObj.Errorln(err)
						return
					}
					c.Flush()
				}
				return
			}

		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func VideoShowHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.VideoShowReq
		err := c.Bind(&req)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		if err == nil {
			l := service.GetVideoSrv()
			resp, err := l.Show(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func VideoCommentHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.VideoCommentReq
		err := c.Bind(&req)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		if err == nil {
			l := service.GetVideoSrv()
			resp, err := l.Comment(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func VideoCommentReplyHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.VideoCommentReq
		err := c.Bind(&req)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		if err == nil {
			l := service.GetVideoSrv()
			resp, err := l.Reply(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}

func VideoDeleteHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.VideoDeleteReq
		err := c.Bind(&req)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}
		if err == nil {
			l := service.GetVideoSrv()
			resp, err := l.Delete(ctx, &req)
			if err != nil {
				log.LogrusObj.Errorln(err)
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
			c.JSON(http.StatusOK, resp)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}
