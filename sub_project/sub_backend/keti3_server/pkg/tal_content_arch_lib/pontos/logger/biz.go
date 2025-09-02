package logger

import (
	plog "keti3/pkg/tal_content_arch_lib/pontos-log"
	"keti3/pkg/tal_content_arch_lib/pontos/constant"
)

type bizTagHook struct{}

var _ plog.Hook = (*bizTagHook)(nil)

func (bizTag *bizTagHook) Levels() []plog.Level {
	return plog.AllLevels
}

func (bizTag *bizTagHook) Fire(e *plog.Entry) error {
	if _, ok := e.Data["x_tag"]; !ok {
		e.Data["x_tag"] = constant.LogBizTag
	}
	return nil
}

// NewBizTagHook return a hook for biz tag
func NewBizTagHook() *bizTagHook {
	return &bizTagHook{}
}
