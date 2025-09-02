package errcode

// UnLoginError 未登录错误码，前端依赖
const UnLoginError = 29999

// UnknownCode 未知错误兜底
const UnknownCode = 60000
const DataEncryptionError = 60001

// CodeMsg defines the map of code to msg
var CodeMsg = map[int]string{
	UnknownCode:         "系统错误，请稍后再试",
	UnLoginError:        "您尚未登录~",
	DataEncryptionError: "数据加密错误",
}

// ErrMSG code -> msg
func ErrMSG(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}

	return CodeMsg[UnknownCode]
}
