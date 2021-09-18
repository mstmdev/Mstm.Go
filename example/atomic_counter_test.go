package example

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {
	var ops uint64 = 0
	for i := 0; i < 50; i++ {
		go func() {
			for {
				atomic.AddUint64(&ops, 1)
				//使其他的goroutine可以对pos进行操作  不会被阻塞
				runtime.Gosched()
			}
		}()
	}
	time.Sleep(time.Second * 3)
	//获取ops当前的值，防止在操作时ops被改变  导致数据不一致
	result := atomic.LoadUint64(&ops)
	fmt.Println(result)
}
