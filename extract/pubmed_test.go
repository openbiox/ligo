package extract

import (
	"strings"
	"testing"
)

func TestGetSimplePubmedFields(t *testing.T) {
	keywords := []string{"virus", "Hubei", "China", "outbreak", "Policy", "Heath", "ACE2", "[^0-9a-zA-Z]SARS[^0-9a-zA-Z]"}
	strings.Join(keywords, "|")
	//pubMedFields, _ := GetSimplePubmedFields("/cluster/home/ljf/repositories/github/openanno/COVID-19-connections/data/pubmed/xml/covid19.XML.tmp49.json", nil, &pat, true, true, true, 60)
	//dat, _ := json.MarshalIndent(pubMedFields, "", "    ")
	//fmt.Println(string(dat))
}
