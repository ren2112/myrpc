package main

import (
	"fmt"
	"reflect"
)

type Geek struct {
	A int `tag1:"First Tag" tag2:"Second Tag"`
	B string
}

// Main function
func main() {
	f := Geek{A: 10, B: "Number"}

	fType := reflect.TypeOf(f)
	fVal := reflect.New(fType)
	fmt.Println(fVal, fVal.Elem(), reflect.TypeOf(fVal), reflect.TypeOf(fVal.Elem()))
	fVal.Elem().Field(0).SetInt(20)
	fVal.Elem().Field(1).SetString("Number")
	f2 := fVal.Elem().Interface().(Geek)
	fmt.Printf("%+v, %d, %s\n", f2, f2.A, f2.B)
}
