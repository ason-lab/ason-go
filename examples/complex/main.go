package main

import (
	"encoding/json"
	"fmt"
	"log"

	asun "github.com/asun-lab/asun-go"
)

type Department struct {
	Title string `asun:"title" json:"title"`
}

type Employee struct {
	ID     int64      `asun:"id" json:"id"`
	Name   string     `asun:"name" json:"name"`
	Dept   Department `asun:"dept" json:"dept"`
	Skills []string   `asun:"skills" json:"skills"`
	Active bool       `asun:"active" json:"active"`
}

type AttrEntry struct {
	Key   string `asun:"key" json:"key"`
	Value int64  `asun:"value" json:"value"`
}

type WithEntries struct {
	Name  string      `asun:"name" json:"name"`
	Attrs []AttrEntry `asun:"attrs" json:"attrs"`
}

type Person struct {
	Name string `asun:"name" json:"name"`
	Age  int64  `asun:"age" json:"age"`
}

type GroupEntry struct {
	Key   string   `asun:"key" json:"key"`
	Value []Person `asun:"value" json:"value"`
}

type WithGroupEntries struct {
	Groups []GroupEntry `asun:"groups" json:"groups"`
}

type Nested struct {
	Name string  `asun:"name" json:"name"`
	Addr Address `asun:"addr" json:"addr"`
}

type Address struct {
	City string `asun:"city" json:"city"`
	Zip  int64  `asun:"zip" json:"zip"`
}

type AllTypes struct {
	B         bool      `asun:"b" json:"b"`
	I8v       int8      `asun:"i8v" json:"i8v"`
	I16v      int16     `asun:"i16v" json:"i16v"`
	I32v      int32     `asun:"i32v" json:"i32v"`
	I64v      int64     `asun:"i64v" json:"i64v"`
	U8v       uint8     `asun:"u8v" json:"u8v"`
	U16v      uint16    `asun:"u16v" json:"u16v"`
	U32v      uint32    `asun:"u32v" json:"u32v"`
	U64v      uint64    `asun:"u64v" json:"u64v"`
	F32v      float32   `asun:"f32v" json:"f32v"`
	F64v      float64   `asun:"f64v" json:"f64v"`
	S         string    `asun:"s" json:"s"`
	OptSome   *int64    `asun:"opt_some" json:"opt_some"`
	OptNone   *int64    `asun:"opt_none" json:"opt_none"`
	VecInt    []int64   `asun:"vec_int" json:"vec_int"`
	VecStr    []string  `asun:"vec_str" json:"vec_str"`
	NestedVec [][]int64 `asun:"nested_vec" json:"nested_vec"`
}

type Building struct {
	Name        string  `asun:"name" json:"name"`
	Floors      int64   `asun:"floors" json:"floors"`
	Residential bool    `asun:"residential" json:"residential"`
	HeightM     float64 `asun:"height_m" json:"height_m"`
}

type Street struct {
	Name      string     `asun:"name" json:"name"`
	LengthKm  float64    `asun:"length_km" json:"length_km"`
	Buildings []Building `asun:"buildings" json:"buildings"`
}

type District struct {
	Name       string   `asun:"name" json:"name"`
	Population int64    `asun:"population" json:"population"`
	Streets    []Street `asun:"streets" json:"streets"`
}

type City struct {
	Name       string     `asun:"name" json:"name"`
	Population int64      `asun:"population" json:"population"`
	AreaKm2    float64    `asun:"area_km2" json:"area_km2"`
	Districts  []District `asun:"districts" json:"districts"`
}

type Region struct {
	Name   string `asun:"name" json:"name"`
	Cities []City `asun:"cities" json:"cities"`
}

type Country struct {
	Name        string   `asun:"name" json:"name"`
	Code        string   `asun:"code" json:"code"`
	Population  int64    `asun:"population" json:"population"`
	GdpTrillion float64  `asun:"gdp_trillion" json:"gdp_trillion"`
	Regions     []Region `asun:"regions" json:"regions"`
}

type State struct {
	Name       string `asun:"name" json:"name"`
	Capital    string `asun:"capital" json:"capital"`
	Population int64  `asun:"population" json:"population"`
}

type Nation struct {
	Name   string  `asun:"name" json:"name"`
	States []State `asun:"states" json:"states"`
}

