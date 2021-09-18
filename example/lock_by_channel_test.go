package example

import (
	"fmt"
	"testing"
)

func TestLockByChan(t *testing.T) {
	num := 100
	lock := make(chan bool)

	go func() {

		for {
			//保证在num输出时不会进行写操作
			lock <- true
			num = num + 1
		}

	}()

	for i := 0; i < 10; i++ {
		fmt.Println("-->", i, num)
		fmt.Println("-->", i, num)
		<-lock

	}
}
