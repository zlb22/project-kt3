package config

import (
	"strings"
	"sync"

	"keti3/pkg/tal_content_arch_lib/pontos/env"

	"github.com/nacos-group/nacos-sdk-go/clients"
	cfgClient "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	client cfgClient.IConfigClient
	once   sync.Once
)

func newTcmClient() cfgClient.IConfigClient {
	once.Do(func() {
		cli, _ := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  nacosEnvConfigs[env.Environment],
				ServerConfigs: nacosServers[env.Environment],
			},
		)

		client = cli
	})

	return client
}

// listenTCM 监听tcm配置同步到viper
func (v *Viper) listenTCM() error {
	if v.TcmCli == nil {
		return NewTCMClientErr
	}

	if err := v.TcmCli.ListenConfig(vo.ConfigParam{
		DataId: v.file,
		Group:  v.appName,
		OnChange: func(namespace, group, dataId, data string) {
			if v.TcmLogger != nil {
				v.TcmLogger.Infof("config changed,service:%v,file:%v", group, dataId)
			}

			if err := v.ReadConfig(strings.NewReader(data)); err != nil {
				if v.TcmLogger != nil {
					v.TcmLogger.Errorf("tcm config changed,but parse to viper error:%v", err)
				}
				return
			}

			if v.TcmNoticeFunc != nil {
				v.TcmNoticeFunc(map[string]string{
					"namespace": namespace,
					"service":   group,
					"file":      dataId,
				})
			}
		},
	}); err != nil {
		return err
	}

	return nil
}

// SetTcmClient ...
func (v *Viper) SetTcmClient() *Viper {
	v.TcmCli = newTcmClient()
	return v
}

// GetTCMConfig 读取TCM配置
func (v *Viper) GetTCMConfig(dataId, appName string) (string, error) {
	if v.TcmCli == nil {
		return "", NewTCMClientErr
	}
	v.file = dataId
	v.appName = appName

	return v.TcmCli.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  appName,
	})
}
