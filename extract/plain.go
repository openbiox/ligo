package extract

import (
	"sort"

	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "github.com/jdkato/prose/v2"
	xurls "mvdan.cc/xurls/v2"
)

type PlainTextFields struct {
	Filename    string
	Correlation map[string][]string
	URLs        []string
	Keywords    []string
}

func GetPlainFields(filename string, dat *[]byte, keywordsPat *string, callCor bool, callURLs bool, thread int) (plain PlainTextFields, err error) {
	var doc *prose.Document
	var urls []string
	if dat == nil {
		dat = &[]byte{}
	}
	if filename != "" {
		if *dat, err = readDocFile(filename); err != nil {
			return plain, err
		}
	}
	datStr := string(*dat)
	if doc, err = prose.NewDocument(datStr); err != nil {
		return plain, err
	}
	if callURLs {
		urls = slice.DropSliceDup(xurls.Relaxed().FindAllString(datStr, -1))
	}
	key := stringo.StrExtract(datStr, *keywordsPat, -1)
	for k := range key {
		key[k] = formartKey(key[k])
	}
	key = slice.DropSliceDup(key)
	sort.Sort(sort.StringSlice(key))

	var corela map[string][]string
	if callCor {
		corela = getKeywordsCorleations(doc, keywordsPat, thread)
	}
	return PlainTextFields{
		filename,
		corela,
		urls,
		key,
	}, err
}
