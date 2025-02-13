package controllers

import (
	"CRMsystemproject/util"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	redisUtil := util.NewRedisCache()
	username, _ := redisUtil.Get("logged_user")
	if username != nil {
		// 用户名存在，返回给前端页面
		c.Data["Username"] = username
	} else {
		// 用户名不存在，可以返回错误信息或者处理匿名用户的情况
		c.Data["Username"] = "Anonymous"
	}

	//c.Data["Website"] = "beego.vip"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
