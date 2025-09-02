package curl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos/constant"
	"keti3/pkg/tal_content_arch_lib/pontos/kit"
)

var (
	RetryDoTimeoutErr = errors.New("curl service RetryDo timeout error")
	HttpRespStatusErr = errors.New("response not http status ok")
)

type ZRequest struct {
	*http.Client
	*http.Request
	params interface{}
	Logger
	err   error
	ctx   context.Context
	trace bool
	*clientTrace
}

func NewZRequest(ctx context.Context, c *http.Client, req *http.Request, logger Logger) *ZRequest {
	zReq := &ZRequest{
		Client:      c,
		Request:     req,
		params:      nil,
		Logger:      logger,
		err:         nil,
		ctx:         ctx,
		trace:       false,
		clientTrace: nil,
	}

	if ctx != nil {
		zReq.Request = req.WithContext(ctx)
	} else {
		zReq.Request = req.WithContext(context.Background())
	}
	zReq.parseCtx(ctx)
	return zReq
}

func (h *ZRequest) EnableClientTrace() *ZRequest {
	h.trace = true
	h.clientTrace = &clientTrace{}
	return h
}

func (h *ZRequest) SetError(err error) {
	h.err = err
}

func (h *ZRequest) SetUrl(reqUrl string) *ZRequest {
	u, err := url.Parse(reqUrl)
	if err != nil {
		h.err = err
		return h
	}
	h.URL = u
	return h
}

func (h *ZRequest) Do() *ZResponse {
	var (
		response = &ZResponse{}
		err      error
		start    = time.Now()
	)

	defer func() {
		h.doLog(start, response, err)
	}()

	if h.err != nil {
		err, response.err = h.err, h.err
		return response
	}

	if h.trace {
		h.Request = h.Request.WithContext(newOtelClientTrace(h.Context()))
	}

	resp, err := h.Client.Do(h.Request)
	if err != nil {
		response.err = err
		return response
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.err = err
		return response
	}

	if resp.StatusCode != http.StatusOK {
		err = HttpRespStatusErr
		response.responseCode = resp.StatusCode
		response.err = err
		response.bytes = body
		return response
	}

	response.responseCode = resp.StatusCode
	response.responseCookie = resp.Cookies()
	response.responseHeader = resp.Header
	response.bytes = body
	return response
}

func (h *ZRequest) RetryDo(retryNum int, timeout time.Duration) *ZResponse {
	start := time.Now()
	var (
		response = &ZResponse{}
		err      error
		resp     *http.Response
		num      = 0
		tCh      = time.After(timeout)
		ch       = make(chan struct{}, 1)
		reqBody  []byte
		respBody []byte
	)

	if h.Body != nil && h.GetBody != nil {
		reqBody, _ = io.ReadAll(h.Request.Body)
	}

	defer func() {
		h.doLog(start, response, err)
	}()

	if h.err != nil {
		err, response.err = h.err, h.err
		return response
	}

	if h.trace {
		h.Request = h.Request.WithContext(newOtelClientTrace(h.Context()))
	}

	go func(reqBody []byte) {
		defer func() {
			ch <- struct{}{}
		}()
		for {
			num++

			// v1.1.14 每次重试需要set body 或者 http.NewRequest
			if reqBody != nil {
				h.setBody(reqBody)
			}

			resp, err = h.Client.Do(h.Request)
			if err == nil {
				return
			}

			err = fmt.Errorf("wrap: %w", err)
			if num >= retryNum {
				err = fmt.Errorf("%w,retry number: %d", err, num)
				return
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}
	}(reqBody)

	select {
	case <-tCh:
		err = RetryDoTimeoutErr
		response.err = err
		return response
	case <-ch:
	}
	if err != nil {
		response.err = err
		return response
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		response.err = err
		return response
	}

	if resp.StatusCode != http.StatusOK {
		err = HttpRespStatusErr
		response.responseCode = resp.StatusCode
		response.err = err
		response.bytes = respBody
		return response
	}

	response.responseCode = resp.StatusCode
	response.responseCookie = resp.Cookies()
	response.responseHeader = resp.Header
	response.bytes = respBody
	return response
}

// SetParams ...
func (h *ZRequest) SetParams(params map[string]string) *ZRequest {
	h.params = params
	args := &url.Values{}
	for key, value := range params {
		args.Add(key, value)
	}
	h.setBody([]byte(args.Encode()))
	return h
}

// SetJsonParams ...
func (h *ZRequest) SetJsonParams(jsonParams interface{}) *ZRequest {
	h.params = jsonParams
	requestBody, err := json.Marshal(jsonParams)
	if err != nil {
		h.err = errors.New(fmt.Sprintf("json params marshal failed,err:%s", err.Error()))
		return h
	}
	h.setBody(requestBody)
	return h
}

func (h *ZRequest) SetBytesParams(params []byte) *ZRequest {
	h.params = params
	h.setBody(params)
	return h
}

func (h *ZRequest) doLog(start time.Time, response *ZResponse, err error) {
	ctx := h.Request.Context()
	if ctx == nil {
		return
	}

	var (
		query strings.Builder
	)
	query.WriteString(h.URL.Path)

	if h.URL.RawQuery != "" {
		query.WriteString("?")
		query.WriteString(h.URL.RawQuery)
	}

	bodyByte, _ := json.Marshal(h.params)
	fields := map[string]interface{}{
		"p_method":        h.Method,
		"p_host":          h.URL.Host,
		"p_path":          h.URL.Path,
		"p_raw_query":     query.String(),
		"p_request_body":  string(bodyByte),
		"p_response_body": string(response.bytes),
		"request_header":  kit.Impl.LogHeader(h.Request.Header),
		"response_header": kit.Impl.LogHeader(response.responseHeader),
		"p_cost_ms":       float64(time.Now().Sub(start).Microseconds()) / 1e3,
		"x_tag":           constant.LogHttpRequest,
		"p_status_code":   response.responseCode,
	}

	if err == nil {
		h.Info(ctx, "request service ok", fields)
	} else {
		h.Error(ctx, err.Error(), fields)
	}
}
