package controllers

import (
	"CRMsystemproject/models"
	"CRMsystemproject/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

const (
	USERNAME = "username"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	// 生成随机的4位数字captcha
	randomCaptcha := generateRandomCaptcha()
	c.Data["RandomCaptcha"] = randomCaptcha

	c.TplName = "login.html"
}

/*
 * @summary 处理登录请求
 * @description 获取用户输入，验证验证码，验证用户名和密码，设置 Session，跳转页面
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *LoginController) Post() {
	// 获取用户输入
	userInput := c.GetString("captcha")
	// 从服务器获取生成的captcha
	generatedCaptcha := c.GetString("randomCaptcha")

	if userInput != generatedCaptcha {
		errCode := 407
		errorResponse := util.NewError(errCode)
		c.Ctx.Output.JSON(errorResponse, false, false)
		//c.Data["json"] = errorResponse
		c.Ctx.WriteString("验证码错误！")
		c.TplName = "login.html"
		return
	}

	//1.拿到数据
	useName := c.GetString("username")
	pwd := c.GetString("password")
	//2.判断数据是否合法
	if useName == "" || pwd == "" {
		c.Ctx.WriteString("请正确输入数据！")
		c.TplName = "login.html"
		return
	}
	//3.查询账号密码是否正确
	o := orm.NewOrm()
	user := models.Users{}
	user.Username = useName

	// 注意这里不再直接将用户输入的密码赋值给user.Password
	err := o.Read(&user, "username")
	if err != nil {
		errCode := 501
		errorResponse := util.NewError(errCode)
		//c.Data["json"] = errorResponce
		c.Ctx.Output.JSON(errorResponse, false, false)
		return
		//c.Ctx.WriteString("用户名输入错误！")
		//c.TplName = "login.html"
	}
	fmt.Println("数据库中的密码哈希：", user.Password)
	fmt.Println("用户输入的密码：", pwd)

	// 使用bcrypt进行密码验证  用bcrypt.CompareHashAndPassword来验证密码是否匹配
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd)); err != nil {
		fmt.Println("密码比较错误：", err)
		// 密码比较错误，设置错误提示
		errCode := 502
		errorResponse := util.NewError(errCode)
		c.Ctx.Output.JSON(errorResponse, false, false)
		//c.Data["json"] = errorResponce
		//c.Data["Error"] = template.JSEscapeString("密码输入错误！")
		c.TplName = "login.html"
		return
	}

	// 假设验证成功，将用户名存储在Redis中
	redisCache := util.NewRedisCache()           // 实例化Redis缓存工具类
	redisCache.Set("logged_user", user.Username) // 将用户名存储在Redis中
	//登录验证通过，设置 Session
	c.SetSession(USERNAME, user.Username)
	fmt.Println("设置session，用户名:", user.Username)
	//4.跳转
	//c.Ctx.WriteString("欢迎您，登陆成功！")
	//c.Redirect("/userlist", 302)
	c.TplName = "index.html"
}

// generateRandomCaptcha 生成随机验证码
/*
 * @summary 生成随机验证码
 * @return 返回生成的随机验证码字符串
 */
func generateRandomCaptcha() string {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	//生成随机的4位数字captcha
	randomCaptcha := rand.Intn(10000)
	return strconv.Itoa(randomCaptcha)
}
