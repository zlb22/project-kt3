package schema

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"time"

	"keti3/internal/application/constant"
)

// TimeInt 自定义时间类型
// 对应数据库timestamp类型
// 接口JSON返回时对应时间戳的int类型（单位秒）
type TimeInt time.Time

type GeneralField struct {
	ID        int64   `gorm:"id" json:"-"`
	CreatedAt TimeInt `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt TimeInt `gorm:"column:updated_at;autoUpdateTime" json:"updated_at,omitempty"`
}

// NewTimeIntByString for time.ParseInLocation
func NewTimeIntByString(dt string) TimeInt {
	t, _ := time.ParseInLocation(constant.TimeFormatDef, dt, time.Local)
	ti := TimeInt(t)
	return ti
}

// NowTimeInt for time.Now
func NowTimeInt() *TimeInt {
	t := TimeInt(time.Now())
	return &t
}

// ToInt64 return timestamp for second
func (ti TimeInt) ToInt64() int64 {
	t := time.Time(ti)
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local).Unix()
}

// Scan implements the Scanner interface.
func (ti *TimeInt) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("can't convert to time.Time")
	}

	*ti = TimeInt(t)
	return nil
}

// Value implements the driver Valuer interface.
func (ti TimeInt) Value() (driver.Value, error) {
	return time.Time(ti), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ti TimeInt) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ti.ToInt64(), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ti *TimeInt) UnmarshalJSON(b []byte) error {
	stamp, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*ti = TimeInt(time.Unix(stamp, 0))
	return nil
}

// ToTime return time.Time
func (ti TimeInt) ToTime() time.Time {
	return time.Time(ti)
}