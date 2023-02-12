# LightRPC

A lightweight RPC framework based on protocol buffer

LightRPC是一个基于`protocol buffer`序列化协议的轻量级RPC框架，它基于`net/rpc`构建，支持`gzip`、`zlib`、`snappy`压缩格式。

## 实现细节

请求头部结构如下：

| CompressType |    Method     |   ID   | RequestLen | Checksum |
| :----------: | :-----------: | :----: | :--------: | :------: |
|    uint16    | uint64+string | uint64 |   uint64   |  uint32  |

响应头部结构如下：

| CompressType |   ID   |     Error     | ResponseLen | Checksum |
| :----------: | :----: | :-----------: | :---------: | :------: |
|    uint16    | uint64 | uint64+string |   uint64    |  uint32  |

为了实现Header重用，为Header构建了缓冲池：

```go
var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{New: func() interface{} {
		return &RequestHeader{}
	}}
	ResponsePool = sync.Pool{New: func() interface{} {
		return &ResponseHeader{}
	}}
}
```


在使用时get出来，生命周期结束后放回缓冲池，并且在put之前进行重置：

```go
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
```

