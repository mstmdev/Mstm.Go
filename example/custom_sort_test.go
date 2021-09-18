package example

import (
	"fmt"
	"sort"
	"testing"
)

func TestCustomSort(t *testing.T) {
	array := newStringArray{"cs", "aaaa", "bb", "aaaaaaa"}
	sort.Sort(array)
	fmt.Println(array)
}

type newStringArray []string

func (str newStringArray) Less(l, r int) bool {
	return len(str[l]) < len(str[r])
}

func (str newStringArray) Len() int {
	return len(str)
}

func (str newStringArray) Swap(l, r int) {
	str[r], str[l] = str[l], str[r]
}
