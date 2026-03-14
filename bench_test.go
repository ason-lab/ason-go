package ason

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// cd /home/X/dev/ason/ason-go
// go test -run '^$' -bench 'BenchmarkCompareFlat8' -benchmem
// go test -run '^$' -bench 'BenchmarkCompareAllTypes' -benchmem
// go test -run '^$' -bench 'BenchmarkCompareDeep' -benchmem
// go test -run '^$' -bench 'BenchmarkRoundtripSingle|BenchmarkAnnotatedComparison|BenchmarkBinaryComparison' -benchmem

var (
	benchFlatUser           = User{ID: 42, Name: "Alice", Active: true}
	benchFlatUsers100       = makeFlatUsers(100)
	benchFlatUsers1000      = makeFlatUsers(1000)
	benchFlatStructUntyped  = []byte("{id,name,active}:(42,Alice,true)")
	benchFlatStructTyped    = []byte("{id@int,name@str,active@bool}:(42,Alice,true)")
	benchFlatStructJSON     = []byte(`{"id":42,"name":"Alice","active":true}`)
	benchFlatVec100         = []byte(makeFlatUserVecText(100, false))
	benchFlatVec100Typed    = []byte(makeFlatUserVecText(100, true))
	benchFlatVec100JSON     = []byte(makeFlatUserVecJSON(100))
	benchFlatVec1000        = []byte(makeFlatUserVecText(1000, false))
	benchFlatVec1000Typed   = []byte(makeFlatUserVecText(1000, true))
	benchFlatVec1000JSON    = []byte(makeFlatUserVecJSON(1000))
	benchFlatStructJSONOut  = mustJSONMarshal(benchFlatUser)
	benchFlatStructASONOut  = mustAsonEncode(benchFlatUser, false)
	benchFlatStructATypOut  = mustAsonEncode(benchFlatUser, true)
	benchFlatVec100JSONOut  = mustJSONMarshal(benchFlatUsers100)
	benchFlatVec100ASONOut  = mustAsonEncode(benchFlatUsers100, false)
	benchFlatVec100ATypOut  = mustAsonEncode(benchFlatUsers100, true)
	benchFlatVec1000JSONOut = mustJSONMarshal(benchFlatUsers1000)
	benchFlatVec1000ASONOut = mustAsonEncode(benchFlatUsers1000, false)
	benchFlatVec1000ATypOut = mustAsonEncode(benchFlatUsers1000, true)
)

func makeFlatUsers(n int) []User {
	users := make([]User, n)
	for i := 0; i < n; i++ {
		users[i] = User{
			ID:     int64(i + 1),
			Name:   "user" + strconv.Itoa(i),
			Active: i&1 == 0,
		}
	}
	return users
}

func mustJSONMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustAsonEncode(v any, typed bool) []byte {
	var (
		b   []byte
		err error
	)
	if typed {
		b, err = EncodeTyped(v)
	} else {
		b, err = Encode(v)
	}
	if err != nil {
		panic(err)
	}
	return b
}

