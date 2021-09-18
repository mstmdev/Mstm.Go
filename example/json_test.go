package example

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson2(t *testing.T) {
	//Json序列化
	p := Person{Name: "张三", Age: 13}
	data, _ := json.Marshal(p)
	fmt.Println(string(data))

	//Json反序列化
	bytes := []byte(`{"NameString":"王五","AgeInt":99}`)
	json.Unmarshal(bytes, &p)
	fmt.Println(p)
}

type Person struct {
	Name string `json:"NameString"`
	Age  int    `json:"AgeInt"`
}
