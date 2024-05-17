package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Result struct {
	Values []interface{}  // 返回值
	Types  []reflect.Type // 返回值类型
}

type A struct {
	Name string
}

func main() {
	// 创建一个A类型的实例
	aInstance := A{"zjp"}

	// 将A的实例转换为interface{}类型并放入切片
	values := []interface{}{aInstance}

	// 获取A的类型
	aType := reflect.TypeOf(aInstance)

	// 创建Result实例
	result := Result{
		Values: values,                // 正确的[]interface{}类型
		Types:  []reflect.Type{aType}, // 类型切片，包含A的类型
	}
	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
	var from Result
	json.Unmarshal([]byte(res), &from)
	fmt.Println(from)

	firstValue := from.Values[0]
	firstValueJSON, err := json.Marshal(firstValue)
	if err != nil {
		fmt.Println("Error marshaling first value to JSON:", err)
		return
	}
	fmt.Println("First value as JSON:", string(firstValueJSON))

	// 将JSON字符串反序列化为A类型的结构体实例
	var aInstance2 A
	err = json.Unmarshal(firstValueJSON, &aInstance2)
	if err != nil {
		fmt.Println("Error unmarshaling JSON to A:", err)
		return
	}
	fmt.Printf("Deserialized A instance: %+v\n", aInstance2)
}
