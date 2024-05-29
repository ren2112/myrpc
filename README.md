# myrpc
纯go语言手写rpc框架，实现服务端，客户端，注册中心（使用http请求）
## go八股考点
- 反射
- 方法集：对指针结构体的方法集但是对普通结构体不适用（当然普通调用是可以的，因为编译会**自动解引用**！）
## 项目结构
- myRpc
  - Consumer（客户端的 main 包代码）
  - Provider（服务端的接口具体实现）
  - Provider-Common（服务端的接口定义）
  - Registry（注册中心 main 包代码）
  - server（服务端 main 包代码）
  - Zjrpc（rpc 相关自定义协议）
    - common（定义通用函数和结构体）
    - loadbalance（负载均衡）
    - protocol（对服务端客户端的 rpc 进行实现）
    - proxy（实现服务端的代理）
    - register（实现注册中心具体代码）


## 难点
- 反射使用（服务端获得传参,根据方法名调用方法）
- go程编写（服务端发送心跳go程，注册中心检测心跳go程）
- 线程安全问题，注册中心的注册表需要使用读写锁
## 可以改进的地方
- http比tcp更加占用资源，发起请求消耗资源大，其实可以用tcp
