package protocol

import (
	"flag"
	"fmt"
	"myRpc/Zjprpc/common"
)

func ParseServerArgs() (*common.ServerConfig, error) {
	config := &common.ServerConfig{}

	flag.StringVar(&config.IP, "l", "0.0.0.0", "服务监听 IP 地址")
	flag.IntVar(&config.Port, "p", 0, "服务端监听端口")
	flag.StringVar(&config.ReIP, "r", "", "注册中心 IP 地址")
	flag.IntVar(&config.RePort, "rp", 0, "注册中心端口")
	flag.BoolVar(&config.Help, "h", false, "帮助参数")

	// Parse命令行参数
	flag.Parse()

	if config.Help {
		//flag.PrintDefaults()
		return nil, fmt.Errorf("\n1）-l，服务端监听的 ip 地址(不需要协议)，支持 IPv4 和 IPv6，可以为空，默认监听所有 ip 地址，即 0.0.0.0\n2）-p，服务端监听的端口号，不得为空" +
			"\n3) -r，注册中心的ip地址，不可以为空，需要协议\n4) -rp，注册中心的端口，不可以为空\n" +
			"启动实例：go run <fileName>.go -l 127.0.0.1 -p 8081 -r http://127.0.0.1 -rp 8082")
	}

	if config.Port == 0 {
		//flag.PrintDefaults()
		return nil, fmt.Errorf("端口不能为空")
	}
	if config.ReIP == "" || config.RePort == 0 {
		return nil, fmt.Errorf("注册中心的ip或者端口不能为空")
	}

	return config, nil
}
