package main

import (
	"fmt"
	"log"
	"myRpc/Zjprpc/register"
	"time"
)

func main() {
	fmt.Println("注册中心于端口8082监听")
	if err := register.StartHTTPRegisterServer(":8082", 10*time.Second); err != nil {
		log.Fatalf("Start HTTP register server failed: %v", err)
	}
}
