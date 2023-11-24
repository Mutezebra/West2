package api

import (
	"five/pkg/log"
	"five/service"
	"five/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.WsConnectReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetWsSrv()
			l.Connect(c.Request.Context(), &req, c.Writer, c.Request)
			c.JSON(http.StatusOK, nil)
			return
		}
		log.LogrusObj.Errorln(err)
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
}
