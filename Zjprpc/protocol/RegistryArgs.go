package protocol

import (
	"flag"
	"fmt"
	"myRpc/Zjprpc/common"
)

func ParseRegistryArgs() (*common.RegistryConfig, error) {
	config := &common.RegistryConfig{}

	flag.StringVar(&config.IP, "i", "0.0.0.0", "注册中心 IP 地址")
	flag.IntVar(&config.Port, "p", 0, "注册中心端口")
	flag.BoolVar(&config.Help, "h", false, "帮助参数")

	// Parse命令行参数
	flag.Parse()

	if config.Help {
		return nil, fmt.Errorf("\n1）-i，注册中心的 ip 地址，同时支持 IPv4 和 IPv6，默认0.0.0.0(此处ip输入不需要协议)\n2）-p，注册中心端口，不得为空。\n" +
			"启动实例：go run <fileName>.go -i 127.0.0.1 -p 8082")
	}

	if config.Port == 0 {
		//flag.PrintDefaults()
		return nil, fmt.Errorf("端口不能为空")
	}

	return config, nil
}
