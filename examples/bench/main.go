package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	asun "github.com/asunLab/asun-go"
)

type User struct {
	ID     int64   `asun:"id" json:"id"`
	Name   string  `asun:"name" json:"name"`
	Email  string  `asun:"email" json:"email"`
	Age    int64   `asun:"age" json:"age"`
	Score  float64 `asun:"score" json:"score"`
	Active bool    `asun:"active" json:"active"`
	Role   string  `asun:"role" json:"role"`
	City   string  `asun:"city" json:"city"`
}

type AllTypes struct {
	B       bool     `asun:"b" json:"b"`
	I8v     int8     `asun:"i8v" json:"i8v"`
	I16v    int16    `asun:"i16v" json:"i16v"`
	I32v    int32    `asun:"i32v" json:"i32v"`
	I64v    int64    `asun:"i64v" json:"i64v"`
	U8v     uint8    `asun:"u8v" json:"u8v"`
	U16v    uint16   `asun:"u16v" json:"u16v"`
	U32v    uint32   `asun:"u32v" json:"u32v"`
	U64v    uint64   `asun:"u64v" json:"u64v"`
	F32v    float32  `asun:"f32v" json:"f32v"`
	F64v    float64  `asun:"f64v" json:"f64v"`
	S       string   `asun:"s" json:"s"`
	OptSome *int64   `asun:"opt_some" json:"opt_some"`
	OptNone *int64   `asun:"opt_none" json:"opt_none"`
	VecInt  []int64  `asun:"vec_int" json:"vec_int"`
	VecStr  []string `asun:"vec_str" json:"vec_str"`
}

type Task struct {
	ID       int64   `asun:"id" json:"id"`
	Title    string  `asun:"title" json:"title"`
	Priority int64   `asun:"priority" json:"priority"`
	Done     bool    `asun:"done" json:"done"`
	Hours    float64 `asun:"hours" json:"hours"`
}

type Project struct {
	Name   string  `asun:"name" json:"name"`
	Budget float64 `asun:"budget" json:"budget"`
	Active bool    `asun:"active" json:"active"`
	Tasks  []Task  `asun:"tasks" json:"tasks"`
}

type Team struct {
	Name     string    `asun:"name" json:"name"`
	Lead     string    `asun:"lead" json:"lead"`
	Size     int64     `asun:"size" json:"size"`
	Projects []Project `asun:"projects" json:"projects"`
}

type Division struct {
	Name      string `asun:"name" json:"name"`
	Location  string `asun:"location" json:"location"`
	Headcount int64  `asun:"headcount" json:"headcount"`
	Teams     []Team `asun:"teams" json:"teams"`
}

type Company struct {
	Name      string     `asun:"name" json:"name"`
	Founded   int64      `asun:"founded" json:"founded"`
	RevenueM  float64    `asun:"revenue_m" json:"revenue_m"`
	Public    bool       `asun:"public" json:"public"`
	Divisions []Division `asun:"divisions" json:"divisions"`
	Tags      []string   `asun:"tags" json:"tags"`
}

type benchResult struct {
	name      string
	jsonSerMs float64
	asunSerMs float64
	binSerMs  float64
	jsonDeMs  float64
	asunDeMs  float64
	binDeMs   float64
	jsonBytes int
	asunBytes int
	binBytes  int
}

func i64ptr(v int64) *int64 { return &v }

func generateUsers(n int) []User {
	names := []string{"Alice", "Bob", "Carol", "David", "Eve", "Frank", "Grace", "Hank"}
	roles := []string{"engineer", "designer", "manager", "analyst"}
	cities := []string{"NYC", "LA", "Chicago", "Houston", "Phoenix"}
	users := make([]User, n)
	for i := 0; i < n; i++ {
		users[i] = User{
			ID:     int64(i),
			Name:   names[i%len(names)],
			Email:  fmt.Sprintf("%s@example.com", strings.ToLower(names[i%len(names)])),
			Age:    int64(25 + i%40),
			Score:  50.0 + float64(i%50) + 0.5,
			Active: i%3 != 0,
			Role:   roles[i%len(roles)],
			City:   cities[i%len(cities)],
		}
	}
	return users
}

