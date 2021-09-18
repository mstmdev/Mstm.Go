package example

import (
	"fmt"
	"testing"
)

func TestMap3(t *testing.T) {
	//定义一个Map
	firstMap := make(map[string]int)
	firstMap["test"] = 99
	firstMap["test2"] = 993
	firstMap["age"] = 13
	fmt.Println(firstMap)

	//输出指定的值
	fmt.Println(firstMap["age"])

	//输出Map的长度
	fmt.Println("firstMap的长度为：", len(firstMap))

	//删除Map中的项
	delete(firstMap, "age")
	fmt.Println("firstMap的长度为：", len(firstMap))
	fmt.Println(firstMap)

	//检查对应的键是否存在  age已删除 返回false
	_, ageMap := firstMap["age"]
	fmt.Println("ageMap:", ageMap)

	//test存在  返回true
	_, testMap := firstMap["test"]
	fmt.Println("testMap:", testMap)

	//快读定义
	quickMap := map[string]int{"tel": 10086, "age": 45}
	fmt.Println(quickMap)
}
