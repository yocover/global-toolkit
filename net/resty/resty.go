package resty

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

// DefaultTimeout 默认的 HTTP 请求超时时间（秒）
const DefaultTimeout = 10

// HTTP 请求相关的常量定义
const (
	// ContentType 请求头的 Content-Type 字段名
	ContentType = "Content-Type"
	// ContentTypeJson JSON 格式的 Content-Type 值
	ContentTypeJson = "application/json"
	// ContentTypeForm 表单格式的 Content-Type 值
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeMultipartForm 多部分表单格式的 Content-Type 值
	ContentTypeMultipartForm = "multipart/form-data"
)

// GetRequest 创建一个基础的 HTTP 请求客户端
//
// 参数:
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - *resty.Request: resty 请求对象，可用于进一步配置和发送请求
//
// 示例:
//
//	req := GetRequest(30)
//	resp, err := req.Get("https://api.example.com")
func GetRequest(timout int64) *resty.Request {
	client := resty.New()
	client.SetTimeout(time.Duration(timout) * time.Second)
	return client.R()
}

// GetHttpsRequest 创建一个支持 HTTPS 的 HTTP 请求客户端，会跳过 TLS 证书验证
//
// 参数:
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - *resty.Request: 配置了 TLS 跳过验证和请求追踪的 resty 请求对象
//
// 示例:
//
//	req := GetHttpsRequest(30)
//	resp, err := req.Get("https://api.example.com")
func GetHttpsRequest(timout int64) *resty.Request {
	// 创建新的resty客户端
	client := resty.New()
	// 配置TLS ，跳过证书验证
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// 设置超时时间
	client.SetTimeout(time.Duration(timout) * time.Second)
	// 创建请求对象并启用追踪
	return client.R().EnableTrace()
}

// Get 发送一个简单的 HTTP GET 请求
//
// 参数:
//   - url: 目标请求地址
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	resp, err := Get("https://api.example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(resp))
func Get(url string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).Get(url)
	if err != nil {
		return nil, err
	}
	resp = request.Body()
	return resp, nil
}

