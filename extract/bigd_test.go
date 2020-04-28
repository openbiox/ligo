package extract

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestGetBigdFields(t *testing.T) {
	keywords := []string{"virus", "Hubei", "China", "outbreak", "Policy", "Heath", "ACE2", "[^0-9a-zA-Z]SARS[^0-9a-zA-Z]"}
	pat := strings.Join(keywords, "|")
	articleFields, err := GetBigdFields("_demo/bigd.json", nil, &pat, true, 60)
	dat, _ := json.MarshalIndent(articleFields, "", "    ")
	fmt.Println(err)
	fmt.Println(string(dat))
}
