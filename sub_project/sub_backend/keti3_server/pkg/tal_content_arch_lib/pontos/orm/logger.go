package orm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"

	pontosLogger "keti3/pkg/tal_content_arch_lib/pontos/logger"

	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

/**
 * 实现GORM的Logger(接收context)，用于自定义输出可追踪的日志
type Interface interface {
   LogMode(LogLevel) Interface
   Info(context.Context, string, ...interface{})
   Warn(context.Context, string, ...interface{})
   Error(context.Context, string, ...interface{})
   Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
*/

var (
	once      sync.Once
	SQLLogger logger.Interface
)

func NewSQLLogger(zLogger *pontosLogger.ZLogger) logger.Interface {
	once.Do(func() {
		conf := logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Info,
		}
		SQLLogger = &sqlLog{Config: conf,
			Writer: zLogger,
		}
	})

	return SQLLogger
}

type sqlLog struct {
	Config logger.Config
	Writer *pontosLogger.ZLogger
}

// LogMode ...
func (l *sqlLog) LogMode(logger.LogLevel) logger.Interface {
	return SQLLogger
}

// Info ...
func (l *sqlLog) Info(ctx context.Context, msg string, extra ...interface{}) {
	l.Writer.WithContext(ctx).WithFields(map[string]interface{}{
		"x_tag":   constant.LogSQLTag,
		"p_extra": convertExtraLog(extra...),
	}).Info(msg)
}

// Warn ...
func (l *sqlLog) Warn(ctx context.Context, msg string, extra ...interface{}) {
	l.Writer.WithContext(ctx).WithFields(map[string]interface{}{
		"x_tag":   constant.LogSQLTag,
		"p_extra": convertExtraLog(extra...),
	}).Warn(msg)
}

// Error ...
func (l *sqlLog) Error(ctx context.Context, msg string, extra ...interface{}) {
	l.Writer.WithContext(ctx).WithFields(map[string]interface{}{
		"x_tag":   constant.LogSQLTag,
		"p_extra": convertExtraLog(extra...),
	}).Error(msg)
}

// Trace ...
func (l *sqlLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var (
		elapsed   = time.Since(begin)
		sql, rows = fc()
		file      = utils.FileWithLineNum()
		logField  = map[string]interface{}{
			"x_tag":           constant.LogSQLTag,
			"p_model_file":    file,
			"p_sql":           sql,
			"p_affected_rows": rows,
			"p_cost_ms":       float64(elapsed.Microseconds()) / 1e3,
		}
	)

	switch {
	case err != nil && l.Config.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		logField["p_err_msg"] = err.Error()
		l.Writer.WithContext(ctx).WithFields(logField).Error("mysql error")

	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Config.LogLevel >= logger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.Config.SlowThreshold)

		logField["slow_log"] = slowLog
		l.Writer.WithContext(ctx).WithFields(logField).Warn("slow sql")

	case l.Config.LogLevel == logger.Info:
		l.Writer.WithContext(ctx).WithFields(logField).Info("mysql info")
	}
}

// convertExtraLog ...
func convertExtraLog(extra ...interface{}) string {
	var log []string
	for _, val := range extra {
		log = append(log, cast.ToString(val))
	}
	return cast.ToString(log)
}
