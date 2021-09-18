package example

import (
	"fmt"
	"testing"
)

func TestSlice01(t *testing.T) {

	//Slices的使用
	firstSlices := make([]string, 3)
	fmt.Println("first:", firstSlices)

	firstSlices[0] = "test1"
	firstSlices[1] = "test2"
	firstSlices[2] = "test3"

	fmt.Println("second:", firstSlices)

	fmt.Println("当前Slices的长度为：", len(firstSlices))

	//动态添加元素 长度改变
	firstSlices = append(firstSlices, "test01")
	fmt.Println("当前Slices的长度为：", len(firstSlices))
	fmt.Println(firstSlices)

	copySlices := make([]string, len(firstSlices))
	copy(copySlices, firstSlices)
	fmt.Println("复制后的数据：", copySlices)

	//将firstSlices[1] firstSlices[2]复制到copySlices2中
	copySlices2 := firstSlices[1:3]
	//从0开始复制到2
	copySlices3 := firstSlices[:3]
	//从1复制到最后
	copySlices4 := firstSlices[1:]
	fmt.Println(copySlices2)
	fmt.Println(copySlices3)
	fmt.Println(copySlices4)

	//快速定义Slices
	//不能指定长度 否则就编程Array了
	newSlices := []string{"a", "b", "c", "d"}
	fmt.Println(newSlices)

	//二维Slices
	doubleSlices := make([][]int, 4)
	for i := 0; i < 4; i++ {
		doubleSlices[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			doubleSlices[i][j] = j
		}
	}
	fmt.Println(doubleSlices)
}
