package extract

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "github.com/jdkato/prose/v2"
	xurls "mvdan.cc/xurls/v2"
)

// SraFields defines extracted Sra fields
type SraFields struct {
	Title        string
	StudyTitle   string
	Abstract     string
	Type         string
	SRX          string
	SRA          string
	SRAFile      parse.SRAFileJSON
	Correlation  map[string][]string
	AbstractURLs []string
	Keywords     []string
}

func GetSimpleSraFields(filename string, dat *[]byte, keywordsPat *string, callCor bool, callURLs bool, keepAbs bool, thread int) (sraFields []SraFields, err error) {
	var sraJSON []parse.ExperimentPkgJSON
	var lock sync.Mutex
	if dat == nil {
		dat = &[]byte{}
	}
	if filename != "" {
		if *dat, err = readDocFile(filename); err != nil {
			return nil, err
		}
	}
	if err := json.Unmarshal(*dat, &sraJSON); err != nil {
		return nil, err
	}
	sem := make(chan bool, thread)
	done := make(map[string]int)
	for _, sra := range sraJSON {
		sem <- true
		go func(sra parse.ExperimentPkgJSON) {
			defer func() {
				<-sem
			}()
			var urls []string
			titleAbs := sra.EXPERIMENT.TITLE + "\n" + sra.STUDY.DESCRIPTOR.STUDYTITLE +
				"\n" + sra.STUDY.DESCRIPTOR.STUDYABSTRACT
			doc, err := prose.NewDocument(titleAbs)
			if done[sra.EXPERIMENT.TITLE+sra.STUDY.DESCRIPTOR.STUDYTITLE] == 1 {
				return
			}
			if callURLs {
				urls = xurls.Relaxed().FindAllString(titleAbs, -1)
			}
			abs := sra.STUDY.DESCRIPTOR.STUDYABSTRACT
			if !keepAbs {
				abs = ""
			}
			key := stringo.StrExtract(titleAbs, *keywordsPat, -1)
			for k := range key {
				key[k] = formartKey(key[k])
			}
			key = slice.DropSliceDup(key)
			sort.Sort(sort.StringSlice(key))
			var corela map[string][]string
			if callCor {
				corela = getKeywordsCorleations(doc, keywordsPat, 10)
			}
			if err != nil {
				log.Warn(err)
			}
			lock.Lock()
			sraFields = append(sraFields, SraFields{
				Title:        sra.EXPERIMENT.TITLE,
				StudyTitle:   sra.STUDY.DESCRIPTOR.STUDYTITLE,
				Abstract:     abs,
				Type:         sra.STUDY.DESCRIPTOR.STUDYTYPE.ExistingStudyType,
				SRX:          sra.EXPERIMENT.Accession,
				SRA:          sra.RUNSET.RUN.Accession,
				SRAFile:      sra.RUNSET.RUN.SRAFiles.SRAFile,
				Correlation:  corela,
				AbstractURLs: urls,
				Keywords:     key,
			})
			done[sra.EXPERIMENT.TITLE+sra.STUDY.DESCRIPTOR.STUDYTITLE] = 1
			lock.Unlock()
		}(sra)
	}
	return sraFields, err
}
