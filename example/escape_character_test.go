package example

import (
	"fmt"
	"testing"
)

func TestEscape(t *testing.T) {

	//转义字符说明
	//系统警告声音\a
	// fmt.Println("\a")
	//退格符\b
	fmt.Println("abcd\b\befgh")
	//\b放在末尾无效
	fmt.Println("abcd\b\b")
	//换页符  \f
	fmt.Println("abc\f123")
	//换行符 \n
	fmt.Printf("换行符\n")
	fmt.Printf("换行符\n")
	//回车符 \r
	fmt.Printf("回车符\r")
	fmt.Printf("回车符\r")

	fmt.Println("水平制表符\taaaaaa")
	fmt.Println("水平制表符\tbbbbbb")
	fmt.Println("水平制表符\tcccccc")

	fmt.Println("垂直制表符\vdddddd")
	fmt.Println("垂直制表符\veeeeee")
	fmt.Println("垂直制表符\vffffff")

	//特殊字符转义
	fmt.Println("\\---\"")

}
