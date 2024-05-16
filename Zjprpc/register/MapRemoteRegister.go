package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myRpc/Zjprpc/common"
	"net/http"
	"sync"
)

// HTTPRegisterServer 使用HTTP作为服务注册中心
type HTTPRegisterServer struct {
	mu       sync.RWMutex
	registry map[string][]common.URL
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
	s.registry[url.InterfaceName] = append(s.registry[url.InterfaceName], url)
	s.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// 处理客户端的服务发现
func (s *HTTPRegisterServer) handleQuery(w http.ResponseWriter, r *http.Request) {
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
	return http.ListenAndServe(addr, server)
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
