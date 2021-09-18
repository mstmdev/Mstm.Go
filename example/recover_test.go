package example

import (
	"errors"
	"fmt"
	"testing"
)

func TestRecoverError(t *testing.T) {

	defer func() {
		if v := recover(); v != nil {
			fmt.Println(v)
		}
	}()
	panic(errors.New("MyException"))
}