func generateAllTypes(n int) []AllTypes {
	items := make([]AllTypes, n)
	for i := 0; i < n; i++ {
		var optSome *int64
		if i%2 == 0 {
			optSome = i64ptr(int64(i))
		}
		items[i] = AllTypes{
			B:       i%2 == 0,
			I8v:     int8(i % 128),
			I16v:    -int16(i),
			I32v:    int32(i) * 1000,
			I64v:    int64(i) * 100000,
			U8v:     uint8(i % 255),
			U16v:    uint16(i % 65535),
			U32v:    uint32(i) * 7919,
			U64v:    uint64(i) * 1000000007,
			F32v:    float32(i) * 1.5,
			F64v:    float64(i)*0.25 + 0.5,
			S:       fmt.Sprintf("item_%d", i),
			OptSome: optSome,
			OptNone: nil,
			VecInt:  []int64{int64(i), int64(i + 1), int64(i + 2)},
			VecStr:  []string{fmt.Sprintf("tag%d", i%5), fmt.Sprintf("cat%d", i%3)},
		}
	}
	return items
}

func generateCompanies(n int) []Company {
	locs := []string{"NYC", "London", "Tokyo", "Berlin"}
	leads := []string{"Alice", "Bob", "Carol", "David"}
	companies := make([]Company, n)
	for i := 0; i < n; i++ {
		divisions := make([]Division, 2)
		for d := 0; d < 2; d++ {
			teams := make([]Team, 2)
			for t := 0; t < 2; t++ {
				projects := make([]Project, 3)
				for p := 0; p < 3; p++ {
					tasks := make([]Task, 4)
					for tk := 0; tk < 4; tk++ {
						tasks[tk] = Task{
							ID:       int64(i*100 + d*10 + t*5 + tk),
							Title:    fmt.Sprintf("Task_%d", tk),
							Priority: int64(tk%3 + 1),
							Done:     tk%2 == 0,
							Hours:    2.0 + float64(tk)*1.5,
						}
					}
					projects[p] = Project{
						Name:   fmt.Sprintf("Proj_%d_%d", t, p),
						Budget: 100.0 + float64(p)*50.5,
						Active: p%2 == 0,
						Tasks:  tasks,
					}
				}
				teams[t] = Team{
					Name:     fmt.Sprintf("Team_%d_%d_%d", i, d, t),
					Lead:     leads[t%len(leads)],
					Size:     int64(5 + t*2),
					Projects: projects,
				}
			}
			divisions[d] = Division{
				Name:      fmt.Sprintf("Div_%d_%d", i, d),
				Location:  locs[d%len(locs)],
				Headcount: int64(50 + d*20),
				Teams:     teams,
			}
		}
		companies[i] = Company{
			Name:      fmt.Sprintf("Corp_%d", i),
			Founded:   int64(1990 + i%35),
			RevenueM:  10.0 + float64(i)*5.5,
			Public:    i%2 == 0,
			Divisions: divisions,
			Tags:      []string{"enterprise", "tech", fmt.Sprintf("sector_%d", i%5)},
		}
	}
	return companies
}

func formatRatio(base, target float64) string {
	if target <= 0 {
		return "infx"
	}
	s := fmt.Sprintf("%.1f", base/target)
	s = strings.TrimSuffix(strings.TrimSuffix(s, "0"), ".")
	return s + "x"
}

func formatPercent(part, whole int) string {
	if whole <= 0 {
		return "0%"
	}
	s := fmt.Sprintf("%.1f", float64(part)*100.0/float64(whole))
	s = strings.TrimSuffix(strings.TrimSuffix(s, "0"), ".")
	return s + "%"
}

func mustAsunEncode(v any, typed bool) []byte {
	var (
		b   []byte
		err error
	)
	if typed {
		b, err = asun.EncodeTyped(v)
	} else {
		b, err = asun.Encode(v)
	}
	if err != nil {
		panic(err)
	}
	return b
}

func printSection(title string, width int) {
	line := strings.Repeat("─", width-2)
	fmt.Printf("┌%s┐\n", line)
	fmt.Printf("│ %-*s │\n", width-4, title)
	fmt.Printf("└%s┘\n", line)
}

