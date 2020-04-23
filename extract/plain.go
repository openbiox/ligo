package extract

import (
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

type PlainTextFields struct {
	Filename    string
	Correlation map[string][]string
	URLs        []string
	Keywords    []string
}

func GetPlainFields(filename string, dat *[]byte, keywordsPat *string, callCor bool, thread int) (plain PlainTextFields, err error) {
	var doc *prose.Document
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
	urls := slice.DropSliceDup(xurls.Relaxed().FindAllString(datStr, -1))
	key := stringo.StrExtract(datStr, *keywordsPat, -1)
	for k := range key {
		key[k] = formartKey(key[k])
	}
	key = slice.DropSliceDup(key)
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
