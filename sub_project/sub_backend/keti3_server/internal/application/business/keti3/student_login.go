package keti3

import (
	"context"
	"keti3/internal/application/schema"

	"gorm.io/gorm"
)

type StudentLoginParam struct {
	SchoolID    int    `json:"school_id" binding:"required"`
	StudentName string `json:"student_name" binding:"required"`
	StudentNum  string `json:"student_num" binding:"required"`
	GradeID     int    `json:"grade_id" binding:"required"`
	ClassName   string `json:"class_name" binding:"required"`
}
type StudentLoginRet struct {
	UID int64 `json:"uid"`
}

func (b business) StudentLogin(ctx context.Context, param StudentLoginParam) (ret StudentLoginRet, err error) {
	userParam := &schema.Keti3StudentInfo{
		SchoolID:    param.SchoolID,
		StudentName: param.StudentName,
		StudentNum:  param.StudentNum,
		GradeID:     param.GradeID,
		ClassName:   param.ClassName,
	}
	// 先查询
	data, err := b.stuModel.GetStudent(ctx, userParam)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if data != nil && data.ID > 0 {
		return StudentLoginRet{UID: data.ID}, nil
	}

	uid, err := b.stuModel.Create(ctx, userParam)
	if err != nil {
		return
	}
	return StudentLoginRet{UID: uid}, nil
}