func (r benchResult) print() {
	fmt.Printf("  %s\n", r.name)
	fmt.Printf("    Serialize:   JSON %.2fms/%dB | ASUN %.2fms(%s)/%dB(%s) | BIN %.2fm(%s)/%dB(%s)\n",
		r.jsonSerMs, r.jsonBytes,
		r.asunSerMs, formatRatio(r.jsonSerMs, r.asunSerMs), r.asunBytes, formatPercent(r.asunBytes, r.jsonBytes),
		r.binSerMs, formatRatio(r.jsonSerMs, r.binSerMs), r.binBytes, formatPercent(r.binBytes, r.jsonBytes))
	fmt.Printf("    Deserialize: JSON %8.2fms | ASUN %8.2fms(%s) | BIN %8.2fms(%s)\n",
		r.jsonDeMs, r.asunDeMs, formatRatio(r.jsonDeMs, r.asunDeMs), r.binDeMs, formatRatio(r.jsonDeMs, r.binDeMs))
}

func benchFlat(count, iterations int) benchResult {
	users := generateUsers(count)

	var jsonData []byte
	start := time.Now()
	for i := 0; i < iterations; i++ {
		jsonData, _ = json.Marshal(users)
	}
	jsonSer := time.Since(start)

	var asunData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		asunData, _ = asun.Encode(users)
	}
	asunSer := time.Since(start)

	var binData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		binData, _ = asun.EncodeBinary(users)
	}
	binSer := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []User
		_ = json.Unmarshal(jsonData, &out)
	}
	jsonDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []User
		_ = asun.Decode(asunData, &out)
	}
	asunDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []User
		_ = asun.DecodeBinary(binData, &out)
	}
	binDe := time.Since(start)

	var decoded []User
	if err := asun.Decode(asunData, &decoded); err != nil || len(decoded) != count {
		panic("flat text roundtrip failed")
	}
	if err := asun.DecodeBinary(binData, &decoded); err != nil || len(decoded) != count {
		panic("flat binary roundtrip failed")
	}

	return benchResult{
		name:      fmt.Sprintf("Flat struct × %d (8 fields, vec)", count),
		jsonSerMs: float64(jsonSer.Nanoseconds()) / 1e6,
		asunSerMs: float64(asunSer.Nanoseconds()) / 1e6,
		binSerMs:  float64(binSer.Nanoseconds()) / 1e6,
		jsonDeMs:  float64(jsonDe.Nanoseconds()) / 1e6,
		asunDeMs:  float64(asunDe.Nanoseconds()) / 1e6,
		binDeMs:   float64(binDe.Nanoseconds()) / 1e6,
		jsonBytes: len(jsonData),
		asunBytes: len(asunData),
		binBytes:  len(binData),
	}
}

func benchAllTypes(count, iterations int) benchResult {
	items := generateAllTypes(count)

	var jsonData []byte
	start := time.Now()
	for i := 0; i < iterations; i++ {
		jsonData, _ = json.Marshal(items)
	}
	jsonSer := time.Since(start)

	var asunData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		asunData, _ = asun.Encode(items)
	}
	asunSer := time.Since(start)

	var binData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		binData, _ = asun.EncodeBinary(items)
	}
	binSer := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []AllTypes
		_ = json.Unmarshal(jsonData, &out)
	}
	jsonDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []AllTypes
		_ = asun.Decode(asunData, &out)
	}
	asunDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []AllTypes
		_ = asun.DecodeBinary(binData, &out)
	}
	binDe := time.Since(start)

	return benchResult{
		name:      fmt.Sprintf("All-types struct × %d (16 fields, vec)", count),
		jsonSerMs: float64(jsonSer.Nanoseconds()) / 1e6,
		asunSerMs: float64(asunSer.Nanoseconds()) / 1e6,
		binSerMs:  float64(binSer.Nanoseconds()) / 1e6,
		jsonDeMs:  float64(jsonDe.Nanoseconds()) / 1e6,
		asunDeMs:  float64(asunDe.Nanoseconds()) / 1e6,
		binDeMs:   float64(binDe.Nanoseconds()) / 1e6,
		jsonBytes: len(jsonData),
		asunBytes: len(asunData),
		binBytes:  len(binData),
	}
}

