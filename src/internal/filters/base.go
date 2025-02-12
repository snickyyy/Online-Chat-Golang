package filters

type Filter interface {
	ApplyFilter()
}

type EmptyFilter struct {}

func (EmptyFilter) ApplyFilter(){}
