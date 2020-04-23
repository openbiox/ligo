package extract

import (
	"fmt"
	"testing"
)

func TestGetPlainFields(t *testing.T) {
	keywords := []string{"virus", "Hubei", "outbreak"}
	dat, _ := GetPlainFields("_demo/plain", nil, &keywords)
	fmt.Println(dat)
}
