package keti3

import (
	"context"
	"time"

	"keti3/internal/application/constant"
	"keti3/internal/application/service/base_video"
)

type GetSTSAuthRet struct {
	AccessKeyID     string    `json:"AccessKeyId"`
	AccessKeySecret string    `json:"AccessKeySecret"`
	Expiration      time.Time `json:"Expiration"`
	SecurityToken   string    `json:"SecurityToken"`
	Bucket          string    `json:"Bucket"`
	Domain          string    `json:"Domain"`
	HTTPSDomain     string    `json:"HttpsDomain"`
	EndPoint        string    `json:"EndPoint"`
	RootPath        string    `json:"RootPath"`
}

// GetSTSAuth 获取上传 oss 配置请求
func (b business) GetSTSAuth(ctx context.Context) (ret GetSTSAuthRet, err error) {
	authCfg, err := b.baseVideoSrv.GetSTSAuth(ctx, &base_video.GetSTSAuthReq{AppID: constant.STSAuthAppID})
	if err != nil {
		return
	}
	ret = GetSTSAuthRet{
		AccessKeyID:     authCfg.AccessKeyID,
		AccessKeySecret: authCfg.AccessKeySecret,
		Expiration:      authCfg.Expiration,
		SecurityToken:   authCfg.SecurityToken,
		Bucket:          authCfg.Bucket,
		Domain:          authCfg.Domain,
		HTTPSDomain:     authCfg.HTTPSDomain,
		EndPoint:        authCfg.EndPoint,
		RootPath:        "/keti/uid",
	}
	return
}
