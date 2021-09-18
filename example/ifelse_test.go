package example

import (
	"fmt"
	"testing"
)

func TestIfElse2(t *testing.T) {
	num1 := 100
	num2 := 200

	//if else使用
	if num1 > num2 {
		fmt.Println("the max is num1")
	} else {
		fmt.Println("the max is num2")
	}

}
