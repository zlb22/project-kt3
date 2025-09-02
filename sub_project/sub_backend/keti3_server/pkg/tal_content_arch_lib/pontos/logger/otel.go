package logger

import (
	"reflect"
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	plog "keti3/pkg/tal_content_arch_lib/pontos-log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")
	logGwTraceKey  = attribute.Key("log.gateway.trace") //大鱼网关trace
)

// otel is a hook that adds logs to the active span as events.
type otel struct{}

var _ plog.Hook = (*otel)(nil)

// Fire is a logrus hook that is fired on a new log entry.
func (hook *otel) Fire(entry *plog.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	//Compatible with gin.Context
	ctx = kit.Impl.ConvertGinCtx(ctx)

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(entry.Data)+3+3)

	attrs = append(attrs, logSeverityKey.String(levelString(entry.Level)))
	attrs = append(attrs, logMessageKey.String(entry.Message))
	if t := kit.Impl.GetGatewayTraceId(ctx); t != "" {
		attrs = append(attrs, logGwTraceKey.String(t))
	}

	if entry.Caller != nil {
		if entry.Caller.Function != "" {
			attrs = append(attrs, semconv.CodeFunctionKey.String(entry.Caller.Function))
		}
		if entry.Caller.File != "" {
			attrs = append(attrs, semconv.CodeFilepathKey.String(entry.Caller.File))
			attrs = append(attrs, semconv.CodeLineNumberKey.Int(entry.Caller.Line))
		}
	}

	for k, v := range entry.Data {
		if k == "error" {
			if err, ok := v.(error); ok {
				typ := reflect.TypeOf(err).String()
				attrs = append(attrs, semconv.ExceptionTypeKey.String(typ))
				attrs = append(attrs, semconv.ExceptionMessageKey.String(err.Error()))
				continue
			}
		}

		attrs = append(attrs, plog.Attribute(k, v))
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))

	if entry.Level <= plog.ErrorLevel {
		span.SetStatus(codes.Error, entry.Message)
	}

	return nil
}

// Levels returns levels on which this hook is fired.
func (hook *otel) Levels() []plog.Level {
	return plog.AllLevels
}

func levelString(lvl plog.Level) string {
	s := lvl.String()
	if s == "warning" {
		s = "warn"
	}
	return strings.ToUpper(s)
}

// NewOTELHook returns a logrus hook.
func NewOTELHook() *otel {
	return &otel{}
}
