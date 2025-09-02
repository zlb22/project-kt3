package keti3

import (
	"context"
	"keti3/internal/application/constant"
	"sort"
)

type ConfigListRet struct {
	School []cfgItem `json:"school"`
	Grade  []cfgItem `json:"grade"`
}

type cfgItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b *business) ConfigList(ctx context.Context) (ret ConfigListRet, err error) {
	for _, sid := range constant.SchoolList {
		if sname, ok := constant.SchoolMap[sid]; ok {
			ret.School = append(ret.School, cfgItem{ID: sid, Name: sname})
		}
	}
	// for k, v := range constant.SchoolMap {
	// 	ret.School = append(ret.School, cfgItem{ID: k, Name: v})
	// }
	for k, v := range constant.GradeMap {
		ret.Grade = append(ret.Grade, cfgItem{ID: k, Name: v})
	}

	// sort.Slice(ret.School, func(i, j int) bool {
	// 	return ret.School[i].ID < ret.School[j].ID
	// })
	sort.Slice(ret.Grade, func(i, j int) bool {
		return ret.Grade[i].ID < ret.Grade[j].ID
	})

	return
}