func benchDeep(count, iterations int) benchResult {
	companies := generateCompanies(count)

	var jsonData []byte
	start := time.Now()
	for i := 0; i < iterations; i++ {
		jsonData, _ = json.Marshal(companies)
	}
	jsonSer := time.Since(start)

	var asunData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		asunData, _ = asun.Encode(companies)
	}
	asunSer := time.Since(start)

	var binData []byte
	start = time.Now()
	for i := 0; i < iterations; i++ {
		binData, _ = asun.EncodeBinary(companies)
	}
	binSer := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []Company
		_ = json.Unmarshal(jsonData, &out)
	}
	jsonDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []Company
		_ = asun.Decode(asunData, &out)
	}
	asunDe := time.Since(start)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		var out []Company
		_ = asun.DecodeBinary(binData, &out)
	}
	binDe := time.Since(start)

	return benchResult{
		name:      fmt.Sprintf("5-level deep × %d (Company>Division>Team>Project>Task)", count),
		jsonSerMs: float64(jsonSer.Nanoseconds()) / 1e6,
		asunSerMs: float64(asunSer.Nanoseconds()) / 1e6,
		binSerMs:  float64(binSer.Nanoseconds()) / 1e6,
		jsonDeMs:  float64(jsonDe.Nanoseconds()) / 1e6,
		asunDeMs:  float64(asunDe.Nanoseconds()) / 1e6,
		binDeMs:   float64(binDe.Nanoseconds()) / 1e6,
		jsonBytes: len(jsonData),
		asunBytes: len(asunData),
		binBytes:  len(binData),
	}
}

func benchSingleRoundtrip(iterations int) (asunMs, jsonMs float64) {
	user := User{
		ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30,
		Score: 95.5, Active: true, Role: "engineer", City: "NYC",
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		s, _ := asun.Encode(&user)
		var out User
		_ = asun.Decode(s, &out)
	}
	asunMs = float64(time.Since(start).Nanoseconds()) / 1e6

	start = time.Now()
	for i := 0; i < iterations; i++ {
		s, _ := json.Marshal(&user)
		var out User
		_ = json.Unmarshal(s, &out)
	}
	jsonMs = float64(time.Since(start).Nanoseconds()) / 1e6
	return
}

