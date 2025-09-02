package pontos

import (
	"fmt"
	"net/http"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/env"

	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/errcode"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// JsonResponse ...
type JsonResponse struct {
	ErrCode     int         `json:"errcode"` //兼容之前字段
	ErrMSG      string      `json:"errmsg"`  //兼容之前字段
	Data        interface{} `json:"data,omitempty"`
	SystemMs    int64       `json:"system_ms,omitempty"`
	Debug       *debugInfo  `json:"debug,omitempty"`
	Trace       string      `json:"trace,omitempty"`
	DataType    int         `json:"data_type"`
	OtelTraceId string      `json:"otel_trace_id,omitempty"`
	OtelSpanId  string      `json:"otel_span_id,omitempty"`
}

type debugInfo struct {
	Err   string `json:"err,omitempty"`
	Stack string `json:"stack,omitempty"`
}

// GenResponse ...
func GenResponse(code int, msg string, opt ...Option) JsonResponse {
	data := JsonResponse{
		ErrCode: code,
		ErrMSG:  msg,
	}

	for _, fn := range opt {
		fn(&data)
	}

	return data
}

type Option func(resp *JsonResponse)

// WithData ...
func WithData(v interface{}) Option {
	return func(resp *JsonResponse) {
		if v != nil {
			resp.Data = v
		}
	}
}

// WithTime ...
func WithTime() Option {
	return func(resp *JsonResponse) {
		resp.SystemMs = time.Now().UnixMilli()
	}
}

// WithDebug contains the stack trace
func WithDebug(e error, stack bool) Option {
	return func(resp *JsonResponse) {
		d := &debugInfo{
			Err: e.Error(),
		}

		if stack {
			d.Stack = fmt.Sprintf("%+v", errors.WithStack(e))
		}

		resp.Debug = d
	}
}

// WithTrace contains the trace in context
func WithTrace(t string) Option {
	return func(resp *JsonResponse) {
		if t != "" {
			resp.Trace = t
		}
	}
}

// WithOtelTraceID set OpenTelemetry trace id
func WithOtelTraceID(t string) Option {
	return func(resp *JsonResponse) {
		resp.OtelTraceId = t
	}
}

// WithOtelSpanID set OpenTelemetry span id
func WithOtelSpanID(span string) Option {
	return func(resp *JsonResponse) {
		resp.OtelSpanId = span
	}
}

func WithDataType(t int) Option {
	return func(resp *JsonResponse) {
		resp.DataType = t
	}
}

// RawData 返回JSON结构体
func RawData(ctx *gin.Context, code int, msg string, opts ...Option) interface{} {
	return GenResponse(code, msg, opts...)
}

// Success 成功响应
func Success(ctx *gin.Context, d interface{}, option ...Option) {
	opts := []Option{
		WithData(d),
		WithTrace(ToolKit.GetGatewayTraceId(ctx)),
		WithOtelTraceID(ToolKit.GetOTELTraceId(ctx)),
		WithOtelSpanID(ToolKit.GetOTELSpanId(ctx)),
		WithDataType(constant.RAW),
	}

	if len(option) > 0 {
		opts = append(opts, option...)
	}

	ctx.JSON(http.StatusOK, RawData(ctx, 0, "ok", opts...))
}

// EncryptSuccess 加密响应
func EncryptSuccess(ctx *gin.Context, d interface{}) {
	if !Config.GetBool("app_sign.is_open") {
		Success(ctx, d)
		return
	}

	secret := Config.GetString("app_sign.app_secret")
	data, err := kit.Encrypt(d, secret)
	if err != nil {
		sysError := NewAppError(err, errcode.DataEncryptionError, errcode.ErrMSG(errcode.DataEncryptionError))
		Error(ctx, sysError, nil)
		return
	}

	opts := []Option{
		WithData(data),
		WithTrace(ToolKit.GetGatewayTraceId(ctx)),
		WithOtelTraceID(ToolKit.GetOTELTraceId(ctx)),
		WithOtelSpanID(ToolKit.GetOTELSpanId(ctx)),
		WithDataType(constant.Encrypt),
	}

	ctx.JSON(http.StatusOK, RawData(ctx, 0, "ok", opts...))
}

// Error 错误响应
func Error(ctx *gin.Context, err error, d interface{}) {
	var (
		code   = errcode.UnknownCode
		msg    string
		stack  bool
		status = http.StatusOK
	)

	// if is aborted, preserve gin's state.
	// eg: ctx.Bind() ctx.MustBindXXX() will abort the request with HTTP 400 if any error occurs.
	if ctx.IsAborted() {
		status = ctx.Writer.Status()
	}

	if err == nil {
		ctx.JSON(status, GenResponse(0, "error nil"))
		return
	}

	opts := []Option{
		WithData(d),
		WithTrace(ToolKit.GetGatewayTraceId(ctx)),
		WithOtelTraceID(ToolKit.GetOTELTraceId(ctx)),
		WithOtelSpanID(ToolKit.GetOTELSpanId(ctx)),
	}

	if !env.IsOnlineEnv {
		opts = append(opts, WithDebug(err, stack))
	}

	//判断error类型
	switch raw := (err).(type) {
	case AppError:
		code, msg = raw.CodeWithMsg()

	default:
		msg = err.Error()
	}

	ctx.JSON(status, GenResponse(code, msg, opts...))
}
