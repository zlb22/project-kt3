package logger

import (
	"os"
	"sync"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	plog "keti3/pkg/tal_content_arch_lib/pontos-log"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	// topList 兼容之前允许的顶层字段
	topList = []string{
		"x_msg", "x_tag", "p_response_time", "p_status_code", "p_path", "p_host", "p_raw_query", "p_comment", "p_request_body",
		"p_response_body", "p_latency", "p_cost_ms", "p_request_time", "p_client_ip", "p_method", "p_ua", "p_pressure_test",
		"x_department", "x_env", "x_hostname", "x_name", "p_affected_rows", "p_sql", "p_model_file", "p_retry_num", "x_event",
		"x_duration", "x_server_ip", "x_node_name", "x_module", "otel_trace_id", "otel_span_id", "x_trace_id", "x_rpc_id",
		"x_cmd", "x_num_cmd", "x_func", "x_file",
	}
	withFieldsKey = "x_extra"
	topFieldsMap  map[string]struct{}
	once          sync.Once
)

func init() {
	once.Do(func() {
		topFieldsMap = make(map[string]struct{}, len(topList))

		for _, v := range topList {
			topFieldsMap[v] = struct{}{}
		}
	})
}

var _ plog.Hook = (*fieldHook)(nil)

type fieldHook struct {
	app, env string

	topField map[string]struct{}
}

func (f *fieldHook) Levels() []plog.Level {
	return plog.AllLevels
}

func (f *fieldHook) Fire(e *plog.Entry) error {
	hn, _ := os.Hostname()
	e.Data["x_env"] = f.env
	e.Data["x_name"] = f.app
	e.Data["x_hostname"] = hn
	e.Data["x_server_ip"] = os.Getenv(constant.K8SPodIPEnvKey)
	e.Data["x_node_name"] = os.Getenv(constant.K8SNodeNameEnvKey)

	if e.Context == nil {
		return nil
	}

	ctx := kit.Impl.ConvertGinCtx(e.Context)

	e.Data[constant.TraceIDCtxKey] = kit.Impl.GetGatewayTraceId(ctx)
	e.Data[constant.RpcIdCtxKey] = kit.Impl.GetCustomRPCId(ctx)

	e.Data[constant.OTELTraceIdKey] = kit.Impl.GetOTELTraceId(ctx)
	e.Data[constant.OTELSpanIdKey] = kit.Impl.GetOTELSpanId(ctx)

	extra := map[string]interface{}{}
	for k, v := range e.Data {
		if _, ok := f.topField[k]; !ok {
			extra[k] = v
			delete(e.Data, k)
		}
	}

	if len(extra) > 0 {
		if bys, err := json.Marshal(extra); err == nil {
			e.Data[withFieldsKey] = string(bys)
		}
	}
	return nil
}

// NewDefaultFieldHook return a hook for default field to log record
func NewDefaultFieldHook(app, env string, fields ...string) plog.Hook {
	for _, field := range fields {
		topFieldsMap[field] = struct{}{}
	}

	return &fieldHook{
		app:      app,
		env:      env,
		topField: topFieldsMap,
	}
}
