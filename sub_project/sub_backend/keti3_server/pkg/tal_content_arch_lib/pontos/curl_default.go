package pontos

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	routerError = errors.New("curl router error")
)

type CurlMap struct {
	clientMap map[string]*CurlClient
	router    map[string]map[string]string
}

var DefaultCurlClient *CurlMap

type defConfig struct {
	Url     string            `mapstructure:"url"`
	Path    map[string]string `mapstructure:"path"`
	Timeout time.Duration     `mapstructure:"timeout"`
}

type Result struct {
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

func initDefaultCurlClient() error {
	if !Config.IsSet(curlSection) {
		return nil
	}
	conf := map[string]*defConfig{}
	if err := Config.UnmarshalKey(curlSection, &conf); err != nil {
		return err
	}

	clientMap := make(map[string]*CurlClient, len(conf))
	router := make(map[string]map[string]string, len(conf))
	for key, config := range conf {
		clientMap[key] = NewCurlClient(config.Timeout).SetLogger(CurlLogger)
		for k, v := range config.Path {
			config.Path[k] = fmt.Sprintf("%s%s", config.Url, v)
		}
		router[key] = config.Path
	}
	DefaultCurlClient = &CurlMap{
		clientMap: clientMap,
		router:    router,
	}
	return nil
}

func (c *CurlMap) Call(ctx *gin.Context, group, method string, headers map[string]string, params interface{}, out interface{}) (*Result, error) {
	url, ok := c.router[group][method]
	if !ok {
		return nil, routerError
	}
	client, ok := c.clientMap[group]
	if !ok {
		return nil, routerError
	}
	zr := client.JsonPostCtx(ctx, url, params)
	if len(headers) > 0 {
		zr.SetHeaders(headers)
	}
	z := zr.Do()
	if out == nil {
		result := &Result{}
		if err := z.JsonDecode(result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	} else {
		result := &Result{Data: out}
		if err := z.JsonDecode(result); err != nil {
			return nil, err
		} else {
			if result.Errcode == 0 {
				return nil, nil
			} else {
				return nil, errors.New(result.Errmsg)
			}
		}
	}
}
