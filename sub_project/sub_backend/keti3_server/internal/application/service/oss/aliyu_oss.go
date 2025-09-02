package oss

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"keti3/internal/application/service/base_video"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Client 请求阿里云存储
func Client(config *base_video.AliOssConfig) *oss.Client {
	var err error
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New(config.EndPoint, config.AccessKeyID, config.AccessKeySecret, oss.SecurityToken(config.SecurityToken))
	if err != nil {
		pontos.Logger.Errorf("oss.New err:%v", err)
		return nil
	}

	return client
}

// PutObject 上传byte数组文件
func PutObject(ctx context.Context, config *base_video.AliOssConfig, fileData []byte, remoteFile string) (fileURL string, err error) {
	client := Client(config)
	// 填写存储空间名称
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("client.Bucket err:%v", err)
		return
	}

	// 将Byte数组上传至exampledir目录下的exampleobject.txt文件。
	err = bucket.PutObject(remoteFile, bytes.NewReader(fileData))
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("bucket.PutObject err:%v", err)
		return
	}

	fileURL = fmt.Sprintf("%s/%s", config.HTTPSDomain, url.QueryEscape(remoteFile))
	return fileURL, nil
}

// PutObjectFromFile 上传本地文件
func PutObjectFromFile(ctx context.Context, config *base_video.AliOssConfig, localPath string, remoteFile string) (fileURL string, err error) {
	client := Client(config)
	// 填写存储空间名称
	bucket, err := client.Bucket(config.Bucket)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("client.Bucket err:%v", err)
		return
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径
	err = bucket.PutObjectFromFile(remoteFile, localPath)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("bucket.PutObjectFromFile err:%v", err)
		return
	}

	fileURL = fmt.Sprintf("%s/%s", config.HTTPSDomain, remoteFile)
	return fileURL, nil
}
