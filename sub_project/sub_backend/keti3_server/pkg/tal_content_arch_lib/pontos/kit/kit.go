package kit

type kit struct{}

func NewKit() *kit {
	return &kit{}
}

var Impl = NewKit()
