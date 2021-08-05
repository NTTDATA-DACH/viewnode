package srv

type LoadAndFilter interface {
	LoadAndFilter() (result []ViewNode, err error)
}
