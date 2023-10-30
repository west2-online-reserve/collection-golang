package api

import (
	"context"
	"four/pkg/log"
	"four/service"
	"four/types"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

func Search() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.SearchReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetSearchSrv()
			resp, err := l.Search(ctx, &req)
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

func AuthSearchHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.SearchReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetSearchSrv()
			resp, err := l.Search(ctx, &req)
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

func FilterHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.FilterReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetSearchSrv()
			resp, err := l.Filter(ctx, &req)
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

func AuthFilterHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req types.FilterReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetSearchSrv()
			resp, err := l.Filter(ctx, &req)
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

func HistorySearchItemsHandle() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		l := service.GetSearchSrv()
		resp, err := l.HistorySearchItems(ctx)
		if err != nil {
			log.LogrusObj.Errorln(err)
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}
