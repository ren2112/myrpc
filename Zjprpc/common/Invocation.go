package common

type Invocation struct {
	InterfaceName string
	MethodName    string
	//记录传参类型
	ParameterTypes []string
	//	记录调用传参
	Parameters []interface{}
}

func NewInvocation(interfaceName string, methodName string, parameterTypes []string, parameters []interface{}) Invocation {
	return Invocation{InterfaceName: interfaceName, MethodName: methodName, ParameterTypes: parameterTypes, Parameters: parameters}
}
