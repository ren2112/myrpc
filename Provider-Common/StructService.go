package Provider_Common

type Person struct {
	Name string
	Age  int
}

type StructService interface {
	StructFun() Person
}
