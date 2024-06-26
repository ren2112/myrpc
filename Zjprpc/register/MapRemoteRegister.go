package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myRpc/Zjprpc/common"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPRegisterServer 使用HTTP作为服务注册中心
type HTTPRegisterServer struct {
	mu       sync.RWMutex            //接口url注册表的读写锁
	registry map[string][]common.URL //记录接口名称对应的url数组的接口url表
	timeout  time.Duration           //心跳超时时间
}

func NewHTTPRegisterServer() *HTTPRegisterServer {
	return &HTTPRegisterServer{
		registry: make(map[string][]common.URL),
	}
}

func (s *HTTPRegisterServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleRegister(w, r)
	case http.MethodGet:
		s.handleQuery(w, r)
	case http.MethodPut:
		s.handleHeartbeat(w, r)
	default:
		http.Error(w, "不支持这种请求方式", http.StatusMethodNotAllowed)
	}
}

func (s *HTTPRegisterServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	//获得请求的ip地址
	ipPort := strings.Split(r.RemoteAddr, ":")
	ip := ipPort[0]

	//获得服务端发起注册的url结构体
	var url common.URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//更新url的ip，因为有可能ip在服务端眼里是本机0.0.0.0，但是在注册中心却并不知道服务端的本机ip
	url.HostName = ip

	s.mu.Lock()
	//注册
	s.registry[url.InterfaceName] = append(s.registry[url.InterfaceName], url)
	s.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// 处理客户端的服务发现
func (s *HTTPRegisterServer) handleQuery(w http.ResponseWriter, r *http.Request) {
	//获取url部分query的interfaceName
	interfaceName := r.URL.Query().Get("interface")
	s.mu.RLock()
	if urls, ok := s.registry[interfaceName]; ok {
		json.NewEncoder(w).Encode(urls)
	} else {
		http.NotFound(w, r)
	}
	s.mu.RUnlock()
}

// 启动服务中心
func StartHTTPRegisterServer(addr string, timeout time.Duration) error {
	server := NewHTTPRegisterServer()
	server.timeout = timeout

	// 启动定期检查心跳的 goroutine
	go server.checkHeartBeat()

	// 监听并处理 HTTP 请求，一直阻塞直到服务器关闭
	err := http.ListenAndServe(addr, server)
	if err != nil {
		return err
	}

	// 服务器关闭
	fmt.Println("服务器已关闭")

	return nil
}

// 定期检查心跳
func (s *HTTPRegisterServer) checkHeartBeat() {
	for {
		time.Sleep(s.timeout) // 先等待 timeout 时间
		s.mu.Lock()
		for interfaceName, urls := range s.registry {
			updatedURLs := make([]common.URL, 0) // 创建一个新的切片，用于保存未失效的服务
			for i := 0; i < len(urls); i++ {
				if time.Since(urls[i].LastHeartbeat) <= s.timeout {
					// 如果服务未失效，则将其添加到新切片中
					updatedURLs = append(updatedURLs, urls[i])
				}
			}
			// 更新注册表中的切片
			s.registry[interfaceName] = updatedURLs
		}
		s.mu.Unlock()
	}
}

// 服务中心注册，给服务端用
func RegisterServiceToHTTP(url common.URL, registerAddr string) error {
	//序列化
	body, err := json.Marshal(url)
	if err != nil {
		return err
	}
	//发送请求给注册中心
	resp, err := http.Post(registerAddr, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("注册服务失败，状态码: %d", resp.StatusCode)
	}
	return nil
}

// 服务中心发现,给客户端用
func QueryServicesFromHTTP(interfaceName, registerAddr string) ([]common.URL, error) {
	//发起http请求
	resp, err := http.Get(registerAddr + "/query?interface=" + interfaceName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务发现调用失败，状态码: %d", resp.StatusCode)
	}

	//对响应url列表反序列化
	var urls []common.URL
	err = json.NewDecoder(resp.Body).Decode(&urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// 处理心跳请求逻辑
func (s *HTTPRegisterServer) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	//获得真实ip地址
	ipPort := strings.Split(r.RemoteAddr, ":")
	ip := ipPort[0]

	var heartbeatData common.HeartBeatData
	err := json.NewDecoder(r.Body).Decode(&heartbeatData)
	if err != nil {
		http.Error(w, "无法反序列化心跳数据", http.StatusBadRequest)
		return
	}
	heartbeatData.URL.HostName = ip

	//	更新服务中心的接口名对应url的时间戳
	s.mu.Lock()
	urls, ok := s.registry[heartbeatData.URL.InterfaceName]
	if ok {
		for i := range urls {
			if urls[i].HostName == heartbeatData.URL.HostName &&
				urls[i].Port == heartbeatData.URL.Port {
				//更新时间
				urls[i].LastHeartbeat = heartbeatData.HeartbeatTime
				break
			}
		}
		s.registry[heartbeatData.URL.InterfaceName] = urls
	}
	s.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}
