package util

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrCodeMsgMap = map[int]string{
	400: "错误请求",
	401: "未经 授权",
	403: "禁止",
	404: "未找到页面",
	405: "解析查询数据失败",
	407: "验证码输入错误",
	408: "查询数据失败",
	409: "数据不能为空",
	410: "解析搜索请求失败",
	411: "搜索失败",
	412: "数据添加失败",
	413: "数据删除失败",
	414: "数据修改失败",
	415: "数据不存在",
	416: "数据插入失败",
	500: "Internal Server Error",
	406: "密码加密失败",
	501: "用户不存在",
	502: "密码输入错误",

	//继续添加其他的错误

}

func NewError(errCode int) ErrorResponse {
	msg, ok := ErrCodeMsgMap[errCode]
	if !ok {
		msg = "Unknown Error"
	}
	return ErrorResponse{
		Code:    errCode,
		Message: msg,
	}
}
