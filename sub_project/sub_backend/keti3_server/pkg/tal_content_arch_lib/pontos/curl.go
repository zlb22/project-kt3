package pontos

import (
	"context"
	"io"
	"net/http"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/curl"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric/global"
)

const (
	defaultTimeout = time.Second * 3
)

const (
	METHODGET  = "GET"
	METHODPOST = "POST"
)

type CurlClient struct {
	*http.Client
	curl.Logger
}

// NewOTELTransport return a http Transport using otelhttp
func NewOTELTransport() *otelhttp.Transport {
	return otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Host + r.URL.Path
		}),
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithMeterProvider(global.MeterProvider()),
	)
}

// CurlOption for extension
type CurlOption interface {
	apply(client *CurlClient)
}

type curlOptionFunc func(client *CurlClient)

func (o curlOptionFunc) apply(cli *CurlClient) {
	o(cli)
}

// WithTimeout set timeout for curl client
func WithTimeout(t time.Duration) CurlOption {
	return curlOptionFunc(func(client *CurlClient) {
		client.Timeout = t
	})
}

// NewCurlClientWithOption return a curl client with options
func NewCurlClientWithOption(opts ...CurlOption) *CurlClient {
	cc := &CurlClient{
		Client: &http.Client{
			Transport: NewOTELTransport(),
			Timeout:   defaultTimeout,
		},
		Logger: curl.NewEmptyLogger(),
	}

	for _, opt := range opts {
		opt.apply(cc)
	}

	return cc
}

func NewCurlClient(timeout time.Duration) *CurlClient {
	return &CurlClient{
		Client: &http.Client{
			Transport: NewOTELTransport(),
			Timeout:   timeout,
		},
		Logger: curl.NewEmptyLogger(),
	}
}

// SetLogger 更改logger实现
func (h *CurlClient) SetLogger(logger curl.Logger) *CurlClient {
	h.Logger = logger
	return h
}

// 内部接口请求请使用带ctx的方法(带链路追踪),请求外部接口用不带ctx的

func (h *CurlClient) GetCtx(ctx context.Context, reqUrl string) *curl.ZRequest {
	return h.defaultZReq(ctx, METHODGET).SetHeader("Content-Type", constant.FORMContentType).SetUrl(reqUrl)
}

func (h *CurlClient) PostCtx(ctx context.Context, reqUrl string, params map[string]string) *curl.ZRequest {
	return h.defaultZReq(ctx, METHODPOST).SetHeader("Content-Type", constant.FORMContentType).SetUrl(reqUrl).SetParams(params)
}

func (h *CurlClient) JsonPostCtx(ctx context.Context, reqUrl string, params interface{}) *curl.ZRequest {
	return h.defaultZReq(ctx, METHODPOST).SetHeader("Content-Type", constant.JSONContentType).SetUrl(reqUrl).SetJsonParams(params)
}

func (h *CurlClient) IoPostCtx(ctx context.Context, reqUrl string, body io.Reader) *curl.ZRequest {
	req, err := http.NewRequest(METHODPOST, reqUrl, body)
	zReq := curl.NewZRequest(ctx, h.Client, req, h.Logger)
	if err != nil {
		zReq.SetError(err)
	}
	return zReq
}

// 以下方法不携带ctx信息

func (h *CurlClient) Get(reqUrl string) *curl.ZRequest {
	return h.defaultZReq(context.Background(), METHODGET).SetHeader("Content-Type", constant.FORMContentType).SetUrl(reqUrl)
}

func (h *CurlClient) Post(reqUrl string, params map[string]string) *curl.ZRequest {
	return h.defaultZReq(context.Background(), METHODPOST).SetHeader("Content-Type", constant.FORMContentType).SetUrl(reqUrl).SetParams(params)
}

func (h *CurlClient) JsonPost(reqUrl string, params interface{}) *curl.ZRequest {
	return h.defaultZReq(context.Background(), METHODPOST).SetHeader("Content-Type", constant.JSONContentType).SetUrl(reqUrl).SetJsonParams(params)
}

func (h *CurlClient) defaultReq(method string) *http.Request {
	req := &http.Request{
		Method: method,
		Header: make(http.Header),
		Close:  true,
	}
	return req
}

func (h *CurlClient) defaultZReq(ctx context.Context, method string) *curl.ZRequest {
	req := h.defaultReq(method)
	return curl.NewZRequest(ctx, h.Client, req, h.Logger)
}

// curlCtxLogger 默认的curl client 日志打印实现
type curlCtxLogger struct {
}

func NewCurlCtxLogger() *curlCtxLogger {
	return &curlCtxLogger{}
}

func (l *curlCtxLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Infof(format, args...)
}

func (l *curlCtxLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Warnf(format, args...)
}

func (l *curlCtxLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Errorf(format, args...)
}

func (l *curlCtxLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	Logger.WithContext(ctx).Fatalf(format, args...)
}

func (l *curlCtxLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	Logger.WithContext(ctx).WithFields(fields).Info(msg)
}

func (l *curlCtxLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	Logger.WithContext(ctx).WithFields(fields).Warn(msg)
}

func (l *curlCtxLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	Logger.WithContext(ctx).WithFields(fields).Error(msg)
}

func (l *curlCtxLogger) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	Logger.WithContext(ctx).WithFields(fields).Fatal(msg)
}
