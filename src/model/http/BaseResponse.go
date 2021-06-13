package http

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GenSuccessResponse() BaseResponse {
	return BaseResponse{
		Code:    0,
		Message: "处理成功",
	}
}

func GenErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Code:    -1,
		Message: message,
	}
}

func GenResponse(code int, message string) BaseResponse {
	return BaseResponse{
		Code:    code,
		Message: message,
	}
}
