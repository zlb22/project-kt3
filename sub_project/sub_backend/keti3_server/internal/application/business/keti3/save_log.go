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

	"github.com/go-redis/redis/v8"
	"github.com/tealeg/xlsx"
)

type SaveLogParam struct {
	LogList       []logItem `json:"log_list" binding:"required"`
	UID           int64     `json:"uid" binding:"required"`
	VoiceURL      string    `json:"voice_url"`
	ScreenshotURL string    `json:"screenshot_url"`
	OpTime        int       `json:"op_time" binding:"required"`
}
type logItem struct {
	OpTime     int    `json:"op_time" binding:"required"`
	OpType     string `json:"op_type" binding:"required"`
	OpObject   string `json:"op_object" binding:"required"`
	ObjectName string `json:"object_name" binding:"required"`
	ObjectNo   int    `json:"object_no" binding:"required"`
	DataBefore string `json:"data_before" binding:"required"`
	DataAfter  string `json:"data_after" binding:"required"`
}

// SaveLog 保存日志
func (b business) SaveLog(ctx context.Context, logParam SaveLogParam) error {
	data := make([]schema.Keti3OpLog, 0, len(logParam.LogList))
	for _, item := range logParam.LogList {
		// 保存日志
		data = append(data, schema.Keti3OpLog{
			UID:        logParam.UID,
			OpTime:     item.OpTime,
			OpType:     item.OpType,
			OpObject:   item.OpObject,
			ObjectName: item.ObjectName,
			ObjectNo:   item.ObjectNo,
			DataBefore: item.DataBefore,
			DataAfter:  item.DataAfter,
		})
	}
	// 查询用户信息
	userInfo, err := b.stuModel.GetStudent(ctx, &schema.Keti3StudentInfo{
		GeneralField: schema.GeneralField{ID: logParam.UID},
	})
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("查询用户信息失败: %v", err)
		return err
	}
	// 查询redis中次数
	dayCount, err := b.GetDayCount(ctx, logParam.UID)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("查询次数失败: %v", err)
		return err
	}
	// 生成日志oss文件
	localFile, filename, err := b.genLogFile(ctx, logParam.UID, dayCount, data)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("生成日志文件失败: %v", err)
		return err
	}
	// 上传文件
	ossFile, err := b.genOSSFile(ctx, localFile, filename)
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("上传日志文件失败: %v", err)
		return err
	}

	// 保存提交
	submitID, err := b.stuSumitModel.Create(ctx, &schema.Keti3StudentSubmit{
		UID:           logParam.UID,
		SchoolID:      userInfo.SchoolID,
		StudentName:   userInfo.StudentName,
		StudentNum:    userInfo.StudentNum,
		GradeID:       userInfo.GradeID,
		ClassName:     userInfo.ClassName,
		VoiceURL:      logParam.VoiceURL,
		ScreenshotURL: logParam.ScreenshotURL,
		OplogURL:      ossFile,
		Date:          time.Now().Format(time.DateOnly),
	})
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("保存日志失败: %v", err)
		return err
	}
	// 保存日志
	for i, item := range data {
		item.SubmitID = submitID
		data[i] = item
	}
	return b.logModel.BatchCreate(ctx, data)
}

// 上传文件
func (b business) genOSSFile(ctx context.Context, localFile, filename string) (string, error) {
	authCfg, err := b.baseVideoSrv.GetSTSAuth(ctx, &base_video.GetSTSAuthReq{AppID: constant.STSAuthAppID})
	if err != nil {
		return "", err
	}
	dir := ""
	for _, v := range authCfg.RootPath {
		dir = v
		break
	}
	remoteFile, err := oss.PutObjectFromFile(ctx, authCfg, localFile, fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("PutObjectFromFile failed: %v", err)
		return "", err
	}
	return remoteFile, nil
}

func (b business) genLogFile(ctx context.Context, uid, count int64, logList []schema.Keti3OpLog) (string, string, error) {
	// 生成文件 命名方式：日期_用户ID_PID.xlsx
	date := time.Now().Format(time.DateOnly)
	localPath, err := helper.MakeTempDir(fmt.Sprintf("oplog-file/%s", date))
	if err != nil {
		return "", "", err
	}
	// 按照操作时间排序
	sort.Slice(logList, func(i, j int) bool {
		return logList[i].OpTime < logList[j].OpTime
	})
	filename := fmt.Sprintf("%s_%d_%d.xlsx", date, uid, count)
	localFile := fmt.Sprintf("%s/%s", localPath, filename)

	// 数据写入文件
	xlsxFile := xlsx.NewFile()
	//生成表单对象
	xlsxSheet, err := xlsxFile.AddSheet("sheet1")
	if err != nil {
		return "", "", err
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
	}
	// 保存文件
	err = xlsxFile.Save(localFile)
	if err != nil {
		return "", "", err
	}
	pontos.Logger.WithContext(ctx).Infof("生成文件：%s", localFile)

	return localFile, filename, nil
}

func (b business) GetDayCount(ctx context.Context, uid int64) (int64, error) {
	// 有效期到明天 0 点
	now := time.Now()
	expiredAt := now.Truncate(time.Hour * 24).Add(time.Hour * 24)
	expire := expiredAt.Sub(now)

	rdsKey := fmt.Sprintf("base-keti:op-count:uid:%d", uid)
	// 查询redis中次数
	count, err := pontos.Redis.Get(ctx, rdsKey).Int64()
	if err != nil && err != redis.Nil {
		return -1, err
	}
	count++
	err = pontos.Redis.Set(ctx, rdsKey, count, expire).Err()
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("redis set failed: %v", err)
	}
	return count, nil
}
