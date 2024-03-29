package extract

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "github.com/jdkato/prose/v2"
)

var log = clog.Logger

func readDocFile(filename string) (dat []byte, err error) {
	var of *os.File
	if filename, err = filepath.Abs(filename); err != nil {
		return nil, err
	}
	if of, err = os.Open(filename); err != nil {
		return nil, err
	}
	if dat, err = io.ReadAll(of); err != nil {
		return nil, err
	}
	defer of.Close()
	return dat, nil
}

func getKeywordsCorleations(doc *prose.Document, keywordsPat *string, sentThread int) map[string][]string {
	corela := make(map[string][]string)
	sem := make(chan bool, sentThread)
	var lock sync.Mutex
	for _, sent := range doc.Sentences() {
		sem <- true
		go func(sent prose.Sentence) {
			defer func() {
				<-sem
			}()
			kStr := stringo.StrExtract(sent.Text, *keywordsPat, -1)
			sort.Sort(sort.StringSlice(kStr))
			for k := range kStr {
				kStr[k] = formartKey(kStr[k])
			}
			kStr = slice.DropSliceDup(kStr)
			if len(kStr) >= 2 {
				key := strings.Join(kStr, "+")
				key = stringo.StrReplaceAll(key, " [+]", "+")
				lock.Lock()
				corela[key] = append(corela[key], formatSentText(sent.Text))
				lock.Unlock()
			}
		}(sent)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return corela
}

func removeDuplicatesAndEmpty(a []string) (ret []string) {
	sort.Sort(sort.StringSlice(a))
	alen := len(a)
	for i := 0; i < alen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return ret
}

func formartKey(key string) string {
	key = stringo.StrRemoveAll(key, "\n|\n")
	key = stringo.StrRemoveAll(key, "\\n|\\n")
	key = stringo.StrReplaceAll(key, "\t|\\t", " ")
	key = stringo.StrRemoveAll(key, "^[)]|[(]$")
	key = stringo.StrRemoveAll(key, "^[(]|[)]$")
	key = stringo.StrRemoveAll(key, "^[-]|[-]$")
	key = stringo.StrRemoveAll(key, "^[% ,./:=]|[% ,./:=]$")
	key = stringo.StrRemoveAll(key, "^[[]|[]]$")
	key = stringo.StrRemoveAll(key, "^\"|\"$")
	key = stringo.StrRemoveAll(key, "^['（; ?]|[; ?）']$")
	key = stringo.StrRemoveAll(key, "^[“_-><”&~!@#$%*+-`•]|[•“_-><”&~!@#$%*+-`]$")
	key = stringo.StrRemoveAll(key, "^-|-$")
	key = stringo.StrRemoveAll(key, "^‐|‐$")

	return key
}

func formatSentText(text string) string {
	text = stringo.StrReplaceAll(text, "\n\n", " ")
	text = stringo.StrReplaceAll(text, "\n", " ")
	return text
}
