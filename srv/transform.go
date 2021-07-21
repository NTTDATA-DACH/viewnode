package srv

type Transform interface {
	Transform() (result []ViewNode, err error)
}

type Filter interface {
	Filter(nodes []ViewNode) (result []ViewNode, err error)
}
