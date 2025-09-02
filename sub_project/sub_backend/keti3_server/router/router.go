package router

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter register all router
func RegisterRouter(router *gin.Engine) {
	registerKeti3Router(router.Group("web/keti3"))
}
