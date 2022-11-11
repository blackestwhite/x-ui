package controller

import (
	"fmt"
	"strconv"
	"time"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/web/job"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type APIController struct {
	userService    service.UserService
	inboundService service.InboundService
	xrayService    service.XrayService
}

func (a *APIController) auth(c *gin.Context) {
	u := c.GetHeader("x-api-username")
	p := c.GetHeader("x-api-password")

	user := a.userService.CheckUser(u, p)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	if user == nil {
		job.NewStatsNotifyJob().UserLoginNotify(u, getRemoteIp(c), timeStr, 0)
		logger.Infof("wrong username or password: \"%s\" \"%s\"", u, p)
		c.JSON(501, gin.H{
			"message": "wrong username or password",
		})
		c.Abort()
		return
	} else {
		logger.Infof("%s login success,Ip Address:%s\n", u, getRemoteIp(c))
		job.NewStatsNotifyJob().UserLoginNotify(p, getRemoteIp(c), timeStr, 1)
	}

	err := session.SetLoginUser(c, user)
	if err != nil {
		c.JSON(501, gin.H{
			"message": "couldn't set login user",
		})
		c.Abort()
	}

	logger.Info("user", user.Id, "login success")
	c.Next()
}

func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

func (a *APIController) initRouter(g *gin.RouterGroup) {
	api := g.Group("/api")
	api.Use(a.auth)

	i := g.Group("/inbound")
	i.POST("/list", a.getInbounds)
	i.POST("/add", a.addInbound)
	i.POST("/del/:id", a.delInbound)
	i.POST("/update/:id", a.updateInbound)
}

func (a *APIController) getInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)
	if err != nil {
		jsonMsg(c, "获取", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *APIController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	err := c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "添加", err)
		return
	}
	user := session.GetLoginUser(c)
	inbound.UserId = user.Id
	inbound.Enable = true
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	err = a.inboundService.AddInbound(inbound)
	jsonMsg(c, "添加", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *APIController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "删除", err)
		return
	}
	err = a.inboundService.DelInbound(id)
	jsonMsg(c, "删除", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *APIController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	inbound := &model.Inbound{
		Id: id,
	}
	err = c.ShouldBind(inbound)
	if err != nil {
		jsonMsg(c, "修改", err)
		return
	}
	err = a.inboundService.UpdateInbound(inbound)
	jsonMsg(c, "修改", err)
	if err == nil {
		a.xrayService.SetToNeedRestart()
	}
}
