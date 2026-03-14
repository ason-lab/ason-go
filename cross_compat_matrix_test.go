package ason

import "testing"

type MatrixPerson struct {
	ID   int64  `ason:"id"`
	Name string `ason:"name"`
}

type MatrixPersonWithActive struct {
	ID     int64  `ason:"id"`
	Name   string `ason:"name"`
	Active bool   `ason:"active"`
}

type MatrixInnerThin struct {
	X int64 `ason:"x"`
	Y int64 `ason:"y"`
}

type MatrixOuterThin struct {
	Name  string          `ason:"name"`
	Inner MatrixInnerThin `ason:"inner"`
}

type MatrixTaskThin struct {
	Title string `ason:"title"`
	Done  bool   `ason:"done"`
}

type MatrixProjectThin struct {
	Name  string           `ason:"name"`
	Tasks []MatrixTaskThin `ason:"tasks"`
}

type MatrixDstFewerOptionals struct {
	ID    int64   `ason:"id"`
	Label *string `ason:"label"`
}

type MatrixL3Thin struct {
	A int64 `ason:"a"`
}

type MatrixL2Thin struct {
	Name string       `ason:"name"`
	Sub  MatrixL3Thin `ason:"sub"`
}

type MatrixL1Thin struct {
	ID    int64        `ason:"id"`
	Child MatrixL2Thin `ason:"child"`
}

type MatrixPersonScore struct {
	ID    int64   `ason:"id"`
	Score float64 `ason:"score"`
}

type MatrixNoOverlap struct {
	Foo int64  `ason:"foo"`
	Bar string `ason:"bar"`
}

type MatrixNestedOptionalThin struct {
	Name string  `ason:"name"`
	Nick *string `ason:"nick"`
}

type MatrixUserWithNestedOptional struct {
	ID      int64                   `ason:"id"`
	Profile MatrixNestedOptionalThin `ason:"profile"`
}

