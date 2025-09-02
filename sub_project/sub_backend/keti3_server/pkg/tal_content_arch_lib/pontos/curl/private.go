package curl

import (
	"bytes"
	"context"
	"io"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
)

func (h *ZRequest) parseCtx(ctx context.Context) *ZRequest {
	if ctx == nil {
		return h
	}
	h.setTraceId(ctx)
	return h
}

// 设置链路追踪,trace_id相关
func (h *ZRequest) setTraceId(ctx context.Context) *ZRequest {
	if ctx == nil {
		return h
	}
	traceId := ctx.Value(constant.TraceIDCtxKey)
	if traceId == nil {
		return h
	}

	if traceId, ok := traceId.(string); !ok {
		return h
	} else {
		h.Header.Set(constant.TraceIDCtxKey, traceId)
		h.Header.Set(constant.GWTraceKey, traceId) //设置大鱼网关trace
	}
	return h
}

func (h *ZRequest) setBody(body []byte) *ZRequest {
	bf := bytes.NewBuffer(body)
	h.Body = io.NopCloser(bf)
	h.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bf), nil
	}
	h.ContentLength = int64(len(body))
	return h
}
