package main

import (
	"fmt"
	"log"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/register"
	"time"
)

func main() {
	var config *common.RegistryConfig
	var err error
	config, err = protocol.ParseRegistryArgs()
	if err != nil {
		fmt.Println("提示：", err)
		return
	}

	fmt.Printf("注册中心于端口%d监听\n", config.Port)
	if err := register.StartHTTPRegisterServer(fmt.Sprintf("%s:%d", config.IP, config.Port), 10*time.Second); err != nil {
		log.Fatalf("Start HTTP register server failed: %v", err)
	}
}
