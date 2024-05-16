package main

import (
	"log"
	"myRpc/Provider"
	Provider_Common "myRpc/Provider-Common"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/register"
	"reflect"
)

func main() {
	//获取接口名字1
	interfaceName := reflect.TypeOf((*Provider_Common.AddService)(nil)).Elem().Name()
	//新建本地注册
	register.InitLocalRegister()
	localRegister := register.GetInstance()
	localRegister.Regist(interfaceName, "1.0", &Provider.AddServiceImpl{})

	//注册中心注册：
	url := common.URL{interfaceName, "localhost", 8083}
	err := register.RegisterServiceToHTTP(url, "http://localhost:8082")
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	httpServer := new(protocol.HttpServer)
	httpServer.Start(url.HostName, url.Port)
}