func makeFlatUserVecText(n int, typed bool) string {
	var b strings.Builder
	b.Grow(n * 24)
	if typed {
		b.WriteString("[{id@int,name@str,active@bool}]:")
	} else {
		b.WriteString("[{id,name,active}]:")
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteByte(',')
		b.WriteString("user")
		b.WriteString(strconv.Itoa(i))
		b.WriteByte(',')
		if i&1 == 0 {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteByte(')')
	}
	return b.String()
}

func makeFlatUserVecJSON(n int) string {
	var b strings.Builder
	b.Grow(n * 40)
	b.WriteByte('[')
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(`{"id":`)
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(`,"name":"user`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`","active":`)
		if i&1 == 0 {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteByte('}')
	}
	b.WriteByte(']')
	return b.String()
}

func BenchmarkDecodeFlatStruct(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		var u User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructJSON)))
		for i := 0; i < b.N; i++ {
			u = User{}
			if err := json.Unmarshal(benchFlatStructJSON, &u); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		var u User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructUntyped)))
		for i := 0; i < b.N; i++ {
			u = User{}
			if err := Decode(benchFlatStructUntyped, &u); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		var u User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructTyped)))
		for i := 0; i < b.N; i++ {
			u = User{}
			if err := Decode(benchFlatStructTyped, &u); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecodeFlatVec100(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100JSON)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := json.Unmarshal(benchFlatVec100JSON, &users); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := Decode(benchFlatVec100, &users); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100Typed)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := Decode(benchFlatVec100Typed, &users); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecodeFlatVec1000(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000JSON)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := json.Unmarshal(benchFlatVec1000JSON, &users); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := Decode(benchFlatVec1000, &users); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		var users []User
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000Typed)))
		for i := 0; i < b.N; i++ {
			users = nil
			if err := Decode(benchFlatVec1000Typed, &users); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEncodeFlatStruct(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructJSONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(benchFlatUser); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructASONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(benchFlatUser); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatStructATypOut)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeTyped(benchFlatUser); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEncodeFlatVec100(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100JSONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(benchFlatUsers100); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100ASONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(benchFlatUsers100); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec100ATypOut)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeTyped(benchFlatUsers100); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEncodeFlatVec1000(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000JSONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(benchFlatUsers1000); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Untyped", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000ASONOut)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(benchFlatUsers1000); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Typed", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(benchFlatVec1000ATypOut)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeTyped(benchFlatUsers1000); err != nil {
				b.Fatal(err)
			}
		}
	})
}

type BenchUser8 struct {
	ID     int64   `ason:"id" json:"id"`
	Name   string  `ason:"name" json:"name"`
	Email  string  `ason:"email" json:"email"`
	Age    int64   `ason:"age" json:"age"`
	Score  float64 `ason:"score" json:"score"`
	Active bool    `ason:"active" json:"active"`
	Role   string  `ason:"role" json:"role"`
	City   string  `ason:"city" json:"city"`
}

type BenchAllTypes struct {
	B       bool     `ason:"b" json:"b"`
	I8v     int8     `ason:"i8v" json:"i8v"`
	I16v    int16    `ason:"i16v" json:"i16v"`
	I32v    int32    `ason:"i32v" json:"i32v"`
	I64v    int64    `ason:"i64v" json:"i64v"`
	U8v     uint8    `ason:"u8v" json:"u8v"`
	U16v    uint16   `ason:"u16v" json:"u16v"`
	U32v    uint32   `ason:"u32v" json:"u32v"`
	U64v    uint64   `ason:"u64v" json:"u64v"`
	F32v    float32  `ason:"f32v" json:"f32v"`
	F64v    float64  `ason:"f64v" json:"f64v"`
	S       string   `ason:"s" json:"s"`
	OptSome *int64   `ason:"opt_some" json:"opt_some"`
	OptNone *int64   `ason:"opt_none" json:"opt_none"`
	VecInt  []int64  `ason:"vec_int" json:"vec_int"`
	VecStr  []string `ason:"vec_str" json:"vec_str"`
}

type BenchTask struct {
	ID       int64   `ason:"id" json:"id"`
	Title    string  `ason:"title" json:"title"`
	Priority int64   `ason:"priority" json:"priority"`
	Done     bool    `ason:"done" json:"done"`
	Hours    float64 `ason:"hours" json:"hours"`
}

type BenchProject struct {
	Name   string      `ason:"name" json:"name"`
	Budget float64     `ason:"budget" json:"budget"`
	Active bool        `ason:"active" json:"active"`
	Tasks  []BenchTask `ason:"tasks" json:"tasks"`
}

type BenchTeam struct {
	Name     string         `ason:"name" json:"name"`
	Lead     string         `ason:"lead" json:"lead"`
	Size     int64          `ason:"size" json:"size"`
	Projects []BenchProject `ason:"projects" json:"projects"`
}

type BenchDivision struct {
	Name      string      `ason:"name" json:"name"`
	Location  string      `ason:"location" json:"location"`
	Headcount int64       `ason:"headcount" json:"headcount"`
	Teams     []BenchTeam `ason:"teams" json:"teams"`
}

