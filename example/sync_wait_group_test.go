package example

import (
	"fmt"
	"sync"
	"testing"
	_ "time"
)

var waitGroup sync.WaitGroup

// TestExec
//普通的方式执行三个go routine
//所有输出语句都来不及执行
func TestExec(t *testing.T) {
	go process1()
	go process2()
	go process3()
}

// TestExec2
//使用WaitGroup对象让程序等待三个goroutine执行完毕
func TestExec2(t *testing.T) {
	waitGroup.Add(3)
	go process1()
	go process2()
	go process3()
	waitGroup.Wait()
}

func process1() {
	fmt.Println("process1")
	waitGroup.Done()
}

func process2() {
	fmt.Println("process2")
	waitGroup.Done()
}

func process3() {
	fmt.Println("process3")
	waitGroup.Done()
}
