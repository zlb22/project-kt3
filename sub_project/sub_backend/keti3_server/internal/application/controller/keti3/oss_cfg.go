package keti3

import (
	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) OssAuth(ctx *gin.Context) {
	var (
		err error
	)
	ret, err := c.biz.GetSTSAuth(ctx)
	if err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	pontos.Success(ctx, ret)

}
