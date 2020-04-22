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
	urls := []string{"https://github.com/openanno/bget"}
	destDir := []string{os.TempDir(), os.TempDir(), os.TempDir()}
	param := &Params{}
	param.Retries = 5
	param.Engine = "git"
	param.Timeout = 35
	param.Overwrite = true
	param.Ignore = true
	param.TaskID = "test"
	param.RetSleepTime = 2
	param.Thread = 2
	param.ThreadQuery = 3
	param.LogDir = os.TempDir()
	param.Pbar = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	HTTPGetURLs(urls, destDir, param)
}
