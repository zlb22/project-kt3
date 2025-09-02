package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"keti3/internal/application/errcode"
	"net/http"

	"keti3/internal/application/service/base_user"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func CheckSignature() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取验证参数
		appId := ctx.GetHeader("X-Auth-AppId")
		timestamp := ctx.GetHeader("X-Auth-TimeStamp")
		sign := ctx.GetHeader("X-Sign")
		appSecret := pontos.Config.GetString(fmt.Sprintf("secret.%s", appId))

		if sign != getSign(appSecret, timestamp) {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"errcode": errcode.SignNotMatch,
				"errmsg":  errcode.ErrMSG(errcode.SignNotMatch),
				"trace":   ctx.GetString("x_trace_id"),
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// getSign 生成签名
func getSign(secret, timestamp string) string {
	str := secret + timestamp
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// CheckUserToken 用户token校验中间件
func CheckUserToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		appId := ctx.GetHeader("X-Auth-AppId")
		timestamp := ctx.GetHeader("X-Auth-TimeStamp")
		sign := ctx.GetHeader("X-Sign")

		pontos.Logger.Infof("CheckUserToken,token:%s,appId:%s,timestamp:%s,sign:%s", token, appId, timestamp, sign)
		if token == "" {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"errcode": errcode.UserNotLogin,
				"errmsg":  errcode.ErrMSG(errcode.UserNotLogin),
				"trace":   ctx.GetString("x_trace_id"),
			})
			ctx.Abort()
			return
		}

		headers := map[string]string{
			"X-Auth-AppId":     appId,
			"X-Auth-TimeStamp": timestamp,
			"X-Sign":           sign,
			"Authorization":    token,
		}
		client := base_user.NewClient()
		if client == nil {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"errcode": errcode.UnknownCode,
				"errmsg":  errcode.ErrMSG(errcode.UnknownCode),
				"trace":   ctx.GetString("x_trace_id"),
			})
			ctx.Abort()
			return
		}
		userInfo, errCode, err := client.CheckGetUserInfo(ctx, headers)
		if err != nil {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"errcode": errCode,
				"errmsg":  err.Error(),
				"trace":   ctx.GetString("x_trace_id"),
			})
			ctx.Abort()
			return
		}
		ctx.Set("user_info", userInfo)
		ctx.Next()
	}
}
