package curl

import (
	"context"
)

type Logger interface {
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})

	Info(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, fields map[string]interface{})
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
}

type logger struct {
}

func NewEmptyLogger() *logger {
	return &logger{}
}

func (l *logger) Infof(ctx context.Context, format string, args ...interface{})        {}
func (l *logger) Warnf(ctx context.Context, format string, args ...interface{})        {}
func (l *logger) Errorf(ctx context.Context, format string, args ...interface{})       {}
func (l *logger) Fatalf(ctx context.Context, format string, args ...interface{})       {}
func (l *logger) Info(ctx context.Context, msg string, fields map[string]interface{})  {}
func (l *logger) Warn(ctx context.Context, msg string, fields map[string]interface{})  {}
func (l *logger) Error(ctx context.Context, msg string, fields map[string]interface{}) {}
func (l *logger) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {}
