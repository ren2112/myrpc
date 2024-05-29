package Provider

type HelloServiceImpl struct {
}

func (h *HelloServiceImpl) SayHello(name string) string {
	return "hello" + name
}
