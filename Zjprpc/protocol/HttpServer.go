package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"myRpc/Zjprpc/common"
	"net/http"
	"time"
)

type HttpServer struct {
	server            *http.Server
	HeartbeatInterval time.Duration //心跳间隔
	RegisterAddr      string        //注册中心地址
	Url               common.URL    //提供服务的url
	stopChan          chan struct{} //关闭进程的管道
}

func NewHttpServer(heartbeatInterval time.Duration, registerAddr string, url common.URL) HttpServer {
	return HttpServer{
		HeartbeatInterval: heartbeatInterval,
		RegisterAddr:      registerAddr,
		Url:               url,
		stopChan:          make(chan struct{})}
}

func (server *HttpServer) Start(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", address)

	//创建http服务器
	server.server = &http.Server{
		Addr:    address,
		Handler: nil,
	}

	//添加handler
	http.HandleFunc("/", new(DispatcherServlet).service)

	//发送心跳
	go server.startHeartbeat()

	//监听停止信号
	go func() {
		<-server.stopChan
		fmt.Println("收到停止信号，准备关闭HTTP服务...")
		if err := server.server.Shutdown(context.Background()); err != nil {
			fmt.Printf("关闭HTTP服务失败！: %s\n", err)
		}
	}()

	//开启监听
	if err := server.server.ListenAndServe(); err != nil {
		fmt.Printf("服务被关闭！: %s\n", err)
	}
}

func (server *HttpServer) startHeartbeat() {
	ticker := time.NewTicker(server.HeartbeatInterval)
	defer func() {
		ticker.Stop()
		if r := recover(); r != nil {
			fmt.Printf("发送心跳失败！请启动注册中心或者重启服务端！ %v\n", r)
			close(server.stopChan)
		}
	}()
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
	heartbeatData := common.HeartBeatData{
		url, time.Now(),
	}

	//序列化心跳数据
	requestBody, err := json.Marshal(heartbeatData)
	if err != nil {
		return errors.New("序列化心跳数据失败！")
	}

	//发送put请求
	req, err := http.NewRequest("PUT", registerAddr, bytes.NewReader(requestBody))
	if err != nil {
		return errors.New("设置心跳请求失败！")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送心跳失败，状态码为: %d", resp.StatusCode)
	}
	return nil
}
