package controllers

import (
	"CRMsystemproject/models"
	"CRMsystemproject/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type ChurnCustomerController struct {
	beego.Controller
}

func (c *ChurnCustomerController) Get() {
	c.TplName = "churn_customer_management.html"
}

// GetAll 从数据库获取所有流失客户信息
/*
 * @summary 从数据库获取所有流失客户信息
 * @return 返回所有流失客户信息的 JSON 格式响应数据给前端
 */
func (c *ChurnCustomerController) GetAll() {
	// 创建 ORM 实例
	o := orm.NewOrm()

	// 定义存储流失客户信息的切片
	var churncustomers []models.Churncustomers

	// 查询所有流失客户信息
	_, err := o.QueryTable("Churncustomers").All(&churncustomers)

	// 处理查询过程中的错误
	if err != nil {
		// 定义错误代码
		errCode := 408

		// 创建错误响应实例
		errorResponse := util.NewError(errCode)

		// 将错误信息设置到控制器的响应数据中
		c.Data["json"] = errorResponse

		// 打印错误信息到控制台
		fmt.Println("查询流失客户失败：", err)
	}

	// 将流失客户信息设置到控制器的响应数据中
	c.Data["json"] = churncustomers

	// 提供 JSON 格式的响应数据给前端
	c.ServeJSON()
}

// SearchChurnCustomer 执行根据搜索词查询流失客户的操作
/*
 * @summary 执行根据搜索词查询流失客户的操作
 * @param searchRequest 包含搜索词的结构体
 * @return 返回符合搜索词的流失客户信息的 JSON 格式响应数据给前端
 */
func (c *ChurnCustomerController) SearchChurnCustomer() {
	// 从请求正文中解析搜索词
	var searchRequest struct {
		SearchTerm string `json:"searchTerm"`
	}

	// 解析 JSON 请求体并检查错误
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &searchRequest)
	if err != nil {
		// 定义错误代码
		errCode := 410

		// 创建错误响应实例
		errorResponse := util.NewError(errCode)

		// 将错误信息设置到控制器的响应数据中
		c.Data["json"] = errorResponse
		//c.Data["json"] = CustomerResponse{Success: false, Message: "解析搜索请求失败"}
		// 提供 JSON 格式的错误响应给前端
		c.ServeJSON()
		return
	}

	// 调用数据库查询方法来检索符合搜索词的流失客户信息
	churncustomers, err := models.GetChurnCustomersBySearchTerm(searchRequest.SearchTerm)
	if err != nil {
		// 定义错误代码
		errCode := 411

		// 创建错误响应实例
		errorResponse := util.NewError(errCode)

		// 将错误信息设置到控制器的响应数据中
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{"success": false, "message": "搜索客户失败"}
		// 提供 JSON 格式的错误响应给前端
		c.ServeJSON()
		return
	}

	// 将成功的响应数据设置到控制器的响应数据中
	c.Data["json"] = map[string]interface{}{"success": true, "data": churncustomers}

	// 提供 JSON 格式的成功响应给前端
	c.ServeJSON()
}
