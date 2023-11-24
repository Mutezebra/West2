package middleware

import (
	"errors"
	"five/consts"
	"five/pkg/ctl"
	"five/pkg/e"
	"five/pkg/log"
	"five/pkg/myutils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	EmptyToken = errors.New("token is empty")
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code = e.SUCCESS
		aToken := c.GetHeader("access_token")
		rToken := c.GetHeader("refresh_token")
		if aToken == "" {
			code = e.InvalidParams
			err := EmptyToken
			log.LogrusObj.Error(err)
			resp := ctl.RespError(code, err)
			c.JSON(http.StatusBadRequest, resp)
			c.Abort()
			return
		}

		newAToken, newRToken, err := myutils.CheckToken(aToken, rToken)
		if err != nil {
			code = e.CheckTokenFailed
			log.LogrusObj.Error(err)
			resp := ctl.RespError(code, err)
			c.JSON(http.StatusInternalServerError, resp)
			c.Abort()
			return
		}

		claims, err, _ := myutils.ParseToken(newAToken)
		if err != nil {
			code = e.ParseTokenFailed
			log.LogrusObj.Error(err)
			resp := ctl.RespError(code, err)
			c.JSON(http.StatusInternalServerError, resp)
			c.Abort()
			return
		}
		setHeader(c, newAToken, newRToken)
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{UID: claims.UID, UserName: claims.UserName}))
		c.Next()
	}
}

func setHeader(c *gin.Context, aToken, rToken string) {
	c.Header(consts.AccessToken, aToken)
	c.Header(consts.RefreshToken, rToken)
}
