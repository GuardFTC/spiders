// Package client @Author:冯铁城 [17615007230@163.com] 2025-10-11 11:13:33
package client

import "time"

// MongoConfig mongodb配置
type MongoConfig struct {
	Uri            string        `json:"uri"`             // uri
	MaxPoolSize    uint64        `json:"max_pool_size"`   // 最大连接池大小
	MinPoolSize    uint64        `json:"min_pool_size"`   // 最小连接池大小
	ConnectTimeout time.Duration `json:"connect_timeout"` // 连接超时
	SocketTimeout  time.Duration `json:"socket_timeout"`  // Socket 超时
}
