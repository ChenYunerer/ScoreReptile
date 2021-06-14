package http

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenSuccessResponse() BaseResponse {
	return BaseResponse{
		Code:    0,
		Message: "处理成功",
	}
}

func GenSuccessResponseWithData(data interface{}) BaseResponse {
	return BaseResponse{
		Code:    0,
		Message: "处理成功",
		Data:    data,
	}
}

func GenErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    -1,
		Message: message,
	}
}

func GenResponse(code int, message string, data interface{}) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
