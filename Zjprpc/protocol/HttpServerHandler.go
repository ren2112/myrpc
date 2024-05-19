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
		http.Error(resp, "读取请求体失败！", http.StatusBadRequest)
		return
	}

	//反序列化
	err = json.Unmarshal([]byte(requestBody), &invocation)
	if err != nil {
		http.Error(resp, "反序列化传参失败!", http.StatusBadRequest)
		return
	}

	// 获取实现类
	interfaceName := invocation.InterfaceName
	localRegister := register.GetInstance()
	classImpl := localRegister.Get(interfaceName, "1.0")
	if classImpl == nil {
		http.Error(resp, "没有找到这个接口！", http.StatusBadRequest)
		return
	}
	classImplValue := reflect.ValueOf(classImpl)

	// 查找方法
	method := classImplValue.MethodByName(invocation.MethodName)
	if !method.IsValid() {
		http.Error(resp, "没有找到对应服务！", http.StatusBadRequest)
		return
	}

	// 构造参数列表
	inputs := make([]reflect.Value, len(invocation.Parameters))
	for i, param := range invocation.Parameters {
		// 获取参数类型
		paramType, ok := register.GetTypeByName(invocation.ParameterTypes[i])
		if !ok {
			http.Error(resp, "输入类型不支持！", http.StatusInternalServerError)
			return
		}
		inputIJSON, err := json.Marshal(param)
		if err != nil {
			http.Error(resp, "输入类型不支持！", http.StatusInternalServerError)
			return
		}

		inputIval := reflect.New(paramType).Elem()
		ptr := inputIval.Addr().Interface()
		err = json.Unmarshal(inputIJSON, ptr)
		if err != nil {
			http.Error(resp, "输入类型不支持！", http.StatusInternalServerError)
			return
		}

		// 根据参数类型创建 reflect.Value 对象
		inputs[i] = inputIval
	}

	// 调用方法并检查返回值
	results := method.Call(inputs)
	if len(results) == 0 {
		http.Error(resp, "这个方法没有返回值！", http.StatusInternalServerError)
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
	//响应报文写入数据
	if err = json.NewEncoder(resp).Encode(result); err != nil {
		http.Error(resp, "响应结果序列化失败！", http.StatusInternalServerError)
	}
}
