package file

import (
	"fmt"
	"testing"
)

func TestLineCounterByNameSlice(t *testing.T) {
	f, _, _ := LineCounterByNameSlice([]string{"./count.go"})
	fmt.Println(f)
}
