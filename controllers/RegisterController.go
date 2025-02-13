package controllers

import (
	"CRMsystemproject/models"
	"CRMsystemproject/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {

	c.TplName = "register.html"

}

// Post 处理注册请求
/*
 * @summary 处理用户注册请求
 * @description 从请求中获取用户名和密码，对数据进行校验，然后插入数据库中，最后重定向到登录界面
 */
func (c *RegisterController) Post() {

	//1.拿到数据
	username := c.GetString("username")
	password := c.GetString("password")
	//2.对数据进行校验
	if username == "" || password == "" {
		//打印出:数据不能为空
		errCode := 406
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("数据不能为空！")
		//c.Redirect("/register", 302) //重新访问
		return
	}
	//3.插入数据库
	o := orm.NewOrm() //要有ORM对象

	user := models.Users{} // 要有一个插入数据的结构体对象
	//加密前密码
	fmt.Println("加密前密码：", password)
	// 生成密码哈希
	passwordHash, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		fmt.Println(err1)
		errCode := 406
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		c.ServeJSON()

		return

	}
	// 对结构体赋值
	user.Username = username
	user.Password = string(passwordHash)
	fmt.Println("加密后的密码：", user.Password)
	_, err2 := o.Insert(&user) //取地址
	if err2 != nil {
		// 打印出：插入数据失败
		fmt.Println("插入数据失败:", err2)
		//c.Redirect("/register", 302)
		errCode := 400
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		c.ServeJSON()

		return
	}
	//4.进入登录界面
	c.Redirect("/login", 302)
	//c.TplName = "login.html"
	//c.Ctx.WriteString("注册成功")
}
