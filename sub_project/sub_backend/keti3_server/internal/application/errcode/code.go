package errcode

// 参数类错误码 1xx
const (
	ParamMissing = iota + 10000
	ParamInvalid
	SignNotMatch
)

// Business 错误码 2xx
const (
	DefaultBizError = iota + 20000
	JsonMarshalErr
	JsonUnmarshalErr
	OperationFrequently
	GetUserBalanceErr
	UserNotLogin = 20004
)

// Redis类错误码 30000
const (
	RedisQueryError = iota + 30000
	RedisExecuteError
)

// MySQL类错误码 40000
const (
	MySQLCreateError = iota + 40000
	MySQLUpdateError
	MySQLQueryError
)

// UnknownCode 未知错误兜底
const UnknownCode = 50000

// CodeMsg defines the map of code to msg
var CodeMsg = map[int]string{
	ParamMissing:        "参数缺失~",
	ParamInvalid:        "参数有误~",
	OperationFrequently: "操作频繁,请稍候再试",
	UnknownCode:         "系统错误，请稍后再试",
	SignNotMatch:        "签名不匹配",
	UserNotLogin:        "用户未登录",
}

// ErrMSG code -> msg
func ErrMSG(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}

	return CodeMsg[UnknownCode]
}