func TestMatrix_A2_TypedSingleExtraFieldDropped(t *testing.T) {
	input := []byte("{id,name,active}:(42,Alice,true)")
	var dst MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A1_TypedSingleExactMatch(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A1_UntypedSingleExactMatch(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A2_UntypedSingleExtraFieldDropped(t *testing.T) {
	input := []byte("{id,name,active}:(42,Alice,true)")
	var dst MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A3_TypedSingleTargetExtraFieldDefaulted(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" || dst.Active {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A3_UntypedSingleTargetExtraFieldDefaulted(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" || dst.Active {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A4_TypedSingleFieldReorder(t *testing.T) {
	input := []byte("{active,id,name}:(true,42,Alice)")
	var dst MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" || !dst.Active {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A4_UntypedSingleFieldReorder(t *testing.T) {
	input := []byte("{active,id,name}:(true,42,Alice)")
	var dst MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Name != "Alice" || !dst.Active {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_A5_TypedVecExtraFieldDropped(t *testing.T) {
	input := []byte("[{id,name,active}]:(42,Alice,true),(7,Bob,false)")
	var dst []MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 42 || dst[0].Name != "Alice" {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 7 || dst[1].Name != "Bob" {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_A5_UntypedVecExtraFieldDropped(t *testing.T) {
	input := []byte("[{id,name,active}]:(42,Alice,true),(7,Bob,false)")
	var dst []MatrixPerson
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 42 || dst[0].Name != "Alice" {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 7 || dst[1].Name != "Bob" {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_N1_TypedNestedExtraFieldsDropped(t *testing.T) {
	input := []byte("{name,inner@{x,y,z,w},flag}:(test,(10,20,3.14,true),true)")
	var dst MatrixOuterThin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.Name != "test" || dst.Inner.X != 10 || dst.Inner.Y != 20 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_N1_UntypedNestedExtraFieldsDropped(t *testing.T) {
	input := []byte("{name,inner@{x,y,z,w},flag}:(test,(10,20,3.14,true),true)")
	var dst MatrixOuterThin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.Name != "test" || dst.Inner.X != 10 || dst.Inner.Y != 20 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_N2_TypedNestedVecExtraFieldsDropped(t *testing.T) {
	input := []byte("[{name,tasks@[{title,done,priority,weight}]}]:(Alpha,[(Design,true,1,0.5),(Code,false,2,0.8)]),(Beta,[(Test,false,3,1.0)])")
	var dst []MatrixProjectThin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 || dst[0].Name != "Alpha" || len(dst[0].Tasks) != 2 || dst[1].Name != "Beta" || len(dst[1].Tasks) != 1 {
		t.Fatalf("mismatch: %+v", dst)
	}
	if dst[0].Tasks[0].Title != "Design" || !dst[0].Tasks[0].Done {
		t.Fatalf("task0 mismatch: %+v", dst[0].Tasks[0])
	}
}

func TestMatrix_N2_UntypedNestedVecExtraFieldsDropped(t *testing.T) {
	input := []byte("[{name,tasks@[{title,done,priority,weight}]}]:(Alpha,[(Design,true,1,0.5),(Code,false,2,0.8)]),(Beta,[(Test,false,3,1.0)])")
	var dst []MatrixProjectThin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 || dst[0].Name != "Alpha" || len(dst[0].Tasks) != 2 || dst[1].Name != "Beta" || len(dst[1].Tasks) != 1 {
		t.Fatalf("mismatch: %+v", dst)
	}
	if dst[0].Tasks[1].Title != "Code" || dst[0].Tasks[1].Done {
		t.Fatalf("task1 mismatch: %+v", dst[0].Tasks[1])
	}
}

func TestMatrix_O1_TypedOptionalSkipTrailing(t *testing.T) {
	input := []byte("[{id,label,score,flag}]:(1,hello,95.5,true),(2,,,false)")
	var dst []MatrixDstFewerOptionals
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 1 || dst[0].Label == nil || *dst[0].Label != "hello" {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 2 || dst[1].Label != nil {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_A6_TypedVecTargetExtraFieldDefaulted(t *testing.T) {
	input := []byte("[{id,name}]:(42,Alice),(7,Bob)")
	var dst []MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 42 || dst[0].Name != "Alice" || dst[0].Active {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 7 || dst[1].Name != "Bob" || dst[1].Active {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_A6_UntypedVecTargetExtraFieldDefaulted(t *testing.T) {
	input := []byte("[{id,name}]:(42,Alice),(7,Bob)")
	var dst []MatrixPersonWithActive
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 42 || dst[0].Name != "Alice" || dst[0].Active {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 7 || dst[1].Name != "Bob" || dst[1].Active {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_N3_TypedDeepNestedExtraFieldsDropped(t *testing.T) {
	input := []byte("{id,child@{name,sub@{a,b,c},code,tags@[str]},extra}:(7,(leaf,(11,hello,true),99,[x,y]),tail)")
	var dst MatrixL1Thin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 7 || dst.Child.Name != "leaf" || dst.Child.Sub.A != 11 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_N3_UntypedDeepNestedExtraFieldsDropped(t *testing.T) {
	input := []byte("{id,child@{name,sub@{a,b,c},code,tags},extra}:(7,(leaf,(11,hello,true),99,[x,y]),tail)")
	var dst MatrixL1Thin
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 7 || dst.Child.Name != "leaf" || dst.Child.Sub.A != 11 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_O1_UntypedOptionalSkipTrailing(t *testing.T) {
	input := []byte("[{id,label,score,flag}]:(1,hello,95.5,true),(2,,,false)")
	var dst []MatrixDstFewerOptionals
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(dst))
	}
	if dst[0].ID != 1 || dst[0].Label == nil || *dst[0].Label != "hello" {
		t.Fatalf("row0 mismatch: %+v", dst[0])
	}
	if dst[1].ID != 2 || dst[1].Label != nil {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_P1_TypedPartialOverlap(t *testing.T) {
	input := []byte("{id,name,score,active}:(42,Alice,9.5,true)")
	var dst MatrixPersonScore
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Score != 9.5 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_P1_UntypedPartialOverlap(t *testing.T) {
	input := []byte("{id,name,score,active}:(42,Alice,9.5,true)")
	var dst MatrixPersonScore
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.ID != 42 || dst.Score != 9.5 {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_P2_TypedNoOverlapDefaults(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixNoOverlap
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.Foo != 0 || dst.Bar != "" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_P2_UntypedNoOverlapDefaults(t *testing.T) {
	input := []byte("{id,name}:(42,Alice)")
	var dst MatrixNoOverlap
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if dst.Foo != 0 || dst.Bar != "" {
		t.Fatalf("mismatch: %+v", dst)
	}
}

func TestMatrix_N4_TypedNestedOptionalSubset(t *testing.T) {
	input := []byte("[{id,profile@{name,nick,score},active}]:(1,(Alice,ally,9.5),true),(2,(Bob,,),false)")
	var dst []MatrixUserWithNestedOptional
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 || dst[0].ID != 1 || dst[0].Profile.Name != "Alice" || dst[0].Profile.Nick == nil || *dst[0].Profile.Nick != "ally" {
		t.Fatalf("row0 mismatch: %+v", dst)
	}
	if dst[1].ID != 2 || dst[1].Profile.Name != "Bob" || dst[1].Profile.Nick != nil {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}

func TestMatrix_N4_UntypedNestedOptionalSubset(t *testing.T) {
	input := []byte("[{id,profile@{name,nick,score},active}]:(1,(Alice,ally,9.5),true),(2,(Bob,,),false)")
	var dst []MatrixUserWithNestedOptional
	if err := Decode(input, &dst); err != nil {
		t.Fatal(err)
	}
	if len(dst) != 2 || dst[0].Profile.Nick == nil || *dst[0].Profile.Nick != "ally" {
		t.Fatalf("row0 mismatch: %+v", dst)
	}
	if dst[1].Profile.Nick != nil {
		t.Fatalf("row1 mismatch: %+v", dst[1])
	}
}
