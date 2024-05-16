package protocol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myRpc/Zjprpc/common"
	"myRpc/Zjprpc/register"
	"net/http"
	"reflect"
)

type HttpServerHandler struct {
}

func (h *HttpServerHandler) Handler(resp http.ResponseWriter, req *http.Request) {
	//	处理请求:接口，方法，方法参数
	var invocation common.Invocation
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Failed to read request body", http.StatusBadRequest)
		return
	}

	//反序列化
	err = json.Unmarshal([]byte(requestBody), &invocation)
	if err != nil {
		http.Error(resp, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// 获取实现类
	interfaceName := invocation.InterfaceName
	localRegister := register.GetInstance()
	classImpl := localRegister.Get(interfaceName, "1.0")
	if classImpl == nil {
		http.Error(resp, "Class implementation not found", http.StatusBadRequest)
		return
	}
	classImplValue := reflect.ValueOf(classImpl)

	// 查找方法
	method := classImplValue.MethodByName(invocation.MethodName)
	if !method.IsValid() {
		http.Error(resp, "Method not found", http.StatusBadRequest)
		return
	}

	// 构造参数列表
	inputs := make([]reflect.Value, len(invocation.Parameters))
	for i, param := range invocation.Parameters {
		// 获取参数类型
		paramType, ok := register.GetTypeByName(invocation.ParameterTypes[i])
		if !ok {
			fmt.Println("输入类型不支持！")
			return
		}
		// 根据参数类型创建 reflect.Value 对象
		inputs[i] = reflect.ValueOf(param).Convert(paramType)
	}

	// 调用方法并检查返回值
	results := method.Call(inputs)
	if len(results) == 0 {
		http.Error(resp, "Method call returned no results", http.StatusInternalServerError)
		return
	}

	// 组装返回值及类型信息
	var returnValues []interface{}
	var returnTypes []string
	for _, result := range results {
		returnValues = append(returnValues, result.Interface())
		returnTypes = append(returnTypes, reflect.TypeOf(result.Interface()).String())
	}

	// 返回结果
	resp.Header().Set("Content-Type", "application/json")
	result := common.Result{
		Values: returnValues,
		Types:  returnTypes,
	}
	if err = json.NewEncoder(resp).Encode(result); err != nil {
		http.Error(resp, "Failed to encode response", http.StatusInternalServerError)
	}
}
