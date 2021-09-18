package example

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i++
	}

	//for循环
	for j := 9; j >= 0; j-- {
		fmt.Println(j)
	}

	fmt.Println("-------------------------")
	tag := 6
	for {
		fmt.Println("for loop")
		tag--
		if tag <= 0 {
			break
		}
	}
}
