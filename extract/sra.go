package parse

import (
	"strings"

	"github.com/openbiox/ligo/parse"
	"github.com/openbiox/ligo/slice"
	"github.com/openbiox/ligo/stringo"
	prose "gopkg.in/jdkato/prose.v2"
	xurls "mvdan.cc/xurls/v2"
)

// SraFields defines extracted Sra fields
type SraFields struct {
	Title        *string
	StudyTitle   *string
	Abstract     *string
	Type         *string
	SRX          *string
	SRA          *string
	SRAFile      *parse.SRAFile
	Correlation  *map[string]string
	AbstractURLs *[]string
	Keywords     *[]string
}

func GetSimpleSraFields(keywords *[]string, sra *parse.ExperimentPkg, callCor bool, done map[string]int) *SraFields {
	titleAbs := sra.EXPERIMENT.TITLE + "\n" + sra.STUDY.DESCRIPTOR.STUDYTITLE +
		"\n" + sra.STUDY.DESCRIPTOR.STUDYABSTRACT
	doc, err := prose.NewDocument(titleAbs)
	if done[sra.EXPERIMENT.TITLE+sra.STUDY.DESCRIPTOR.STUDYTITLE] == 1 {
		return &SraFields{
			Title:      &sra.EXPERIMENT.TITLE,
			StudyTitle: &sra.STUDY.DESCRIPTOR.STUDYTITLE,
			Type:       &sra.STUDY.DESCRIPTOR.STUDYTYPE.ExistingStudyType,
			SRX:        &sra.EXPERIMENT.Accession,
			SRA:        &sra.RUNSET.RUN.Accession,
			SRAFile:    &sra.RUNSET.RUN.SRAFiles.SRAFile,
		}
	}
	corela := make(map[string]string)
	urls := xurls.Relaxed().FindAllString(titleAbs, -1)
	keywordsPat := strings.Join(*keywords, "|")
	key := stringo.StrExtract(titleAbs, keywordsPat, 1000000)
	key = slice.DropSliceDup(key)
	if len(key) >= 2 && callCor {
		getKeywordsCorleations(doc, &keywordsPat, &corela)
	}
	if err != nil {
		log.Warn(err)
	}
	return &SraFields{
		Title:        &sra.EXPERIMENT.TITLE,
		StudyTitle:   &sra.STUDY.DESCRIPTOR.STUDYTITLE,
		Abstract:     &sra.STUDY.DESCRIPTOR.STUDYABSTRACT,
		Type:         &sra.STUDY.DESCRIPTOR.STUDYTYPE.ExistingStudyType,
		SRX:          &sra.EXPERIMENT.Accession,
		SRA:          &sra.RUNSET.RUN.Accession,
		SRAFile:      &sra.RUNSET.RUN.SRAFiles.SRAFile,
		Correlation:  &corela,
		AbstractURLs: &urls,
		Keywords:     &key,
	}
}
