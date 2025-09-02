package curl

import (
	"net/http"
)

// SetHeaders 主动设置header头,可以覆盖之前的配置,批量添加
func (h *ZRequest) SetHeaders(header map[string]string) *ZRequest {
	if len(header) > 0 {
		for k, v := range header {
			h.Header.Set(k, v)
		}
	}

	return h
}

// SetHeader 主动设置header头,可以覆盖之前的配置，单个添加
func (h *ZRequest) SetHeader(k, v string) *ZRequest {
	if k != "" {
		h.Header.Set(k, v)
	}
	return h
}

// SetCookies cookie批量添加
func (h *ZRequest) SetCookies(cookies map[string]string) *ZRequest {
	if len(cookies) > 0 {
		for k, v := range cookies {
			h.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}
	return h
}

//SetCookie cookie 单个添加
func (h *ZRequest) SetCookie(k, v string) *ZRequest {
	if k != "" {
		h.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return h
}

//SetReferer 添加referer
func (h *ZRequest) SetReferer(referer string) *ZRequest {
	h.Header.Add("referer", referer)
	return h
}
