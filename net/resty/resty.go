package resty

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// Default timeout 10s
const DefaultTimeout = 10

const (
	ContentType              = "Content-Type"
	ContentTypeJson          = "application/json"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm = "multipart/form-data"
)

// 提供一个统一的方式来创建 HTTP 请求客户端
func GetRequest(timout int64) *resty.Request {
	client := resty.New()
	client.SetTimeout(time.Duration(timout) * time.Second)
	return client.R()
}

func GetHttpRequest(timout int64) *resty.Request {
	// 创建新的resty客户端
	client := resty.New()
	// 配置TLS ，跳过证书验证
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// 设置超时时间
	client.SetTimeout(time.Duration(timout) * time.Second)
	// 创建请求对象并启用追踪
	return client.R().EnableTrace()
}

func Get(url string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).Get(url)
	if err != nil {
		return nil, err
	}
	resp = request.Body()
	return resp, nil
}
