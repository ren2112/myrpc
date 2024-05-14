package main

import (
	"fmt"
	"reflect"
)

type HelloService interface {
	SayHello(name string) string
}

type HelloServiceImpl struct {
}

func (h HelloServiceImpl) SayHello(name string) string {
	fmt.Println("get it!")
	return "hello" + name
}

func set(imp interface{}) {
	fmt.Println(reflect.TypeOf(imp))
	if helloService, ok := imp.(HelloService); ok {
		helloService.SayHello("") // 这里可以安全调用SayHello方法了
	} else {
		fmt.Println("implClass does not implement HelloService")
	}
}

func main() {
	//mp := map[string]interface{}{}
	imp := HelloServiceImpl{}
	set(imp)
}
