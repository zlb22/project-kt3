package base_video

import (
	"context"
	"errors"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type GetPrivateURLReq struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
	AppID    string `json:"app_id"`
}

type GetPrivateURLResp struct {
	CommonResp
	Data *AliOssPrivateURL `json:"data"`
}
type AliOssPrivateURL struct {
	URL string `json:"url"`
}

const (
	ossPrivateURL = "/video-service/api/v1/oss/secret/getUrl"
)

// GetPrivateURL 设备基础信息
func (c Client) GetPrivateURL(ctx context.Context, req *GetPrivateURLReq) (*AliOssPrivateURL, error) {
	var resp *GetPrivateURLResp

	err := c.JsonPostCtx(ctx, c.Host+ossPrivateURL, req).EnableClientTrace().RetryDo(3, c.Timeout).JsonDecode(&resp)
	if err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return nil, pontos.NewAppError(errors.New(resp.Errmsg), resp.Errcode, resp.Errmsg)
	}
	return resp.Data, nil
}
