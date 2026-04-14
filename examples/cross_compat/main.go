package main

import (
	"fmt"

	"github.com/asun-lab/asun-go"
)

type User struct {
	Details []Detail `asun:"details" json:"details"`
}

type Detail struct {
	ID     int64
	Name   string
	Age    int
	Gender bool
}

type Human struct {
	Details []Person `asun:"details" json:"persons"`
}
type Person struct {
	ID   int64
	Name string
	Age  int
}

func main() {
	users := []User{{Details: []Detail{{ID: 1, Name: "Alice", Age: 30, Gender: true}, {ID: 2, Name: "Bob", Age: 25, Gender: false}}}}
	buf, err := asun.Encode(users)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
	var decoded []Human
	if err := asun.Decode(buf, &decoded); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", decoded)
}
