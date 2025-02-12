package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"math"
)

// Users 用户表
type Users struct {
	Userid   int `orm:"pk"`
	Username string
	Password string
}

// Customers 客户表
type Customers struct {
	Customerid      int `orm:"pk"`
	Customername    string
	Contactid       int
	Customeraddress string
	Customerinfo    string
}

// Contacts 联系人表
type Contacts struct {
	Contactid   int `orm:"pk"`
	Contactname string
	Contactinfo string
}

// Salesorders 销售订单表
type Salesorders struct {
	Orderid      int `orm:"pk"`
	Ordername    string
	Orderstatus  string
	Customername string
	Orderaddress string
	Orderprice   float64
	Ordertime    string
}

// Churncustomers 流失客户表
type Churncustomers struct {
	Churnid         int `orm:"pk"`
	Customerid      int
	Customername    string
	Customermanager string
	Ordertime       string
	Churnreasons    string
}

// Services 服务表
type Services struct {
	Serviceid      int    `orm:"pk" json:"serviceid"`
	Servicename    string `json:"servicename"`
	Customername   string `json:"customername"`
	Servicestatus  string `json:"servicestatus"`
	Servicecreator string `json:"servicecreator"`
}

type UnifiedResponse struct {
	Code    int         `json:"code"`    // 响应状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

func init() {
	//设置数据库基本信息
	orm.RegisterDataBase("default", "mysql", ":@tcp()/test?charset=utf8")
	// 映射model数据
	orm.RegisterModel(new(Users))
	orm.RegisterModel(new(Customers))
	orm.RegisterModel(new(Salesorders))
	orm.RegisterModel(new(Churncustomers))
	orm.RegisterModel(new(Services))
	orm.RegisterModel(new(Contacts))
	// 生成表
	//orm.RunSyncdb("default", false, true)
	//打印orm日志
	orm.Debug = true
}

func (user *Users) UserToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"userid":   user.Userid,
		"username": user.Username,
	}
	return respInfo
}

func (customer *Customers) CustomerToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"customerid":      customer.Customerid,
		"customername":    customer.Customername,
		"contactid":       customer.Contactid,
		"customeraddress": customer.Customeraddress,
		"customerinfo":    customer.Customerinfo,
	}
	return respInfo
}

func (service *Services) ServiceToRespDesc() interface{} {
	respInfo := map[string]interface{}{
		"serviceid":      service.Serviceid,
		"servicename":    service.Servicename,
		"customername":   service.Customername,
		"servicestatus":  service.Servicestatus,
		"Servicecreator": service.Servicecreator,
	}
	return respInfo
}

func GetUsers(page, pageSize int) ([]Users, error) {
	o := orm.NewOrm()
	count, _ := o.QueryTable("users").Count()
	fmt.Println("用户总数量：", count)
	// 计算总页数
	totalPage1 := int(math.Ceil(float64(count) / float64(pageSize)))
	// 确保当前页码在有效范围内
	if page < 1 {
		page = 1
	} else if page > totalPage1 {
		page = totalPage1
	}
	var users []Users
	_, err := o.QueryTable("users").Limit(pageSize, (page-1)*pageSize).All(&users)
	return users, err
}

func GetCustomers(page, pageSize int) ([]Customers, error) {
	o := orm.NewOrm()
	count, _ := o.QueryTable("customers").Count()
	fmt.Println("客户总数量：", count)
	// 计算总页数
	totalPage2 := int(math.Ceil(float64(count) / float64(pageSize)))
	// 确保当前页码在有效范围内
	if page < 1 {
		page = 1
	} else if page > totalPage2 {
		page = totalPage2
	}
	var customers []Customers
	_, err := o.QueryTable("customers").Limit(pageSize, (page-1)*pageSize).All(&customers)
	return customers, err
}

func GetService(page, pageSize int) ([]Services, error) {
	o := orm.NewOrm()
	count, _ := o.QueryTable("services").Count()
	fmt.Println("服务总数量：", count)
	// 计算总页数
	totalPage3 := int(math.Ceil(float64(count) / float64(pageSize)))
	// 确保当前页码在有效范围内
	if page < 1 {
		page = 1
	} else if page > totalPage3 {
		page = totalPage3
	}
	var services []Services
	_, err := o.QueryTable("services").Limit(pageSize, (page-1)*pageSize).All(&services)
	return services, err
}

