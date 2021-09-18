package example

import (
	"fmt"
	"testing"
)

func TestValues(t *testing.T) {
	fmt.Println("go" + "lang")
	fmt.Println("1+87=", 1+87)
	fmt.Println("6.2/3.1=", 7.2/3.1)
	fmt.Println(true && 1 > 0)
	fmt.Println(12 > 14 || false)
	fmt.Println(false || true)
	fmt.Println(!false)
}
