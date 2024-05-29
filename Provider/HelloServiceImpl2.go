package Provider

type HelloServiceImpl2 struct {
}

func (h *HelloServiceImpl2) SayHello(name string) string {
	return "hello_from 2" + name
}
