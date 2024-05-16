package main

import (
	"fmt"
	"myRpc/Zjprpc/proxy"
)

func main() {
	//interfaceName := reflect.TypeOf((*Provider_Common.HelloService)(nil)).Elem().Name()
	//invocation := common.NewInvocation(
	//	interfaceName,
	//	"SayHello",
	//	[]string{reflect.TypeOf("").String()},
	//	[]interface{}{" World"},
	//)
	//httpClient := new(protocol.HttpClient)
	//result, err := httpClient.Send("localhost", 8081, invocation)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	// 创建 RPC 代理
	rpcProxy := proxy.NewRpcProxy()

	// 进行 RPC 调用
	result, err := rpcProxy.Invoke("HelloService", "SayHello", []interface{}{"World"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	rpcProxy2 := proxy.NewRpcProxy()
	// 进行 RPC 调用
	result2, err := rpcProxy2.Invoke("AddService", "Add", []interface{}{1, 2})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result2)
}
