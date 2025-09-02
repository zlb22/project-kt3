package logger

import (
	plog "keti3/pkg/tal_content_arch_lib/pontos-log"

	"github.com/petermattis/goid"
	uuid "github.com/satori/go.uuid"
)

var prefix int64 = 10000000000000000

type ZLogger struct {
	*plog.Logger
}

func NewPLogger(filename string) (*ZLogger, error) {
	newLogger := &ZLogger{
		Logger: plog.New(),
	}
	newLogger.SetFormatter(plog.NewJsonFormatter())
	newLogger.SetOutput(plog.NewFileWriter(filename))
	newLogger.SetReportCaller(true)

	return newLogger, nil
}

// GenTraceId generate trace id
func (l *ZLogger) GenTraceId() string {
	v4 := uuid.NewV4()
	return v4.String()
}

// GenLoggerId generate logger id
func (l *ZLogger) GenLoggerId() int64 {
	id := goid.Get()
	return id + prefix
}
