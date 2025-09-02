package plog

import (
	"reflect"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

// Option applies a configuration to the given config.
type Option func(h *OTEL)

// WithLevels sets the logging levels on which the hook is fired.
//
// The default is all levels between PanicLevel and WarnLevel inclusive.
func WithLevels(levels ...Level) Option {
	return func(h *OTEL) {
		h.levels = levels
	}
}

// WithErrorStatusLevel sets the minimal logging level on which
// the span status is set to codes.Error.
//
// The default is <= ErrorLevel.
func WithErrorStatusLevel(level Level) Option {
	return func(h *OTEL) {
		h.errorStatusLevel = level
	}
}

var (
	logSeverityKey = attribute.Key("log.severity")
	logMessageKey  = attribute.Key("log.message")
)

// OTEL is a hook that adds logs to the active span as events.
type OTEL struct {
	levels           []Level
	errorStatusLevel Level
}

var _ Hook = (*OTEL)(nil)

// NewOTELHook returns a logrus hook.
func NewOTELHook(opts ...Option) *OTEL {
	hook := &OTEL{
		levels: []Level{
			PanicLevel,
			FatalLevel,
			ErrorLevel,
			WarnLevel,
			InfoLevel,
		},
		errorStatusLevel: ErrorLevel,
	}

	for _, fn := range opts {
		fn(hook)
	}

	return hook
}

// Fire is a logrus hook that is fired on a new log entry.
func (hook *OTEL) Fire(entry *Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(entry.Data)+2+3)

	attrs = append(attrs, logSeverityKey.String(levelString(entry.Level)))
	attrs = append(attrs, logMessageKey.String(entry.Message))

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

		attrs = append(attrs, Attribute(k, v))
	}

	span.AddEvent("log", trace.WithAttributes(attrs...))

	if entry.Level <= hook.errorStatusLevel {
		span.SetStatus(codes.Error, entry.Message)
	}

	return nil
}

// Levels returns levels on which this hook is fired.
func (hook *OTEL) Levels() []Level {
	return hook.levels
}

func levelString(lvl Level) string {
	s := lvl.String()
	if s == "warning" {
		s = "warn"
	}
	return strings.ToUpper(s)
}
