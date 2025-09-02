package orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"keti3/pkg/tal_content_arch_lib/pontos/kit"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type internalCtxKey string

var (
	dbSystemKeyValue = semconv.DBSystemMySQL

	dbNameKey         = semconv.DBNameKey
	dbUserKey         = semconv.DBUserKey
	dbHostKey         = semconv.DBConnectionStringKey
	dbTableKey        = attribute.Key("db.sql.table")
	dbRowsAffectedKey = attribute.Key("db.rows_affected")
	dbOperationKey    = semconv.DBOperationKey
	dbStatementKey    = semconv.DBStatementKey
	omitVarsKey       = internalCtxKey("omit_vars")
)

func dbTable(name string) attribute.KeyValue {
	return dbTableKey.String(name)
}

func dbStatement(stmt string) attribute.KeyValue {
	return dbStatementKey.String(stmt)
}

func dbCount(n int64) attribute.KeyValue {
	return dbRowsAffectedKey.Int64(n)
}

func dbOperation(op string) attribute.KeyValue {
	return dbOperationKey.String(op)
}

func (op *OtelPlugin) spanName(tx *gorm.DB, operation string) string {
	target := op.dbName

	if tx.Statement != nil && tx.Statement.Table != "" {
		target += "." + tx.Statement.Table
	}

	return fmt.Sprintf("%s %s", operation, target)
}

func operationForQuery(query, op string) string {
	if op != "" {
		return op
	}

	return strings.ToUpper(strings.Split(query, " ")[0])
}

func (op *OtelPlugin) before(operation string) gormHookFunc {
	return func(tx *gorm.DB) {
		tx.Statement.Context = kit.Impl.ConvertGinCtx(tx.Statement.Context)
		tx.Statement.Context, _ = op.tracer.
			Start(tx.Statement.Context, op.spanName(tx, operation), trace.WithSpanKind(trace.SpanKindClient))
	}
}

func extractQuery(tx *gorm.DB) string {
	if shouldOmit, _ := tx.Statement.Context.Value(omitVarsKey).(bool); shouldOmit {
		return tx.Statement.SQL.String()
	}
	return tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
}

const (
	eventMaxSize = 250
	maxChunks    = 4
)

func chunkBy(val string, size int, callback func(string, ...trace.EventOption)) {
	if len(val) > maxChunks*size {
		return
	}

	for i := 0; i < maxChunks*size; i += size {
		end := len(val)
		if end > size {
			end = size
		}
		callback(val[0:end])
		if end > len(val)-1 {
			break
		}
		val = val[end:]
	}
}

func (op *OtelPlugin) after(operation string) gormHookFunc {
	return func(tx *gorm.DB) {
		span := trace.SpanFromContext(tx.Statement.Context)
		if !span.IsRecording() {
			// skip the reporting if not recording
			return
		}
		defer span.End()

		switch tx.Error {
		case nil,
			gorm.ErrRecordNotFound,
			sql.ErrNoRows,
			driver.ErrSkip,
			io.EOF:
			//ignore
		default:
			span.RecordError(tx.Error)
			span.SetStatus(codes.Error, tx.Error.Error())
		}

		//fill attributes
		span.SetAttributes(op.attrs...)

		// extract the db operation
		query := strings.ToValidUTF8(extractQuery(tx), "")
		span.SetAttributes(dbStatement(query))

		opt := operationForQuery(query, operation)
		if tx.Statement.Table != "" {
			span.SetAttributes(dbTable(tx.Statement.Table))
		}

		span.SetName(op.spanName(tx, opt))
		span.SetAttributes(
			dbOperation(opt),
			dbCount(tx.Statement.RowsAffected),
		)
	}
}

func WithOmitVariablesFromTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, omitVarsKey, true)
}
