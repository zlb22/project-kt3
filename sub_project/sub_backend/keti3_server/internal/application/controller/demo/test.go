package demo

import (
	"errors"
	"time"

	"keti3/internal/application/errcode"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

// TestLog 日志使用
func (c *Ctrl) TestLog(ctx *gin.Context) {
	pontos.Logger.WithContext(ctx).Info("测试日志")
	pontos.Logger.WithContext(ctx).Errorf("error:%s", errors.New("测试错误"))

	pontos.Success(ctx, nil)
}

// TestRedis redis使用
func (c *Ctrl) TestRedis(ctx *gin.Context) {
	err := pontos.Redis.Set(ctx, "test_key", "test_value", time.Hour).Err()
	if err != nil {
		errCode := errcode.ParamInvalid
		pontos.Error(ctx, pontos.NewAppError(err, errCode, errcode.ErrMSG(errCode)), nil)
		return
	}

	data := pontos.Redis.Get(ctx, "test_key").Val()
	pontos.Success(ctx, data)
}
