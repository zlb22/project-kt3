package middleware

import (
	"keti3/pkg/tal_content_arch_lib/pontos/constant"

	"github.com/gin-gonic/gin"
)

// MarkPressure 标记压测流量
func MarkPressure() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var i int
		if ctx.GetHeader(constant.PressureTestKey) == constant.PressureTestValue {
			i = 1
		}

		ctx.Set(constant.IsPressureTestKey, i)
	}
}
