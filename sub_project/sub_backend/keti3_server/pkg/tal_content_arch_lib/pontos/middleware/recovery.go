package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/logger"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery 异常恢复
func Recovery(logger *logger.ZLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)

				var (
					status    = http.StatusInternalServerError //500
					end       = time.Now()
					comment   = c.Errors.ByType(gin.ErrorTypePrivate).String()
					logParams = map[string]interface{}{
						"x_tag":           constant.LogResponseTag,
						"p_response_time": end.Format(constant.TimeFormatWithMS),
						"p_status_code":   status,
						"p_path":          c.Request.URL.Path,
						"p_comment":       comment,
						"p_response_body": http.StatusText(status),
						"error":           fmt.Sprintf("%v", err),
						"stacks":          string(stack),
					}
				)

				if startValue, ok := c.Get("start"); ok {
					if start, ok := startValue.(time.Time); ok {
						latency := end.Sub(start)
						logParams["p_latency"] = latency.String()
						logParams["p_cost_ms"] = float64(latency.Microseconds()) / 1e3
					}
				}

				logger.WithContext(c).WithFields(logParams).Error("recovery form panic error")
				c.AbortWithStatus(status)
			}
		}()
		c.Next()
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}
