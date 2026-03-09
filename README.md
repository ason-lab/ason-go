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
[{id:int,name:str,active:bool}]:(1,Alice,true),(2,Bob,false)
```

That reduces repeated keys, payload size, and often parsing cost.

## Highlights

- Standard library only
- Current API uses `Encode` / `Decode`, not the older `Marshal` / `Unmarshal` names
- Text, pretty text, and binary formats
- Struct tags via `ason:"..."`, with `json` tag fallback
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
    // {id:int,name:str,active:bool}:(1,Alice,true)

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

## Benchmark Snapshot

Measured on this machine with:

```bash
go run ./examples/bench
```

Headline numbers:

- Flat 1,000-record dataset: ASON serialize `15.71ms` vs JSON `76.92ms`, deserialize `77.32ms` vs JSON `380.56ms`
- Flat 10,000-record dataset: ASON serialize `75.82ms` vs JSON `173.69ms`, deserialize `283.21ms` vs JSON `770.39ms`
- Deep 100-record dataset: ASON serialize `425.39ms` vs JSON `699.25ms`, deserialize `1136.91ms` vs JSON `3001.79ms`
- Throughput summary on 1,000 records: ASON text was `3.37x` faster than JSON for serialize and `4.40x` faster for deserialize
- Binary summary on 1,000 flat records: BIN serialize `22.8ms` vs JSON `32.1ms`, BIN deserialize `15.2ms` vs JSON `186.2ms`, with binary size `74,450 B`

## License

MIT
