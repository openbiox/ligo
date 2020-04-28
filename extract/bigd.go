package extract

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

func GetBigdFields(filename string, dat *[]byte, keywordsPat *string, callCor bool, thread int) (articleFields []parse.BigdArticle, err error) {
	var dataJSON parse.BigdArticleList
	var lock sync.Mutex
	if dat == nil {
		dat = &[]byte{}
	}
	if filename != "" {
		if *dat, err = readDocFile(filename); err != nil {
			return nil, err
		}
	}
	if err := json.Unmarshal(*dat, &dataJSON); err != nil {
		return nil, err
	}
	sem := make(chan bool, thread)
	done := make(map[int]int)
	for _, article := range dataJSON.Data {
		sem <- true
		go func(article parse.BigdArticle) {
			defer func() {
				<-sem
			}()
			if done[article.ID] == 1 {
				return
			}
			article.Abst = stringo.StrReplaceAll(article.Abst, "\n  *", "\n")
			article.Abst = stringo.StrReplaceAll(article.Abst, "\n", "\\n")
			titleAbs := article.Title + "\n" + article.Abst
			article.URLs = xurls.Relaxed().FindAllString(titleAbs, -1)
			article.Keywords = stringo.StrExtract(titleAbs, *keywordsPat, -1)
			for k := range article.Keywords {
				article.Keywords[k] = formartKey(article.Keywords[k])
			}
			article.Keywords = slice.DropSliceDup(article.Keywords)
			sort.Sort(sort.StringSlice(article.Keywords))

			doc, err := prose.NewDocument(titleAbs)
			if callCor {
				article.Correlation = getKeywordsCorleations(doc, keywordsPat, 10)
			}
			if err != nil {
				log.Warn(err)
			}
			lock.Lock()
			articleFields = append(articleFields, article)
			done[article.ID] = 1
			lock.Unlock()
		}(article)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return articleFields, err
}
