package Provider

import Provider_Common "myRpc/Provider-Common"

type StructServiceImpl struct {
}

func (st *StructServiceImpl) StructFun() Provider_Common.Person {
	return Provider_Common.Person{Name: "zhoujp", Age: 20}
}