// GetWithHeaders 发送带自定义请求头的 HTTP GET 请求
//
// 参数:
//   - url: 目标请求地址
//   - header: 自定义的 HTTP 请求头，key-value 键值对形式
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	headers := map[string]string{
//	    "Content-Type": "application/json",
//	    "Authorization": "Bearer token123"
//	}
//	resp, err := GetWithHeaders("https://api.example.com", headers)
func GetWithHeaders(url string, header map[string]string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).SetHeaders(header).Get(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsGetWithHeaders 发送带自定义请求头的 HTTPS GET 请求，会跳过 TLS 证书验证
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//   - header: 自定义的 HTTP 请求头，key-value 键值对形式
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 在生产环境中使用时需要注意安全风险
//
// 示例:
//
//	headers := map[string]string{
//	    "Content-Type": "application/json",
//	    "Authorization": "Bearer token123"
//	}
//	resp, err := HttpsGetWithHeaders("https://api.example.com", headers)
func HttpsGetWithHeaders(url string, header map[string]string) (resp []byte, err error) {
	request, err := GetHttpsRequest(DefaultTimeout).SetHeaders(header).Get(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsGet 发送一个简单的 HTTPS GET 请求，会跳过 TLS 证书验证
//
// 这是一个简化版的 HTTPS GET 请求函数，适用于不需要自定义请求头的场景。
// 内部使用 GetHttpsRequest 实现，继承了其跳过证书验证的特性。
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 使用默认超时时间 DefaultTimeout（10秒）
//   - 在生产环境中使用时需要注意安全风险
//   - 如果需要自定义请求头，请使用 HttpsGetWithHeaders 函数
//
// 示例:
//
//	resp, err := HttpsGet("https://api.example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(resp))
func HttpsGet(url string) (resp []byte, err error) {
	request, err := GetHttpsRequest(DefaultTimeout).Get(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// GetWithEntity 发送 GET 请求并将响应解析为指定的实体对象
//
// 这个函数会自动将响应的 JSON 内容解析到提供的实体对象中。
//
// 参数:
//   - url: 目标请求地址
//   - entity: 用于存储响应数据的目标对象指针
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - error: JSON 解析错误或请求错误，如果成功则为 nil
//
// 示例:
//
//	var user User
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	err := GetWithEntity("https://api.example.com/user", &user, headers, 30)
func GetWithEntity(url string, entity interface{}, header map[string]string, timeout int64) error {
	request, err := GetRequest(timeout).SetHeaders(header).Get(url)
	if err != nil {
		return err
	}
	resp := request.Body()

	err = json.Unmarshal(resp, &entity)
	if err != nil {
		zap.L().Error("Json Transform Error", zap.Error(err))
		return err
	}
	return err
}

// GetWithTimeOut 发送带超时设置的 HTTP GET 请求
//
// 参数:
//   - url: 目标请求地址
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	resp, err := GetWithTimeOut("https://api.example.com", headers, 30)
func GetWithTimeOut(url string, header map[string]string, timeout int64) (resp []byte, err error) {
	request, err := GetRequest(timeout).SetHeaders(header).Get(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsGetWithTimeOut 发送带超时设置的 HTTPS GET 请求，会跳过 TLS 证书验证
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 在生产环境中使用时需要注意安全风险
//
// 示例:
//
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	resp, err := HttpsGetWithTimeOut("https://api.example.com", headers, 30)
func HttpsGetWithTimeOut(url string, header map[string]string, timeout int64) (resp []byte, err error) {
	request, err := GetHttpsRequest(timeout).SetHeaders(header).Get(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// Post 发送一个简单的 HTTP POST 请求
//
// 参数:
//   - url: 目标请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	resp, err := Post("https://api.example.com", body, headers)
func Post(url string, body interface{}, header map[string]string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// PostWithTimeOut 发送带超时设置的 HTTP POST 请求
//
// 参数:
//   - url: 目标请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	resp, err := PostWithTimeOut("https://api.example.com", body, headers, 30)
func PostWithTimeOut(url string, body interface{}, header map[string]string, timeout int64) (resp []byte, err error) {
	request, err := GetRequest(timeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsPost 发送一个 HTTPS POST 请求，会跳过 TLS 证书验证
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 在生产环境中使用时需要注意安全风险
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	resp, err := HttpsPost("https://api.example.com", body, headers)
func HttpsPost(url string, body interface{}, header map[string]string) (resp []byte, err error) {
	request, err := GetHttpsRequest(DefaultTimeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsPostWithTimeOut 发送带超时设置的 HTTPS POST 请求，会跳过 TLS 证书验证
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 在生产环境中使用时需要注意安全风险
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	resp, err := HttpsPostWithTimeOut("https://api.example.com", body, headers, 30)
func HttpsPostWithTimeOut(url string, body interface{}, header map[string]string, timeout int64) (resp []byte, err error) {
	request, err := GetHttpsRequest(timeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// PostWithEntity 发送 POST 请求并将响应解析为指定的实体对象
//
// 这个函数会自动将响应的 JSON 内容解析到提供的实体对象中。
//
// 参数:
//   - url: 目标请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//   - entity: 用于存储响应数据的目标对象指针
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - error: JSON 解析错误或请求错误，如果成功则为 nil
//
// 示例:
//
//	var response Response
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	err := PostWithEntity("https://api.example.com", body, headers, &response, 30)
func PostWithEntity(url string, body interface{}, header map[string]string, entity interface{}, timeout int64) error {
	request, err := GetRequest(timeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return err
	}
	resp := request.Body()

	err = json.Unmarshal(resp, &entity)
	if err != nil {
		zap.L().Error("Json Transform Error", zap.Error(err))
		return err
	}
	return err
}

// Json 发送 JSON 格式的 POST 请求
//
// 自动设置 Content-Type 为 application/json。
//
// 参数:
//   - url: 目标请求地址
//   - body: 请求体内容，将被序列化为 JSON
//   - header: 自定义的 HTTP 请求头
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	resp, err := Json("https://api.example.com", body, headers)
func Json(url string, body interface{}, header map[string]string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).SetHeaders(header).SetHeader(ContentType, ContentTypeJson).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// Form 发送 x-www-form-urlencoded 格式的 POST 请求
//
// 自动设置 Content-Type 为 application/x-www-form-urlencoded。
//
// 参数:
//   - url: 目标请求地址
//   - FormData: 表单数据，key-value 键值对形式
//   - header: 自定义的 HTTP 请求头
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	formData := map[string]string{"name": "test", "age": "18"}
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	resp, err := Form("https://api.example.com", formData, headers)
func Form(url string, FormData map[string]string, header map[string]string) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).SetHeaders(header).SetHeader(ContentType, ContentTypeForm).SetFormData(FormData).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// File 发送带文件的 POST 请求
//
// 自动设置 Content-Type 为 application/x-www-form-urlencoded。
//
// 参数:
//   - url: 目标请求地址
//   - FormData: 表单数据，key-value 键值对形式
//   - header: 自定义的 HTTP 请求头
//   - param: 文件参数名
//   - fileName: 文件名
//   - reader: 文件内容读取器
//
// 返回值:
//   - resp: 响应体的字节数组
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 示例:
//
//	file, _ := os.Open("test.txt")
//	formData := map[string]string{"name": "test"}
//	headers := map[string]string{"Authorization": "Bearer token123"}
//	resp, err := File("https://api.example.com", formData, headers, "file", "test.txt", file)
func File(url string, FormData map[string]string, header map[string]string, param, fileName string, reader io.Reader) (resp []byte, err error) {
	request, err := GetRequest(DefaultTimeout).SetHeaders(header).SetHeader(ContentType, ContentTypeForm).SetFormData(FormData).SetFileReader(param, fileName, reader).Post(url)
	if err != nil {
		return
	}
	resp = request.Body()
	return
}

// HttpsPostWithTimeOutResHeader 发送带超时设置的 HTTPS POST 请求，并返回响应头
//
// 参数:
//   - url: 目标 HTTPS 请求地址
//   - body: 请求体内容，可以是任意类型
//   - header: 自定义的 HTTP 请求头
//   - timeout: 请求超时时间（秒）
//
// 返回值:
//   - resp: 响应体的字节数组
//   - resHeader: 响应头
//   - err: 请求过程中的错误信息，如果请求成功则为 nil
//
// 注意:
//   - 此方法会跳过 TLS 证书验证，主要用于自签名证书或测试环境
//   - 在生产环境中使用时需要注意安全风险
//
// 示例:
//
//	body := map[string]interface{}{"name": "test"}
//	headers := map[string]string{"Content-Type": "application/json"}
//	resp, resHeaders, err := HttpsPostWithTimeOutResHeader("https://api.example.com", body, headers, 30)
func HttpsPostWithTimeOutResHeader(url string, body interface{}, header map[string]string, timeout int64) (resp []byte, resHeader http.Header, err error) {
	res, err := GetHttpsRequest(timeout).SetHeaders(header).SetBody(body).Post(url)
	if err != nil {
		return
	}
	resp = res.Body()
	resHeader = res.Header()
	return
}
