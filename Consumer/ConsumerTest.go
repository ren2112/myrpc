package main

import (
	"fmt"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/proxy"
	"sync"
)

func main() {
	var config *common.ClientConfig
	var err error
	config, err = protocol.ParseClientArgs()
	if err != nil {
		fmt.Println("提示：", err)
		return
	}
	// 创建 RPC 代理
	rpcProxy := proxy.NewRpcProxy()

	wg := sync.WaitGroup{}
	// 进行 RPC 调用
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result, err := rpcProxy.Invoke("TestService", fmt.Sprintf("Test%d", i), []interface{}{}, config)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(result)
		}(i)
	}
	wg.Wait()
}
