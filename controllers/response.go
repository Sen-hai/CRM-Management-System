package controllers

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONResponse(code int, message string, data ...interface{}) APIResponse {
	var respData interface{}
	if len(data) > 0 {
		respData = data[0]
	}
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    respData,
	}
}
