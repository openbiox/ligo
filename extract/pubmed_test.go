package extract

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestGetSimplePubmedFields(t *testing.T) {
	var pubMedFields []PubmedFields
	keywords := []string{"virus", "Hubei", "China", "outbreak", "Policy", "Heath", "ACE2", "[^0-9a-zA-Z]SARS[^0-9a-zA-Z]"}
	pat := strings.Join(keywords, "|")
	pubMedFields, _ = GetSimplePubmedFields("/cluster/home/ljf/repositories/github/openanno/covid-19/xml/covid19.XML.tmp49.json", nil, &pat, true, 60)
	dat, _ := json.MarshalIndent(pubMedFields, "", "    ")
	fmt.Println(string(dat))
}
