package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OptionsRequest 处理跨域的options请求
func OptionsRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.ToUpper(c.Request.Method) == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// Cors 跨域+response header
func Cors(headers ...string) gin.HandlerFunc {
	allowHeaders := "Content-Type, Authorization, X-Requested-With, X-TAL-AppId, X-TAL-DeviceId, X-TAL-Timestamp, X-TAL-Nonce, X-TAL-Signature,Authtoken,X-TAL-ClientId,X-TAL-Version"

	return func(ctx *gin.Context) {
		if len(headers) > 0 {
			allowHeaders = fmt.Sprintf("%s,%s", allowHeaders, strings.Join(headers, ","))
		}

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT,PATCH,DELETE")
		ctx.Header("Access-Control-Allow-Headers", allowHeaders)

		if strings.ToUpper(ctx.Request.Method) == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
