# ason-go

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.24+-00ADD8.svg)](https://go.dev)

High-performance Go support for [ASON](https://github.com/ason-lab/ason), a schema-driven format for compact structured data.

[中文文档](README_CN.md)

## Why ASON

ASON writes field names once and stores rows positionally:

```json
[
  {"id": 1, "name": "Alice", "active": true},
  {"id": 2, "name": "Bob", "active": false}
]
```

```text
[{id@int,name@str,active@bool}]:(1,Alice,true),(2,Bob,false)
```

That reduces repeated keys, payload size, and often parsing cost.

## Highlights

- Standard library only
- Current API uses `Encode` / `Decode`, not the older `Marshal` / `Unmarshal` names
- Text, pretty text, and binary formats
- Struct tags via `ason:"..."`, with `json` tag fallback
- No native `map` / dictionary field syntax; model key-value data as slices of entry structs
- Good fit for LLM payloads, internal services, logs, and fixtures

## Install

```bash
go get github.com/ason-lab/ason-go
```

## Quick Start

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

### Encode a slice

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

### Pretty and binary output

```go
pretty, _ := ason.EncodePretty(users)
prettyTyped, _ := ason.EncodePrettyTyped(users)
bin, _ := ason.EncodeBinary(users)

var decoded []User
_ = ason.DecodeBinary(bin, &decoded)
```

### Model key-value data with entry structs

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

Typed ASON output:

```text
{name@str,env@[{key@str,value@str}]}:(api,[(RUST_LOG,debug),(PORT,8080)])
```

## Current API

| Function | Purpose |
| --- | --- |
| `Encode` / `EncodeTyped` | Encode to text |
| `Decode` | Decode from text |
| `EncodePretty` / `EncodePrettyTyped` | Pretty text output |
| `EncodeBinary` | Encode to binary |
| `DecodeBinary` | Decode from binary |

## Run Examples

```bash
go test ./...
go run ./examples/basic
go run ./examples/complex
go run ./examples/bench
```

## Contributors

- [Athan](https://github.com/athxx)

## Benchmarks

Run:

```bash
go run ./examples/bench
```

The benchmark output now follows the same layout as the C and C++ versions:

```text
Serialize:   JSON    16.22ms | ASON    16.80ms (1x) | BIN    15.02ms (1.1x)
Deserialize: JSON   111.90ms | ASON    35.50ms (3.2x) | BIN    35.10ms (3.2x)
Size:        JSON   218737 B | ASON    84861 B (39%) | BIN    85282 B (39%)
```

## License

MIT
