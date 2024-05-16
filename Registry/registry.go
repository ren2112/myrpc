package main

import (
	"fmt"
	"log"
	"myRpc/Zjprpc/register"
)

func main() {
	fmt.Println("注册中心于端口8082监听")
	if err := register.StartHTTPRegisterServer(":8082"); err != nil {
		log.Fatalf("Start HTTP register server failed: %v", err)
	}
}
