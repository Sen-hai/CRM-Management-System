package controllers

import (
	"CRMsystemproject/models"
	"CRMsystemproject/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
	"math"
)

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type PageParam struct {
	Pagesize int `json:"pagesize"` //每页显示多少条
	Pagenum  int `json:"pagenum"`  //第几页
}

type UserController struct {
	beego.Controller
}

/*
 * IsUserLoggedIn 检查用户是否已登录
 * @param c beego.Controller类型的指针，用于获取用户的session和控制器信息
 * @return 如果用户已登录，返回true；否则，返回false
 */
func IsUserLoggedIn(c *beego.Controller) bool {
	// 获取 session 中的用户名
	username := c.GetSession(USERNAME)
	fmt.Println("Session中的用户名:", username)
	// 获取当前控制器和操作的名称
	controllerName, actionName := c.GetControllerAndAction()
	fmt.Println("当前控制器和操作:", controllerName, actionName)
	// 不需要登录检查的操作列表
	nonAuthActions := map[string][]string{
		"userListController": {"user_management.html", "/customer_management", "/getall", "/getuser/:userid"},
	}
	// 检查当前控制器和操作是否在不需要身份验证的操作列表中
	exemptActions, exists := nonAuthActions[controllerName]
	if exists {
		for _, exemptAction := range exemptActions {
			if exemptAction == actionName {
				fmt.Println("当前操作无需身份验证")
				return true
			}
		}
	}
	//如果不在列表中，则执行登录检查
	// 判断用户名是否存在
	if username == nil {
		fmt.Println("用户未登录，执行重定向")
		return false
	}
	fmt.Println("用户已登录")
	return true
}

/*
 * Get 处理用户请求的方法，检查用户登录状态，未登录则重定向到登录页面
 */
func (c *UserController) Get() {
	if !IsUserLoggedIn(&c.Controller) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "user_management.html"
}

/*
 * GetAll 获取所有用户信息的方法
 */
func (c *UserController) GetAll() {
	var users []models.Users
	o := orm.NewOrm()
	o.QueryTable("users").All(&users)
	var respList []interface{}
	for _, user := range users {
		respList = append(respList, user.UserToRespDesc())
	}
	c.Data["json"] = respList
	c.ServeJSON()
}

/*
 * GetUserInfo 获取指定用户信息的方法
 * @param userid 用户ID
 */
func (c *UserController) GetUserInfo() {
	// 获取 URL 中的用户 ID
	userid, err := c.GetInt(":userid")
	if err != nil {
		fmt.Println(err)
		c.Data["json"] = LoginResponse{Success: false, Message: "获取用户ID失败"}
		c.ServeJSON()
		return
	}
	// 查询数据库获取用户信息
	o := orm.NewOrm()
	var user models.Users
	err = o.QueryTable("users").Filter("Userid", userid).One(&user)
	if err != nil {
		errCode := 407
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "获取用户信息失败"}
		c.ServeJSON()
		return
	}

	// 返回用户信息
	c.Data["json"] = user.UserToRespDesc()
	c.ServeJSON()
}

/*
 * AddUsers 添加用户的方法
 * @param userid 用户ID
 * @param username 用户名
 * @param password 用户密码
 */
