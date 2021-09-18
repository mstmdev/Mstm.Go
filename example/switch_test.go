package example

import (
	"fmt"
	"testing"
	"time"
)

func TestSwitch(t *testing.T) {
	digit := 100
	fmt.Print("test", digit, digit, digit, "test2")

	switch digit {
	case 1:
		fmt.Println("to0 Lower")

	case 2:
		fmt.Println("this is best digit")
	}

	fmt.Println("---------------")
	fmt.Println(time.Now().Weekday())

	switch time.Now().Weekday() {
	case time.Monday:
		fmt.Println("周一")
	case time.Tuesday:
		fmt.Println("周二")
	case time.Wednesday:
		fmt.Println("周三")
	case time.Thursday:
		fmt.Println("周四")
	default:
		fmt.Println("Rest")
	}

	fmt.Println("当前时间为：", time.Now())

	//没有表达式的Switch
	switch {
	case 1 > 2:
		fmt.Println("1>2")
	case true:
		fmt.Println("true")
	//case 4 > 1:
	//	fmt.Println("4>1")
	default:
		fmt.Println("error")
	}

	test := 2
	switch test {
	case 1, 2:
		fmt.Println("1...2")
	case 3, 4:
		fmt.Println("3...4")
	default:
		fmt.Println("other")
	}

}
