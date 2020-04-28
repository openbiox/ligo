package extract

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

// PubmedFields defines extracted Pubmed fields
type PubmedFields struct {
	Pmid, Doi, Title, Abs, Journal, Issue, Volume, Date, Issn string
	Author                                                    []string
	Affiliation                                               []string
	Correlation                                               map[string][]string
	URLs                                                      []string
	Keywords                                                  []string
}

type SimplePubmedFieldsOpt struct {
	Title, Abs, Journal, Issue, Volume, Date, Issn, Author, Affiliation, URLs, Keywords bool
}

func GetSimplePubmedFields(filename string, dat *[]byte, keywordsPat *string, callCor bool, thread int) (pubMedFields []PubmedFields, err error) {
	var pubmedJSON []parse.PubmedArticleJSON
	var lock sync.Mutex
	if dat == nil {
		dat = &[]byte{}
	}
	if filename != "" {
		if *dat, err = readDocFile(filename); err != nil {
			return nil, err
		}
	}
	if err := json.Unmarshal(*dat, &pubmedJSON); err != nil {
		return nil, err
	}
	sem := make(chan bool, thread)
	done := make(map[string]int)

	for _, article := range pubmedJSON {
		sem <- true
		go func(article parse.PubmedArticleJSON) {
			defer func() {
				<-sem
			}()
			year := article.MedlineCitation.Article.ArticleDate.Year
			month := article.MedlineCitation.Article.ArticleDate.Month
			day := article.MedlineCitation.Article.ArticleDate.Day
			date := fmt.Sprintf("%s/%s/%s", year, month, day)
			var pmid, doi, abs string
			for _, v := range article.PubmedData.ArticleIDList.ArticleID {
				if v.IDType == "pubmed" {
					pmid = v.Text
				} else if v.IDType == "doi" {
					doi = v.Text
				}
			}
			if done[pmid] == 1 {
				return
			}
			abs = stringo.StrReplaceAll(article.MedlineCitation.Article.Abstract.AbstractText.Text, "\n  *", "\n")
			abs = stringo.StrReplaceAll(abs, "(<[/]AbstractText.*>)|(^\n)|(\n$)", "")
			abs = stringo.StrReplaceAll(abs, "\n", "\\n")
			title := article.MedlineCitation.Article.ArticleTitle.Text
			title = stringo.StrReplaceAll(title, "\n", "\\n")
			titleAbs := title + "\n" + abs
			urls := xurls.Relaxed().FindAllString(titleAbs, -1)
			key := stringo.StrExtract(titleAbs, *keywordsPat, -1)
			for k := range key {
				key[k] = formartKey(key[k])
			}
			key = slice.DropSliceDup(key)
			sort.Sort(sort.StringSlice(key))

			doc, err := prose.NewDocument(titleAbs)
			var corela map[string][]string
			if callCor {
				corela = getKeywordsCorleations(doc, keywordsPat, 10)
			}
			if err != nil {
				log.Warn(err)
			}
			var author, affiliation []string
			for _, v := range article.MedlineCitation.Article.AuthorList.Author {
				author = append(author, v.ForeName+" "+v.LastName)
				affiliationTmp := ""
				for _, j := range v.AffiliationInfo {
					if affiliationTmp == "" {
						affiliationTmp = j.Affiliation
					} else {
						affiliationTmp = affiliationTmp + "; " + j.Affiliation
					}
				}
				affiliation = append(affiliation, affiliationTmp)
			}
			lock.Lock()
			pubMedFields = append(pubMedFields, PubmedFields{
				Pmid:        pmid,
				Doi:         doi,
				Title:       title,
				Abs:         abs,
				Journal:     article.MedlineCitation.Article.Journal.ISOAbbreviation,
				Issn:        article.MedlineCitation.Article.Journal.ISSN.Text,
				Date:        date,
				Issue:       article.MedlineCitation.Article.Journal.JournalIssue.Issue,
				Volume:      article.MedlineCitation.Article.Journal.JournalIssue.Volume,
				Author:      author,
				Affiliation: affiliation,
				Correlation: corela,
				URLs:        urls,
				Keywords:    key,
			})
			done[pmid] = 1
			lock.Unlock()
		}(article)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return pubMedFields, err
}
