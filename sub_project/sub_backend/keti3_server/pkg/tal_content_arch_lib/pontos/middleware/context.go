package middleware

import (
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// Context 处理上下文相关
func Context(log *logger.ZLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToUpper(ctx.Request.Method) == "OPTIONS" {
			ctx.Set(constant.SKIPLOGCtxKey, 1)
		}

		otelTraceId := trace.SpanContextFromContext(ctx.Request.Context()).TraceID().String()

		//历史逻辑保留
		traceId := ctx.GetHeader("x_trace_id")
		if traceId == "" {
			traceId = ctx.GetHeader("trace_id")
		}
		if traceId == "" {
			traceId = ctx.GetHeader("Traceid") //兼容大鱼网关trace
		}
		if traceId == "" {
			traceId = otelTraceId //网关没透传直接用 opentelemetry 的trace
		}

		rpcId := ctx.GetHeader("x_rpc_id")
		if rpcId == "" {
			rpcId = ctx.GetHeader("rpc_id")
		}

		if rpcId == "" {
			rpcId = "1"
		}
		ctx.Set(constant.TraceIDCtxKey, traceId)
		ctx.Set(constant.RpcIdCtxKey, rpcId)

		//opentelemetry
		ctx.Set(constant.OTELTraceIdKey, otelTraceId)
		ctx.Set(constant.OTELSpanIdKey, trace.SpanContextFromContext(ctx.Request.Context()).SpanID().String())
	}
}