type Continent struct {
	Name    string   `asun:"name" json:"name"`
	Nations []Nation `asun:"nations" json:"nations"`
}

type Planet struct {
	Name       string      `asun:"name" json:"name"`
	RadiusKm   float64     `asun:"radius_km" json:"radius_km"`
	HasLife    bool        `asun:"has_life" json:"has_life"`
	Continents []Continent `asun:"continents" json:"continents"`
}

type SolarSystem struct {
	Name     string   `asun:"name" json:"name"`
	StarType string   `asun:"star_type" json:"star_type"`
	Planets  []Planet `asun:"planets" json:"planets"`
}

type Galaxy struct {
	Name              string        `asun:"name" json:"name"`
	StarCountBillions float64       `asun:"star_count_billions" json:"star_count_billions"`
	Systems           []SolarSystem `asun:"systems" json:"systems"`
}

type Universe struct {
	Name            string   `asun:"name" json:"name"`
	AgeBillionYears float64  `asun:"age_billion_years" json:"age_billion_years"`
	Galaxies        []Galaxy `asun:"galaxies" json:"galaxies"`
}

type DbConfig struct {
	Host           string  `asun:"host" json:"host"`
	Port           int64   `asun:"port" json:"port"`
	MaxConnections int64   `asun:"max_connections" json:"max_connections"`
	SSL            bool    `asun:"ssl" json:"ssl"`
	TimeoutMs      float64 `asun:"timeout_ms" json:"timeout_ms"`
}

type CacheConfig struct {
	Enabled    bool  `asun:"enabled" json:"enabled"`
	TTLSeconds int64 `asun:"ttl_seconds" json:"ttl_seconds"`
	MaxSizeMb  int64 `asun:"max_size_mb" json:"max_size_mb"`
}

type LogConfig struct {
	Level  string  `asun:"level" json:"level"`
	File   *string `asun:"file" json:"file"`
	Rotate bool    `asun:"rotate" json:"rotate"`
}

type StringEntry struct {
	Key   string `asun:"key" json:"key"`
	Value string `asun:"value" json:"value"`
}

type ServiceConfig struct {
	Name     string        `asun:"name" json:"name"`
	Version  string        `asun:"version" json:"version"`
	Db       DbConfig      `asun:"db" json:"db"`
	Cache    CacheConfig   `asun:"cache" json:"cache"`
	Log      LogConfig     `asun:"log" json:"log"`
	Features []string      `asun:"features" json:"features"`
	Env      []StringEntry `asun:"env" json:"env"`
}

func i64ptr(v int64) *int64   { return &v }
func strptr(v string) *string { return &v }

