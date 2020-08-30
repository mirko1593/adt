package main

import (
	"fmt"
	"reflect"
)

func main() {
}

func hello(val interface{}) {
	if v, ok := val.(string); ok {
		fmt.Println("string", v, reflect.ValueOf(v).Kind(), reflect.TypeOf(val))
	}
}
