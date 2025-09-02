package kit

import (
	"net/http"
	"os"
	"strings"
	"time"
)

// TimeToInt64 convert time to int timestamp
func (k *kit) TimeToInt64(t time.Time) int64 {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local).Unix()
}

// FileExist 文件是否存在
func (k *kit) FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return os.IsExist(err)
	}
	return true
}

// LogHeader format header for http request or response
func (k *kit) LogHeader(headers http.Header) map[string]string {
	logHeaders := make(map[string]string)
	if len(headers) > 0 {
		for k, v := range headers {
			value := strings.Join(v, ",")
			logHeaders[k] = value
		}
	}

	return logHeaders
}
