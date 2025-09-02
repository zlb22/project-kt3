package curl

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ZResponse struct {
	responseCode   int
	responseHeader http.Header
	responseCookie []*http.Cookie
	bytes          []byte
	err            error
}

//Byte return []byte
func (zr *ZResponse) Byte() ([]byte, error) {
	if zr.err != nil {
		return zr.bytes, zr.err
	}
	return zr.bytes, nil
}

// JsonDecode unmarshal to struct
func (zr *ZResponse) JsonDecode(ret interface{}) error {
	if zr.err != nil {
		return zr.err
	}
	if zr.bytes == nil {
		return errors.New("body empty")
	}
	err := json.Unmarshal(zr.bytes, ret)
	return err
}

func (zr *ZResponse) GetRespBody() []byte {
	return zr.bytes
}

//GetRespHeader 获取响应的header
func (zr *ZResponse) GetRespHeader() http.Header {
	return zr.responseHeader
}

//GetRespCookie 获取响应的cookie
func (zr *ZResponse) GetRespCookie() []*http.Cookie {
	return zr.responseCookie
}

//GetRespStatus 获取响应的状态码
func (zr *ZResponse) GetRespStatus() int {
	return zr.responseCode
}
