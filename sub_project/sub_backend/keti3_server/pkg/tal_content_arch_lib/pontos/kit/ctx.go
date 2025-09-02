package kit

import (
	"context"
	"net/http"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// IsPressureTest 是否是压测流量
func (k *kit) IsPressureTest(ctx *gin.Context) bool {
	v, ok := ctx.Get(constant.IsPressureTestKey)
	if !ok {
		return false
	}

	return cast.ToInt(v) == 1
}

// NewCtxWithTrace 生成一个携带trace的root节点的span
func (k *kit) NewCtxWithTrace(spanName string) (context.Context, oteltrace.Span) {
	tracer := otel.GetTracerProvider().Tracer(
		"keti3/pkg/tal_content_arch_lib/pontos/kit",
		oteltrace.WithInstrumentationVersion(constant.Version()),
	)

	eCtx := context.Background()
	ctx, span := tracer.Start(eCtx, spanName,
		oteltrace.WithNewRoot(),
		oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
	)

	traceId := span.SpanContext().TraceID().String()

	ctx = context.WithValue(ctx, constant.TraceIDCtxKey, traceId)
	ctx = context.WithValue(ctx, constant.OTELTraceIdKey, traceId)
	ctx = context.WithValue(ctx, constant.OTELSpanIdKey, span.SpanContext().SpanID().String())
	return ctx, span
}

// NewCtxWithParentTrace 集成parent的trace信息，生成新的context。用于协程操作串联trace
func (k *kit) NewCtxWithParentTrace(parent context.Context) context.Context {
	spanContext := trace.SpanContextFromContext(parent)
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)

	traceId := spanContext.TraceID().String()
	spanId := spanContext.SpanID().String()

	if t, ok := parent.Value(constant.TraceIDCtxKey).(string); ok {
		ctx = context.WithValue(ctx, constant.TraceIDCtxKey, t)
	} else {
		ctx = context.WithValue(ctx, constant.TraceIDCtxKey, traceId)
	}

	ctx = context.WithValue(ctx, constant.OTELTraceIdKey, traceId)
	ctx = context.WithValue(ctx, constant.OTELSpanIdKey, spanId)
	return ctx
}

// GetOTELTraceId 适配 OpenTelemetry 的trace id
func (k *kit) GetOTELTraceId(ctx context.Context) string {
	if t, ok := ctx.Value(constant.OTELTraceIdKey).(string); ok {
		if traceId, err := oteltrace.TraceIDFromHex(t); err == nil {
			return traceId.String()
		}
	}

	return ""
}

// GetOTELSpanId 适配 OpenTelemetry 的 span id
func (k *kit) GetOTELSpanId(ctx context.Context) string {
	if t, ok := ctx.Value(constant.OTELSpanIdKey).(string); ok {
		if spanId, err := oteltrace.SpanIDFromHex(t); err == nil {
			return spanId.String()
		}
	}

	return ""
}

// GetGatewayTraceId 获取大鱼网关trace
func (k *kit) GetGatewayTraceId(ctx context.Context) string {
	if t, ok := ctx.Value(constant.TraceIDCtxKey).(string); ok {
		if t != "00000000000000000000000000000000" {
			return t
		}
	}

	return ""
}

// GetCustomRPCId 获取rpcid
func (k *kit) GetCustomRPCId(ctx context.Context) string {
	if r, ok := ctx.Value(constant.RpcIdCtxKey).(string); ok {
		return r
	}
	return ""
}

func (k *kit) GetRequestCtxFromGinCtx(ctx *gin.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx.Request.Context()
}

// ConvertGinCtx convert to context.Context when ctx is gin.Context
func (k *kit) ConvertGinCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	switch (ctx).(type) {
	case *gin.Context:
		if ginCtx, ok := ctx.(*gin.Context); ok {
			if ginCtx.Request == nil {
				ctx = context.Background()
			} else {
				ctx = ginCtx.Request.Context()
			}

			for k, v := range ginCtx.Keys {
				ctx = context.WithValue(ctx, k, v)
			}
		}
	}

	return ctx
}

// CopyGinCtxWithTrace !!! Deprecate: use NewCtxWithTrace or NewCtxWithParentTrace instead.
// 返回的gin.Context 携带参数，Trace等信息，但是清理了timeout,cancel等context的上下文，适用于协程中使用
func (k *kit) CopyGinCtxWithTrace(ctx *gin.Context) *gin.Context {
	ginCtx := ctx.Copy()

	if ctx.Request != nil {
		ginCtx.Request = ctx.Request.WithContext(detachCtx(ctx.Request.Context()))
	} else {
		ginCtx.Request = &http.Request{}
	}

	return ginCtx
}

type detached struct {
	ctx context.Context
}

func (detached) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (detached) Done() <-chan struct{} {
	return nil
}

func (detached) Err() error {
	return nil
}

func (d detached) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

func detachCtx(ctx context.Context) context.Context {
	return detached{ctx: ctx}
}
