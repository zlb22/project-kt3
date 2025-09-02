package keti3

import (
	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) ExportLog(ctx *gin.Context) {
	var (
		err error
	)
	ret, err := c.biz.ExportLog(ctx)
	if err != nil {
		pontos.Error(ctx, err, nil)
		return
	}
	ctx.String(200, ret)
	// pontos.Success(ctx, ret)

}
