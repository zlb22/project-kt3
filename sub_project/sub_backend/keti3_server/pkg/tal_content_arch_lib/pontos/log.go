package pontos

import (
	"os"
	"path"
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/env"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"

	plog "keti3/pkg/tal_content_arch_lib/pontos-log"
)

type logConfig struct {
	OutPath   string `mapstructure:"path"`
	BizStdout bool   `mapstructure:"biz_stdout"`
	SysStdout bool   `mapstructure:"sys_stdout"`
}

var (
	systemLogger *logger.ZLogger
)

func InitLogger(app string) error {
	var logConf logConfig
	if err := Config.UnmarshalKey(logSection, &logConf); err != nil {
		return err
	}

	if logConf.OutPath == "" {
		logConf.OutPath = path.Join(AppRootPath, "logs")
	}

	//初始化业务日志
	newZLogger, err := logger.NewPLogger(path.Join(logConf.OutPath, logFile(app, "biz.log")))
	if err != nil {
		return err
	}

	if logConf.BizStdout && env.Environment != env.OnlineEnv {
		newZLogger.SetOutput(os.Stdout)
		newZLogger.SetFormatter(plog.NewTextFormatter())
	} else {
		newZLogger.AddHook(logger.NewDefaultFieldHook(app, env.Environment))
		newZLogger.AddHook(logger.NewBizTagHook())
	}
	newZLogger.AddHook(logger.NewOTELHook())
	Logger = newZLogger

	//初始化系统日志(request/response/mysql...)
	systemLogger, err = logger.NewPLogger(path.Join(logConf.OutPath, logFile(app, "sys.log")))
	if err != nil {
		return err
	}

	if logConf.SysStdout && env.Environment != env.OnlineEnv {
		systemLogger.SetOutput(os.Stdout)
		systemLogger.SetFormatter(plog.NewTextFormatter())
	} else {
		systemLogger.AddHook(logger.NewDefaultFieldHook(app, env.Environment, constant.DeviceSNCtxKey, constant.TalIDCtxKey))
	}
	systemLogger.AddHook(logger.NewOTELHook())

	//traceLogger, err = logger.NewPLogger(path.Join(logConf.OutPath, logFile(app, "tracing.log")))
	//if err != nil {
	//	return err
	//}
	return nil
}

func logFile(app, suffix string) string {
	var c strings.Builder

	hn, _ := os.Hostname()
	if hn == "" {
		hn = os.Getenv(constant.K8SPodIPEnvKey)
	}

	if hn == "" {
		c.WriteString(app)
		c.WriteString("_")
		c.WriteString(suffix)
		return c.String()
	}

	c.WriteString(hn)
	c.WriteString("_")
	c.WriteString(suffix)
	return c.String()
}
