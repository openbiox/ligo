package net

import (
	"os"
	"testing"
	"time"

	mpb "github.com/vbauerster/mpb/v5"
)

func TestHttpGetURLs(t *testing.T) {
	urls := []string{"http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20150330.txt.gz", "http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20170130.txt.gz",
		"http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20180603.txt.gz"}
	destDir := []string{os.TempDir(), os.TempDir(), os.TempDir()}
	param := &Params{}
	param.Retries = 5
	param.Engine = "default"
	param.Timeout = 35
	param.Overwrite = false
	param.Ignore = true
	param.TaskID = "test"
	param.Thread = 2
	param.ThreadQuery = 3
	param.LogDir = os.TempDir()
	param.Pbar = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	HTTPGetURLs(urls, destDir, param)
}
