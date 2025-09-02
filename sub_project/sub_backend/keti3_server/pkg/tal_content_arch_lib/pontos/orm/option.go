package orm

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type PluginOption func(p *OtelPlugin)

// WithTracerProvider configures a tracer traceProvider that is used to create a tracer.
func WithTracerProvider(provider trace.TracerProvider) PluginOption {
	return func(p *OtelPlugin) {
		p.tp = provider
	}
}

// WithMetricProvider configure a metric provider that is used to create a meter
func WithMetricProvider(provider metric.MeterProvider) PluginOption {
	return func(p *OtelPlugin) {
		p.mp = provider
	}
}

// WithAttributes configures attributes that are used to create a span.
func WithAttributes(attrs ...attribute.KeyValue) PluginOption {
	return func(p *OtelPlugin) {
		p.attrs = append(p.attrs, attrs...)
	}
}

// WithDBName configures a db.name attribute.
func WithDBName(name string) PluginOption {
	return func(p *OtelPlugin) {
		p.dbName = name
		p.attrs = append(p.attrs, dbNameKey.String(name))
	}
}

// WithDBUser configures a db.user attribute.
func WithDBUser(user string) PluginOption {
	return func(p *OtelPlugin) {
		p.attrs = append(p.attrs, dbUserKey.String(user))
	}
}

// WithDBHost configures a db.connection_string attribute
func WithDBHost(host string) PluginOption {
	return func(p *OtelPlugin) {
		p.attrs = append(p.attrs, dbHostKey.String(host))
	}
}