type BenchCompany struct {
	Name      string          `ason:"name" json:"name"`
	Founded   int64           `ason:"founded" json:"founded"`
	RevenueM  float64         `ason:"revenue_m" json:"revenue_m"`
	Public    bool            `ason:"public" json:"public"`
	Divisions []BenchDivision `ason:"divisions" json:"divisions"`
	Tags      []string        `ason:"tags" json:"tags"`
}

func benchI64Ptr(v int64) *int64 { return &v }

func generateBenchUsers(n int) []BenchUser8 {
	names := []string{"Alice", "Bob", "Carol", "David", "Eve", "Frank", "Grace", "Hank"}
	roles := []string{"engineer", "designer", "manager", "analyst"}
	cities := []string{"NYC", "LA", "Chicago", "Houston", "Phoenix"}
	users := make([]BenchUser8, n)
	for i := 0; i < n; i++ {
		users[i] = BenchUser8{
			ID:     int64(i),
			Name:   names[i%len(names)],
			Email:  strings.ToLower(names[i%len(names)]) + "@example.com",
			Age:    int64(25 + i%40),
			Score:  50.0 + float64(i%50) + 0.5,
			Active: i%3 != 0,
			Role:   roles[i%len(roles)],
			City:   cities[i%len(cities)],
		}
	}
	return users
}

func generateBenchAllTypes(n int) []BenchAllTypes {
	items := make([]BenchAllTypes, n)
	for i := 0; i < n; i++ {
		var optSome *int64
		if i%2 == 0 {
			optSome = benchI64Ptr(int64(i))
		}
		items[i] = BenchAllTypes{
			B: i%2 == 0, I8v: int8(i % 256), I16v: -int16(i),
			I32v: int32(i) * 1000, I64v: int64(i) * 100000,
			U8v: uint8(i % 256), U16v: uint16(i % 65536),
			U32v: uint32(i) * 7919, U64v: uint64(i) * 1000000007,
			F32v: float32(i) * 1.5, F64v: float64(i)*0.25 + 0.5,
			S: fmt.Sprintf("item_%d", i), OptSome: optSome, OptNone: nil,
			VecInt: []int64{int64(i), int64(i + 1), int64(i + 2)},
			VecStr: []string{fmt.Sprintf("tag%d", i%5), fmt.Sprintf("cat%d", i%3)},
		}
	}
	return items
}

func generateBenchCompanies(n int) []BenchCompany {
	locs := []string{"NYC", "London", "Tokyo", "Berlin"}
	leads := []string{"Alice", "Bob", "Carol", "David"}
	companies := make([]BenchCompany, n)
	for i := 0; i < n; i++ {
		divisions := make([]BenchDivision, 2)
		for d := 0; d < 2; d++ {
			teams := make([]BenchTeam, 2)
			for t := 0; t < 2; t++ {
				projects := make([]BenchProject, 3)
				for p := 0; p < 3; p++ {
					tasks := make([]BenchTask, 4)
					for tk := 0; tk < 4; tk++ {
						tasks[tk] = BenchTask{
							ID: int64(i*100 + d*10 + t*5 + tk), Title: fmt.Sprintf("Task_%d", tk),
							Priority: int64(tk%3 + 1), Done: tk%2 == 0, Hours: 2.0 + float64(tk)*1.5,
						}
					}
					projects[p] = BenchProject{
						Name: fmt.Sprintf("Proj_%d_%d", t, p), Budget: 100.0 + float64(p)*50.5,
						Active: p%2 == 0, Tasks: tasks,
					}
				}
				teams[t] = BenchTeam{
					Name: fmt.Sprintf("Team_%d_%d_%d", i, d, t), Lead: leads[t%4],
					Size: int64(5 + t*2), Projects: projects,
				}
			}
			divisions[d] = BenchDivision{
				Name: fmt.Sprintf("Div_%d_%d", i, d), Location: locs[d%4],
				Headcount: int64(50 + d*20), Teams: teams,
			}
		}
		companies[i] = BenchCompany{
			Name: fmt.Sprintf("Corp_%d", i), Founded: int64(1990 + i%35),
			RevenueM: 10.0 + float64(i)*5.5, Public: i%2 == 0, Divisions: divisions,
			Tags: []string{"enterprise", "tech", fmt.Sprintf("sector_%d", i%5)},
		}
	}
	return companies
}

