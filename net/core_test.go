package net

import (
	"os"
	"testing"
	"time"

	mpb "github.com/vbauerster/mpb/v5"
)

func TestHttpGetURLs(t *testing.T) {
	// urls := []string{"https://dldir1.qq.com/qqfile/qq/PCQQ9.1.6/25786/QQ9.1.6.25786.exe",
	//	"https://dldir1.qq.com/weixin/Windows/WeChatSetup.exe"}
	//urls := []string{"http://61.129.70.139:3030/api/viewfile/?path=/tmp/hiplot-It8jqd/a7yqy740q78.tar"}
	urls := []string{
		"http://61.129.70.139:3030/api/viewfile/?path=/tmp/hiplot-ox4t9X/cob38g8rgqv.tar",
		"http://61.129.70.139:3030/api/viewfile/?path=/tmp/hiplot-Jcw5fC/s1051on90mh.tar",
		"http://61.129.70.139:3030/api/viewfile/?path=/tmp/hiplot-It8jqd/a7yqy740q78.tar",
		"http://61.129.70.139:3030/api/view_user_file/?path=/62dcd780-0db6-11e9-855b-bb4c4b386613/data/readnew_1.fq.gz",
		"http://61.129.70.139:3030/api/view_user_file/?path=/62dcd780-0db6-11e9-855b-bb4c4b386613/upload/1587281476978-PTBP1-1.csv",
		"http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20131105.txt.idx.gz",
		"http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20170130.txt.gz",
		"http://www.openbioinformatics.org/annovar/download/hg19_clinvar_20170905.txt.gz"}
	//urls := []string{"https://github.com/openanno/bget"}
	urls = []string{"https://www.sciencedirect.com/science/article/pii/S2211034820301747/pdfft?md5=2fe69b9687518895596f0c2c1c55c8ed&pid=1-s2.0-S2211034820301747-main.pdf"}
	destDir := []string{}
	for range urls {
		destDir = append(destDir, "/cluster/home/ljf/repositories/github/openbiox/ligo/net/a")
	}
	param := &Params{}
	param.Retries = 5
	param.Engine = "default"
	param.Timeout = 35
	param.Overwrite = true
	param.Quiet = false
	param.Ignore = true
	param.TaskID = "test"
	param.RetSleepTime = 2
	param.Thread = 3
	param.ThreadQuery = 2
	param.LogDir = os.TempDir()
	param.Pbar = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	HTTPGetURLs(urls, destDir, param)
}
