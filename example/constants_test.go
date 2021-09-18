package example

import (
	"fmt"
	"math"
	"testing"
)

const outConst string = "constants"

func TestConstants(t *testing.T) {

	const strConst = "first Const"
	fmt.Println(outConst)
	fmt.Println(strConst)

	//常量  不可修改
	//strConst = "test"

	const digit = 3e20 / 5000
	fmt.Println(digit)

	//类型转换
	fmt.Println(int64(digit))

	//使用math
	fmt.Println(math.Sin(digit))

}