func benchDeepSingleRoundtrip(iterations int) (asunMs, jsonMs float64) {
	company := generateCompanies(1)[0]

	start := time.Now()
	for i := 0; i < iterations; i++ {
		s, _ := asun.Encode(&company)
		var out Company
		_ = asun.Decode(s, &out)
	}
	asunMs = float64(time.Since(start).Nanoseconds()) / 1e6

	start = time.Now()
	for i := 0; i < iterations; i++ {
		s, _ := json.Marshal(&company)
		var out Company
		_ = json.Unmarshal(s, &out)
	}
	jsonMs = float64(time.Since(start).Nanoseconds()) / 1e6
	return
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║            ASUN vs JSON Comprehensive Benchmark              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("\nSystem: %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Iterations per test: 100")

	fmt.Println()
	printSection("Section 1: Flat Struct (schema-driven vec)", 47)
	fmt.Println()
	for _, count := range []int{100, 500, 1000, 5000} {
		benchFlat(count, 100).print()
		fmt.Println()
	}

	printSection("Section 2: All-Types Struct (16 fields)", 48)
	fmt.Println()
	for _, count := range []int{100, 500} {
		benchAllTypes(count, 100).print()
		fmt.Println()
	}

	printSection("Section 3: 5-Level Deep Nesting (Company hierarchy)", 60)
	fmt.Println()
	for _, count := range []int{10, 50, 100} {
		benchDeep(count, 50).print()
		fmt.Println()
	}

	printSection("Section 4: Single Struct Roundtrip (10000x)", 48)
	fmt.Println()
	asunFlat, jsonFlat := benchSingleRoundtrip(10000)
	fmt.Printf("  Flat:  ASUN %8.2fms | JSON %8.2fms | ratio %.2fx\n", asunFlat, jsonFlat, jsonFlat/asunFlat)
	asunDeep, jsonDeep := benchDeepSingleRoundtrip(10000)
	fmt.Printf("  Deep:  ASUN %8.2fms | JSON %8.2fms | ratio %.2fx\n", asunDeep, jsonDeep, jsonDeep/asunDeep)

	fmt.Println()
	printSection("Section 5: Large Payload (10k records)", 48)
	fmt.Println()
	large := benchFlat(10000, 10)
	fmt.Println("  (10 iterations for large payload)")
	large.print()

	fmt.Println()
	printSection("Section 6: Annotated vs Unannotated Schema (deserialize)", 64)
	fmt.Println()
	{
		users := generateUsers(1000)
		untyped := mustAsunEncode(users, false)
		typed := mustAsunEncode(users, true)
		deIters := 200

		start := time.Now()
		for i := 0; i < deIters; i++ {
			var out []User
			_ = asun.Decode(untyped, &out)
		}
		untypedMs := float64(time.Since(start).Nanoseconds()) / 1e6

		start = time.Now()
		for i := 0; i < deIters; i++ {
			var out []User
			_ = asun.Decode(typed, &out)
		}
		typedMs := float64(time.Since(start).Nanoseconds()) / 1e6

		fmt.Printf("  Flat struct x 1000 (%d iters, deserialize only)\n", deIters)
		fmt.Printf("    Unannotated: %8.2fms  (%d B)\n", untypedMs, len(untyped))
		fmt.Printf("    Annotated:   %8.2fms  (%d B)\n", typedMs, len(typed))
		fmt.Printf("    Ratio: %.3fx (unannotated / annotated)\n", untypedMs/typedMs)
	}

	fmt.Println()
	printSection("Section 7: Annotated vs Unannotated Schema (serialize)", 62)
	fmt.Println()
	{
		users := generateUsers(1000)
		serIters := 200

		start := time.Now()
		var untyped []byte
		for i := 0; i < serIters; i++ {
			untyped, _ = asun.Encode(users)
		}
		untypedMs := float64(time.Since(start).Nanoseconds()) / 1e6

		start = time.Now()
		var typed []byte
		for i := 0; i < serIters; i++ {
			typed, _ = asun.EncodeTyped(users)
		}
		typedMs := float64(time.Since(start).Nanoseconds()) / 1e6

		fmt.Printf("  Flat struct x 1000 (%d iters, serialize only)\n", serIters)
		fmt.Printf("    Unannotated: %8.2fms  (%d B)\n", untypedMs, len(untyped))
		fmt.Printf("    Annotated:   %8.2fms  (%d B)\n", typedMs, len(typed))
		fmt.Printf("    Ratio: %.3fx (unannotated / annotated)\n", untypedMs/typedMs)
	}

	fmt.Println()
	printSection("Section 8: Throughput Summary", 48)
	fmt.Println()
	{
		users := generateUsers(1000)
		jsonData, _ := json.Marshal(users)
		asunData, _ := asun.Encode(users)
		iters := 100

		start := time.Now()
		for i := 0; i < iters; i++ {
			_, _ = json.Marshal(users)
		}
		jsonSerDur := time.Since(start).Seconds()

		start = time.Now()
		for i := 0; i < iters; i++ {
			_, _ = asun.Encode(users)
		}
		asunSerDur := time.Since(start).Seconds()

		start = time.Now()
		for i := 0; i < iters; i++ {
			var out []User
			_ = json.Unmarshal(jsonData, &out)
		}
		jsonDeDur := time.Since(start).Seconds()

		start = time.Now()
		for i := 0; i < iters; i++ {
			var out []User
			_ = asun.Decode(asunData, &out)
		}
		asunDeDur := time.Since(start).Seconds()

		totalRecords := 1000.0 * float64(iters)
		fmt.Printf("  Serialize throughput (1000 records x %d iters):\n", iters)
		fmt.Printf("    JSON: %.0f records/s\n", totalRecords/jsonSerDur)
		fmt.Printf("    ASUN: %.0f records/s\n", totalRecords/asunSerDur)
		fmt.Printf("    Speed: %.2fx\n", jsonSerDur/asunSerDur)
		fmt.Println("  Deserialize throughput:")
		fmt.Printf("    JSON: %.0f records/s\n", totalRecords/jsonDeDur)
		fmt.Printf("    ASUN: %.0f records/s\n", totalRecords/asunDeDur)
		fmt.Printf("    Speed: %.2fx\n", jsonDeDur/asunDeDur)
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Benchmark Complete                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
