package pontos

import (
	"errors"
	"fmt"
	"net/http"

	"keti3/pkg/tal_content_arch_lib/pontos/bootstrap"
	"keti3/pkg/tal_content_arch_lib/pontos/config"
	"keti3/pkg/tal_content_arch_lib/pontos/kit"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"
	"keti3/pkg/tal_content_arch_lib/pontos/middleware"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	_ "go.uber.org/automaxprocs"
)

var (
	Logger     *logger.ZLogger
	Server     *bootstrap.Server
	Config     *config.Viper
	Redis      = new(RedisClient)
	Mysql      = new(MysqlClient)
	MysqlMap   = make(map[string]*MysqlClient, 2)
	RedisMap   = make(map[string]*RedisClient, 2)
	ToolKit    = kit.NewKit()
	CurlLogger = NewCurlCtxLogger()
)

var (
	svrSection       = "server"
	logSection       = "log"
	sqlSection       = "mysql"
	sqlMapSection    = "mysql_map"
	redisSection     = "redis"
	redisMapSection  = "redis_map"
	curlSection      = "curl"
	tcmNoticeSection = "tcm_notice"
)

// ChainOptions 组件加载配置
type ChainOptions struct {
	// cfgWithTCM 加载智能学习配置中心namespace
	cfgWithTCM bool

	// ConfigFile //只初始化配置文件
	ConfigFile string

	// IsInitDefault 是否初始化默认组件：logger,mysql,mysql_map,redis,redis_map
	IsInitDefault bool

	// AppVersion 携带app version(git tag)
	AppVersion string
}

type InitAppOption interface {
	apply(opt *ChainOptions)
}

type appOptionFunc func(opt *ChainOptions)

func (o appOptionFunc) apply(opt *ChainOptions) {
	o(opt)
}

// WithAppVersion ...
func WithAppVersion(version string) InitAppOption {
	return appOptionFunc(func(opt *ChainOptions) {
		opt.AppVersion = version
	})
}

// initChain 初始化组件
func initChain(chain *ChainOptions) (*bootstrap.ServerOptions, error) {
	// 初始化应用配置
	if chain.cfgWithTCM {
		if err := initConfig(); err != nil {
			return nil, errors.New(fmt.Sprintf("init config error :%v", err))
		}
	} else {
		if chain.ConfigFile == "" {
			return nil, errors.New("config file of ChainOptions is required.")
		}

		var err error
		Config, err = ReadConfig(chain.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	//获取服务配置
	option := &bootstrap.ServerOptions{}
	if err := Config.UnmarshalKey(svrSection, option); err != nil {
		return nil, err
	}

	//是否初始化各个组件
	if !chain.IsInitDefault {
		return option, nil
	}

	//初始化日志
	if err := InitLogger(option.AppName); err != nil {
		return nil, errors.New(fmt.Sprintf("init logger error :%v", err))
	}

	//初始化Mysql,如有
	if err := InitMysql(systemLogger); err != nil {
		return nil, errors.New(fmt.Sprintf("init mysql error :%v", err))
	}

	//初始化MysqlMap 多实例,如有
	if err := InitMysqlMap(systemLogger); err != nil {
		return nil, errors.New(fmt.Sprintf("init mysql map error :%v", err))
	}

	//初始化redis,如果有
	if err := InitRedis(Logger); err != nil {
		return nil, errors.New(fmt.Sprintf("init redis error :%v", err))
	}

	//初始化redisMap 多实例，如果有
	if err := InitRedisMap(Logger); err != nil {
		return nil, errors.New(fmt.Sprintf("init redis map error :%v", err))
	}

	//初始化默认curl client
	if err := initDefaultCurlClient(); err != nil {
		return nil, errors.New(fmt.Sprintf("init curl error :%v", err))
	}
	return option, nil
}

// NewApp 自定义组件初始化服务
func NewApp(options *ChainOptions) error {
	setEnv()

	option, err := initChain(options)
	if err != nil {
		return err
	}

	//init server
	Server = bootstrap.NewServer()
	Server.Init(option)

	registerRouter()
	return nil
}

// InitApp 默认初始化各组件和监听中间件(智能学习服务端专用)
func InitApp(opts ...InitAppOption) error {
	setEnv()

	chain := &ChainOptions{
		cfgWithTCM:    true,
		IsInitDefault: true,
	}

	for _, opt := range opts {
		opt.apply(chain)
	}

	option, err := initChain(chain)
	if err != nil {
		return err
	}

	//tracing.log
	//tracer, err := initLogTracer(option.AppName)
	//if err != nil {
	//	return err
	//}

	//init server
	Server = bootstrap.NewServer()

	Server.Init(option)
	Server.UseMiddleware(
		otelgin.Middleware(option.AppName, otelgin.WithFilter(func(request *http.Request) bool {
			if isHitFilterUri(request.RequestURI) {
				return false
			}
			return true
		})),

		middleware.Recovery(systemLogger),
		middleware.Context(systemLogger),
		middleware.MarkPressure(),
		middleware.Logger(systemLogger),
	)
	Server.AddAfterServerFunc(
		CloseMysql(),
		CloseMysqlMap(),
		CloseRedis(),
		CloseRedisMap(),
	)

	registerRouter()
	Server.RegisterShutdown(func() {
	})
	return nil
}

// InitCommand 初始化命令行
func InitCommand(opts ...InitAppOption) error {
	setEnv()

	chain := &ChainOptions{
		cfgWithTCM:    true,
		IsInitDefault: true,
	}

	for _, opt := range opts {
		opt.apply(chain)
	}

	_, err := initChain(chain)
	if err != nil {
		return err
	}

	return nil
}

// CloseCommand 命令行退出各个组件
func CloseCommand() {
	CloseMysql()
	CloseMysqlMap()
	CloseRedis()
	CloseRedisMap()

}
