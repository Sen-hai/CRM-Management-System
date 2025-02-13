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

type CustomerController struct {
	beego.Controller
}

type PageParam2 struct {
	Pagesize int `json:"pagesize"` //每页显示多少条
	Pagenum  int `json:"pagenum"`  //第几页
}

type CustomerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *CustomerController) Get() {
	c.TplName = "customer_management.html"
}

/*
 * @summary 查询数据库中的所有客户数据
 * @return 返回所有客户数据的 JSON 格式响应数据给前端
 */
func (c *CustomerController) GetAll() {
	// 查询数据库中的客户数据
	o := orm.NewOrm()
	var customers []models.Customers
	_, err := o.QueryTable("Customers").All(&customers)

	if err != nil {
		errCode := 405
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("查询客户数据失败:", err)
		return
	}
	// 将查询结果传递给前端页面
	c.Data["json"] = customers
	c.ServeJSON()
}

// GetCustomer 处理 GET 请求以获取客户详细信息的方法
/*
 * @summary 处理 GET 请求以获取客户详细信息的方法
 * @param customerID 客户ID，从 URL 参数中获取
 * @return 返回客户详细信息作为 JSON 给前端
 */
func (c *CustomerController) GetCustomer() {
	// 从 URL 参数中获取客户ID
	customerID := c.Ctx.Input.Param(":customerId")

	// 查询数据库，获取客户详细信息
	o := orm.NewOrm()
	var customer models.Customers
	err := o.QueryTable("customers").Filter("Customerid", customerID).One(&customer)
	if err != nil {
		// 处理错误
		errCode := 407
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{
		//	"error": "无法获取客户详细信息！",
		//}
		c.ServeJSON()
		return
	}

	// 返回客户详细信息作为 JSON
	c.Data["json"] = customer
	c.ServeJSON()
}

// AddCustomer  添加客户功能
/*
 * @summary 添加客户功能
 * @description 从前端获取客户数据，对数据进行校验后插入数据库
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) AddCustomer() {
	fmt.Println("收到添加客户的 POST 请求")
	// 拿到前端输入数据
	customerid, _ := c.GetInt("customerid")
	customername := c.GetString("customername")
	contactid, _ := c.GetInt("contactid")
	customeraddress := c.GetString("customeradress")
	customerinfo := c.GetString("customerinfo")

	fmt.Println("customerid:", customerid)
	// 对数据进行校验
	if customerid == 0 || customername == "" || contactid == 0 || customeraddress == "" || customerinfo == "" {
		errCode := 409
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("数据不能为空！")
		c.Redirect("/index", 302)
		return
	}
	//插入数据库
	o := orm.NewOrm()
	customers := models.Customers{} // 插入数据的结构体对象
	// 对结构体赋值
	customers.Customerid = customerid
	customers.Customername = customername
	customers.Contactid = contactid
	customers.Customeraddress = customeraddress
	customers.Customerinfo = customerinfo

	_, err := o.Insert(&customers)
	if err != nil {
		// 打印出添加失败
		fmt.Println("客户添加失败：", err)
		errCode := 412
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{
		//	"success": false,
		//	"message": "客户添加失败：数据不能为空",
		//}
		c.ServeJSON()
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "客户添加成功！",
	}
	c.ServeJSON()

}

// DeleteCustomer  删除功能
/*
 * @summary 删除客户功能
 * @description 从 URL 获取客户ID，调用 ORM 的 Delete 方法删除指定客户
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) DeleteCustomer() {
	customerid, err := c.GetInt(":customerid") // 获取到 url 当中 id 变量的值
	if err != nil {                            // 有错误就返回数据：获取参数失败
		errCode := 413
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("deleteerr:", err)
		//c.Data["json"] = CustomerResponse{Success: false, Message: "获取参数失败"}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm() // 创建一个orm对象
	// 调用 orm 的 Delete 方法，&models.Userinfo{Id: id} 表示删除的是哪一个跟数据库相关的模型以及限制条件
	_, err = o.Delete(&models.Customers{Customerid: customerid})
	fmt.Println("deleteerr:", err)
	if err != nil { // 如果删除错误就返回信息："删除数据失败"
		errCode := 413
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = CustomerResponse{Success: false, Message: "删除数据失败"}
		fmt.Println("删除失败", err)
		c.ServeJSON()
		return
	}
	c.Data["json"] = CustomerResponse{Success: true, Message: "删除成功"}
	fmt.Println("删除成功", customerid)
	c.ServeJSON()
}

// UpdateCustomer 修改功能
/*
 * @summary 修改客户信息功能
 * @description 从请求正文解析客户信息，调用 ORM 的 Update 方法更新客户信息
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) UpdateCustomer() {
	var customers models.Customers
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &customers)
	if err != nil {
		fmt.Println("解析请求正文时出错:", err)
		c.Data["json"] = CustomerResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	customer := models.Customers{Customerid: customers.Customerid}
	if o.Read(&customer) == nil {
		fmt.Println("找到了存在的客户:", customer)
		customer.Customername = customers.Customername
		customer.Customeraddress = customers.Customeraddress
		customer.Customerinfo = customers.Customerinfo
		customer.Contactid = customers.Contactid
		fmt.Println()
		if _, err := o.Update(&customer); err == nil {
			c.Data["json"] = CustomerResponse{Success: true, Message: "更新成功"}
		} else {
			errCode := 414
			errorResponse := util.NewError(errCode)
			c.Data["json"] = errorResponse
			fmt.Println("更新失败:", err)
			//c.Data["json"] = CustomerResponse{Success: false, Message: "更新失败"}
		}
	} else {
		errCode := 415
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		fmt.Println("该客户不存在 Customerid:", customers.Customerid)
		//c.Data["json"] = CustomerResponse{Success: false, Message: "该客户不存在"}
	}

	c.ServeJSON()
}

// SearchCustomer 查找功能
/*
 * @summary 查找客户功能
 * @description 从请求正文中解析搜索词，调用数据库查询方法检索搜索结果
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) SearchCustomer() {
	// 从请求正文中解析搜索词
	var searchRequest struct {
		SearchTerm string `json:"searchTerm"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &searchRequest)
	if err != nil {
		errCode := 410
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = CustomerResponse{Success: false, Message: "解析搜索请求失败"}
		c.ServeJSON()
		return
	}

	// 调用数据库查询方法来检索搜索结果
	customers, err := models.GetCustomersBySearchTerm(searchRequest.SearchTerm)
	if err != nil {
		errCode := 411
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{"success": false, "message": "搜索客户失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "data": customers}
	c.ServeJSON()
}

// GetContactDetails 查看联系人
/*
 * @summary 查看联系人详情
 * @description 从 URL 获取联系人 ID，查询数据库获取联系人详细信息
 * @return 返回 JSON 格式的联系人详细信息给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) GetContactDetails() {
	// 从 URL 获取联系人 ID
	contactID := c.Ctx.Input.Param(":contactId")

	// 查询数据库以获取联系人详细信息
	o := orm.NewOrm()
	var contact models.Contacts
	err := o.QueryTable("Contacts").Filter("Contactid", contactID).One(&contact)

	if err != nil {
		errCode := 408
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{"error": "无法获取联系人详细信息！"}
		c.ServeJSON()
		return
	}

	// 以 JSON 格式返回联系人详细信息
	c.Data["json"] = contact
	c.ServeJSON()
}

// MarkCustomerAsLost 流失客户功能
/*
 * @summary 流失客户功能
 * @description 从前端获取流失客户数据，插入数据库
 * @return 返回 JSON 格式的响应数据给前端，包括成功或失败的状态和消息
 */
func (c *CustomerController) MarkCustomerAsLost() {
	fmt.Println("流失客户插入表")
	// 从前端获取数据
	churnid, _ := c.GetInt("churnid")
	customerid, _ := c.GetInt("customerid")
	customername := c.GetString("customername")
	customermanager := c.GetString("customermanager")
	ordertime := c.GetString("ordertime")
	churnreasons := c.GetString("churnreasons")

	// 插入数据库
	o := orm.NewOrm()
	//// 检查订单编号是否已存在
	//existingChurn := models.Churncustomers{Churnid: churnid}
	//err := o.QueryTable("Churncustomers").Filter("Churnid", churnid).One(&existingChurn)
	//fmt.Println("churnid", churnid)
	//fmt.Println("addchurn:", err)
	//if err != nil {
	//	// 处理客户编号已存在的情况，例如返回错误响应给客户端
	//	c.Data["json"] = map[string]interface{}{
	//		"success": false,
	//		"message": "客户编号已存在",
	//	}
	//	c.ServeJSON()
	//	return
	//}

	churncustomers := models.Churncustomers{}
	// 对结构体赋值
	churncustomers.Churnid = churnid
	churncustomers.Customerid = customerid
	churncustomers.Customername = customername
	churncustomers.Customermanager = customermanager
	churncustomers.Ordertime = ordertime
	churncustomers.Churnreasons = churnreasons

	_, err1 := o.Insert(&churncustomers)
	if err1 != nil {
		fmt.Println("流失客户失败！")
		errCode := 416
		errorResponse := util.NewError(errCode)
		c.Data["json"] = errorResponse
		//c.Data["json"] = map[string]interface{}{
		//	"success": false,
		//	"message": "流失客户失败，数据不能为空！",
		//}
	}
	// 返回成功的 JSON 响应
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "流失客户成功！",
	}
	c.ServeJSON()
}

