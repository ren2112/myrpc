package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"myRpc/Zjprpc/common"
	"net/http"
)

type HttpClient struct {
}

func (h *HttpClient) Send(hostname string, port int, invocation common.Invocation) (string, error) {
	// 将 invocation 转换为 JSON 格式的字节切片
	requestBody, err := json.Marshal(invocation)
	if err != nil {
		return "", err
	}

	// 构建请求 URL
	url := fmt.Sprintf("http://%s:%d/", hostname, port)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}
