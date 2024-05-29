package register

import (
	"reflect"
	"sync"
)

// 服务注册结构体
type LocalRegister struct {
	localMap map[string]interface{}
}

var TypeRegistry map[string]reflect.Type //数据类型注册表
var localRegisterInstance *LocalRegister //服务注册结构体实例
var muRegist sync.RWMutex                //读写锁给服务注册表
var muType sync.RWMutex                  //读写锁给类型注册表

// 获得服务注册实例，为了给服务端调用注册服务用
func GetInstance() *LocalRegister {
	return localRegisterInstance
}

func InitLocalRegister() {
	localMap := map[string]interface{}{}
	TypeRegistry = make(map[string]reflect.Type)
	muType.Lock()
	TypeRegistry["int"] = reflect.TypeOf(0)
	TypeRegistry["string"] = reflect.TypeOf("")
	TypeRegistry["float64"] = reflect.TypeOf(0.0)
	muType.Unlock()
	localRegisterInstance = &LocalRegister{localMap: localMap}
}

func GetTypeByName(name string) (reflect.Type, bool) {
	muType.RLock()
	typ, ok := TypeRegistry[name]
	muType.RUnlock()
	return typ, ok
}

func RegisterType(name string, t reflect.Type) {
	muType.Lock()
	TypeRegistry[name] = t
	muType.Unlock()
}

func (l *LocalRegister) Regist(interfaceName string, implClass interface{}) {
	muRegist.Lock()
	l.localMap[interfaceName] = implClass
	muRegist.Unlock()
}

func (l *LocalRegister) Get(interfaceName string) interface{} {
	muRegist.RLock()
	ret := l.localMap[interfaceName]
	muRegist.RUnlock()
	return ret
}
