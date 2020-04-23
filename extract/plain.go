package extract

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

type PlainTextFields struct {
	Filename    string
	Correlation map[string]string
	URLs        []string
	Keywords    []string
}

func GetPlainFields(filename string, dat *[]byte, keywords *[]string, callCor bool) (plain *PlainTextFields, err error) {
	var of *os.File
	var dat2 []byte
	var doc *prose.Document
	if filename != "" {
		if filename, err = filepath.Abs(filename); err != nil {
			return nil, err
		}
		if of, err = os.Open(filename); err != nil {
			return nil, err
		}
		if dat2, err = ioutil.ReadAll(of); err != nil {
			return nil, err
		}
		dat = &dat2
	}
	datStr := string(*dat)
	if doc, err = prose.NewDocument(datStr); err != nil {
		return nil, err
	}
	urls := slice.DropSliceDup(xurls.Relaxed().FindAllString(datStr, -1))
	keywordsPat := strings.Join(*keywords, "|")
	key := stringo.StrExtract(datStr, keywordsPat, 1000000)
	key = slice.DropSliceDup(key)
	cor := make(map[string]string)
	if callCor {
		for _, sent := range doc.Sentences() {
			kStr := stringo.StrExtract(sent.Text, keywordsPat, 1000000)
			kStr = slice.DropSliceDup(kStr)
			if len(kStr) >= 2 {
				cor[strings.Join(kStr, "+")] = sent.Text
			}
		}
	}
	return &PlainTextFields{
		filename,
		cor,
		urls,
		key,
	}, nil
}
