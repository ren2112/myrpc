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
	interfaceName := reflect.TypeOf((*Provider_Common.StructService)(nil)).Elem().Name()
	//新建本地注册
	register.InitLocalRegister()
	localRegister := register.GetInstance()
	localRegister.Regist(interfaceName, "1.0", &Provider.StructServiceImpl{})
	//注册结构体
	register.RegisterType(reflect.TypeOf((*Provider_Common.Person)(nil)).Elem().String(), reflect.TypeOf(Provider_Common.Person{}))

	//注册中心注册：
	url := common.URL{interfaceName, "localhost", 8084}
	err := register.RegisterServiceToHTTP(url, "http://localhost:8082")
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	httpServer := new(protocol.HttpServer)
	httpServer.Start(url.HostName, url.Port)
}
