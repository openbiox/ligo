package extract

import (
	"fmt"
	"strings"

	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

var log = clog.Logger

// PubmedFields defines extracted Pubmed fields
type PubmedFields struct {
	Pmid, Doi, Title, Abs, Journal, Issue, Volume, Date, Issn *string
	Author                                                    *[]string
	Affiliation                                               *[]string
	Correlation                                               *map[string]string
	URLs                                                      *[]string
	Keywords                                                  *[]string
}

func GetSimplePubmedFields(keywords *[]string, article *parse.PubmedArticleJSON, callCor bool) *PubmedFields {
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
	abs = stringo.StrReplaceAll(article.MedlineCitation.Article.Abstract.AbstractText.Text, "\n  *", "\n")
	abs = stringo.StrReplaceAll(abs, "(<[/]AbstractText.*>)|(^\n)|(\n$)", "")
	titleAbs := article.MedlineCitation.Article.ArticleTitle.Text + "\n" + abs
	urls := xurls.Relaxed().FindAllString(titleAbs, -1)
	keywordsPat := strings.Join(*keywords, "|")
	key := stringo.StrExtract(titleAbs, keywordsPat, 1000000)
	key = slice.DropSliceDup(key)

	doc, err := prose.NewDocument(titleAbs)
	corela := make(map[string]string)
	if len(key) >= 2 && callCor {
		getKeywordsCorleations(doc, &keywordsPat, &corela)
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
	return &PubmedFields{
		Pmid:        &pmid,
		Doi:         &doi,
		Title:       &article.MedlineCitation.Article.ArticleTitle.Text,
		Abs:         &abs,
		Journal:     &article.MedlineCitation.Article.Journal.ISOAbbreviation,
		Issn:        &article.MedlineCitation.Article.Journal.ISSN.Text,
		Date:        &date,
		Issue:       &article.MedlineCitation.Article.Journal.JournalIssue.Issue,
		Volume:      &article.MedlineCitation.Article.Journal.JournalIssue.Volume,
		Author:      &author,
		Affiliation: &affiliation,
		Correlation: &corela,
		URLs:        &urls,
		Keywords:    &key,
	}
}

func getKeywordsCorleations(doc *prose.Document, keywordsPat *string, corela *map[string]string) {
	for _, sent := range doc.Sentences() {
		kStr := stringo.StrExtract(sent.Text, *keywordsPat, 1000000)
		kStr = slice.DropSliceDup(kStr)
		if len(kStr) >= 2 {
			(*corela)[strings.Join(kStr, "+")] = sent.Text
		}
	}
}
