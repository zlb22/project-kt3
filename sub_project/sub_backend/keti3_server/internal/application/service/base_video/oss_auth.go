package base_video

import (
	"context"
	"errors"
	"time"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type GetSTSAuthReq struct {
	AppID string `json:"app_id"`
}

type GetSTSAuthResp struct {
	CommonResp
	Data *AliOssConfig `json:"data"`
}
type AliOssConfig struct {
	AccessKeyID     string            `json:"AccessKeyId"`
	AccessKeySecret string            `json:"AccessKeySecret"`
	Expiration      time.Time         `json:"Expiration"`
	SecurityToken   string            `json:"SecurityToken"`
	Bucket          string            `json:"Bucket"`
	Domain          string            `json:"Domain"`
	HTTPSDomain     string            `json:"HttpsDomain"`
	EndPoint        string            `json:"EndPoint"`
	RootPath        map[string]string `json:"RootPath"`
}

const (
	stsAuthURL = "/video-service/api/v1/oss/upload/sts/auth"
)

// SNBasic 设备基础信息
func (c Client) GetSTSAuth(ctx context.Context, req *GetSTSAuthReq) (*AliOssConfig, error) {
	var resp *GetSTSAuthResp

	err := c.JsonPostCtx(ctx, c.Host+stsAuthURL, req).EnableClientTrace().RetryDo(3, c.Timeout).JsonDecode(&resp)
	if err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return nil, pontos.NewAppError(errors.New(resp.Errmsg), resp.Errcode, resp.Errmsg)
	}
	return resp.Data, nil
}
