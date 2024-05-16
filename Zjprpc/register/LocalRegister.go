package register

import "reflect"

type LocalRegister struct {
	localMap map[string]interface{}
}

var TypeRegistry map[string]reflect.Type
var localRegisterInstance *LocalRegister

func GetInstance() *LocalRegister {
	return localRegisterInstance
}

func InitLocalRegister() {
	localMap := map[string]interface{}{}
	TypeRegistry = make(map[string]reflect.Type)
	TypeRegistry["int"] = reflect.TypeOf(0)
	TypeRegistry["string"] = reflect.TypeOf("")
	TypeRegistry["float64"] = reflect.TypeOf(0.0)
	localRegisterInstance = &LocalRegister{localMap: localMap}
}

func GetTypeByName(name string) (reflect.Type, bool) {
	typ, ok := TypeRegistry[name]
	return typ, ok
}

func RegisterType(name string, t reflect.Type) {
	TypeRegistry[name] = t
}

func (l *LocalRegister) Regist(interfaceName string, version string, implClass interface{}) {
	l.localMap[interfaceName+version] = implClass
}

func (l *LocalRegister) Get(interfaceName, version string) interface{} {
	return l.localMap[interfaceName+version]
}
