package example

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort01(t *testing.T) {
	//数字排序
	nums := []int{3, 41, 34, 577, 4, 26, 78, 5, 64}
	sort.Ints(nums)
	fmt.Println(nums)

	//字符串排序
	strs := []string{"这", "a", "c", "b", "ac", "ab", "+"}
	sort.Strings(strs)
	fmt.Println(strs)
}
