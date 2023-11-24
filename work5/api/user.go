package api

import (
	"five/pkg/log"
	"five/service"
	"five/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegisterHandle 用户注册
func UserRegisterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserRegisterReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Register(c.Request.Context(), &req)
			if err != nil {
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

// UserLoginHandle 用户登录
func UserLoginHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserLoginReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.Login(c.Request.Context(), &req)
			if err != nil {
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

func CreateGroupHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CreateGroupReq
		err := c.Bind(&req)
		if err == nil {
			l := service.GetUserSrv()
			resp, err := l.CreateGroup(c.Request.Context(), &req)
			if err != nil {
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
