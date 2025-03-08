package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRPCHeader(t *testing.T) {
	// 测试用例
	tests := []struct {
		name     string
		setup    func() context.Context
		key      string
		expected string
		hasValue bool
	}{
		{
			name: "empty context",
			setup: func() context.Context {
				return context.Background()
			},
			key:      "test-key",
			expected: "",
			hasValue: false,
		},
		{
			name: "context with value",
			setup: func() context.Context {
				ctx := context.Background()
				return SetRPCHeader(ctx, "test-key", "test-value")
			},
			key:      "test-key",
			expected: "test-value",
			hasValue: true,
		},
		{
			name: "context with non-existent key",
			setup: func() context.Context {
				ctx := context.Background()
				return SetRPCHeader(ctx, "other-key", "test-value")
			},
			key:      "test-key",
			expected: "",
			hasValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setup()
			value, ok := GetRPCHeader(ctx, tt.key)

			assert.Equal(t, tt.hasValue, ok, "hasValue mismatch")
			assert.Equal(t, tt.expected, value, "value mismatch")
		})
	}
}

func TestSetRPCHeader(t *testing.T) {
	// 基础上下文
	ctx := context.Background()

	// 测试用例
	tests := []struct {
		name     string
		key      string
		value    string
		expected string
	}{
		{
			name:     "set simple header",
			key:      "test-key",
			value:    "test-value",
			expected: "test-value",
		},
		{
			name:     "set empty value",
			key:      "empty-key",
			value:    "",
			expected: "",
		},
		{
			name:     "set special characters",
			key:      "special-key",
			value:    "!@#$%^&*()",
			expected: "!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置 header
			newCtx := SetRPCHeader(ctx, tt.key, tt.value)

			// 验证设置的值
			value, ok := GetRPCHeader(newCtx, tt.key)
			assert.True(t, ok, "header should be set")
			assert.Equal(t, tt.expected, value, "value mismatch")

			// 验证原始上下文未被修改
			originalValue, ok := GetRPCHeader(ctx, tt.key)
			assert.False(t, ok, "original context should not be modified")
			assert.Empty(t, originalValue, "original context should not have value")
		})
	}
}

func TestGetRPCHeaders(t *testing.T) {
	// 创建带有多个 header 的上下文
	ctx := context.Background()
	ctx = SetRPCHeader(ctx, "key1", "value1")
	ctx = SetRPCHeader(ctx, "key2", "value2")
	ctx = SetRPCHeader(ctx, "key3", "value3")

	// 获取所有 headers
	headers := GetRPCHeaders(ctx)

	// 验证结果
	assert.Equal(t, 3, len(headers), "should have 3 headers")
	assert.Equal(t, "value1", headers["key1"], "value1 mismatch")
	assert.Equal(t, "value2", headers["key2"], "value2 mismatch")
	assert.Equal(t, "value3", headers["key3"], "value3 mismatch")

	// 测试空上下文
	emptyHeaders := GetRPCHeaders(context.Background())
	assert.Equal(t, 0, len(emptyHeaders), "empty context should have no headers")
}

func TestSetRPCHeaders(t *testing.T) {
	// 准备测试数据
	headers := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// 设置 headers
	ctx := SetRPCHeaders(context.Background(), headers)

	// 验证所有 header 都被正确设置
	for key, expectedValue := range headers {
		value, ok := GetRPCHeader(ctx, key)
		assert.True(t, ok, "header should be set: %s", key)
		assert.Equal(t, expectedValue, value, "value mismatch for key: %s", key)
	}

	// 验证获取的所有 headers 与设置的一致
	retrievedHeaders := GetRPCHeaders(ctx)
	assert.Equal(t, headers, retrievedHeaders, "retrieved headers should match set headers")
}

func TestRPCHeadersOverwrite(t *testing.T) {
	// 初始上下文
	ctx := context.Background()
	ctx = SetRPCHeader(ctx, "key", "value1")

	// 覆盖值
	ctx = SetRPCHeader(ctx, "key", "value2")

	// 验证新值
	value, ok := GetRPCHeader(ctx, "key")
	assert.True(t, ok, "header should be set")
	assert.Equal(t, "value2", value, "value should be overwritten")
}