func main() {
	fmt.Println("=== ASUN Complex Examples ===")
	fmt.Println()

	// 1. Nested struct
	fmt.Println("1. Nested struct:")
	var emp Employee
	if err := asun.Decode([]byte("{id,name,dept@{title},skills,active}:(1,Alice,(Manager),[rust],true)"), &emp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   %+v\n\n", emp)

	// 2. Vec with nested structs
	fmt.Println("2. Vec with nested structs:")
	input2 := []byte(`[{id,name,dept@{title},skills@[str],active}]:
  (1, Alice, (Manager), [Rust, Go], true),
  (2, Bob, (Engineer), [Python], false),
  (3, "Carol Smith", (Director), [Leadership, Strategy], true)`)
	var employees []Employee
	if err := asun.Decode(input2, &employees); err != nil {
		log.Fatal(err)
	}
	for _, e := range employees {
		fmt.Printf("   %+v\n", e)
	}

	// 3. Entry-list field
	fmt.Println("\n3. Entry-list field:")
	var wm WithEntries
	if err := asun.Decode([]byte("{name,attrs@[{key,value}]}:(Alice,[(age,30),(score,95)])"), &wm); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   %+v\n", wm)

	// 3b. Nested entry-list field
	fmt.Println("\n3b. Nested entry-list field:")
	var groups WithGroupEntries
	if err := asun.Decode([]byte("{groups@[{key,value@[{name,age}]}]}:([(teamA,[(Alice,30),(Bob,28)]),(teamB,[(Carol,41)])])"), &groups); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   %+v\n", groups)

	// 4. Nested struct roundtrip
	fmt.Println("\n4. Nested struct roundtrip:")
	nested := Nested{Name: "Alice", Addr: Address{City: "NYC", Zip: 10001}}
	s, err := asun.Encode(&nested)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized:   %s\n", s)
	var nested2 Nested
	if err := asun.Decode(s, &nested2); err != nil {
		log.Fatal(err)
	}
	if nested != nested2 {
		log.Fatal("nested roundtrip mismatch")
	}
	fmt.Println("   ✓ roundtrip OK")

	// 5. Escaped strings
	fmt.Println("\n5. Escaped strings:")
	type Note struct {
		Text string `asun:"text"`
	}
	note := Note{Text: "say \"hi\", then (wave)\tnewline\nend"}
	s, err = asun.Encode(&note)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized:   %s\n", s)
	var note2 Note
	if err := asun.Decode(s, &note2); err != nil {
		log.Fatal(err)
	}
	if note != note2 {
		log.Fatal("escape roundtrip mismatch")
	}
	fmt.Println("   ✓ escape roundtrip OK")

	// 6. Float fields
	fmt.Println("\n6. Float fields:")
	type Measurement struct {
		ID    int64   `asun:"id"`
		Value float64 `asun:"value"`
		Label string  `asun:"label"`
	}
	m := Measurement{ID: 2, Value: 95.0, Label: "score"}
	s, err = asun.Encode(&m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized: %s\n", s)
	var m2 Measurement
	if err := asun.Decode(s, &m2); err != nil {
		log.Fatal(err)
	}
	if m != m2 {
		log.Fatal("float roundtrip mismatch")
	}
	fmt.Println("   ✓ float roundtrip OK")

	// 7. Negative numbers
	fmt.Println("\n7. Negative numbers:")
	type Nums struct {
		A int64   `asun:"a"`
		B float64 `asun:"b"`
		C int64   `asun:"c"`
	}
	n := Nums{A: -42, B: -3.14, C: -9223372036854775807}
	s, err = asun.Encode(&n)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized:   %s\n", s)
	var n2 Nums
	if err := asun.Decode(s, &n2); err != nil {
		log.Fatal(err)
	}
	if n != n2 {
		log.Fatal("negative roundtrip mismatch")
	}
	fmt.Println("   ✓ negative roundtrip OK")

	// 8. All types struct
	fmt.Println("\n8. All types struct:")
	optVal := int64(42)
	all := AllTypes{
		B: true, I8v: -128, I16v: -32768, I32v: -2147483648, I64v: -9223372036854775807,
		U8v: 255, U16v: 65535, U32v: 4294967295, U64v: 18446744073709551615,
		F32v: 3.15, F64v: 2.718281828459045,
		S:         "hello, world (test) [arr]",
		OptSome:   &optVal,
		OptNone:   nil,
		VecInt:    []int64{1, 2, 3, -4, 0},
		VecStr:    []string{"alpha", "beta gamma", "delta"},
		NestedVec: [][]int64{{1, 2}, {3, 4, 5}},
	}
	s, err = asun.Encode(&all)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized (%d bytes):\n", len(s))
	fmt.Printf("   %s\n", s)
	var all2 AllTypes
	if err := asun.Decode(s, &all2); err != nil {
		log.Fatal(err)
	}
	if all.B != all2.B || all.I64v != all2.I64v || all.U64v != all2.U64v || all.S != all2.S {
		log.Fatal("all-types roundtrip mismatch")
	}
	if all2.OptSome == nil || *all2.OptSome != 42 {
		log.Fatal("all-types opt_some mismatch")
	}
	if all2.OptNone != nil {
		log.Fatal("all-types opt_none should be nil")
	}
	fmt.Println("   ✓ all-types roundtrip OK")

	// 9. 5-level deep: Country > Region > City > District > Street > Building
	fmt.Println("\n9. Five-level nesting (Country>Region>City>District>Street>Building):")
	country := Country{
		Name: "Rustland", Code: "RL", Population: 50000000, GdpTrillion: 1.5,
		Regions: []Region{
			{Name: "Northern", Cities: []City{
				{Name: "Ferriton", Population: 2000000, AreaKm2: 350.5, Districts: []District{
					{Name: "Downtown", Population: 500000, Streets: []Street{
						{Name: "Main St", LengthKm: 2.5, Buildings: []Building{
							{Name: "Tower A", Floors: 50, Residential: false, HeightM: 200.0},
							{Name: "Apt Block 1", Floors: 12, Residential: true, HeightM: 40.5},
						}},
						{Name: "Oak Ave", LengthKm: 1.2, Buildings: []Building{
							{Name: "Library", Floors: 3, Residential: false, HeightM: 15.0},
						}},
					}},
					{Name: "Harbor", Population: 150000, Streets: []Street{
						{Name: "Dock Rd", LengthKm: 0.8, Buildings: []Building{
							{Name: "Warehouse 7", Floors: 1, Residential: false, HeightM: 8.0},
						}},
					}},
				}},
			}},
			{Name: "Southern", Cities: []City{
				{Name: "Crabville", Population: 800000, AreaKm2: 120.0, Districts: []District{
					{Name: "Old Town", Population: 200000, Streets: []Street{
						{Name: "Heritage Ln", LengthKm: 0.5, Buildings: []Building{
							{Name: "Museum", Floors: 2, Residential: false, HeightM: 12.0},
							{Name: "Town Hall", Floors: 4, Residential: false, HeightM: 20.0},
						}},
					}},
				}},
			}},
		},
	}
	s, err = asun.Encode(&country)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized (%d bytes)\n", len(s))
	preview := string(s)
	if len(preview) > 200 {
		preview = preview[:200] + "..."
	}
	fmt.Printf("   first 200 chars: %s\n", preview)
	var country2 Country
	if err := asun.Decode(s, &country2); err != nil {
		log.Fatal(err)
	}
	if country.Name != country2.Name || country.Population != country2.Population {
		log.Fatal("5-level roundtrip mismatch")
	}
	fmt.Println("   ✓ 5-level ASUN-text roundtrip OK")

	// ASUN binary roundtrip
	binBytes, err := asun.EncodeBinary(&country)
	if err != nil {
		log.Fatal(err)
	}
	var country3 Country
	if err := asun.DecodeBinary(binBytes, &country3); err != nil {
		log.Fatal(err)
	}
	if country.Name != country3.Name || country.Population != country3.Population {
		log.Fatal("5-level binary roundtrip mismatch")
	}
	fmt.Println("   ✓ 5-level ASUN-bin roundtrip OK")

	jsonBytes, _ := json.Marshal(&country)
	fmt.Printf("   ASUN text: %d B | ASUN bin: %d B | JSON: %d B\n",
		len(s), len(binBytes), len(jsonBytes))
	fmt.Printf("   BIN vs JSON: %.0f%% smaller | TEXT vs JSON: %.0f%% smaller\n",
		(1.0-float64(len(binBytes))/float64(len(jsonBytes)))*100.0,
		(1.0-float64(len(s))/float64(len(jsonBytes)))*100.0)

	// 10. 7-level deep
	fmt.Println("\n10. Seven-level nesting (Universe>Galaxy>SolarSystem>Planet>Continent>Nation>State):")
	universe := Universe{
		Name: "Observable", AgeBillionYears: 13.8,
		Galaxies: []Galaxy{{
			Name: "Milky Way", StarCountBillions: 250.0,
			Systems: []SolarSystem{{
				Name: "Sol", StarType: "G2V",
				Planets: []Planet{
					{Name: "Earth", RadiusKm: 6371.0, HasLife: true, Continents: []Continent{
						{Name: "Asia", Nations: []Nation{
							{Name: "Japan", States: []State{
								{Name: "Tokyo", Capital: "Shinjuku", Population: 14000000},
								{Name: "Osaka", Capital: "Osaka City", Population: 8800000},
							}},
							{Name: "China", States: []State{
								{Name: "Beijing", Capital: "Beijing", Population: 21500000},
							}},
						}},
						{Name: "Europe", Nations: []Nation{
							{Name: "Germany", States: []State{
								{Name: "Bavaria", Capital: "Munich", Population: 13000000},
								{Name: "Berlin", Capital: "Berlin", Population: 3600000},
							}},
						}},
					}},
					{Name: "Mars", RadiusKm: 3389.5, HasLife: false, Continents: []Continent{}},
				},
			}},
		}},
	}
	s, err = asun.Encode(&universe)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized (%d bytes)\n", len(s))
	var universe2 Universe
	if err := asun.Decode(s, &universe2); err != nil {
		log.Fatal(err)
	}
	if universe.Name != universe2.Name {
		log.Fatal("7-level roundtrip mismatch")
	}
	fmt.Println("   ✓ 7-level ASUN-text roundtrip OK")

	// ASUN binary roundtrip
	uniBin, err := asun.EncodeBinary(&universe)
	if err != nil {
		log.Fatal(err)
	}
	var universe3 Universe
	if err := asun.DecodeBinary(uniBin, &universe3); err != nil {
		log.Fatal(err)
	}
	if universe.Name != universe3.Name {
		log.Fatal("7-level binary roundtrip mismatch")
	}
	fmt.Println("   ✓ 7-level ASUN-bin roundtrip OK")

	jsonBytes, _ = json.Marshal(&universe)
	fmt.Printf("   ASUN text: %d B | ASUN bin: %d B | JSON: %d B\n",
		len(s), len(uniBin), len(jsonBytes))
	fmt.Printf("   BIN vs JSON: %.0f%% smaller | TEXT vs JSON: %.0f%% smaller\n",
		(1.0-float64(len(uniBin))/float64(len(jsonBytes)))*100.0,
		(1.0-float64(len(s))/float64(len(jsonBytes)))*100.0)

	// 11. Service config
	fmt.Println("\n11. Complex config struct (nested + entry list + optional):")
	config := ServiceConfig{
		Name: "my-service", Version: "2.1.0",
		Db:       DbConfig{Host: "db.example.com", Port: 5432, MaxConnections: 100, SSL: true, TimeoutMs: 3000.5},
		Cache:    CacheConfig{Enabled: true, TTLSeconds: 3600, MaxSizeMb: 512},
		Log:      LogConfig{Level: "info", File: strptr("/var/log/app.log"), Rotate: true},
		Features: []string{"auth", "rate-limit", "websocket"},
		Env: []StringEntry{
			{Key: "RUST_LOG", Value: "debug"},
			{Key: "DATABASE_URL", Value: "postgres://localhost:5432/mydb"},
			{Key: "SECRET_KEY", Value: "abc123!@#"},
		},
	}
	s, err = asun.Encode(&config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   serialized (%d bytes):\n", len(s))
	fmt.Printf("   %s\n", s)
	var config2 ServiceConfig
	if err := asun.Decode(s, &config2); err != nil {
		log.Fatal(err)
	}
	if config.Name != config2.Name || config.Db.Port != config2.Db.Port {
		log.Fatal("config roundtrip mismatch")
	}
	fmt.Println("   ✓ config ASUN-text roundtrip OK")
	jsonBytes, _ = json.Marshal(&config)
	fmt.Printf("   ASUN text: %d B | JSON: %d B | TEXT vs JSON: %.0f%% smaller\n",
		len(s), len(jsonBytes), (1.0-float64(len(s))/float64(len(jsonBytes)))*100.0)
	// Binary roundtrip
	cfgBin, err := asun.EncodeBinary(&config)
	if err != nil {
		log.Fatal(err)
	}
	var config3 ServiceConfig
	if err := asun.DecodeBinary(cfgBin, &config3); err != nil {
		log.Fatal(err)
	}
	if config.Name != config3.Name || config.Db.Port != config3.Db.Port {
		log.Fatal("config binary roundtrip mismatch")
	}
	fmt.Println("   ✓ config ASUN-bin roundtrip OK")
	fmt.Printf("   ASUN bin: %d B | BIN vs JSON: %.0f%% smaller\n",
		len(cfgBin), (1.0-float64(len(cfgBin))/float64(len(jsonBytes)))*100.0)

	// 12. Large structure — 100 countries
	fmt.Println("\n12. Large structure (100 countries × nested regions):")
	countries := make([]Country, 100)
	for i := range countries {
		regions := make([]Region, 3)
		for r := 0; r < 3; r++ {
			cities := make([]City, 2)
			for c := 0; c < 2; c++ {
				cities[c] = City{
					Name: fmt.Sprintf("City_%d_%d_%d", i, r, c), Population: int64(100000 + c*50000),
					AreaKm2: 50.0 + float64(c)*25.5,
					Districts: []District{{
						Name: fmt.Sprintf("Dist_%d", c), Population: int64(50000 + c*10000),
						Streets: []Street{{
							Name: fmt.Sprintf("St_%d", c), LengthKm: 1.0 + float64(c)*0.5,
							Buildings: []Building{
								{Name: fmt.Sprintf("Bldg_%d_0", c), Floors: 5, Residential: true, HeightM: 15.0},
								{Name: fmt.Sprintf("Bldg_%d_1", c), Floors: 8, Residential: false, HeightM: 25.5},
							},
						}},
					}},
				}
			}
			regions[r] = Region{Name: fmt.Sprintf("Region_%d_%d", i, r), Cities: cities}
		}
		countries[i] = Country{
			Name: fmt.Sprintf("Country_%d", i), Code: fmt.Sprintf("C%02d", i%100),
			Population: int64(1000000 + i*500000), GdpTrillion: float64(i) * 0.5, Regions: regions,
		}
	}
	totalASUN, totalJSON, totalBIN := 0, 0, 0
	for i := range countries {
		as, _ := asun.Encode(&countries[i])
		js, _ := json.Marshal(&countries[i])
		bs, _ := asun.EncodeBinary(&countries[i])
		// Verify text roundtrip
		var c2 Country
		if err := asun.Decode(as, &c2); err != nil {
			log.Fatalf("country %d roundtrip failed: %v", i, err)
		}
		if countries[i].Name != c2.Name {
			log.Fatalf("country %d name mismatch", i)
		}
		// Verify binary roundtrip
		var c3 Country
		if err := asun.DecodeBinary(bs, &c3); err != nil {
			log.Fatalf("country %d binary roundtrip failed: %v", i, err)
		}
		if countries[i].Name != c3.Name {
			log.Fatalf("country %d binary name mismatch", i)
		}
		totalASUN += len(as)
		totalJSON += len(js)
		totalBIN += len(bs)
	}
	fmt.Println("   100 countries with 5-level nesting:")
	fmt.Printf("   Total ASUN text: %d bytes (%.1f KB)\n", totalASUN, float64(totalASUN)/1024.0)
	fmt.Printf("   Total ASUN bin:  %d bytes (%.1f KB)\n", totalBIN, float64(totalBIN)/1024.0)
	fmt.Printf("   Total JSON:      %d bytes (%.1f KB)\n", totalJSON, float64(totalJSON)/1024.0)
	fmt.Printf("   TEXT vs JSON: %.0f%% smaller | BIN vs JSON: %.0f%% smaller\n",
		(1.0-float64(totalASUN)/float64(totalJSON))*100.0,
		(1.0-float64(totalBIN)/float64(totalJSON))*100.0)
	fmt.Println("   ✓ all 100 countries roundtrip OK (text + bin)")

	// 13. Deserialize with nested schema type hints
	fmt.Println("\n13. Deserialize with nested schema type hints:")
	deepInput := []byte("{name,code,population,gdp_trillion,regions@[{name,cities@[{name,population,area_km2,districts@[{name,population,streets@[{name,length_km,buildings@[{name,floors,residential,height_m}]}]}]}]}]}:(TestLand,TL,1000000,0.5,[(TestRegion,[(TestCity,500000,100.0,[(Central,250000,[(Main St,2.5,[(HQ,10,false,45.0)])])])])])")
	var dc Country
	if err := asun.Decode(deepInput, &dc); err != nil {
		log.Fatal(err)
	}
	if dc.Name != "TestLand" {
		log.Fatal("deep schema parse name mismatch")
	}
	bld := dc.Regions[0].Cities[0].Districts[0].Streets[0].Buildings[0]
	if bld.Name != "HQ" {
		log.Fatal("deep schema parse building name mismatch")
	}
	fmt.Println("   ✓ deep schema type-hint parse OK")
	fmt.Printf("   Building at depth 6: %+v\n", bld)

	// 14. Typed serialization
	fmt.Println("\n14. Typed serialization (EncodeTyped):")
	empForTyped := Employee{ID: 1, Name: "Alice", Dept: Department{Title: "Engineering"}, Skills: []string{"Rust", "Go"}, Active: true}
	typedBytes, err := asun.EncodeTyped(&empForTyped)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   nested struct: %s\n", typedBytes)
	var empBack Employee
	if err := asun.Decode(typedBytes, &empBack); err != nil {
		log.Fatal(err)
	}
	if empBack.Name != "Alice" {
		log.Fatal("typed nested struct roundtrip mismatch")
	}
	fmt.Println("   ✓ typed nested struct roundtrip OK")

	allTyped, err := asun.EncodeTyped(&all)
	if err != nil {
		log.Fatal(err)
	}
	allTypedStr := string(allTyped)
	p := allTypedStr
	if len(p) > 80 {
		p = p[:80]
	}
	fmt.Printf("   all-types (%d bytes): %s...\n", len(allTyped), p)
	var allBack AllTypes
	if err := asun.Decode(allTyped, &allBack); err != nil {
		log.Fatal(err)
	}
	if allBack.S != all.S || allBack.B != all.B {
		log.Fatal("typed all-types roundtrip mismatch")
	}
	fmt.Println("   ✓ typed all-types roundtrip OK")

	configTyped, err := asun.EncodeTyped(&config)
	if err != nil {
		log.Fatal(err)
	}
	configTypedStr := string(configTyped)
	p = configTypedStr
	if len(p) > 100 {
		p = p[:100]
	}
	fmt.Printf("   config (%d bytes): %s...\n", len(configTyped), p)
	var configBack ServiceConfig
	if err := asun.Decode(configTyped, &configBack); err != nil {
		log.Fatal(err)
	}
	if configBack.Name != config.Name {
		log.Fatal("typed config roundtrip mismatch")
	}
	fmt.Println("   ✓ typed config roundtrip OK")
	untyped, _ := asun.Encode(&config)
	fmt.Printf("   untyped: %d bytes | typed: %d bytes | overhead: %d bytes\n",
		len(untyped), len(configTyped), len(configTyped)-len(untyped))

	// 15. Edge cases
	fmt.Println("\n15. Edge cases:")
	type WithVec struct {
		Items []int64 `asun:"items"`
	}
	wv := WithVec{Items: []int64{}}
	s, _ = asun.Encode(&wv)
	fmt.Printf("   empty vec: %s\n", s)
	var wv2 WithVec
	if err := asun.Decode(s, &wv2); err != nil {
		log.Fatal(err)
	}

	type Special struct {
		Val string `asun:"val"`
	}
	sp := Special{Val: "tabs\there, newlines\nhere, quotes\"and\\backslash"}
	s, _ = asun.Encode(&sp)
	fmt.Printf("   special chars: %s\n", s)
	var sp2 Special
	if err := asun.Decode(s, &sp2); err != nil {
		log.Fatal(err)
	}
	if sp != sp2 {
		log.Fatal("special chars roundtrip mismatch")
	}

	sp3 := Special{Val: "true"}
	s, _ = asun.Encode(&sp3)
	fmt.Printf("   bool-like string: %s\n", s)
	var sp4 Special
	asun.Decode(s, &sp4)
	if sp3 != sp4 {
		log.Fatal("bool-like string roundtrip mismatch")
	}

	sp5 := Special{Val: "12345"}
	s, _ = asun.Encode(&sp5)
	fmt.Printf("   number-like string: %s\n", s)
	var sp6 Special
	asun.Decode(s, &sp6)
	if sp5 != sp6 {
		log.Fatal("number-like string roundtrip mismatch")
	}
	fmt.Println("   ✓ all edge cases OK")

	// 16. Triple-nested arrays
	fmt.Println("\n16. Triple-nested arrays:")
	type Matrix3D struct {
		Data [][][]int64 `asun:"data"`
	}
	m3 := Matrix3D{Data: [][][]int64{{{1, 2}, {3, 4}}, {{5, 6, 7}, {8}}}}
	s, _ = asun.Encode(&m3)
	fmt.Printf("   %s\n", s)
	var m3b Matrix3D
	if err := asun.Decode(s, &m3b); err != nil {
		log.Fatal(err)
	}
	if len(m3b.Data) != 2 || len(m3b.Data[0]) != 2 || m3b.Data[0][0][0] != 1 {
		log.Fatal("triple-nested array roundtrip mismatch")
	}
	fmt.Println("   ✓ triple-nested array roundtrip OK")

	// 17. Comments
	fmt.Println("\n17. Comments:")
	var empComment Employee
	if err := asun.Decode([]byte("{id,name,dept@{title},skills,active}:/* inline */ (1,Alice,(HR),[rust],true)"), &empComment); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   with inline comment: %+v\n", empComment)
	fmt.Println("   ✓ comment parsing OK")

	fmt.Printf("\n=== All %d complex examples passed! ===\n", 17)
}
