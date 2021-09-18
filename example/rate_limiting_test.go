package example

import (
	"fmt"
	"testing"
	"time"
)

// TestBatch1 以500毫秒匀速执行
func TestBatch1(t *testing.T) {
	works := make(chan int, 10)
	for i := 1; i < 10; i++ {
		works <- i
	}
	close(works)
	timeout := time.Tick(time.Millisecond * 500)
	for work := range works {
		fmt.Println("batch1:", work, <-timeout)
	}
}

// TestBatch2 前三次同步执行  后面每隔1s执行
func TestBatch2(t *testing.T) {
	works := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		works <- i
	}
	close(works)

	timeout := make(chan time.Time, 10)
	for i := 1; i <= 3; i++ {
		timeout <- time.Now()
	}

	go func() {
		for v := range time.Tick(time.Second) {
			timeout <- v
		}
	}()

	for v := range works {
		<-timeout
		fmt.Println("batch2:", v, time.Now())
	}

}