func benchmarkCompare[T any](b *testing.B, jsonData []byte, asonUntyped []byte, asonTyped []byte, encodeValue T) {
	b.Run("EncodeJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(jsonData)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(encodeValue); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("EncodeASON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(asonUntyped)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(encodeValue); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("EncodeASONTyped", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(asonTyped)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeTyped(encodeValue); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCompareFlat8(b *testing.B) {
	for _, n := range []int{100, 500, 1000, 5000, 10000} {
		users := generateBenchUsers(n)
		jsonData := mustJSONMarshal(users)
		asonUntyped := mustAsonEncode(users, false)
		asonTyped := mustAsonEncode(users, true)

		b.Run(fmt.Sprintf("Count=%d", n), func(b *testing.B) {
			benchmarkCompare(b, jsonData, asonUntyped, asonTyped, users)

			b.Run("DecodeJSON", func(b *testing.B) {
				var out []BenchUser8
				b.ReportAllocs()
				b.SetBytes(int64(len(jsonData)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := json.Unmarshal(jsonData, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASON", func(b *testing.B) {
				var out []BenchUser8
				b.ReportAllocs()
				b.SetBytes(int64(len(asonUntyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonUntyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASONTyped", func(b *testing.B) {
				var out []BenchUser8
				b.ReportAllocs()
				b.SetBytes(int64(len(asonTyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonTyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkCompareAllTypes(b *testing.B) {
	for _, n := range []int{100, 500} {
		items := generateBenchAllTypes(n)
		jsonData := mustJSONMarshal(items)
		asonUntyped := mustAsonEncode(items, false)
		asonTyped := mustAsonEncode(items, true)

		b.Run(fmt.Sprintf("Count=%d", n), func(b *testing.B) {
			benchmarkCompare(b, jsonData, asonUntyped, asonTyped, items)

			b.Run("DecodeJSON", func(b *testing.B) {
				var out []BenchAllTypes
				b.ReportAllocs()
				b.SetBytes(int64(len(jsonData)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := json.Unmarshal(jsonData, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASON", func(b *testing.B) {
				var out []BenchAllTypes
				b.ReportAllocs()
				b.SetBytes(int64(len(asonUntyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonUntyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASONTyped", func(b *testing.B) {
				var out []BenchAllTypes
				b.ReportAllocs()
				b.SetBytes(int64(len(asonTyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonTyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkCompareDeep(b *testing.B) {
	for _, n := range []int{10, 50, 100} {
		companies := generateBenchCompanies(n)
		jsonData := mustJSONMarshal(companies)
		asonUntyped := mustAsonEncode(companies, false)
		asonTyped := mustAsonEncode(companies, true)

		b.Run(fmt.Sprintf("Count=%d", n), func(b *testing.B) {
			benchmarkCompare(b, jsonData, asonUntyped, asonTyped, companies)

			b.Run("DecodeJSON", func(b *testing.B) {
				var out []BenchCompany
				b.ReportAllocs()
				b.SetBytes(int64(len(jsonData)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := json.Unmarshal(jsonData, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASON", func(b *testing.B) {
				var out []BenchCompany
				b.ReportAllocs()
				b.SetBytes(int64(len(asonUntyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonUntyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("DecodeASONTyped", func(b *testing.B) {
				var out []BenchCompany
				b.ReportAllocs()
				b.SetBytes(int64(len(asonTyped)))
				for i := 0; i < b.N; i++ {
					out = nil
					if err := Decode(asonTyped, &out); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkRoundtripSingle(b *testing.B) {
	flat := generateBenchUsers(1)[0]
	deep := generateBenchCompanies(1)[0]

	b.Run("FlatJSON", func(b *testing.B) {
		var out BenchUser8
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s, err := json.Marshal(&flat)
			if err != nil {
				b.Fatal(err)
			}
			if err := json.Unmarshal(s, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatASON", func(b *testing.B) {
		var out BenchUser8
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s, err := Encode(&flat)
			if err != nil {
				b.Fatal(err)
			}
			if err := Decode(s, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepJSON", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s, err := json.Marshal(&deep)
			if err != nil {
				b.Fatal(err)
			}
			if err := json.Unmarshal(s, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepASON", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s, err := Encode(&deep)
			if err != nil {
				b.Fatal(err)
			}
			if err := Decode(s, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAnnotatedComparison(b *testing.B) {
	users := generateBenchUsers(1000)
	untypedFlat := mustAsonEncode(users, false)
	typedFlat := mustAsonEncode(users, true)
	deep := generateBenchCompanies(1)[0]
	untypedDeep := mustAsonEncode(&deep, false)
	typedDeep := mustAsonEncode(&deep, true)

	b.Run("DecodeFlat1000Untyped", func(b *testing.B) {
		var out []BenchUser8
		b.ReportAllocs()
		b.SetBytes(int64(len(untypedFlat)))
		for i := 0; i < b.N; i++ {
			out = nil
			if err := Decode(untypedFlat, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DecodeFlat1000Typed", func(b *testing.B) {
		var out []BenchUser8
		b.ReportAllocs()
		b.SetBytes(int64(len(typedFlat)))
		for i := 0; i < b.N; i++ {
			out = nil
			if err := Decode(typedFlat, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DecodeDeepUntyped", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		b.SetBytes(int64(len(untypedDeep)))
		for i := 0; i < b.N; i++ {
			out = BenchCompany{}
			if err := Decode(untypedDeep, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DecodeDeepTyped", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		b.SetBytes(int64(len(typedDeep)))
		for i := 0; i < b.N; i++ {
			out = BenchCompany{}
			if err := Decode(typedDeep, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkBinaryComparison(b *testing.B) {
	flat := generateBenchUsers(1)[0]
	deep := generateBenchCompanies(1)[0]
	flatJSON := mustJSONMarshal(&flat)
	flatASON := mustAsonEncode(&flat, false)
	flatBIN, _ := EncodeBinary(&flat)
	deepJSON := mustJSONMarshal(&deep)
	deepASON := mustAsonEncode(&deep, false)
	deepBIN, _ := EncodeBinary(&deep)

	b.Run("FlatEncodeJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(flatJSON)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(&flat); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatEncodeASON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(flatASON)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(&flat); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatEncodeBIN", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(flatBIN)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeBinary(&flat); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatDecodeJSON", func(b *testing.B) {
		var out BenchUser8
		b.ReportAllocs()
		b.SetBytes(int64(len(flatJSON)))
		for i := 0; i < b.N; i++ {
			out = BenchUser8{}
			if err := json.Unmarshal(flatJSON, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatDecodeASON", func(b *testing.B) {
		var out BenchUser8
		b.ReportAllocs()
		b.SetBytes(int64(len(flatASON)))
		for i := 0; i < b.N; i++ {
			out = BenchUser8{}
			if err := Decode(flatASON, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("FlatDecodeBIN", func(b *testing.B) {
		var out BenchUser8
		b.ReportAllocs()
		b.SetBytes(int64(len(flatBIN)))
		for i := 0; i < b.N; i++ {
			out = BenchUser8{}
			if err := DecodeBinary(flatBIN, &out); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DeepEncodeJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(deepJSON)))
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(&deep); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepEncodeASON", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(deepASON)))
		for i := 0; i < b.N; i++ {
			if _, err := Encode(&deep); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepEncodeBIN", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(deepBIN)))
		for i := 0; i < b.N; i++ {
			if _, err := EncodeBinary(&deep); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepDecodeJSON", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		b.SetBytes(int64(len(deepJSON)))
		for i := 0; i < b.N; i++ {
			out = BenchCompany{}
			if err := json.Unmarshal(deepJSON, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepDecodeASON", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		b.SetBytes(int64(len(deepASON)))
		for i := 0; i < b.N; i++ {
			out = BenchCompany{}
			if err := Decode(deepASON, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("DeepDecodeBIN", func(b *testing.B) {
		var out BenchCompany
		b.ReportAllocs()
		b.SetBytes(int64(len(deepBIN)))
		for i := 0; i < b.N; i++ {
			out = BenchCompany{}
			if err := DecodeBinary(deepBIN, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}
