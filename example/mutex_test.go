package example

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	mutex := &sync.Mutex{}
	num := 100

	go func() {

		for {
			mutex.Lock()
			num = num + 1
			mutex.Unlock()
		}

	}()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 200)
		//保证在一次输出中num的值不会被改变
		mutex.Lock()
		fmt.Println("-->", i, num)
		fmt.Println("-->", i, num)
		mutex.Unlock()
	}

}