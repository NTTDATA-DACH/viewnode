package srv

type LoadAndFilter interface {
	ResourceName() string
	LoadAndFilter(vns []ViewNode) (result []ViewNode, err error)
}
