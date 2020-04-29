package extract

import (
	"strings"
	"testing"
)

func TestGetPlainFields(t *testing.T) {
	keywords := []string{"virus", "Hubei", "China", "outbreak", "Policy", "Heath", ""}
	pat := strings.Join(keywords, "|")
	GetPlainFields("_demo/covid19.XML.tmp1", nil, &pat, true, true, 60)
}
