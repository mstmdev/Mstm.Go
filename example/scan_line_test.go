package example

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestScanner2(t *testing.T) {
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			if scanner.Err() != nil {
				fmt.Println("we get err")
			}
		}
	}()

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
