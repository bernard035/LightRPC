# LightRPC

A lightweight RPC framework based on protocol buffer

LightRPC是一个基于`protocol buffer`序列化协议的轻量级RPC框架，它基于`net/rpc`构建，支持`gzip`、`zlib`、`snappy`压缩格式。

请求头部结构如下：

| CompressType |    Method     |   ID   | RequestLen | Checksum |
| :----------: | :-----------: | :----: | :--------: | :------: |
|    uint16    | uint64+string | uint64 |   uint64   |  uint32  |

响应头部结构如下：

| CompressType |   ID   |     Error     | ResponseLen | Checksum |
| :----------: | :----: | :-----------: | :---------: | :------: |
|    uint16    | uint64 | uint64+string |   uint64    |  uint32  |
