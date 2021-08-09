package srv

type LoadAndFilter interface {
	LoadAndFilter(vns []ViewNode) (result []ViewNode, err error)
}
