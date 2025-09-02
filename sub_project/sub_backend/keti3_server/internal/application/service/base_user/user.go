package base_user

import (
	"context"
	"encoding/json"
	"errors"
)

type GetUserInfoReq struct {
	UserID string `json:"user_id"`
}

type GetUserInfoResp struct {
	CommonResp
	Data any `json:"data"`
}
type UserInfo struct {
	UserID              string      `json:"user_id"`                // 用户ID
	UserName            string      `json:"user_name"`              // 名字
	Sex                 int         `json:"sex"`                    // 性别 1=男 2=女
	StudentType         int         `json:"student_type"`           // 学生类型 1=小学生 2=高中生
	Nation              int         `json:"nation"`                 // 民族
	Birthday            string      `json:"birthday"`               // 出生日期
	ProvinceCode        string      `json:"province_code"`          // 省份编码
	CityCode            string      `json:"city_code"`              // 城市编码
	CountryCode         string      `json:"country_code"`           // 区县编码
	SchoolName          string      `json:"school_name"`            // 学校名称
	GradeCode           int         `json:"grade_code"`             // 年级编码
	ClassName           string      `json:"class_name"`             // 班级名称
	StudyNo             string      `json:"study_no"`               // 学号
	Height              float64     `json:"height"`                 // 身高
	Weight              float64     `json:"weight"`                 // 体重
	MotherContact       string      `json:"mother_contact"`         // 母亲联系方式
	FatherContact       string      `json:"father_contact"`         // 父亲联系方式
	MotherEduLevel      int         `json:"mother_edu_level"`       // 母亲教育程度
	FatherEduLevel      int         `json:"father_edu_level"`       // 父亲教育程度
	FatherBusiness      int         `json:"father_business"`        // 父亲所属行业
	MotherBusiness      int         `json:"mother_business"`        // 母亲所属行业
	Region              int         `json:"region"`                 // 地区
	ClassRank           int         `json:"class_rank"`             // 班级排名
	MoveHouse           int         `json:"move_house"`             // 近一年搬家次数
	MusicHobby          string      `json:"music_hobby"`            // 音乐爱好
	MusicHobbyLearnTime string      `json:"music_hobby_learn_time"` // 音乐爱好学习时长
	FamilyStruct        int         `json:"family_struct"`          // 家庭结构
	Brothers            int         `json:"brothers"`               // 兄弟姐妹数
	FamilyLevel         int         `json:"family_level"`           // 家庭幸福指数
	GradeDetail         GradeDetail `json:"grade_detail"`           // 高中生分数明细
}

type GradeDetail struct {
	ZkGrade             int `json:"zk_grade"`               // 中考成绩
	ZkFullGrade         int `json:"zk_full_grade"`          // 中考满分分数
	ZkLanguageGrade     int `json:"zk_language_grade"`      // 语文成绩
	ZkLanguageFullGrade int `json:"zk_language_full_grade"` // 语文满分
	ZkMathGrade         int `json:"zk_math_grade"`          // 数学成绩
	ZkMathFullGrade     int `json:"zk_math_full_grade"`     // 数学满分
	ZkEnglishGrade      int `json:"zk_english_grade"`       // 英语成绩
	ZkEnglishFullGrade  int `json:"zk_english_full_grade"`  // 英语满分
}

const (
	getUserInfoURL = "/gtc-ice/api/user/get"
)

// SNBasic 设备基础信息
func (c Client) CheckGetUserInfo(ctx context.Context, header map[string]string) (*UserInfo, int, error) {
	var resp *GetUserInfoResp

	err := c.GetCtx(ctx, c.Host+getUserInfoURL).SetHeaders(header).EnableClientTrace().RetryDo(3, c.Timeout).JsonDecode(&resp)
	if err != nil {
		return nil, 0, err
	}
	if resp.Errcode != 0 {
		return nil, resp.Errcode, errors.New(resp.Errmsg)
	}
	data := &UserInfo{}
	jsonData, _ := json.Marshal(resp.Data)
	_ = json.Unmarshal(jsonData, data)

	return data, resp.Errcode, nil
}
