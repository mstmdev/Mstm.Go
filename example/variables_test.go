package example

import (
	"fmt"
	"testing"
)

func TestVars(t *testing.T) {
	//定义一个字符串
	var str string = "test String"
	fmt.Println(str)

	//第一两个数字
	var numLeft, numRight int = 3, 4
	fmt.Println(numLeft + numRight)

	//定义布尔变量
	var isOpen = true
	fmt.Println(isOpen)

	var isClose bool
	fmt.Println(isClose)

	var defaultInt int
	fmt.Println(defaultInt)

	//快速定义变量  使用字面量进行快速定义
	quickString := "quick String Variable"
	fmt.Println(quickString)

	quickInt := 99
	fmt.Println(quickInt)
}