// ShowCustomerByPage  分页功能
/*
 * @summary 分页功能
 * @description 从前端获取分页信息，查询数据库获取分页结果，返回 JSON 格式的响应数据给前端
 */
func (c *CustomerController) ShowCustomerByPage() {
	var pageinfo2 PageParam2
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo2)
	fmt.Println("err =", err)
	if err != nil {
		c.Data["json"] = CustomerResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := int64(pageinfo2.Pagenum)
	pageSize := int64(pageinfo2.Pagesize)

	// 查询总数并获取分页结果
	totalCount, _ := models.GetCustomerCount()
	pagination2 := util.NewPagination(totalCount, currentPage, pageSize)
	offset := (pagination2.CurrentPage - 1) * pagination2.PageSize
	customers, _ := models.GetCustomers(int(offset), int(pagination2.PageSize))

	var repList []interface{}
	for _, customer := range customers {
		repList = append(repList, customer.CustomerToRespDesc())
	}

	// 将分页结果返回到前端
	pagination2.Data = repList
	c.Data["json"] = map[string]interface{}{
		"pagination":  pagination2,
		"total_pages": totalCount,
		"customers":   repList,
	}
	c.ServeJSON()

}

func (c *CustomerController) GetCustomerList() {

	var pageinfo PageParam2
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	if err != nil {
		c.Data["json"] = CustomerResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := pageinfo.Pagenum
	pageSize := pageinfo.Pagesize
	o := orm.NewOrm()
	totalCount, _ := o.QueryTable("customers").Count()
	fmt.Println("客户总数量：", totalCount)

	//totalCount, _ := models.GetUserCount()
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	if currentPage < 1 {
		currentPage = 1
	} else if currentPage > totalPages {
		currentPage = totalPages
	}
	var customers []models.Customers
	o.QueryTable("customers").Limit(pageSize, (currentPage-1)*pageSize).All(&customers)

	var repList []interface{}
	for _, customer := range customers {
		repList = append(repList, customer.CustomerToRespDesc())
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
func (c *CustomerController) ShowCustomerList() {

	var pageinfo PageParam2
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &pageinfo)
	fmt.Println("err =", err)
	if err != nil {
		c.Data["json"] = CustomerResponse{Success: false, Message: "解析请求失败"}
		c.ServeJSON()
		return
	}
	currentPage := int64(pageinfo.Pagenum)
	pageSize := int64(pageinfo.Pagesize)
	// 查询总记录数
	totalCount, _ := models.GetCustomerCount()
	// 创建分页实例
	pagination := util.NewPagination(totalCount, currentPage, pageSize)
	// 查询当前页的数据
	offset := (pagination.CurrentPage - 1) * pagination.PageSize
	customers, _ := models.GetCustomerlist(offset, pagination.PageSize)
	var repList []interface{}
	for _, customer := range customers {
		repList = append(repList, customer.CustomerToRespDesc())
	}
	pagination.Data = repList
	// 返回分页信息给前端
	c.Data["json"] = pagination
	c.ServeJSON()
}
