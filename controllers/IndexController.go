package controllers

import beego "github.com/beego/beego/v2/server/web"

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	if !IsUserLoggedIn(&c.Controller) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "index.html"
}
