package resty_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/yocover/global-toolkit/net/resty"
)

// 编写测试方法
func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	tests := []struct {
		name    string
		timeout int64
	}{
		{
			name:    "normal timeout 30 seconds",
			timeout: 30,
		},
		{
			name:    "minimum timeout 1 second",
			timeout: 1,
		},
		{
			name:    "large timeout 300 seconds",
			timeout: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetRequest(tt.timeout)
			assert.NotNil(t, req, "GetRequest() should not return nil")
			assert.NotNil(t, req.Header, "Request headers should not be nil")

			resp, err := req.Get(ts.URL)
			assert.NoError(t, err, "Request should not fail")
			assert.Equal(t, http.StatusOK, resp.StatusCode(), "Status code should be 200")
			assert.Equal(t, `{"status":"ok"}`, string(resp.Body()), "Response body should match")
		})
	}
}

func TestGetHttpRequest(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	tests := []struct {
		name    string
		timeout int64
	}{
		{
			name:    "normal timeout 30 seconds",
			timeout: 30,
		},
		{
			name:    "minimum timeout 1 second",
			timeout: 1,
		},
		{
			name:    "large timeout 300 seconds",
			timeout: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := GetHttpsRequest(tt.timeout)
			assert.NotNil(t, req, "GetHttpRequest() should not return nil")
			assert.NotNil(t, req.Header, "Request headers should not be nil")

			resp, err := req.Get(ts.URL)
			assert.NoError(t, err, "Request should not fail")
			assert.Equal(t, http.StatusOK, resp.StatusCode(), "Status code should be 200")
			assert.Equal(t, `{"status":"ok"}`, string(resp.Body()), "Response body should match")
		})
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))

	defer ts.Close()

	resp, err := Get(ts.URL + "/test")
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestGetWithHeaders(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和路径
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		// 验证关键请求头
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-token", r.Header.Get("Authorization"))

		// 返回响应
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	// 设置请求头
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "test-token",
	}

	// 发送请求并验证结果
	resp, err := GetWithHeaders(ts.URL+"/test", headers)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestHttpsGetWithHeaders(t *testing.T) {
	// 创建 HTTPS 测试服务器
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和路径
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		// 验证关键请求头
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-token", r.Header.Get("Authorization"))

		// 返回响应
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	// 设置请求头
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "test-token",
	}

	// 发送请求并验证结果
	resp, err := HttpsGetWithHeaders(ts.URL+"/test", headers)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

// 用于测试的示例结构体
type TestResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func TestGetWithEntity(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok","data":"test"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	var response TestResponse
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	err := GetWithEntity(ts.URL+"/test", &response, headers, 30)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "test", response.Data)
}

func TestGetWithTimeOut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := GetWithTimeOut(ts.URL+"/test", headers, 30)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestHttpsGetWithTimeOut(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := HttpsGetWithTimeOut(ts.URL+"/test", headers, 30)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, `{"name":"test"}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := map[string]string{
		"name": "test",
	}
	resp, err := Post(ts.URL+"/test", body, headers)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, `{"name":"test"}`, string(body))

		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{}
	body := map[string]string{
		"name": "test",
	}
	resp, err := Json(ts.URL+"/test", body, headers)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestForm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := r.ParseForm()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "test", r.Form.Get("name"))
		assert.Equal(t, "18", r.Form.Get("age"))

		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{}
	formData := map[string]string{
		"name": "test",
		"age":  "18",
	}
	resp, err := Form(ts.URL+"/test", formData, headers)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			t.Fatal(err)
		}

		// 验证普通表单字段
		assert.Equal(t, "test", r.FormValue("name"))

		// 验证文件
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		assert.Equal(t, "test.txt", header.Filename)
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "test content", string(content))

		w.WriteHeader(http.StatusOK)
		_, err = io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{}
	formData := map[string]string{
		"name": "test",
	}
	reader := strings.NewReader("test content")
	resp, err := File(ts.URL+"/test", formData, headers, "file", "test.txt", reader)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
}

func TestHttpsPostWithTimeOutResHeader(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.Header().Set("X-Test-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	body := map[string]string{
		"name": "test",
	}
	resp, resHeader, err := HttpsPostWithTimeOutResHeader(ts.URL+"/test", body, headers, 30)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"ok"}`), resp)
	assert.Equal(t, "test-value", resHeader.Get("X-Test-Header"))
}
