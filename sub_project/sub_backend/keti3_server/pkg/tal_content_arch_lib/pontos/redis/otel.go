package redis

import (
	"context"

	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	"github.com/go-redis/redis/extra/rediscmd/v8"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "keti3/pkg/tal_content_arch_lib/pontos/redis" //"github.com/go-redis/redis/extra/redisotel"
)

type TracingHook struct {
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

var _ redis.Hook = (*TracingHook)(nil)

// NewOTELHook return a redis Hook compatible with gin.Context
func NewOTELHook(idx int) redis.Hook {
	return NewTracingHook(
		WithAttributes(semconv.DBRedisDBIndexKey.Int(idx)),
	)
}

// NewTracingHook return a redis Hook using OpenTelemetry
func NewTracingHook(opts ...Option) *TracingHook {
	cfg := &config{
		tp: otel.GetTracerProvider(),
		attrs: []attribute.KeyValue{
			semconv.DBSystemRedis,
		},
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	tracer := cfg.tp.Tracer(
		defaultTracerName,
		trace.WithInstrumentationVersion("semver:"+redis.Version()),
	)
	return &TracingHook{tracer: tracer, attrs: cfg.attrs}
}

func spanName(name string) string {
	return "RedisCMD " + name
}

func (th *TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx = kit.Impl.ConvertGinCtx(ctx)
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(th.attrs...),
		trace.WithAttributes(
			semconv.DBStatementKey.String(rediscmd.CmdString(cmd)),
		),
	}

	ctx, _ = th.tracer.Start(ctx, spanName(cmd.FullName()), opts...)

	return ctx, nil
}

func (th *TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmd.Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func (th *TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx = kit.Impl.ConvertGinCtx(ctx)
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	summary, cmdsString := rediscmd.CmdsString(cmds)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(th.attrs...),
		trace.WithAttributes(
			semconv.DBStatementKey.String(cmdsString),
			attribute.Int("db.redis.num_cmd", len(cmds)),
		),
	}

	ctx, _ = th.tracer.Start(ctx, spanName("pipeline "+summary), opts...)

	return ctx, nil
}

func (th *TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if err := cmds[0].Err(); err != nil {
		recordError(ctx, span, err)
	}
	span.End()
	return nil
}

func recordError(ctx context.Context, span trace.Span, err error) {
	if err != redis.Nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
