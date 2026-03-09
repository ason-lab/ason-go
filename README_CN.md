# ason-go

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24+-00ADD8.svg)](https://go.dev)

面向 [ASON](https://github.com/ason-lab/ason) 的高性能 Go 实现。ASON 是一种适合紧凑结构化数据的 Schema 驱动格式。

[English](README.md)

## 为什么用 ASON

ASON 只写一次字段名，后续数据按位置存储：

```json
[
  {"id": 1, "name": "Alice", "active": true},
  {"id": 2, "name": "Bob", "active": false}
]
```

```text
[{id:int,name:str,active:bool}]:(1,Alice,true),(2,Bob,false)
```

这能减少重复键名、减小体积，并且通常降低解析成本。

## 特性

- 仅依赖 Go 标准库
- 当前 API 是 `Encode` / `Decode`，不再是旧文档里的 `Marshal` / `Unmarshal`
- 同时支持文本、格式化文本和二进制格式
- 通过 `ason:"..."` struct tag 定义字段名，并回退支持 `json` tag
- 适合 LLM 载荷、内部服务、日志和测试数据

## 安装

```bash
go get github.com/ason-lab/ason-go
```

## 快速开始

```go
package main

import (
    "fmt"
    ason "github.com/ason-lab/ason-go"
)

type User struct {
    ID     int64  `ason:"id"`
    Name   string `ason:"name"`
    Active bool   `ason:"active"`
}

func main() {
    user := User{ID: 1, Name: "Alice", Active: true}

    text, _ := ason.Encode(&user)
    fmt.Println(string(text))
    // {id,name,active}:(1,Alice,true)

    typed, _ := ason.EncodeTyped(&user)
    fmt.Println(string(typed))
    // {id:int,name:str,active:bool}:(1,Alice,true)

    var decoded User
    _ = ason.Decode(text, &decoded)
}
```

### 编码切片

```go
users := []User{
    {ID: 1, Name: "Alice", Active: true},
    {ID: 2, Name: "Bob", Active: false},
}

text, _ := ason.Encode(users)
typed, _ := ason.EncodeTyped(users)

var decoded []User
_ = ason.Decode(text, &decoded)
```

### 格式化文本和二进制

```go
pretty, _ := ason.EncodePretty(users)
prettyTyped, _ := ason.EncodePrettyTyped(users)
bin, _ := ason.EncodeBinary(users)

var decoded []User
_ = ason.DecodeBinary(bin, &decoded)
```

## 当前 API

| 函数 | 作用 |
| --- | --- |
| `Encode` / `EncodeTyped` | 编码为文本 |
| `Decode` | 从文本解码 |
| `EncodePretty` / `EncodePrettyTyped` | 生成更易读的文本 |
| `EncodeBinary` | 编码为二进制 |
| `DecodeBinary` | 从二进制解码 |

## 运行示例

```bash
go test ./...
go run ./examples/basic
go run ./examples/complex
go run ./examples/bench
```

## Contributors

- [Athan](https://github.com/athxx)

## Benchmark Snapshot

在当前机器上通过下面命令实测：

```bash
go run ./examples/bench
```

关键结果：

- 扁平 1,000 条记录：ASON 序列化 `15.71ms`，JSON `76.92ms`；反序列化 ASON `77.32ms`，JSON `380.56ms`
- 扁平 10,000 条记录：ASON 序列化 `75.82ms`，JSON `173.69ms`；反序列化 ASON `283.21ms`，JSON `770.39ms`
- 深层 100 条数据：ASON 序列化 `425.39ms`，JSON `699.25ms`；反序列化 ASON `1136.91ms`，JSON `3001.79ms`
- 1,000 条记录吞吐总结：ASON 文本序列化比 JSON 快 `3.37x`，反序列化快 `4.40x`
- 1,000 条扁平记录二进制总结：BIN 序列化 `22.8ms`，JSON `32.1ms`；BIN 反序列化 `15.2ms`，JSON `186.2ms`；二进制大小 `74,450 B`

## 许可证

MIT
