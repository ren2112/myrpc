package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"myRpc/Zjprpc/loadbalance"
	"myRpc/Zjprpc/register"
	"reflect"
	"sync"
	"time"

	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/protocol"
)

type RpcProxy struct {
	client *protocol.HttpClient
}

var mu sync.RWMutex
var TypeRegistry map[string]reflect.Type

// NewRpcProxy 创建一个新的 RpcProxy 实例
func NewRpcProxy() *RpcProxy {
	mu.Lock()
	TypeRegistry = make(map[string]reflect.Type)
	TypeRegistry["int"] = reflect.TypeOf(0)
	TypeRegistry["string"] = reflect.TypeOf("")
	TypeRegistry["float64"] = reflect.TypeOf(0.0)
	mu.Unlock()
	return &RpcProxy{
		client: new(protocol.HttpClient),
	}
}

func GetTypeByName(name string) (reflect.Type, bool) {
	mu.RLock()
	typ, ok := TypeRegistry[name]
	mu.RUnlock()
	return typ, ok
}

func RegisterType(name string, t reflect.Type) {
	mu.Lock()
	TypeRegistry[name] = t
	mu.Unlock()
}

// 服务发现
func callServicecDiscovery(interfaceName string, registerAddr string) ([]common.URL, error) {
	urlList, err := register.QueryServicesFromHTTP(interfaceName, registerAddr)
	if err != nil || len(urlList) == 0 {
		return nil, fmt.Errorf("在接口 %s,没有可以用的服务， 报错: %v", interfaceName, err)
	}
	return urlList, nil
}

// 服务调用
func callServiceSend(p *protocol.HttpClient, hostName string, port int, invocation common.Invocation) ([]byte, error) {
	resultStr, err := p.Send(hostName, port, invocation)
	if err != nil {
		return []byte{}, err
	}
	return resultStr, nil
}

// Invoke 进行 RPC 调用
func (p *RpcProxy) Invoke(interfaceName, methodName string, params []interface{}, config *common.ClientConfig) ([]interface{}, error) {
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

	//	设置超时时间
	timeout := 5 * time.Second

	//服务发现(包含服务发现异常)：
	registerAddr := fmt.Sprintf("%s:%d", config.IP, config.Port)
	var urlList []common.URL
	err := protocol.WithTimeout(func() error {
		var err error
		urlList, err = callServicecDiscovery(interfaceName, registerAddr)
		return err
	}, timeout)
	if err != nil {
		return nil, err
	}

	//负载均衡
	url := loadbalance.Random(urlList)

	// 服务调用
	resultStr, err := p.client.Send(url.HostName, url.Port, invocation)
	err = protocol.WithTimeout(func() error {
		var err error
		resultStr, err = callServiceSend(p.client, url.HostName, url.Port, invocation)
		return err
	}, timeout)
	var result common.Result
	err = json.Unmarshal(resultStr, &result)
	if err != nil {
		return nil, err
	}

	//将结果转换为对应类型
	converted, err := convertResult(result)
	if err != nil {
		return nil, err
	}
	return converted, nil
}

// 将结果转换为对应类型的函数
func convertResult(result common.Result) ([]interface{}, error) {
	// 创建一个空的interface{}切片，长度与Values相同
	converted := make([]interface{}, len(result.Values))

	for i, val := range result.Values {
		// 根据Types中的类型字符串创建相应类型的零值
		instance, ok := TypeRegistry[result.Types[i]]
		if !ok {
			return nil, errors.New("不支持这种类型！")
		}
		valJSON, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}

		convertedVal := reflect.New(instance).Elem()
		ptr := convertedVal.Addr().Interface() // 获取convertedVal的地址并转换为interface{}

		err = json.Unmarshal(valJSON, ptr)
		if err != nil {
			return nil, err
		}
		converted[i] = convertedVal.Interface()
	}

	return converted, nil
}