/*
获取用户总记录数
*/
func GetUserCount() (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("users").Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

/*
获取客户信息总记录数
*/
func GetCustomerCount() (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("customers").Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

/*
获取服务信息总记录数
*/
func GetServiceCount() (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("services").Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

/*
获取用户当前页的数据
*/
func GetUserlist(offset, limit int64) ([]Users, error) {
	o := orm.NewOrm()
	var users []Users
	_, err := o.QueryTable("users").Offset(offset).Limit(limit).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

/*
获取客户当前页的数据
*/
func GetCustomerlist(offset, limit int64) ([]Customers, error) {
	o := orm.NewOrm()
	var customers []Customers
	_, err := o.QueryTable("customers").Offset(offset).Limit(limit).All(&customers)
	if err != nil {
		return nil, err

	}
	return customers, nil
}

/*
获取服务当前页的数据
*/
func GetServicelist(offset, limit int64) ([]Services, error) {
	o := orm.NewOrm()
	var services []Services
	_, err := o.QueryTable("services").Offset(offset).Limit(limit).All(&services)
	if err != nil {
		return nil, err

	}
	return services, nil
}

// GetUsersBySearchTerm 根据搜索词获取用户列表
func GetUsersBySearchTerm(searchTerm string) ([]Users, error) {
	o := orm.NewOrm()

	// 创建一个用于存储查询结果的切片
	var users []Users
	// 使用 ORM 查询
	//username 是数据库表中的列名，表示你希望在哪一列进行查询。
	//__icontains 是 beego ORM 中的一个过滤器，表示执行不区分大小写的模糊查询。
	_, err := o.QueryTable("users").Filter("username__icontains", searchTerm).All(&users)
	if err != nil {
		// 处理错误，例如日志记录或返回错误信息
		return nil, err
	}

	return users, nil

}

// GetCustomersBySearchTerm 根据搜索词获取客户列表
func GetCustomersBySearchTerm(searchTerm string) ([]Customers, error) {
	o := orm.NewOrm()

	// 创建一个用于存储查询结果的切片
	var customers []Customers

	// 使用 ORM 查询
	//customername 是数据库表中的列名，表示在哪一列进行查询。
	//__icontains 是 beego ORM 中的一个过滤器，表示执行不区分大小写的模糊查询。
	_, err := o.QueryTable("customers").Filter("customername__icontains", searchTerm).All(&customers)
	if err != nil {
		// 处理错误，例如日志记录或返回错误信息
		return nil, err
	}

	return customers, nil
}

// GetServicesBySearchTerm 根据搜索词获取服务列表
func GetServicesBySearchTerm(searchTerm string) ([]Services, error) {
	o := orm.NewOrm()

	// 创建一个用于存储查询结果的切片
	var services []Services
	// 使用 ORM 查询
	//customername 是数据库表中的列名，表示你希望在哪一列进行查询。
	//__icontains 是 beego ORM 中的一个过滤器，表示执行不区分大小写的模糊查询。

	_, err := o.QueryTable("services").Filter("customername__icontains", searchTerm).All(&services)
	if err != nil {
		// 处理错误，例如日志记录或返回错误信息
		return nil, err
	}
	fmt.Println("services", services)

	return services, nil
}

// GetChurnCustomersBySearchTerm 根据搜索词获取服务列表
func GetChurnCustomersBySearchTerm(searchTerm string) ([]Churncustomers, error) {
	o := orm.NewOrm()

	// 创建一个用于存储查询结果的切片
	var churncustomers []Churncustomers
	// 使用 ORM 查询
	//customername 是数据库表中的列名，表示你希望在哪一列进行查询。
	//__icontains 是 beego ORM 中的一个过滤器，表示执行不区分大小写的模糊查询。

	_, err := o.QueryTable("churncustomers").Filter("customername__icontains", searchTerm).All(&churncustomers)
	if err != nil {
		// 处理错误，例如日志记录或返回错误信息
		return nil, err
	}
	fmt.Println("services", churncustomers)

	return churncustomers, nil
}

// GetOrdersBySearchTerm  根据搜索词获取订单列表
func GetOrdersBySearchTerm(searchTerm string) ([]Salesorders, error) {
	o := orm.NewOrm()

	// 创建一个用于存储查询结果的切片
	var salesorders []Salesorders
	// 使用 ORM 查询
	//ordername 是数据库表中的列名，表示你希望在哪一列进行查询。
	//__icontains 是 beego ORM 中的一个过滤器，表示执行不区分大小写的模糊查询。

	_, err := o.QueryTable("salesorders").Filter("ordername__icontains", searchTerm).All(&salesorders)
	if err != nil {
		// 处理错误，例如日志记录或返回错误信息
		return nil, err
	}
	fmt.Println("salesorders", salesorders)

	return salesorders, nil
}
