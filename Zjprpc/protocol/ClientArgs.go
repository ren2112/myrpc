package protocol

import (
	"flag"
	"fmt"
	"myRpc/Zjprpc/common"
)

func ParseClientArgs() (*common.ClientConfig, error) {
	config := &common.ClientConfig{}

	flag.StringVar(&config.IP, "i", "", "服务端 IP 地址")
	flag.IntVar(&config.Port, "p", 0, "服务端端口")
	flag.BoolVar(&config.Help, "h", false, "帮助参数")

	// Parse命令行参数
	flag.Parse()

	if config.Help {
		//flag.PrintDefaults()
		return nil, fmt.Errorf("\n1）-i，客户端需要发送的服务端 ip 地址，同时支持 IPv4 和 IPv6，不得为空\n2）-p，客户端需要发送的服务端端口，不得为空。\n" +
			"启动实例：go run <fileName>.go -i http://127.0.0.1 -p 8082")
	}

	if config.IP == "" || config.Port == 0 {
		//flag.PrintDefaults()
		return nil, fmt.Errorf("IP 地址和端口不能为空")
	}

	return config, nil
}
