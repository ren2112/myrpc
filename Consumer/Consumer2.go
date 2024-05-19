package main

import (
	"fmt"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/proxy"
	"sync"
)

func threadTest(wg *sync.WaitGroup) {
	defer wg.Done()
	rpcProxy := proxy.NewRpcProxy()

	// 进行 RPC 调用
	result, err := rpcProxy.Invoke("HelloService", "SayHello", []interface{}{"World"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

var config *common.ClientConfig

func main() {
	var err error
	config, err = protocol.ParseClientArgs()
	if err != nil {
		fmt.Println("错误：", err)
		return
	}

	fmt.Println("客户端启动参数：")
	fmt.Println("- 服务端 IP 地址：", config.IP)
	fmt.Println("- 服务端端口：", config.Port)

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1) // 在启动每个goroutine前增加WaitGroup计数
		go threadTest(&wg)
	}
	wg.Wait()
}
