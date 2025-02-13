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

type ServiceManagementController struct {
	beego.Controller
}

type PageParam3 struct {
	Pagesize int `json:"pagesize"` //每页显示多少条
	Pagenum  int `json:"pagenum"`  //第几页
}

type ServiceManagementResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *ServiceManagementController) Get() {
	c.TplName = "service_management.html"
}

// GetAll 获取所有服务信息
/*
 * @summary 获取所有服务信息
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和服务信息列表
 */
func (c *ServiceManagementController) GetAll() {
	o := orm.NewOrm()
	var services []models.Services
	_, err := o.QueryTable("Services").All(&services)

	if err != nil {
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("查询服务信息失败：", err)
	}
	c.Data["json"] = services
	c.ServeJSON()
}

// GetService  处理 GET 请求以获取客户详细信息的方法
/*
 * @summary 处理 GET 请求以获取客户详细信息
 * @param serviceid 服务ID
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和服务详细信息
 */
func (c *ServiceManagementController) GetService() {
	// 从 URL 参数中获取客户ID
	Serviceid := c.Ctx.Input.Param(":serviceid")

	// 查询数据库，获取服务详细信息
	o := orm.NewOrm()
	var service models.Services
	err := o.QueryTable("Services").Filter("Serviceid", Serviceid).One(&service)
	if err != nil {
		// 处理错误，例如，返回带有错误消息的 JSON 响应
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{
		//	"error": "无法获取服务详细信息！",
		//}
		c.ServeJSON()
		return
	}
	// 返回客户详细信息作为 JSON
	c.Data["json"] = service
	c.ServeJSON()
}

// Post 添加服务
/*
 * @summary 添加服务
 * @param serviceid 服务ID
 * @param servicename 服务名称
 * @param customername 客户名称
 * @param servicestatus 服务状态
 * @param servicecreator 服务创建人
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *ServiceManagementController) Post() {
	fmt.Println("添加服务需求Post")
	// 拿到前端输入数据
	serviceid, _ := c.GetInt("serviceid")
	servicename := c.GetString("servicename")
	customername := c.GetString("customername")
	servicestatus := c.GetString("servicestatus")
	servicecreator := c.GetString("servicecreator")

	fmt.Println("customername", customername)
	// 对数据进行校验
	if serviceid == 0 || servicename == "" || customername == "" || servicestatus == "" || servicecreator == "" {
		errCode := 409
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("数据不能为空！")
		c.Redirect("/service_management", 302)
		return
	}
	// 插入数据库
	o := orm.NewOrm()
	services := models.Services{} // 插入数据库的结构体对象
	// 对结构体赋值
	services.Serviceid = serviceid
	services.Servicename = servicename
	services.Customername = customername
	services.Servicestatus = servicestatus
	services.Servicecreator = servicecreator
	_, err := o.Insert(&services)
	if err != nil {
		fmt.Println("服务添加失败！")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "服务添加失败：数据不能为空",
		}
		c.ServeJSON()
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "服务添加成功！",
	}
	c.ServeJSON()
}

// DeleteService 删除服务功能
/*
 * @summary 删除服务功能
 * @param serviceid 服务ID
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *ServiceManagementController) DeleteService() {
	serviceid, err := c.GetInt(":serviceid") // 获取到 url 当中 id 变量的值
	if err != nil {
		c.Data["json"] = ServiceManagementResponse{Success: false, Message: "获取参数失败"}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm() // 创建一个orm对象

	_, err = o.Delete(&models.Services{Serviceid: serviceid})
	if err != nil { // 如果删除错误就返回信息："删除数据失败"
		errCode := 413
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "删除数据失败"}
		fmt.Println("删除失败")
		c.ServeJSON()
		return
	}
	c.Data["json"] = ServiceManagementResponse{Success: true, Message: "删除成功"}
	fmt.Println("删除成功", serviceid)
	c.ServeJSON()
}

// UpdateService 修改服务功能
/*
 * @summary 修改服务功能
 * @param services 包含服务信息的结构体
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *ServiceManagementController) UpdateService() {
	var services models.Services
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &services)
	if err != nil {
		fmt.Println("解析正文请求时出错：", err)
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	service := models.Services{Serviceid: services.Serviceid}
	if o.Read(&service) == nil {
		fmt.Println("找到了要修改的服务：", service)
		service.Serviceid = services.Serviceid
		service.Servicename = services.Servicename
		service.Customername = services.Customername
		service.Servicestatus = services.Servicestatus
		service.Servicecreator = services.Servicecreator
		fmt.Println()
		if _, err := o.Update(&service); err == nil {
			c.Data["json"] = ServiceManagementResponse{Success: true, Message: "更新成功"}
		} else {
			errCode := 414
			errorResponse := util.NewError(errCode)
			c.Data["json"] = errorResponse
			fmt.Println("更新失败:", err)
			//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "更新失败"}
		}
	} else {
		errCode := 415
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("该服务不存在 Serviceid:", services.Serviceid)
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "该服务不存在"}
	}

	c.ServeJSON()
}

// SearchService 查找功能
/*
 * @summary 查找服务功能
 * @param searchRequest3 包含搜索词字段的结构体
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和服务信息列表
 */
func (c *ServiceManagementController) SearchService() {
	// 从正文中解析搜索词
	fmt.Println("进入查找方法")
	var searchRequest3 struct {
		SearchTerm string `json:"searchTerm"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &searchRequest3)
	if err != nil {
		errCode := 410
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "解析搜索请求失败"}
		c.ServeJSON()
		return
	}
	// 调用数据库查询方法来检索搜索结果
	fmt.Println("searchRequest3.SearchTerm:", searchRequest3.SearchTerm)
	services, err := models.GetServicesBySearchTerm(searchRequest3.SearchTerm)
	if err != nil {
		errCode := 411
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "搜索服务失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"success": true, "data": services}
	c.ServeJSON()
}

// ShowServiceByPage  分页功能
/*
 * @summary 分页功能
 * @param pageinfo3 包含分页信息的结构体
 * @return 返回 JSON 格式的响应数据给前端，包括分页结果和总服务数
 */
func (c *ServiceManagementController) ShowServiceByPage() {
	var pageinfo3 PageParam3
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo3)
	if err != nil {
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = ServiceManagementResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := int64(pageinfo3.Pagenum)
	pageSize := int64(pageinfo3.Pagesize)

	// 查询总数并获取分页结果
	totalCount, _ := models.GetServiceCount()
	pagination3 := util.NewPagination(totalCount, currentPage, pageSize)
	offset := (pagination3.CurrentPage - 1) * pagination3.PageSize
	services, _ := models.GetService(int(offset), int(pagination3.PageSize))

	var repList []interface{}
	for _, service := range services {
		repList = append(repList, service.ServiceToRespDesc())
	}

	// 将分页结果和总服务数一起返回到前端
	pagination3.Data = repList
	c.Data["json"] = map[string]interface{}{
		"pagination":  pagination3,
		"total_pages": totalCount,
		"services":    repList,
	}
	c.ServeJSON()

}
func (c *ServiceManagementController) GetServiceList() {

	var pageinfo PageParam3
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		c.Data["json"] = ServiceManagementResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := pageinfo.Pagenum
	pageSize := pageinfo.Pagesize
	o := orm.NewOrm()
	totalCount, _ := o.QueryTable("services").Count()
	fmt.Println("服务总数量：", totalCount)

	//totalCount, _ := models.GetUserCount()
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}
	var services []models.Services
	o.QueryTable("services").Limit(pageSize, (currentPage-1)*pageSize).All(&services)

	var repList []interface{}
	for _, service := range services {
		repList = append(repList, service.ServiceToRespDesc())
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
func (c *ServiceManagementController) ShowServiceList() {

	var pageinfo PageParam3
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		c.Data["json"] = ServiceManagementResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := int64(pageinfo.Pagenum)
	pageSize := int64(pageinfo.Pagesize)
	// 查询总记录数
	totalCount, _ := models.GetServiceCount()
	// 创建分页实例
	pagination := util.NewPagination(totalCount, currentPage, pageSize)
	// 查询当前页的数据
	offset := (pagination.CurrentPage - 1) * pagination.PageSize
	services, _ := models.GetCustomerlist(offset, pagination.PageSize)
	var repList []interface{}
	for _, service := range services {
		repList = append(repList, service.CustomerToRespDesc())
	}
	pagination.Data = repList
	// 返回分页信息给前端
	c.Data["json"] = pagination
	c.ServeJSON()
}

// CreateOrder 生成订单功能
func (c *ServiceManagementController) CreateOrder() {
	fmt.Println("生成订单CreateOrder")
	// 从前端获取数据
	//orderid, _ := c.GetInt("orderid")
	//orderid := c.GetString("orderid")
	//ordername := c.GetString("ordername")
	//orderstatus := c.GetString("orderstatus")
	//customername := c.GetString("customername")
	//orderaddress := c.GetString("orderaddress")
	//orderprice, _ := c.GetFloat("orderprice")
	//ordertime := c.GetString("ordertime")

	var salesorders models.Salesorders

	err1 := json.Unmarshal(c.Ctx.Input.RequestBody, &salesorders)

	fmt.Println("err", err1)
	fmt.Println("salesorders.orderid", salesorders.Orderid)

	// 插入数据库
	o := orm.NewOrm()
	//检查订单编号是否已存在
	existingOrder := models.Salesorders{Orderid: salesorders.Orderid}
	err := o.QueryTable("Salesorders").Filter("orderid", salesorders.Orderid).One(&existingOrder)

	if err == nil {
		// 处理订单编号已存在的情况，例如返回错误响应给客户端
		fmt.Println("err:", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "订单编号已存在",
		}
		c.ServeJSON()
		return
	}

	_, err2 := o.Insert(&salesorders)
	if err2 != nil {
		fmt.Println("订单添加失败！")
		fmt.Println("err1:", err2)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "订单添加失败，数据不能为空！",
		}
	}
	// 返回成功的 JSON 响应
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "订单生成成功！",
	}
	c.ServeJSON()
}
