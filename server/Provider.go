package main

import (
	"myRpc/Provider"
	"myRpc/Provider-Common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/register"
	"reflect"
)

func main() {
	//获取接口名字1
	interfaceName := reflect.TypeOf((*Provider_Common.HelloService)(nil)).Elem().Name()
	//新建本地注册
	register.InitLocalRegister()
	localRegister := register.GetInstance()
	localRegister.Regist(interfaceName, "1.0", &Provider.HelloServiceImpl{})
	httpServer := new(protocol.HttpServer)
	httpServer.Start("localhost", 8081)
}
