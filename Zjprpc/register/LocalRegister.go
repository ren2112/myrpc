package register

type LocalRegister struct {
	localMap map[string]interface{}
}

var localRegisterInstance *LocalRegister

func GetInstance() *LocalRegister {
	return localRegisterInstance
}

func InitLocalRegister() {
	localMap := map[string]interface{}{}
	localRegisterInstance = &LocalRegister{localMap: localMap}
}

func (l *LocalRegister) Regist(interfaceName string, version string, implClass interface{}) {
	l.localMap[interfaceName+version] = implClass
}

func (l *LocalRegister) Get(interfaceName, version string) interface{} {
	return l.localMap[interfaceName+version]
}
