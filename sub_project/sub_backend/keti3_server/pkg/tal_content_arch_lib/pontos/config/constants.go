package config

import (
	"errors"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

var (
	FileNotFoundErr = errors.New("config file not found")
	NewTCMClientErr = errors.New("new tcm client error")
)

var nacosServers = map[string][]constant.ServerConfig{}

var nacosEnvConfigs = make(map[string]*constant.ClientConfig, 4)
