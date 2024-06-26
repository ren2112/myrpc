package main

import (
	"fmt"
	Provider_Common "myRpc/Provider-Common"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/proxy"
	"reflect"
)

func main() {
	var config *common.ClientConfig
	var err error
	config, err = protocol.ParseClientArgs()
	if err != nil {
		fmt.Println("提示：", err)
		//return
	}
	// 创建 RPC 代理
	rpcProxy := proxy.NewRpcProxy()

	// 进行 RPC 调用
	result, err := rpcProxy.Invoke("HelloService", "SayHello", []interface{}{1}, config)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println(result)

	rpcProxy2 := proxy.NewRpcProxy()
	// 进行 RPC 调用
	result2, err := rpcProxy2.Invoke("AddService", "Add", []interface{}{1, 2}, config)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println(result2)
	if res, ok := result2[0].(int); ok {
		fmt.Println(res, "is int")
	}

	//测试结构体
	rpcProxy3 := proxy.NewRpcProxy()
	//注册结构体
	proxy.RegisterType(reflect.TypeOf((*Provider_Common.Person)(nil)).Elem().String(), reflect.TypeOf(Provider_Common.Person{}))
	// 进行 RPC 调用
	result3, err := rpcProxy3.Invoke("StructService", "StructFun", []interface{}{}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result3)
	if res, ok := result3[0].(Provider_Common.Person); ok {
		fmt.Println(res, "is struct")
	}

	//测试数组
	rpcProxy4 := proxy.NewRpcProxy()
	//注册所有出现的类型
	proxy.RegisterType(reflect.TypeOf((*[]int)(nil)).Elem().String(), reflect.TypeOf([]int{}))
	// 进行 RPC 调用
	result4, err := rpcProxy4.Invoke("SortService", "QuickSort", []interface{}{[]int{3, 1, 2, 6, 2, 10, 7, 5}}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result4)
}
