package base_user

import (
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type config struct {
	Host    string        `mapstructure:"host"`
	Timeout time.Duration `mapstructure:"timeout"`
}

var clientConfig *config

type Client struct {
	config
	*pontos.CurlClient
}

func NewClient() *Client {
	if err := pontos.Config.UnmarshalKey("service.user-host", &clientConfig); err != nil {
		return nil
	}
	return &Client{
		CurlClient: pontos.NewCurlClient(clientConfig.Timeout).SetLogger(pontos.CurlLogger),
		config:     *clientConfig,
	}
}

type CommonResp struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	Trace       string `json:"trace"`
	DataType    int    `json:"data_type"`
	OtelTraceID string `json:"otel_trace_id"`
	OtelSpanID  string `json:"otel_span_id"`
}
