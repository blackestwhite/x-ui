package controller

import (
	"net/http"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	userService service.UserService
}

func (a *BaseController) checkLogin(c *gin.Context) {
	if Auth(a.userService, c) || session.IsLogin(c) {
		c.Next()
	} else {
		if isAjax(c) {
			pureJsonMsg(c, false, "登录时效已过，请重新登录")
		} else {
			c.Redirect(http.StatusTemporaryRedirect, c.GetString("base_path"))
		}
		c.Abort()
	}
}
