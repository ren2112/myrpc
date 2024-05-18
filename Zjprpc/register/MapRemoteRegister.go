package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myRpc/Zjprpc/common"
	"net/http"
	"sync"
	"time"
)

// HTTPRegisterServer 使用HTTP作为服务注册中心
type HTTPRegisterServer struct {
	mu       sync.RWMutex
	registry map[string][]common.URL
	timeout  time.Duration
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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *HTTPRegisterServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var url common.URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	urls := make([]common.URL, len(s.registry[url.InterfaceName]))
	copy(urls, s.registry[url.InterfaceName])

	//urls := s.registry[url.InterfaceName]
	// 向新切片中添加新的 URL
	newURL := common.URL{
		InterfaceName: url.InterfaceName,
		HostName:      url.HostName,
		Port:          url.Port,
		LastHeartbeat: url.LastHeartbeat,
	}
	urls = append(urls, newURL)

	// 更新注册表中的切片
	s.registry[url.InterfaceName] = urls
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
func StartHTTPRegisterServer(addr string) error {
	server := NewHTTPRegisterServer()
	go server.checkHeartBeat()
	return http.ListenAndServe(addr, server)
}

// 定期检查心跳
func (s *HTTPRegisterServer) checkHeartBeat() {
	for {
		time.Sleep(s.timeout) //先等待timeout时间
		s.mu.Lock()
		for interfaceName, urls := range s.registry {
			for i := 0; i < len(urls); i++ {
				if time.Since(urls[i].LastHeartbeat) > s.timeout {
					//	移除失效服务
					urls = append(urls[:i], urls[i+1:]...)
					i--
				}
			}
			s.registry[interfaceName] = urls
		}
		s.mu.Unlock()
	}
}

// 服务中心注册，给服务端用
func RegisterServiceToHTTP(url common.URL, registerAddr string) error {
	body, err := json.Marshal(url)
	if err != nil {
		return err
	}
	resp, err := http.Post(registerAddr+"/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service, status code: %d", resp.StatusCode)
	}
	return nil
}

// 服务中心发现,给客户端用
func QueryServicesFromHTTP(interfaceName, registerAddr string) ([]common.URL, error) {
	resp, err := http.Get(registerAddr + "/query?interface=" + interfaceName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to query services, status code: %d", resp.StatusCode)
	}

	var urls []common.URL
	err = json.NewDecoder(resp.Body).Decode(&urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// 处理心跳请求逻辑
func (s *HTTPRegisterServer) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var heartbeatData common.HeartBeatData
	err := json.NewDecoder(r.Body).Decode(&heartbeatData)
	if err != nil {
		http.Error(w, "无法反序列化心跳数据", http.StatusBadRequest)
		return
	}
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
