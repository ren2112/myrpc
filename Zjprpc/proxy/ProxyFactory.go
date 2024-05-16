package proxy

import (
	"encoding/json"
	"fmt"
	"myRpc/Zjprpc/loadbalance"
	"myRpc/Zjprpc/register"
	"reflect"

	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
)

type RpcProxy struct {
	client *protocol.HttpClient
}

// NewRpcProxy 创建一个新的 RpcProxy 实例
func NewRpcProxy() *RpcProxy {
	return &RpcProxy{
		client: new(protocol.HttpClient),
	}
}

// Invoke 进行 RPC 调用
func (p *RpcProxy) Invoke(interfaceName, methodName string, params []interface{}) (interface{}, error) {
	// 创建参数类型列表
	paramTypes := make([]string, len(params))
	for i, param := range params {
		paramTypes[i] = reflect.TypeOf(param).String()
	}

	// 创建 Invocation 对象
	invocation := common.Invocation{
		InterfaceName:  interfaceName,
		MethodName:     methodName,
		ParameterTypes: paramTypes,
		Parameters:     params,
	}

	//服务发现：
	registerAddr := "http://localhost:8082"
	urlList, err := register.QueryServicesFromHTTP(interfaceName, registerAddr)
	if err != nil || len(urlList) == 0 {
		return "", fmt.Errorf("在接口 %s,没有可以用的服务， 报错: %v", interfaceName, err)
	}

	//负载均衡
	url := loadbalance.Random(urlList)

	// 服务调用
	resultStr, err := p.client.Send(url.HostName, url.Port, invocation)
	var result common.Result
	err = json.Unmarshal(resultStr, &result)
	if err != nil {
		return common.Result{}, err
	}

	return result, nil
}
