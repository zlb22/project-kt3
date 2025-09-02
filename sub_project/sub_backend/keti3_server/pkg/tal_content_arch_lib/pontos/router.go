package pontos

import (
	"net/http"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"

	"github.com/gin-gonic/gin"
)

// registerRouter register router
func registerRouter() {
	Server.GinEngine().GET(constant.PingPath, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Thank you for using Pontos")
	})
}

var filterUrls = []string{
	constant.PingPath,
}

func isHitFilterUri(uri string) bool {
	for _, url := range filterUrls {
		if url == uri {
			return true
		}
	}

	return false
}
