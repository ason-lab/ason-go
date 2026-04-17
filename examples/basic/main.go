package main

import (
	"encoding/json"
	"fmt"
	"log"

	asun "github.com/asunLab/asun-go"
)

type User struct {
	ID     int64  `asun:"id" json:"id"`
	Name   string `asun:"name" json:"name"`
	Active bool   `asun:"active" json:"active"`
}

func main() {
	fmt.Println("=== ASUN Basic Examples ===")
	fmt.Println()

	// 1. Serialize a single struct
	user := User{ID: 1, Name: "Alice", Active: true}
	b, err := asun.Encode(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Serialize single struct:")
	fmt.Printf("  %s\n\n", b)

	// 2. Serialize with type annotations
	typed, err := asun.EncodeTyped(&user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Serialize with type annotations:")
	fmt.Printf("  %s\n\n", typed)

	// 3. Deserialize from ASUN (accepts both annotated and unannotated)
	input := []byte("{id@int,name@str,active@bool}:(1,Alice,true)")
	var u User
	if err := asun.Decode(input, &u); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialize single struct:")
	fmt.Printf("  %+v\n\n", u)

	// 4. Serialize a vec of structs
	users := []User{
		{ID: 1, Name: "Alice", Active: true},
		{ID: 2, Name: "Bob", Active: false},
		{ID: 3, Name: "Carol Smith", Active: true},
	}
	vecBytes, err := asun.Encode(users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Serialize vec (schema-driven):")
	fmt.Printf("  %s\n\n", vecBytes)

	// 5. Serialize vec with type annotations
	typedVec, err := asun.EncodeTyped(users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Serialize vec with type annotations:")
	fmt.Printf("  %s\n\n", typedVec)

	// 6. Deserialize vec
	vecInput := []byte(`[{id@int,name@str,active@bool}]:(1,Alice,true),(2,Bob,false),(3,"Carol Smith",true)`)
	var parsed []User
	if err := asun.Decode(vecInput, &parsed); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deserialize vec:")
	for _, u := range parsed {
		fmt.Printf("  %+v\n", u)
	}

	// 7. Multiline format
	fmt.Println("\nMultiline format:")
	multiline := []byte(`[{id, name, active}]:
  (1, Alice, true),
  (2, Bob, false),
  (3, "Carol Smith", true)`)
	var multi []User
	if err := asun.Decode(multiline, &multi); err != nil {
		log.Fatal(err)
	}
	for _, u := range multi {
		fmt.Printf("  %+v\n", u)
	}

	// 8. Roundtrip (ASUN-text vs ASUN-bin vs JSON)
	fmt.Println("\n8. Roundtrip (ASUN-text vs ASUN-bin vs JSON):")
	original := User{ID: 42, Name: "Test User", Active: true}
	// ASUN text
	asunText, err := asun.Encode(&original)
	if err != nil {
		log.Fatal(err)
	}
	var fromAsun User
	if err := asun.Decode(asunText, &fromAsun); err != nil {
		log.Fatal(err)
	}
	if original != fromAsun {
		log.Fatal("ASUN text roundtrip mismatch")
	}
	// ASUN binary
	asunBin, err := asun.EncodeBinary(&original)
	if err != nil {
		log.Fatal(err)
	}
	var fromBin User
	if err := asun.DecodeBinary(asunBin, &fromBin); err != nil {
		log.Fatal(err)
	}
	if original != fromBin {
		log.Fatal("ASUN binary roundtrip mismatch")
	}
	// JSON
	jsonData, err := json.Marshal(&original)
	if err != nil {
		log.Fatal(err)
	}
	var fromJSON User
	if err := json.Unmarshal(jsonData, &fromJSON); err != nil {
		log.Fatal(err)
	}
	if original != fromJSON {
		log.Fatal("JSON roundtrip mismatch")
	}
	fmt.Printf("  original:     %+v\n", original)
	fmt.Printf("  ASUN text:    %s (%d B)\n", asunText, len(asunText))
	fmt.Printf("  ASUN binary:  %d B\n", len(asunBin))
	fmt.Printf("  JSON:         %s (%d B)\n", jsonData, len(jsonData))
	fmt.Println("  ✓ all 3 formats roundtrip OK")

	// 9. Optional fields
	fmt.Println("\n9. Optional fields:")
	type Item struct {
		ID    int64   `asun:"id" json:"id"`
		Label *string `asun:"label" json:"label"`
	}
	var item Item
	if err := asun.Decode([]byte("{id,label}:(1,hello)"), &item); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  with value: %+v (label=%s)\n", item, *item.Label)

	var item2 Item
	if err := asun.Decode([]byte("{id,label}:(2,)"), &item2); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  with null:  %+v\n", item2)

	// 10. Array fields
	fmt.Println("\n10. Array fields:")
	type Tagged struct {
		Name string   `asun:"name" json:"name"`
		Tags []string `asun:"tags" json:"tags"`
	}
	var t Tagged
	if err := asun.Decode([]byte("{name,tags}:(Alice,[rust,go,python])"), &t); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  %+v\n", t)

	// 11. Comments
	fmt.Println("\n11. With comments:")
	var commented User
	if err := asun.Decode([]byte("/* user list */ {id,name,active}:(1,Alice,true)"), &commented); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  %+v\n", commented)

	fmt.Println("\n=== All examples passed! ===")
}
