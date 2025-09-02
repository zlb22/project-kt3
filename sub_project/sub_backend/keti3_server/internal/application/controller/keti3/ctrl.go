package keti3

import "keti3/internal/application/business/keti3"

// Ctrl TODO delete demo usages.
type Ctrl struct {
	biz keti3.Business
}

// NewCtrl ...
func NewCtrl() *Ctrl {
	return &Ctrl{
		biz: keti3.NewBusiness(),
	}
}
