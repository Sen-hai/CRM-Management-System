package controllers

import (
	"CRMsystemproject/models"
	"CRMsystemproject/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	beego "github.com/beego/beego/v2/server/web"
)

type OrdersController struct {
	beego.Controller
}

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *OrdersController) Get() {
	c.TplName = "order_view.html"
}

/*
 * @summary 获取所有销售订单
 * @description 通过 ORM 查询表 "salesorders" 中的所有数据，并将结果返回给前端
 */
func (c *OrdersController) GetAll() {
	// 创建 ORM 对象
	o := orm.NewOrm()
	// 定义一个切片，用于存储查询到的销售订单数据
	var salesorders []models.Salesorders
	// 使用 ORM 查询表 "salesorders" 中的所有数据，并将结果存储到 salesorders 切片中
	_, err := o.QueryTable("Salesorders").All(&salesorders)
	if err != nil {
		// 如果查询出错，打印错误信息
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("查询销售订单失败：", err)
	}
	// 将查询结果传递到前端界面，通过 c.Data 设置模板变量
	c.Data["json"] = salesorders
	c.ServeJSON()
}

// SearchOrder 查找功能
/*
 * @summary 根据搜索词查找销售订单
 * @param searchRequest 搜索请求结构体，包含搜索词字段
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *OrdersController) SearchOrder() {
	// 从请求正文中解析搜索词
	var searchRequest struct {
		SearchTerm string `json:"searchTerm"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &searchRequest)
	if err != nil {
		errCode := 410
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = OrderResponse{Success: false, Message: "解析搜索请求失败"}
		c.ServeJSON()
		return
	}

	// 调用数据库查询方法来检索搜索结果
	salesorders, err := models.GetOrdersBySearchTerm(searchRequest.SearchTerm)
	if err != nil {
		errCode := 408
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{"success": false, "message": "搜索订单失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "data": salesorders}
	c.ServeJSON()
}
