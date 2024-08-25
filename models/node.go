package models

type Node struct {
	IpAdd string
	Name  string
}

func CreateNode(ip, name string) *Node {
	return &Node{
		IpAdd: ip,
		Name:  name,
	}
}
