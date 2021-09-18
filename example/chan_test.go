package example

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	chan1 := make(chan int, 5)
	select {
	case chan1 <- 1:
	case chan1 <- 2:
	case chan1 <- 3:
	default:
		fmt.Println("default exec")
	}
	close(chan1)
	for item := range chan1 {
		fmt.Println(item)
	}
}
