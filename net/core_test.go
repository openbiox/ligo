package net

import (
	"os"
	"testing"
	"time"

	mpb "github.com/vbauerster/mpb/v4"
)

func TestHttpGetURLs(t *testing.T) {
	urls := []string{"https://dldir1.qq.com/weixin/Windows/WeChatSetup.exe",
		"https://dldir1.qq.com/qqfile/qq/PCQQ9.1.6/25786/QQ9.1.6.25786.exe"}
	destDir := []string{os.TempDir(), os.TempDir(), os.TempDir()}
	param := &Params{}
	param.Retries = 5
	param.Engine = "default"
	param.Timeout = 35
	param.Overwrite = false
	param.Ignore = true
	param.TaskID = "test"
	param.Thread = 2
	param.ThreadQuery = 4
	param.LogDir = os.TempDir()
	param.Pbar = mpb.New(
		mpb.WithWidth(45),
		mpb.WithRefreshRate(180*time.Millisecond),
	)
	HTTPGetURLs(urls, destDir, param)
}
