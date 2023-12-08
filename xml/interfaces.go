package xml

type Name struct {
	ns   *NS
	name string
}

type NS struct {
	uri    string
	prefix string
}

func (ns NS) Equal(other NS) bool {
	return ns.uri == other.uri
}
