package main

type geth struct {
	name       string
	storage    string
	nodeCount  int
	nameTX     string
	id         string
	accountPwd string
}

func newGeth(name string) geth {
	g := geth{
		name:       name,
		storage:    "",
		nodeCount:  0,
		nameTX:     "",
		id:         "",
		accountPwd: "",
	}
	return g
}

func (g *geth) setNodeCount(count int) {
	g.nodeCount = count
}

func (g *geth) setStorageCapacity(storage string) {
	g.storage = storage
}

func (g *geth) setNameOfTX(nameTX string) {
	g.nameTX = nameTX
}

func (g *geth) setid(id string) {
	g.id = id
}

func (g *geth) setAccountPwd(pwd string) {
	g.accountPwd = pwd
}
