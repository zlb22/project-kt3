package keti3

import (
	"keti3/internal/application/business/keti3"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) StudentLogin(ctx *gin.Context) {
	var (
		err  error
		data interface{}
		req  keti3.StudentLoginParam
	)
	if err := ctx.ShouldBind(&req); err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	data, err = c.biz.StudentLogin(ctx, req)
	if err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	pontos.Success(ctx, data)

}
