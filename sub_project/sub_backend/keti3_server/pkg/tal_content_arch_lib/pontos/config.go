package pontos

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos/config"
	"keti3/pkg/tal_content_arch_lib/pontos/env"

	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

// WithCfgOptions ...
type WithCfgOptions func(*config.Viper) *config.Viper

var (
	AppRootPath string
)

func initConfig() error {
	if AppRootPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		AppRootPath = wd
	}

	//Load Config File (app.toml)
	file := fmt.Sprintf("%s/config/%s/app.toml", AppRootPath, env.Environment)
	fmt.Println("-----------------------init config---------------------------------------------")
	fmt.Println("current env:" + env.Environment)
	fmt.Println("trying to use default config file:" + file)

	var (
		fileNotFound bool
		appName      = parseAppName()
		err          error
	)

	// 初始化应用配置
	if Config, err = config.NewConfigByFile(file); err != nil {
		if !errors.Is(err, config.FileNotFoundErr) {
			return err
		}

		fileNotFound = true
		Config, _ = config.NewConfig("app.toml", appName)
	}

	//尝试加载nacos配置，merge
	fmt.Println("trying to get tcm config for app: " + appName)
	tcmConfig, _ := Config.SetTcmClient().GetTCMConfig("app.toml", appName)
	if fileNotFound && tcmConfig == "" {
		return errors.New("init tcm config error")
	}

	fmt.Printf("tcmConfig:%s-------\n\n\n", tcmConfig)

	if tcmConfig != "" {
		fmt.Println("trying to merge tcm config to default config file")
		if err = Config.MergeConfig(strings.NewReader(tcmConfig)); err != nil {
			return err
		}

		//监听app.toml的tcm配置
		if err = Config.ListenConfig(); err != nil {
			return err
		}
	}

	fmt.Println("-----------------------init config done---------------------------------------------")
	return nil
}

// ReadConfig 读取配置文件(存在其他配置文件时可使用此方法读取配置,file传入全路径)
func ReadConfig(file string) (*config.Viper, error) {
	return config.NewConfigByFile(file)
}

// ReadFileConfig 读取单独配置文件，不监听TCM
func ReadFileConfig(file string) (*config.Viper, error) {
	return config.NewConfigByFile(file)
}

// ReadTCMConfig 读取TCM配置中心配置
func ReadTCMConfig(dataId string, opts ...WithCfgOptions) (*config.Viper, error) {
	//mock local env
	if env.Environment == env.LocalEnv {
		file := fmt.Sprintf("%s/config/%s/%s", AppRootPath, env.Environment, dataId)
		fmt.Println("mock get tcm config at local env,using file:" + file)
		return ReadFileConfig(file)
	}

	cfg, err := config.NewConfig(dataId, parseAppName())
	if err != nil {
		return nil, err
	}

	data, _ := cfg.SetTCMLogger(Logger).GetTCMConfig(dataId, parseAppName())
	if data != "" {
		if err = cfg.ReadConfig(strings.NewReader(data)); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}

	return cfg, nil
}

// parseAppName ...
func parseAppName() string {
	appPath := strings.Split(AppRootPath, "/")
	return appPath[len(appPath)-1]
}

func SetTCMLogger(l logger.Logger) WithCfgOptions {
	return func(viper *config.Viper) *config.Viper {
		viper = viper.SetTCMLogger(l)
		return viper
	}
}
