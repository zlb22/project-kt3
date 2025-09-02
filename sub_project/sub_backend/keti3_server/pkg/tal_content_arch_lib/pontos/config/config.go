package config

import (
	"io"

	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	cfgClient "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/spf13/viper"
)

type ChangeNoticeFunc func(map[string]string)

type Viper struct {
	*viper.Viper

	TcmCli        cfgClient.IConfigClient
	TcmLogger     logger.Logger
	TcmNoticeFunc ChangeNoticeFunc
	file          string //config file or tcm dataId
	appName       string //tcm service name
}

// NewConfig new an instance of viper
func NewConfig(fName, appName string) (*Viper, error) {
	v := viper.New()

	v.SetConfigFile(fName)
	return &Viper{
		Viper:   v,
		TcmCli:  newTcmClient(),
		file:    fName,
		appName: appName,
	}, nil
}

// NewConfigByFile read file content to viper
func NewConfigByFile(file string) (*Viper, error) {
	v := viper.New()

	if !kit.Impl.FileExist(file) {
		return nil, FileNotFoundErr
	}

	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Viper{
		Viper:   v,
		TcmCli:  nil,
		file:    file,
		appName: "",
	}, nil
}

func (v *Viper) MergeConfig(in io.Reader) error {
	if v.Viper == nil {
		v.Viper = viper.New()
		v.SetConfigType("toml")
		return v.ReadConfig(in)
	}
	return v.Viper.MergeConfig(in)
}

// ListenConfig 兼容viper监听和TCM配置监听
func (v *Viper) ListenConfig() error {
	if v.TcmCli == nil {
		v.WatchConfig()
		return nil
	}

	return v.listenTCM()
}

// SetTCMLogger set config logger
func (v *Viper) SetTCMLogger(l logger.Logger) *Viper {
	if l != nil {
		v.TcmLogger = l
		_ = logger.InitLogger(logger.Config{CustomLogger: v.TcmLogger})
	}
	return v
}

// SetChangeNoticeFunc ...
func (v *Viper) SetChangeNoticeFunc(f ChangeNoticeFunc) *Viper {
	v.TcmNoticeFunc = f
	return v
}
