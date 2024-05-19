package protocol

import (
	"flag"
	"fmt"
	"myRpc/Zjprpc/common"
)

func ParseClientArgs() (*common.ClientConfig, error) {
	config := &common.ClientConfig{}

	flag.StringVar(&config.IP, "i", "", "服务端 IP 地址")
	flag.StringVar(&config.Port, "p", "", "服务端端口")
	flag.BoolVar(&config.Help, "h", false, "帮助参数")

	flag.Parse()

	if config.Help {
		flag.PrintDefaults()
		return nil, fmt.Errorf("显示帮助信息")
	}

	if config.IP == "" || config.Port == "" {
		flag.PrintDefaults()
		return nil, fmt.Errorf("IP 地址和端口不能为空")
	}

	return config, nil
}
