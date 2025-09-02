package keti3

import (
	"context"
	"strings"
	"sync"
	"time"

	"keti3/internal/application/constant"
	"keti3/internal/application/model"
	"keti3/internal/application/schema"
	"keti3/internal/application/service/base_video"

	"keti3/pkg/tal_content_arch_lib/pontos"
)

type LogListParam struct {
	SchoolID    int    `json:"school_id"`
	GradeID     int    `json:"grade_id"`
	StudentName string `json:"student_name"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
}
type LogListRet struct {
	Total int64         `json:"total"`
	List  []logListItem `json:"list"`
}
type logListItem struct {
	ID          int64  `json:"id"`
	UID         int64  `json:"uid"`
	Date        string `json:"date"`
	SchoolName  string `json:"school_name"`
	StudentName string `json:"student_name"`
	GradeName   string `json:"grade_name"`
	ClassName   string `json:"class_name"`
	StudentNum  string `json:"student_num"`
	VoiceURL    string `json:"voice_url"`
	ScreenURL   string `json:"screenshot_url"`
	OpLogURL    string `json:"oplog_url"`
	CreateTime  string `json:"create_time"`
}

func (b business) LogList(ctx context.Context, param LogListParam) (ret LogListRet, err error) {
	optionList := make([]model.WithOptions, 0)
	if param.SchoolID > 0 {
		optionList = append(optionList, model.AndWhere("school_id = ?", param.SchoolID))
	}
	if param.GradeID > 0 {
		optionList = append(optionList, model.AndWhere("grade_id = ?", param.GradeID))
	}
	stuName := strings.TrimSpace(param.StudentName)
	if stuName != "" {
		optionList = append(optionList, model.AndWhere("student_name = ?", stuName))
	}
	// 总数
	cnt, err := b.stuSumitModel.Count(ctx, optionList)
	if err != nil {
		return
	}
	ret.Total = cnt

	// 列表
	optionList = append(optionList, model.SetOrder("id desc"))
	optionList = append(optionList, model.WithPage(param.Page, param.PageSize))
	list, err := b.stuSumitModel.GetList(ctx, optionList)
	if err != nil {
		return
	}
	// 请求公开地址
	murl := b.getPublicURL(ctx, list)

	ret.List = make([]logListItem, 0, len(list))
	for _, item := range list {
		pburl := murl[item.ID]
		ret.List = append(ret.List, logListItem{
			ID:          item.ID,
			UID:         item.UID,
			CreateTime:  item.CreatedAt.ToTime().Format(time.DateTime),
			ClassName:   item.ClassName,
			Date:        item.Date,
			SchoolName:  constant.SchoolMap[item.SchoolID],
			StudentName: item.StudentName,
			GradeName:   constant.GradeMap[item.GradeID],
			VoiceURL: func() string {
				if pburl == nil {
					return ""
				}
				return pburl.voiceURL
			}(),
			ScreenURL: func() string {
				if pburl == nil {
					return ""
				}
				return pburl.screenURL
			}(),
			OpLogURL: func() string {
				if pburl == nil {
					return ""
				}
				return pburl.oplogURL
			}(),
		})
	}

	return
}

type pburl struct {
	voiceURL  string
	screenURL string
	oplogURL  string
}

func (b business) getPublicURL(ctx context.Context, list []*schema.Keti3StudentSubmit) map[int64]*pburl {
	syncLock := make(chan struct{}, 5)
	ret := make(map[int64]*pburl)

	var (
		wg   sync.WaitGroup
		lock sync.Mutex
	)
	for _, item := range list {
		wg.Add(1)
		go func(item *schema.Keti3StudentSubmit) {
			defer wg.Done()
			syncLock <- struct{}{}
			defer func() {
				<-syncLock
			}()
			voiceurl, err := b.publicURL(ctx, item.VoiceURL)
			if err != nil {
				pontos.Logger.WithContext(ctx).Errorf("voiceurl publicURL failed: %v", err)
				return
			}
			screenurl, err := b.publicURL(ctx, item.ScreenshotURL)
			if err != nil {
				pontos.Logger.WithContext(ctx).Errorf("screenurl publicURL failed: %v", err)
				return
			}
			oplogurl, err := b.publicURL(ctx, item.OplogURL)
			if err != nil {
				pontos.Logger.WithContext(ctx).Errorf("oplogurl publicURL failed: %v", err)
				return
			}

			t := &pburl{
				voiceURL:  voiceurl,
				screenURL: screenurl,
				oplogURL:  oplogurl,
			}

			lock.Lock()
			defer lock.Unlock()
			ret[item.ID] = t
		}(item)
	}
	wg.Wait()
	return ret
}

func (b business) publicURL(ctx context.Context, url string) (string, error) {
	duration := 3600 * 24 * 3
	urlPrivate, err := b.baseVideoSrv.GetPrivateURL(ctx, &base_video.GetPrivateURLReq{URL: url, Duration: duration, AppID: constant.STSAuthAppID})
	if err != nil {
		pontos.Logger.WithContext(ctx).Errorf("GetPrivateURL failed: %v", err)
		return "", err
	}
	return urlPrivate.URL, nil
}
