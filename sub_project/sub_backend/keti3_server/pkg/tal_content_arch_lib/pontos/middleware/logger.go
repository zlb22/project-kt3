package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/kit"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if _, err := w.body.Write(b); err != nil {
		fmt.Printf("bodyLogWriter err:%v", err)
	}

	return w.ResponseWriter.Write(b)
}

// SetSkipLogKey set the key for skipping request and response log
func SetSkipLogKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constant.SKIPLOGCtxKey, 1)
	}
}

func Logger(log *logger.ZLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, skip := c.Get(constant.SKIPLOGCtxKey)
		if skip {
			return
		}

		// Start timer
		var (
			start   = time.Now()
			path    = c.Request.URL.Path
			raw     = c.Request.URL.RawQuery
			body    []byte
			builder strings.Builder
			sn      = c.GetHeader(constant.SnHeaderKey)
			version = c.GetHeader(constant.VersionHeaderKey)
			blw     = &bodyLogWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		)
		c.Writer = blw

		builder.WriteString(path)
		if raw != "" {
			builder.WriteString("?")
			builder.WriteString(raw)
		}

		span := trace.SpanFromContext(c)
		//log request
		logReq := map[string]interface{}{
			"x_tag":           constant.LogRequestTag,
			"p_request_time":  start.Format(constant.TimeFormatWithMS),
			"p_client_ip":     c.ClientIP(),
			"p_method":        c.Request.Method,
			"p_path":          c.Request.URL.Path,
			"p_raw_query":     builder.String(),
			"p_ua":            c.Request.UserAgent(),
			"p_pressure_test": kit.Impl.IsPressureTest(c),
			"request_header":  kit.Impl.LogHeader(c.Request.Header),
		}

		if c.ContentType() == constant.FORMContentType || c.ContentType() == constant.JSONContentType {
			if c.Request.Body != nil {
				body, _ = ioutil.ReadAll(c.Request.Body)
			}
			logReq["p_request_body"] = string(body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		if sn != "" {
			logReq[constant.DeviceSNCtxKey] = sn
			span.SetAttributes(attribute.String(constant.OtelDeviceSNAttr, sn))
		}

		if version != "" {
			logReq[constant.DeviceVersionCtxKey] = version
			span.SetAttributes(attribute.String(constant.OtelDeviceVersionAttr, sn))
		}

		log.WithContext(c).WithFields(logReq).Info("request")

		// Process request
		c.Next()

		// skip log if necessary
		_, skip = c.Get(constant.SKIPLOGCtxKey)
		if skip {
			return
		}

		// Stop timer
		var (
			end        = time.Now()
			latency    = end.Sub(start)
			statusCode = c.Writer.Status()
			comment    = c.Errors.ByType(gin.ErrorTypePrivate).String()
			buf        = blw.body.Bytes()
			logParams  = map[string]interface{}{
				"x_tag":           constant.LogResponseTag,
				"p_response_time": end.Format(constant.TimeFormatWithMS),
				"p_status_code":   statusCode,
				"p_path":          c.Request.URL.Path,
				"p_raw_query":     builder.String(),
				"p_response_body": string(buf),
				"p_latency":       latency.String(),
				"p_cost_ms":       float64(latency.Microseconds()) / 1e3,
				"request_header":  kit.Impl.LogHeader(c.Request.Header),
				"response_header": kit.Impl.LogHeader(c.Writer.Header()),
			}
		)
		if talId := c.GetString(constant.TalId); talId != "" {
			span.SetAttributes(attribute.String(constant.OtelDeviceTalIDAttr, talId))
			logParams[constant.TalIDCtxKey] = talId
		}
		if comment != "" {
			logParams["p_comment"] = comment
		}

		log.WithContext(c).WithFields(logParams).Info("response")
	}
}
