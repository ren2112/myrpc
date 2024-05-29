package Provider

type AddServiceImpl struct {
}

func (add *AddServiceImpl) Add(a, b int) int {
	return a + b
}
