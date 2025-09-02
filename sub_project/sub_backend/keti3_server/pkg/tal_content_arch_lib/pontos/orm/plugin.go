package orm

import (
	"context"
	"database/sql"
	"fmt"

	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	defaultTracerName = "keti3/pkg/tal_content_arch_lib/pontos/orm" // "go.opentelemetry.io/contrib/instrumentation/github.com/go-gorm/gorm/otelgorm"

	callBackBeforeName = "otel:before"
	callBackAfterName  = "otel:after"

	opCreate = "INSERT"
	opQuery  = "SELECT"
	opDelete = "DELETE"
	opUpdate = "UPDATE"
)

type gormHookFunc func(tx *gorm.DB)

type OtelPlugin struct {
	dbName string
	attrs  []attribute.KeyValue
	tp     oteltrace.TracerProvider
	tracer oteltrace.Tracer
	mp     metric.MeterProvider
	meter  metric.Meter
}

func (op *OtelPlugin) Name() string {
	return "OpenTelemetryPlugin"
}

// NewDefaultPlugin init a default gorm.DB plugin for OpenTelemetry.
func NewDefaultPlugin(host, db, user string) *OtelPlugin {
	return NewPlugin(WithAttributes(dbSystemKeyValue), WithDBName(db), WithDBUser(user), WithDBHost(host))
}

// NewPlugin initialize a new gorm.DB plugin that traces queries
// You may pass optional Options to the function
func NewPlugin(opts ...PluginOption) *OtelPlugin {
	p := &OtelPlugin{}
	for _, opt := range opts {
		opt(p)
	}

	if p.tp == nil {
		p.tp = otel.GetTracerProvider()
	}

	p.tracer = p.tp.Tracer(
		defaultTracerName,
		oteltrace.WithInstrumentationVersion(contrib.SemVersion()),
	)

	if p.mp == nil {
		p.mp = global.MeterProvider()
	}
	p.meter = p.mp.Meter(defaultTracerName)

	return p
}

type registerCallback interface {
	Register(name string, fn func(*gorm.DB)) error
}

func beforeName(name string) string {
	return callBackBeforeName + "_" + name
}

func afterName(name string) string {
	return callBackAfterName + "_" + name
}

func (op *OtelPlugin) Initialize(db *gorm.DB) error {
	if sqlDb, ok := db.ConnPool.(*sql.DB); ok {
		op.ReportDBStatsMetrics(sqlDb)
	}

	registerHooks := []struct {
		callback registerCallback
		hook     gormHookFunc
		name     string
	}{
		// before hooks
		{db.Callback().Create().Before("gorm:before_create"), op.before(opCreate), beforeName("create")},
		{db.Callback().Query().Before("gorm:query"), op.before(opQuery), beforeName("query")},
		{db.Callback().Delete().Before("gorm:before_delete"), op.before(opDelete), beforeName("delete")},
		{db.Callback().Update().Before("gorm:before_update"), op.before(opUpdate), beforeName("update")},
		{db.Callback().Row().Before("gorm:row"), op.before(""), beforeName("row")},
		{db.Callback().Raw().Before("gorm:raw"), op.before(""), beforeName("raw")},

		// after hooks
		{db.Callback().Create().After("gorm:after_create"), op.after(opCreate), afterName("create")},
		{db.Callback().Query().After("gorm:after_query"), op.after(opQuery), afterName("select")},
		{db.Callback().Delete().After("gorm:after_delete"), op.after(opDelete), afterName("delete")},
		{db.Callback().Update().After("gorm:after_update"), op.after(opUpdate), afterName("update")},
		{db.Callback().Row().After("gorm:row"), op.after(""), afterName("row")},
		{db.Callback().Raw().After("gorm:raw"), op.after(""), afterName("raw")},
	}

	for _, h := range registerHooks {
		if err := h.callback.Register(h.name, h.hook); err != nil {
			return fmt.Errorf("register %s hook: %w", h.name, err)
		}
	}

	return nil
}

// ReportDBStatsMetrics reports DBStats metrics using OpenTelemetry Metrics API.
func (op *OtelPlugin) ReportDBStatsMetrics(db *sql.DB, opts ...Option) {
	meter := global.MeterProvider().Meter(defaultTracerName)
	labels := op.attrs

	maxOpenConns, _ := meter.AsyncInt64().Gauge(
		"pontos.sql.connections_max_open",
		instrument.WithDescription("Maximum number of open connections to the database"),
	)
	openConns, _ := meter.AsyncInt64().Gauge(
		"pontos.sql.connections_open",
		instrument.WithDescription("The number of established connections both in use and idle"),
	)
	inUseConns, _ := meter.AsyncInt64().Gauge(
		"pontos.sql.connections_in_use",
		instrument.WithDescription("The number of connections currently in use"),
	)
	idleConns, _ := meter.AsyncInt64().Gauge(
		"pontos.sql.connections_idle",
		instrument.WithDescription("The number of idle connections"),
	)
	connsWaitCount, _ := meter.AsyncInt64().Counter(
		"pontos.sql.connections_wait_count",
		instrument.WithDescription("The total number of connections waited for"),
	)
	connsWaitDuration, _ := meter.AsyncInt64().Counter(
		"pontos.sql.connections_wait_duration",
		instrument.WithDescription("The total time blocked waiting for a new connection"),
		instrument.WithUnit("nanoseconds"),
	)
	connsClosedMaxIdle, _ := meter.AsyncInt64().Counter(
		"pontos.sql.connections_closed_max_idle",
		instrument.WithDescription("The total number of connections closed due to SetMaxIdleConns"),
	)
	connsClosedMaxIdleTime, _ := meter.AsyncInt64().Counter(
		"pontos.sql.connections_closed_max_idle_time",
		instrument.WithDescription("The total number of connections closed due to SetConnMaxIdleTime"),
	)
	connsClosedMaxLifetime, _ := meter.AsyncInt64().Counter(
		"pontos.sql.connections_closed_max_lifetime",
		instrument.WithDescription("The total number of connections closed due to SetConnMaxLifetime"),
	)

	if err := meter.RegisterCallback(
		[]instrument.Asynchronous{
			maxOpenConns,

			openConns,
			inUseConns,
			idleConns,

			connsWaitCount,
			connsWaitDuration,
			connsClosedMaxIdle,
			connsClosedMaxIdleTime,
			connsClosedMaxLifetime,
		},
		func(ctx context.Context) {
			stats := db.Stats()

			maxOpenConns.Observe(ctx, int64(stats.MaxOpenConnections), labels...)

			openConns.Observe(ctx, int64(stats.OpenConnections), labels...)
			inUseConns.Observe(ctx, int64(stats.InUse), labels...)
			idleConns.Observe(ctx, int64(stats.Idle), labels...)

			connsWaitCount.Observe(ctx, stats.WaitCount, labels...)
			connsWaitDuration.Observe(ctx, int64(stats.WaitDuration), labels...)
			connsClosedMaxIdle.Observe(ctx, stats.MaxIdleClosed, labels...)
			connsClosedMaxIdleTime.Observe(ctx, stats.MaxIdleTimeClosed, labels...)
			connsClosedMaxLifetime.Observe(ctx, stats.MaxLifetimeClosed, labels...)
		},
	); err != nil {
		panic(err)
	}
}
