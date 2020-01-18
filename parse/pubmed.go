package parse

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"

	clog "github.com/openbiox/ligo/log"
)

var log = clog.Logger

// PubmedXML convert Pubmed XML to json
func PubmedXML(xmlPaths *[]string, stdin *[]byte, outfn string, keywords *[]string, thread int) {
	if len(*xmlPaths) == 1 {
		thread = 1
	}
	if len(*stdin) > 0 {
		*xmlPaths = append(*xmlPaths, "ParsePubmedXMLStdin")
	}
	sem := make(chan bool, thread)

	//|os.O_APPEND
	var of *os.File
	if outfn == "" {
		of = os.Stdout
	} else {
		of, err := os.OpenFile(outfn, os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatal(err)
		}
		defer of.Close()
	}

	var err error
	for i, xmlPath := range *xmlPaths {
		sem <- true
		go func(xmlPath string, i int) {
			defer func() {
				<-sem
			}()
			var pubmed = PubmedArticleSet{}
			if xmlPath != "ParsePubmedXMLStdin" {
				xmlData, err := ioutil.ReadFile(xmlPath)
				if err != nil {
					log.Warnln(err)
				}
				err = xml.Unmarshal(xmlData, &pubmed)
			} else if xmlPath == "ParsePubmedXMLStdin" && len(*stdin) > 0 {
				err = xml.Unmarshal(*stdin, &pubmed)
			}
			if err != nil {
				log.Warnf("%v", err)
				return
			}
			jsonData, _ := json.MarshalIndent(pubmed.PubmedArticle, "", "  ")
			io.Copy(of, bytes.NewBuffer(jsonData))
		}(xmlPath, i)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

type PubmedArticleSet struct {
	XMLName       xml.Name        `xml:"PubmedArticleSet"`
	PubmedArticle []PubmedArticle `xml:"PubmedArticle"`
}

type PubmedArticle struct {
	MedlineCitation struct {
		Status string `xml:"Status,attr"`
		Owner  string `xml:"Owner,attr"`
		PMID   struct {
			Text    string `xml:",chardata"`
			Version string `xml:"Version,attr"`
		} `xml:"PMID"`
		DateRevised struct {
			Year  string `xml:"Year"`
			Month string `xml:"Month"`
			Day   string `xml:"Day"`
		} `xml:"DateRevised"`
		Article struct {
			PubModel string `xml:"PubModel,attr"`
			Journal  struct {
				ISSN struct {
					Text     string `xml:",chardata"`
					IssnType string `xml:"IssnType,attr"`
				} `xml:"ISSN"`
				JournalIssue struct {
					CitedMedium string `xml:"CitedMedium,attr"`
					Volume      string `xml:"Volume"`
					Issue       string `xml:"Issue"`
					PubDate     struct {
						Year  string `xml:"Year"`
						Month string `xml:"Month"`
					} `xml:"PubDate"`
				} `xml:"JournalIssue"`
				Title           string `xml:"Title"`
				ISOAbbreviation string `xml:"ISOAbbreviation"`
			} `xml:"Journal"`
			ArticleTitle struct {
				Text string   `xml:",chardata"`
				I    []string `xml:"i"`
				Sup  string   `xml:"sup"`
			} `xml:"ArticleTitle"`
			Pagination struct {
				MedlinePgn string `xml:"MedlinePgn"`
			} `xml:"Pagination"`
			ELocationID struct {
				Text    string `xml:",chardata"`
				EIdType string `xml:"EIdType,attr"`
				ValidYN string `xml:"ValidYN,attr"`
			} `xml:"ELocationID"`
			Abstract struct {
				AbstractText struct {
					Text string   `xml:",chardata"`
					I    []string `xml:"i"`
					B    string   `xml:"b"`
					Sub  string   `xml:"sub"`
					Sup  []string `xml:"sup"`
				} `xml:"AbstractText"`
			} `xml:"Abstract"`
			AuthorList struct {
				CompleteYN string `xml:"CompleteYN,attr"`
				Author     []struct {
					ValidYN         string `xml:"ValidYN,attr"`
					LastName        string `xml:"LastName"`
					ForeName        string `xml:"ForeName"`
					Initials        string `xml:"Initials"`
					AffiliationInfo []struct {
						Affiliation string `xml:"Affiliation"`
					} `xml:"AffiliationInfo"`
					Identifier struct {
						Text   string `xml:",chardata"`
						Source string `xml:"Source,attr"`
					} `xml:"Identifier"`
				} `xml:"Author"`
			} `xml:"AuthorList"`
			Language  string `xml:"Language"`
			GrantList struct {
				CompleteYN string `xml:"CompleteYN,attr"`
				Grant      []struct {
					GrantID string `xml:"GrantID"`
					Acronym string `xml:"Acronym"`
					Agency  string `xml:"Agency"`
					Country string `xml:"Country"`
				} `xml:"Grant"`
			} `xml:"GrantList"`
			PublicationTypeList struct {
				PublicationType []struct {
					Text string `xml:",chardata"`
					UI   string `xml:"UI,attr"`
				} `xml:"PublicationType"`
			} `xml:"PublicationTypeList"`
			ArticleDate struct {
				DateType string `xml:"DateType,attr"`
				Year     string `xml:"Year"`
				Month    string `xml:"Month"`
				Day      string `xml:"Day"`
			} `xml:"ArticleDate"`
		} `xml:"Article"`
		MedlineJournalInfo struct {
			Country     string `xml:"Country"`
			MedlineTA   string `xml:"MedlineTA"`
			NlmUniqueID string `xml:"NlmUniqueID"`
			ISSNLinking string `xml:"ISSNLinking"`
		} `xml:"MedlineJournalInfo"`
		CitationSubset string `xml:"CitationSubset"`
		KeywordList    struct {
			Owner   string `xml:"Owner,attr"`
			Keyword []struct {
				Text         string `xml:",chardata"`
				MajorTopicYN string `xml:"MajorTopicYN,attr"`
			} `xml:"Keyword"`
		} `xml:"KeywordList"`
	} `xml:"MedlineCitation"`
	PubmedData struct {
		History struct {
			PubMedPubDate []struct {
				PubStatus string `xml:"PubStatus,attr"`
				Year      string `xml:"Year"`
				Month     string `xml:"Month"`
				Day       string `xml:"Day"`
				Hour      string `xml:"Hour"`
				Minute    string `xml:"Minute"`
			} `xml:"PubMedPubDate"`
		} `xml:"History"`
		PublicationStatus string `xml:"PublicationStatus"`
		ArticleIdList     struct {
			ArticleId []ArticleId `xml:"ArticleId"`
		} `xml:"ArticleIdList"`
		ReferenceList struct {
			Reference struct {
				Citation      string `xml:"Citation"`
				ArticleIdList struct {
					ArticleId struct {
						Text   string `xml:",chardata"`
						IdType string `xml:"IdType,attr"`
					} `xml:"ArticleId"`
				} `xml:"ArticleIdList"`
			} `xml:"Reference"`
		} `xml:"ReferenceList"`
	} `xml:"PubmedData"`
}

type ArticleId struct {
	Text   string `xml:",chardata"`
	IdType string `xml:"IdType,attr"`
}
