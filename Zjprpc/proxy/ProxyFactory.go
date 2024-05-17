package proxy

import (
	"encoding/json"
	"errors"
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

var TypeRegistry map[string]reflect.Type

// NewRpcProxy 创建一个新的 RpcProxy 实例
func NewRpcProxy() *RpcProxy {
	TypeRegistry = make(map[string]reflect.Type)
	TypeRegistry["int"] = reflect.TypeOf(0)
	TypeRegistry["string"] = reflect.TypeOf("")
	TypeRegistry["float64"] = reflect.TypeOf(0.0)
	return &RpcProxy{
		client: new(protocol.HttpClient),
	}
}

func GetTypeByName(name string) (reflect.Type, bool) {
	typ, ok := TypeRegistry[name]
	return typ, ok
}

func RegisterType(name string, t reflect.Type) {
	TypeRegistry[name] = t
}

// Invoke 进行 RPC 调用
func (p *RpcProxy) Invoke(interfaceName, methodName string, params []interface{}) ([]interface{}, error) {
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
		return nil, fmt.Errorf("在接口 %s,没有可以用的服务， 报错: %v", interfaceName, err)
	}

	//负载均衡
	url := loadbalance.Random(urlList)

	// 服务调用
	resultStr, err := p.client.Send(url.HostName, url.Port, invocation)
	var result common.Result
	err = json.Unmarshal(resultStr, &result)
	if err != nil {
		return nil, err
	}

	converted, err := convertResult(result)
	if err != nil {
		return nil, err
	}
	return converted, nil
}

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
