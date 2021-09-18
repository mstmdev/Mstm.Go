package example

import (
	"fmt"
	"testing"
)

func TestArrays(t *testing.T) {
	var array [10]int
	fmt.Println("my array", array)

	array[1] = 100
	fmt.Println(array)

	fmt.Println("array的长度为：", len(array))

	//数组必须指定长度 未指定长度则为Slice
	array2 := [5]int{3, 4, 5, 6}
	fmt.Println(array2)
}
