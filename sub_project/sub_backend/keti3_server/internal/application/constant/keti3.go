package constant

// -- op_type操作类型
// -- 新增 add
// -- 移动 move
// -- 删除 delete
// -- 缩放 scale
// -- 宽高 aspect
// -- 旋转 rotate
// -- 翻转 flip
// -- 扭曲 distortion
// -- 撤销 undo
// -- 重做 redo
// -- 下一步 save

const (
	// OpTypeAdd 新增
	OpTypeAdd = "add"
	// OpTypeMove 移动
	OpTypeMove = "move"
	// OpTypeDelete 删除
	OpTypeDelete = "delete"
	// OpTypeScale 缩放
	OpTypeScale = "scale"
	// OpTypeAspect 宽高
	OpTypeAspect = "aspect"
	// OpTypeRotate 旋转
	OpTypeRotate = "rotate"
	// OpTypeFlip 翻转
	OpTypeFlip = "flip"
	// OpTypeDistortion 扭曲
	OpTypeDistortion = "distortion"
	// OpTypeUndo 撤销
	OpTypeUndo = "undo"
	// OpTypeRedo 重做
	OpTypeRedo = "redo"
	// OpTypeSave 下一步
	OpTypeSave   = "save"
	OpTypeUpload = "upload"

	OpTypeUndoDesc = "Undo"
	OpTypeRedoDesc = "Redo"
	OpTypeSaveDesc = "Save"
)

const (
	STSAuthAppID = "xkedj1010"
	// STSAuthAppID = "uyqkz9998"
)

var (
	// OpTypeMap 操作类型map
	OpTypeMap = map[string]string{
		OpTypeAdd:        "新增对象",
		OpTypeMove:       "移动对象",
		OpTypeDelete:     "删除对象",
		OpTypeScale:      "调整缩放比例",
		OpTypeAspect:     "调整宽高比例",
		OpTypeRotate:     "调整旋转角度",
		OpTypeFlip:       "翻转对象",
		OpTypeDistortion: "扭曲对象",
		OpTypeUndo:       "Undo",
		OpTypeRedo:       "Redo",
		OpTypeSave:       "Save",
	}
)

// 学校
// 江北中学
// 鲁能巴蜀中学
// 育才中学
// 珊瑚教育集团
// 双福育才中学
// 四牌坊小学
// 四牌坊尚融小学
// 重庆文德中学校
// 重庆市十一中学校
// 重庆融汇清华实验中学校
// 重庆第二外国语学校
// 其他学校
const (
	SchoolXNDXFSZX    = "西南大学附属中学"
	SchoolJBZX        = "江北中学"
	SchoolLNBSCZX     = "鲁能巴蜀中学"
	SchoolYCZX        = "育才中学"
	SchoolSHJYJT      = "珊瑚教育集团"
	SchoolSFYCZX      = "双福育才中学"
	SchoolSPFX        = "四牌坊小学"
	SchoolSPFSRX      = "四牌坊尚融小学"
	SchoolCQWDZXX     = "重庆文德中学校"
	SchoolCCSSYZXX    = "重庆市十一中学校"
	SchoolCQRHQHSYZXX = "重庆融汇清华实验中学校"
	SchoolCCDEWGYXX   = "重庆第二外国语学校"
	SchoolQTXX        = "其他学校"
)

var (
	SchoolMap = map[int]string{
		1:  SchoolXNDXFSZX,
		2:  SchoolJBZX,
		3:  SchoolLNBSCZX,
		4:  SchoolYCZX,
		5:  SchoolSHJYJT,
		6:  SchoolSFYCZX,
		7:  SchoolSPFX,
		8:  SchoolSPFSRX,
		9:  SchoolQTXX,
		10: SchoolCQWDZXX,
		11: SchoolCCSSYZXX,
		12: SchoolCQRHQHSYZXX,
		13: SchoolCCDEWGYXX,
	}

	SchoolList = []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 9} // 学校排序
)

// 年级
const (
	GradeOne    = "一年级"
	GradeTwo    = "二年级"
	GradeThree  = "三年级"
	GradeFour   = "四年级"
	GradeFive   = "五年级"
	GradeSix    = "六年级"
	GradeSeven  = "七年级"
	GradeEight  = "八年级"
	GradeNine   = "九年级"
	GradeTen    = "高一"
	GradeEleven = "高二"
	GradeTwelve = "高三"
)

var (
	GradeMap = map[int]string{
		1:  GradeOne,
		2:  GradeTwo,
		3:  GradeThree,
		4:  GradeFour,
		5:  GradeFive,
		6:  GradeSix,
		7:  GradeSeven,
		8:  GradeEight,
		9:  GradeNine,
		10: GradeTen,
		11: GradeEleven,
		12: GradeTwelve,
	}
)
