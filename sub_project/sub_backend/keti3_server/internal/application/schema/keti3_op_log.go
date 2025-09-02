package schema

type Keti3OpLog struct {
	SubmitID      int64  `json:"submit_id" gorm:"column:submit_id;type:int(11);not null;default:0;comment:父级ID"`
	UID           int64  `json:"uid" gorm:"column:uid;type:int(11);not null;default:0;comment:用户ID"`
	OpTime        int    `json:"op_time" gorm:"column:op_time;type:int;not null;default:0;comment:操作时间"`
	OpType        string `json:"op_type" gorm:"column:op_type;type:varchar(32);not null;default:'';comment:操作类型"`
	OpObject      string `json:"op_object" gorm:"column:op_object;type:varchar(20);not null;default:'';comment:操作对象"`
	ObjectNo      int    `json:"object_no" gorm:"column:object_no;type:int(11);not null;default:0;comment:操作对象编号"`
	ObjectName    string `json:"object_name" gorm:"column:object_name;type:varchar(32);not null;default:'';comment:操作对象名称"`
	DataBefore    string `json:"data_before" gorm:"column:data_before;type:json;not null;comment:更改之前的数据"`
	DataAfter     string `json:"data_after" gorm:"column:data_after;type:json;not null;comment:更改之后的数据"`
	VoiceURL      string `json:"voice_url" gorm:"column:voice_url;type:varchar(255);not null;default:'';comment:音频地址"`
	ScreenshotURL string `json:"screenshot_url" gorm:"column:screenshot_url;type:varchar(255);not null;default:'';comment:截图地址"`
	DeletedAt     int    `json:"deleted_at" gorm:"column:deleted_at;type:int(10);not null;default:0;comment:删除时间"`
	GeneralField
}

type OpLogDataJson struct {
	StartCoords OpLogValue `json:"start_coords"`
	EndCoords   OpLogValue `json:"end_coords"`
	Scale       OpLogValue `json:"scale"`
	Aspect      OpLogValue `json:"aspect"`
	Rotate      OpLogValue `json:"rotate"`
	Flip        OpLogValue `json:"flip"`
	Distortion  OpLogValue `json:"distortion"`
}

type OpLogValue struct {
	Value string `json:"value"`
}

func (Keti3OpLog) TableName() string {
	return "keti3_op_log"
}
