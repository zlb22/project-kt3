package keti3

import (
	"keti3/internal/application/business/keti3"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) LogList(ctx *gin.Context) {
	var (
		err  error
		data interface{}
		req  keti3.LogListParam
	)

	authToken := ctx.GetHeader("authtoken")
	if err = c.biz.CheckTeacherLogin(ctx, authToken); err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	if err := ctx.ShouldBind(&req); err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	data, err = c.biz.LogList(ctx, req)
	if err != nil {
		pontos.Error(ctx, err, nil)
		return
	}

	pontos.Success(ctx, data)

}
