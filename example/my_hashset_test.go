package example

import (
	"fmt"
	"testing"
)

func TestHashSet(t *testing.T) {
	hash := HashSet{}
	hash.Add("item1")
	hash.Add("item2")
	// hash.Delete("item1")
	// fmt.Println(hash.Contains("item2"))
	// fmt.Println(hash.values)
	fmt.Println(hash.Elements())
}

////////////////////////////
type HashSet struct {
	values map[interface{}]bool
}

//添加一个元素
func (obj *HashSet) Add(item interface{}) {
	if obj.values == nil {
		obj.values = make(map[interface{}]bool)
	}
	obj.values[item] = true
}

//清空所有元素
func (obj *HashSet) Clear() {
	obj.values = make(map[interface{}]bool)
}

//删除一个元素
func (obj *HashSet) Delete(item interface{}) {
	delete(obj.values, item)
}

//是否包含指定的元素
func (obj *HashSet) Contains(item interface{}) bool {
	if obj.values == nil || len(obj.values) == 0 {
		return false
	}
	return obj.values[item]
}

func (obj *HashSet) Elements() []interface{} {

	result := make([]interface{}, 0)
	for key := range obj.values {
		result = append(result, key)
	}
	return result
}
