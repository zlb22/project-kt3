package pontos

import (
	"flag"
	"log"
	"os"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/env"

	"github.com/gin-gonic/gin"
)

var (
	Environment             string
	IsDebugEnv, IsOnlineEnv bool

	LocalEnv       = env.LocalEnv
	DevelopmentEnv = env.DevelopmentEnv
	TestEnv        = env.TestEnv
	GrayEnv        = env.GrayEnv
	OnlineEnv      = env.OnlineEnv
)

func setEnv() {
	//兼容-env=xxx
	flag.StringVar(&env.Environment, "env", "local", "set server environment")
	flag.Parse()

	//未来云环境变量，test/online/gray
	cloudEnv := os.Getenv(constant.K8SEnvKey)
	if cloudEnv != "" {
		//开发环境需要特殊处理，未来云的环境变量只能为test
		if !(cloudEnv == env.TestEnv && env.Environment == env.DevelopmentEnv) {
			env.Environment = cloudEnv
		}
	}

	switch env.Environment {
	case env.LocalEnv:
		gin.SetMode(gin.DebugMode)

	case env.DevelopmentEnv, env.TestEnv:
		env.IsDebugEnv = true
		gin.SetMode(gin.DebugMode)

	case env.GrayEnv, env.OnlineEnv:
		env.IsOnlineEnv = true
		gin.SetMode(gin.ReleaseMode)

	default:
		log.Fatalf("the environment must be the following enumerated values:local/development/test/gray/online")
	}

	// 兼容之前直接用 pontos.Environment
	Environment = env.Environment
	IsDebugEnv = env.IsDebugEnv
	IsOnlineEnv = env.IsOnlineEnv
}
