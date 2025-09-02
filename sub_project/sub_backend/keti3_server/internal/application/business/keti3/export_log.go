package keti3

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"keti3/internal/application/constant"
	"keti3/internal/application/helper"
	"keti3/internal/application/schema"
	"keti3/internal/application/service/base_video"
	"keti3/internal/application/service/oss"

	"keti3/pkg/tal_content_arch_lib/pontos"

	"github.com/tealeg/xlsx"
)

func (b business) ExportLog(ctx context.Context) (url string, err error) {
	// 查询所有日志
	count, err := b.logModel.GetCount(ctx)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("GetCount failed: %v", err)
		return
	}
	logList, err := b.GetLogList(ctx, count)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("GetLogList failed: %v", err)
		return
	}
	mUIDLog := b.categoryLogList(logList)
	// 生成多个 csv 文件
	_, localPath, err := b.genXLSXFile(ctx, mUIDLog)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("genXLSXFile failed: %v", err)
		return
	}
	// 打包文件,将这些文件打包成一个压缩包
	zipFile := fmt.Sprintf("%s.zip", localPath)
	if err := helper.ZipDirectory(localPath, zipFile); err != nil {
		pontos.Logger.WithContext(ctx).Errorf("ZipDirectory failed: %v", err)
		return "", err
	}
	// 上传文件
	authCfg, err := b.baseVideoSrv.GetSTSAuth(ctx, &base_video.GetSTSAuthReq{AppID: constant.STSAuthAppID})
	if err != nil {
		return
	}
	dir := ""
	for _, v := range authCfg.RootPath {
		dir = v
		break
	}
	remoteFile, err := oss.PutObjectFromFile(ctx, authCfg, zipFile, fmt.Sprintf("%s/%s", dir, zipFile))
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("PutObjectFromFile failed: %v", err)
		return
	}
	// 返回 url
	urlPrivate, err := b.baseVideoSrv.GetPrivateURL(ctx, &base_video.GetPrivateURLReq{URL: remoteFile, Duration: 3600, AppID: constant.STSAuthAppID})
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("GetPrivateURL failed: %v", err)
		return
	}
	return urlPrivate.URL, nil
}

func (b business) genXLSXFile(ctx context.Context, mUIDLog map[int64][]schema.Keti3OpLog) (localFileList []string, localPath string, err error) {
	// 生成 csv 文件 命名方式：日期+操作记录表+用户ID
	urlDuration := 3600 * 24 * 100
	date := time.Now().Format(time.DateOnly)
	localPath, err = helper.MakeTempDir(fmt.Sprintf("oplog-file/%s", date))
	if err != nil {
		return
	}
	for uid, logList := range mUIDLog {
		if len(logList) == 0 {
			continue
		}
		// 按照操作时间排序
		sort.Slice(logList, func(i, j int) bool {
			return logList[i].OpTime < logList[j].OpTime
		})

		createDate := logList[0].CreatedAt.ToTime().Format(time.DateOnly)
		filename := fmt.Sprintf("%s_操作记录表_%d.xlsx", createDate, uid)
		localFile := fmt.Sprintf("%s/%s", localPath, filename)

		// 数据写入文件
		xlsxFile := xlsx.NewFile()
		//生成表单对象
		xlsxSheet, err := xlsxFile.AddSheet("sheet1")
		if err != nil {
			return nil, "", err
		}
		rh := xlsxSheet.AddRow()
		rh.AddCell().Value = "操作时间"
		rh.AddCell().Value = "操作"
		rh.AddCell().Value = "操作对象类型"
		rh.AddCell().Value = "对象编号"
		rh.AddCell().Value = "图形起始点坐标"
		rh.AddCell().Value = "图形结束点坐标"
		rh.AddCell().Value = "缩放比例"
		rh.AddCell().Value = "宽高比例"
		rh.AddCell().Value = "旋转角度"
		rh.AddCell().Value = "反转"
		rh.AddCell().Value = "扭曲"
		rh.AddCell().Value = "语音地址"
		rh.AddCell().Value = "截图地址"
		for _, log := range logList {
			row := xlsxSheet.AddRow()
			// 向行中添加单元格并设置值
			row.AddCell().Value = time.Unix(int64(log.OpTime), 0).UTC().Format(time.TimeOnly)
			row.AddCell().Value = log.OpType
			if helper.IsInSlice(log.OpType, []string{constant.OpTypeUndoDesc, constant.OpTypeRedoDesc, constant.OpTypeSaveDesc}) {
				continue
			}
			row.AddCell().Value = log.ObjectName
			row.AddCell().Value = strconv.Itoa(log.ObjectNo)

			// 解析数据
			var dataAfter *schema.OpLogDataJson
			json.Unmarshal([]byte(log.DataAfter), &dataAfter)
			if dataAfter == nil {
				continue
			}

			row.AddCell().Value = dataAfter.StartCoords.Value
			row.AddCell().Value = dataAfter.EndCoords.Value
			row.AddCell().Value = dataAfter.Scale.Value
			row.AddCell().Value = dataAfter.Aspect.Value
			row.AddCell().Value = dataAfter.Rotate.Value
			row.AddCell().Value = dataAfter.Flip.Value
			row.AddCell().Value = dataAfter.Distortion.Value
			row.AddCell().Value = func() string {
				if log.VoiceURL == "" {
					return ""
				}

				// 返回 url
				urlPrivate, err := b.baseVideoSrv.GetPrivateURL(ctx, &base_video.GetPrivateURLReq{URL: log.VoiceURL, Duration: urlDuration, AppID: constant.STSAuthAppID})
				if err != nil {
					pontos.Logger.WithContext(ctx).Errorf("log.VoiceURL GetPrivateURL failed: %v", err)
					return ""
				}
				return urlPrivate.URL
			}()
			row.AddCell().Value = func() string {
				if log.ScreenshotURL == "" {
					return ""
				}
				// 返回 url
				urlPrivate, err := b.baseVideoSrv.GetPrivateURL(ctx, &base_video.GetPrivateURLReq{URL: log.ScreenshotURL, Duration: urlDuration, AppID: constant.STSAuthAppID})
				if err != nil {
					pontos.Logger.WithContext(ctx).Errorf("log.ScreenshotURL  GetPrivateURL failed: %v", err)
					return ""
				}
				return urlPrivate.URL
			}()
		}
		// 保存文件
		err = xlsxFile.Save(localFile)
		if err != nil {
			return nil, "", err
		}
		pontos.Logger.WithContext(ctx).Infof("生成文件：%s", localFile)

		localFileList = append(localFileList, localFile)
	}
	return
}

func (b business) GetLogList(ctx context.Context, total int64) (data []schema.Keti3OpLog, err error) {
	// 分批获取日志
	limit := 2000
	offset := 0
	for {
		logList, err := b.logModel.GetList(ctx, limit, offset)
		if err != nil {
			return nil, err
		}
		data = append(data, logList...)
		if len(logList) < limit {
			break
		}
		offset += limit
	}
	return
}

func (b business) categoryLogList(logList []schema.Keti3OpLog) map[int64][]schema.Keti3OpLog {
	m := make(map[int64][]schema.Keti3OpLog)
	for _, log := range logList {
		if _, ok := m[log.UID]; !ok {
			m[log.UID] = make([]schema.Keti3OpLog, 0)
		}
		m[log.UID] = append(m[log.UID], log)
	}
	return m
}
