package rpc

import (
	"context"
	"sync"
)

// headerKey 用于在 context 中存储 header 的 key
type headerKey string

// headerKeysKey 用于在 context 中存储所有 header keys 的 key
const headerKeysKey headerKey = "__header_keys__"

var (
	headerKeysMutex sync.RWMutex
	headerKeysMap   = make(map[context.Context][]string)
)

// GetRPCHeader 从上下文中获取指定的 header 值
//
// 参数:
//   - ctx: 上下文
//   - key: header 的键名
//
// 返回值:
//   - string: header 的值
//   - bool: 是否存在该 header
func GetRPCHeader(ctx context.Context, key string) (string, bool) {
	if ctx == nil {
		return "", false
	}
	value := ctx.Value(headerKey(key))
	if value == nil {
		return "", false
	}
	strValue, ok := value.(string)
	if !ok {
		return "", false
	}
	return strValue, true
}

// SetRPCHeader 在上下文中设置 header
//
// 参数:
//   - ctx: 原始上下文
//   - key: header 的键名
//   - value: header 的值
//
// 返回值:
//   - context.Context: 新的上下文，包含设置的 header
func SetRPCHeader(ctx context.Context, key, value string) context.Context {
	// 设置 header 值
	newCtx := context.WithValue(ctx, headerKey(key), value)

	// 更新 header keys
	headerKeysMutex.Lock()
	defer headerKeysMutex.Unlock()

	keys := headerKeysMap[ctx]
	found := false
	for _, k := range keys {
		if k == key {
			found = true
			break
		}
	}
	if !found {
		keys = append(keys, key)
	}
	headerKeysMap[newCtx] = keys
	delete(headerKeysMap, ctx) // 清理旧的 context

	return newCtx
}

// GetRPCHeaders 获取上下文中的所有 headers
//
// 参数:
//   - ctx: 上下文
//
// 返回值:
//   - map[string]string: 所有的 headers
func GetRPCHeaders(ctx context.Context) map[string]string {
	if ctx == nil {
		return make(map[string]string)
	}

	headers := make(map[string]string)

	headerKeysMutex.RLock()
	keys := headerKeysMap[ctx]
	headerKeysMutex.RUnlock()

	for _, key := range keys {
		if value, ok := GetRPCHeader(ctx, key); ok {
			headers[key] = value
		}
	}

	return headers
}

// SetRPCHeaders 在上下文中批量设置 headers
//
// 参数:
//   - ctx: 原始上下文
//   - headers: 要设置的 headers 键值对
//
// 返回值:
//   - context.Context: 新的上下文，包含设置的所有 headers
func SetRPCHeaders(ctx context.Context, headers map[string]string) context.Context {
	newCtx := ctx
	for key, value := range headers {
		newCtx = SetRPCHeader(newCtx, key, value)
	}
	return newCtx
}
