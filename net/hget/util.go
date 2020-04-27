package hget

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"github.com/mattn/go-isatty"
)

func FatalCheck(err error) {
	if err != nil {
		Errorf("%v", err)
		panic(err)
	}
}

func FilterIPV4(ips []net.IP) []string {
	var ret = make([]string, 0)
	for _, ip := range ips {
		if ip.To4() != nil {
			ret = append(ret, ip.String())
		}
	}
	return ret
}

func MkdirIfNotExist(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0700); err != nil {
			return err
		}
	}
	return nil
}

func ExistDir(folder string) bool {
	_, err := os.Stat(folder)
	return err == nil
}

func DisplayProgressBar() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) && displayProgress
}

func SciencedirectassetsRed(url string, req *http.Request, client *http.Client) {
	if strings.Contains(url, "www.sciencedirect.com/science/article/pii") &&
		strings.Contains(url, "-main.pdf") {
		resp, err := http.Head(url)
		if err != nil {
			return
		}
		ctype := resp.Header.Get("Content-Type")
		if strings.Contains(ctype, "text/html") {
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			resp, err = client.Do(req)
			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				return
			}
			var wg sync.WaitGroup
			wg.Add(1)
			doc.Find("noscript").SetHtml(doc.Find("noscript").Text()).Find("#redirect-message p a").Each(func(i int, selection *goquery.Selection) {
				url, _ = selection.Attr("href")
				req2, _ := http.NewRequest("GET", url, nil)
				*req = *req2
				req.Header.Set("Connection", "keep-alive")
				req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
				wg.Done()
			})
			wg.Wait()
		}
	}
}
