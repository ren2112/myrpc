package main

import (
	"fmt"
	Provider_Common "myRpc/Provider-Common"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"reflect"
)

func main() {
	interfaceName := reflect.TypeOf((*Provider_Common.HelloService)(nil)).Elem().Name()
	invocation := common.NewInvocation(
		interfaceName,
		"SayHello",
		[]string{reflect.TypeOf("").String()},
		[]interface{}{" World"},
	)
	httpClient := new(protocol.HttpClient)
	result, err := httpClient.Send("localhost", 8081, invocation)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

}
