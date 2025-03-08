package resty_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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
			req := GetHttpRequest(tt.timeout)
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
