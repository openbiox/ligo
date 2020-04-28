package parse

type BigdArticle struct {
	ID                    int                 `json:"id"`
	Title                 string              `json:"title"`
	Author                string              `json:"author"`
	Abst                  string              `json:"abst"`
	Journal               string              `json:"journal"`
	JournalAbbr           string              `json:"journalAbbr"`
	Volume                interface{}         `json:"volume"`
	Issue                 interface{}         `json:"issue"`
	PubYear               int                 `json:"pubYear"`
	PubMonth              interface{}         `json:"pubMonth"`
	PubDay                int                 `json:"pubDay"`
	Pmid                  interface{}         `json:"pmid"`
	Doi                   string              `json:"doi"`
	Citation              int                 `json:"citation"`
	Type                  interface{}         `json:"type"`
	CreateTime            string              `json:"createTime"`
	LastModified          string              `json:"lastModified"`
	PubDate               string              `json:"pubDate"`
	Keyword               interface{}         `json:"keyword"`
	Source                string              `json:"source"`
	Pprid                 interface{}         `json:"pprid"`
	Affiliation           interface{}         `json:"affiliation"`
	Exclude               int                 `json:"exclude"`
	PublicationReviewList []interface{}       `json:"publicationReviewList"`
	URLs                  []string            `json:"urls"`
	Correlation           map[string][]string `json:"correlation"`
	Keywords              []string            `json:"keywords"`
}

type BigdArticleList struct {
	Draw            int           `json:"draw"`
	RecordsTotal    int           `json:"recordsTotal"`
	RecordsFiltered int           `json:"recordsFiltered"`
	Data            []BigdArticle `json:"data"`
	Error           interface{}   `json:"error"`
}
