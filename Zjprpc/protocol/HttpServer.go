package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"myRpc/Zjprpc/common"
	"net/http"
	"time"
)

type HttpServer struct {
	HeartbeatInterval time.Duration //心跳间隔
	RegisterAddr      string        //注册中心地址
	Url               common.URL    //提供服务的url
}

func NewHttpServer(heartbeatInterval time.Duration, registerAddr string, url common.URL) HttpServer {
	return HttpServer{heartbeatInterval, registerAddr, url}
}

func (server *HttpServer) Start(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", address)

	http.HandleFunc("/", new(DispatcherServlet).service)
	go server.startHeartbeat()
	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}

func (server *HttpServer) startHeartbeat() {
	ticker := time.NewTicker(server.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("我发送心跳了")
			//	发送心跳
			err := SendHeartbeat(server.Url, server.RegisterAddr)
			if err != nil {
				fmt.Printf("发送心跳失败: %s\n", err)
			}
		}
	}
}

func SendHeartbeat(url common.URL, registerAddr string) error {
	heartbeatUrl := registerAddr
	heartbeatData := common.HeartBeatData{
		url, time.Now(),
	}
	//序列化
	requestBody, err := json.Marshal(heartbeatData)
	if err != nil {
		return err
	}
	//发送post请求
	req, err := http.NewRequest("PUT", heartbeatUrl, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送心跳失败，状态码为: %d", resp.StatusCode)
	}
	return nil
}
