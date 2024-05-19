package main

import (
	"fmt"
	"log"
	"myRpc/Provider"
	"myRpc/Provider-Common"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
	"myRpc/Zjprpc/register"
	"reflect"
	"time"
)

func main() {
	var config *common.ServerConfig
	var err error
	config, err = protocol.ParseServerArgs()
	if err != nil {
		fmt.Println("提示：", err)
		return
	}

	//获取接口名字1
	interfaceName := reflect.TypeOf((*Provider_Common.HelloService)(nil)).Elem().Name()
	//新建本地注册
	register.InitLocalRegister()
	localRegister := register.GetInstance()
	localRegister.Regist(interfaceName, "1.0", &Provider.HelloServiceImpl{})

	//注册中心注册：
	url := common.URL{interfaceName, config.IP, config.Port, time.Now()}
	err = register.RegisterServiceToHTTP(url, fmt.Sprintf("%s:%d", config.ReIP, config.RePort))
	if err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}

	//启动服务
	httpServer := protocol.NewHttpServer(time.Second, fmt.Sprintf("%s:%d", config.ReIP, config.RePort), url)
	httpServer.Start(url.HostName, url.Port)
}
