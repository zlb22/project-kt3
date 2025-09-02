package schema

type Keti3StudentInfo struct {
	SchoolID    int    `gorm:"column:school_id;not null;default:0;comment:学校id" json:"school_id"`
	StudentName string `gorm:"column:student_name;type:varchar(32);not null;default:'';comment:学生姓名" json:"student_name"`
	StudentNum  string `gorm:"column:student_num;type:varchar(32);not null;default:'';comment:学号" json:"student_num"`
	GradeID     int    `gorm:"column:grade_id;not null;default:0;comment:年级id" json:"grade_id"`
	ClassName   string `gorm:"column:class_name;type:varchar(32);not null;default:'';comment:班级名称" json:"class_name"`
	GeneralField
}

func (Keti3StudentInfo) TableName() string {
	return "keti3_student_info"
}
