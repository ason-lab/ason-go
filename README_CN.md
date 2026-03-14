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
[{id@int,name@str,active@bool}]:(1,Alice,true),(2,Bob,false)
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
    // {id@int,name@str,active@bool}:(1,Alice,true)

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

### 用 entry struct 表达键值集合

```go
type EnvEntry struct {
    Key   string `ason:"key"`
    Value string `ason:"value"`
}

type Config struct {
    Name string     `ason:"name"`
    Env  []EnvEntry `ason:"env"`
}
```

对应的带类型 ASON 文本：

```text
{name@str,env@[{key@str,value@str}]}:(api,[(RUST_LOG,debug),(PORT,8080)])
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

## Benchmarks

可通过下面命令运行：

```bash
go run ./examples/bench
```

输出格式与 C / C++ 版本保持一致，例如：

```text
Serialize:   JSON    16.22ms | ASON    16.80ms (1x) | BIN    15.02ms (1.1x)
Deserialize: JSON   111.90ms | ASON    35.50ms (3.2x) | BIN    35.10ms (3.2x)
Size:        JSON   218737 B | ASON    84861 B (39%) | BIN    85282 B (39%)
```

## 许可证

MIT
