package main

import (
	"fmt"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/proxy"
	"sync"
)

func threadTest(wg *sync.WaitGroup, config *common.ClientConfig) {
	defer wg.Done()
	rpcProxy := proxy.NewRpcProxy()

	// 进行 RPC 调用
	result, err := rpcProxy.Invoke("HelloService", "SayHello", []interface{}{"World"}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func main() {
	var config *common.ClientConfig
	var err error
	config, err = protocol.ParseClientArgs()
	if err != nil {
		fmt.Println("提示：", err)
		return
	}

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1) // 在启动每个goroutine前增加WaitGroup计数
		go threadTest(&wg, config)
	}
	wg.Wait()
}
