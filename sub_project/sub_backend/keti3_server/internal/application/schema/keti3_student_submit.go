package schema

type Keti3StudentSubmit struct {
	UID           int64  `gorm:"column:uid" json:"uid"`
	SchoolID      int    `gorm:"column:school_id" json:"school_id"`
	StudentName   string `gorm:"column:student_name" json:"student_name"`
	StudentNum    string `gorm:"column:student_num" json:"student_num"`
	GradeID       int    `gorm:"column:grade_id" json:"grade_id"`
	ClassName     string `gorm:"column:class_name" json:"class_name"`
	VoiceURL      string `gorm:"column:voice_url" json:"voice_url"`
	OplogURL      string `gorm:"column:oplog_url" json:"oplog_url"`
	ScreenshotURL string `gorm:"column:screenshot_url" json:"screenshot_url"`
	Date          string `gorm:"column:date" json:"date"`
	GeneralField
}

func (Keti3StudentSubmit) TableName() string {
	return "keti3_student_submit"
}
