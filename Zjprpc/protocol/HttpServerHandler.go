package protocol

import (
	"encoding/json"
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
		inputs[i] = reflect.ValueOf(param)
	}

	// 调用方法并检查返回值
	results := method.Call(inputs)
	if len(results) == 0 {
		http.Error(resp, "Method call returned no results", http.StatusInternalServerError)
		return
	}

	result := results[0].Interface()

	// 返回结果
	resp.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(resp).Encode(result); err != nil {
		http.Error(resp, "Failed to encode response", http.StatusInternalServerError)
	}
}