func (c *UserController) AddUsers() {
	//1.拿到数据
	userid, _ := c.GetInt("userid")
	username := c.GetString("username")
	password := c.GetString("password")
	//2.对数据进行校验
	if userid == 0 || username == "" || password == "" {
		errCode := 409
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//打印出:数据不能为空
		fmt.Println("数据不能为空！")
		c.Redirect("/register", 302) //重新访问
		return
	}
	//3.插入数据库
	o := orm.NewOrm() //要有ORM对象

	user := models.Users{} // 要有一个插入数据的结构体对象
	// 对结构体赋值
	user.Userid = userid
	user.Username = username
	user.Password = password
	_, err := o.Insert(&user) //取地址
	if err != nil {
		// 打印出：插入数据失败
		fmt.Println("插入数据失败:", err)
		errCode := 416
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{
		//	"success": false,
		//	"message": "用户添加失败，数据不能为空！",
		//}
		c.ServeJSON()
	}
}

/*
 * SearchUser 用户模糊查找的方法
 * @param searchRequest 包含搜索词的结构体
 */
func (c *UserController) SearchUser() {
	// 从正文中解析搜索词
	fmt.Println("进入用户模糊查找方法")
	var searchRequest struct {
		SearchTerm string `json:"searchTerm"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &searchRequest)
	if err != nil {
		errCode := 411
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "解析搜索请求失败"}
		c.ServeJSON()
		return
	}
	// 调用数据库查询方法来检索搜索结果
	fmt.Println("searchRequest.SearchTerm:", searchRequest.SearchTerm)
	users, err := models.GetUsersBySearchTerm(searchRequest.SearchTerm)
	if err != nil {
		errCode := 411
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "搜索用户失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"success": true, "data": users}
	c.ServeJSON()
}

/*
 * UpdateUser 修改用户信息的方法
 * @param users 包含新用户信息的结构体
 */
func (c *UserController) UpdateUser() {
	var users models.Users
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &users)
	if err != nil {
		errMsg := util.NewError(504)
		c.Data["json"] = LoginResponse{Success: false, Message: errMsg.Message}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	user := models.Users{Userid: users.Userid}
	if o.Read(&user) == nil {
		fmt.Println("找到了存在的用户:", user)
		user.Username = users.Username
		user.Password = users.Password
		fmt.Println()
		if _, err := o.Update(&user); err == nil {
			c.Data["json"] = LoginResponse{Success: true, Message: "更新成功"}
		} else {
			fmt.Println("更新失败:", err)
			errCode := 414
			errorResponse := util.NewError(errCode)
			c.Data["json"] = errorResponse
			//c.Data["json"] = LoginResponse{Success: false, Message: "更新失败"}
		}
	} else {
		fmt.Println("用户不存在 UserID:", users.Userid)
		errCode := 501
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "用户不存在"}
	}

	c.ServeJSON()
}

/*
 * DeleteUser 删除用户的方法
 * @param userid 用户ID
 */
func (c *UserController) DeleteUser() {
	userid, err := c.GetInt(":userid") // 获取到 url 当中 id 变量的值
	if err != nil {                    // 有错误就返回数据：获取参数失败
		c.Data["json"] = LoginResponse{Success: false, Message: "获取参数失败"}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm() // 创建一个orm对象
	// 调用 orm 的 Delete 方法，&models.Userinfo{Id: id} 表示删除的是哪一个跟数据库相关的模型以及限制条件
	_, err = o.Delete(&models.Users{Userid: userid})
	if err != nil { // 如果删除错误就返回信息："删除数据失败"
		errCode := 414
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "删除数据失败"}
		fmt.Println("删除失败")
		c.ServeJSON()
		return
	}
	c.Data["json"] = LoginResponse{Success: true, Message: "删除成功"}
	fmt.Println("删除成功", userid)
	c.ServeJSON()
}

/*
 * ShowUserByPage 分页展示用户信息的方法
 * @param pageinfo 包含分页信息的结构体
 */
func (c *UserController) ShowUserByPage() {
	var pageinfo PageParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	page := pageinfo.Pagenum
	pagesize := pageinfo.Pagesize
	users, err := models.GetUsers(page, pagesize)
	if err != nil {
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = LoginResponse{Success: false, Message: "获取数据失败"}
		c.ServeJSON()
		return
	}
	if len(users) > 0 {
		var repList []interface{}
		for _, user := range users {
			repList = append(repList, user.UserToRespDesc())
		}
		// 构建统一响应结构体
		response := models.UnifiedResponse{
			//Code:    util.RECODE_OK,
			Message: "获取用户成功",
			Data:    repList,
		}
		c.Data["json"] = response

	}
	c.ServeJSON()
}
func (c *UserController) GetUserList() {

	var pageinfo PageParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		c.Data["json"] = LoginResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := pageinfo.Pagenum
	pageSize := pageinfo.Pagesize
	o := orm.NewOrm()
	totalCount, _ := o.QueryTable("users").Count()
	fmt.Println("用户总数量：", totalCount)

	//totalCount, _ := models.GetUserCount()
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}
	var users []models.Users
	o.QueryTable("users").Limit(pageSize, (currentPage-1)*pageSize).All(&users)

	var repList []interface{}
	for _, user := range users {
		repList = append(repList, user.UserToRespDesc())
	}
	// 构建分页结构体
	pagination := &models.Pagination{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    pageSize,
		Data:        repList,
	}

	// 返回分页信息给前端
	c.Data["json"] = pagination
	c.ServeJSON()
}
func (c *UserController) ShowList() {

	var pageinfo PageParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		c.Data["json"] = LoginResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := int64(pageinfo.Pagenum)
	pageSize := int64(pageinfo.Pagesize)
	// 查询总记录数
	totalCount, _ := models.GetUserCount()
	// 创建分页实例
	pagination := util.NewPagination(totalCount, currentPage, pageSize)
	// 查询当前页的数据
	offset := (pagination.CurrentPage - 1) * pagination.PageSize
	users, _ := models.GetUserlist(offset, pagination.PageSize)
	var repList []interface{}
	for _, user := range users {
		repList = append(repList, user.UserToRespDesc())
	}
	pagination.Data = repList
	// 返回分页信息给前端
	c.Data["json"] = pagination
	c.ServeJSON()
}
