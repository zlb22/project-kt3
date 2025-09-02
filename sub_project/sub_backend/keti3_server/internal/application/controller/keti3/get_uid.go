package keti3

import (
	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetUIDByIncrease(ctx *gin.Context) {
	var (
		err  error
		data interface{}
	)
	data, err = c.biz.GetUIDByIncrease(ctx)
	if err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	pontos.Success(ctx, data)
}
