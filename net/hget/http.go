package hget

import (
	"fmt"
	"io"
	"net"
	"net/http"
	stdurl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	mpb "github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var (
	acceptRangeHeader   = "Accept-Ranges"
	contentLengthHeader = "Content-Length"
)

type HttpDownloader struct {
	url       string
	file      string
	par       int64
	len       int64
	ips       []string
	skipTls   bool
	parts     []Part
	resumable bool
}

func NewHttpDownloader(url string, par int, skipTls bool, dest string) (*HttpDownloader, error) {
	var resumable = true
	if !strings.Contains(url, "://") {
		url = "http://" + url
	}
	parsed, _ := stdurl.Parse(url)

	ips, err := net.LookupIP(parsed.Host)

	ipstr := FilterIPV4(ips)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	gCurCookies = gCurCookieJar.Cookies(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Header.Get(acceptRangeHeader) == "" {
		//fallback to par = 1
		par = 1
	}

	//get download range
	clen := resp.Header.Get(contentLengthHeader)
	if clen == "" {
		clen = "1" //set 1 because of progress bar not accept 0 length
		par = 1
		resumable = false
	}

	len, err := strconv.ParseInt(clen, 10, 64)
	FatalCheck(err)

	file := dest
	ret := new(HttpDownloader)
	ret.url = url
	ret.file = file
	ret.par = int64(par)
	ret.len = len
	ret.ips = ipstr
	ret.skipTls = skipTls
	ret.parts = partCalculate(int64(par), len, url)
	ret.resumable = resumable

	return ret, nil
}

func partCalculate(par int64, len int64, url string) []Part {
	ret := make([]Part, 0)
	for j := int64(0); j < par; j++ {
		from := (len / par) * j
		var to int64
		if j < par-1 {
			to = (len/par)*(j+1) - 1
		} else {
			to = len
		}

		file := filepath.Base(url)
		folder := FolderOf(url)
		if err := MkdirIfNotExist(folder); err != nil {
			Errorf("%v", err)
			os.Exit(1)
		}

		fname := fmt.Sprintf("%s.part%d", file, j)
		path := filepath.Join(folder, fname) // ~/.hget/download-file-name/part-name
		ret = append(ret, Part{Url: url, Path: path, RangeFrom: from, RangeTo: to})
	}
	return ret
}

func (d *HttpDownloader) Do(doneChan chan bool, fileChan chan string, errorChan chan error, interruptChan chan bool, stateSaveChan chan Part, bars []*mpb.Bar) {

	var ws sync.WaitGroup

	if DisplayProgressBar() {
		for i, part := range d.parts {
			prefixStr := filepath.Base(d.file)
			prefixStrLen := utf8.RuneCountInString(prefixStr)
			if prefixStrLen >= 18 {
				prefixStr = fmt.Sprintf("%-22s\t", fmt.Sprintf("%s..-%d", prefixStr[0:17], i))
			} else {
				prefixStr = fmt.Sprintf("%-22s\t", fmt.Sprintf("%s-%d", prefixStr, i))
			}
			size := part.RangeTo - part.RangeFrom
			if !d.resumable {
				size = -1
			}
			newbar := pbg.AddBar(size,
				mpb.BarNoPop(),
				mpb.BarStyle("[=>-|"),
				mpb.PrependDecorators(
					decor.Name(prefixStr, decor.WC{W: len(prefixStr) + 1, C: decor.DidentRight}),
					decor.CountersKibiByte("% -.1f / % -.1f\t"),
					decor.OnComplete(decor.Percentage(decor.WC{W: 4}), " "+"√"),
				),
				mpb.AppendDecorators(
					decor.EwmaETA(decor.ET_STYLE_MMSS, float64(part.RangeTo-part.RangeFrom)/2048),
					decor.Name(" ] "),
					decor.AverageSpeed(decor.UnitKiB, "% .1f"),
				),
			)
			bars = append(bars, newbar)
		}
	}

	for i, p := range d.parts {
		ws.Add(1)
		go func(d *HttpDownloader, loop int64, part Part) {
			defer ws.Done()
			var bar *mpb.Bar

			if DisplayProgressBar() {
				bar = bars[loop]
			}

			var ranges string
			if part.RangeTo != d.len {
				ranges = fmt.Sprintf("bytes=%d-%d", part.RangeFrom, part.RangeTo)
			} else {
				ranges = fmt.Sprintf("bytes=%d-", part.RangeFrom) //get all
			}

			//send request
			req, err := http.NewRequest("GET", d.url, nil)
			req.Header.Set("Connection", "keep-alive")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
			if err != nil {
				errorChan <- err
				return
			}

			if d.par > 1 { //support range download just in case parallel factor is over 1
				req.Header.Add("Range", ranges)
				if err != nil {
					errorChan <- err
					return
				}
			}
			gCurCookies = gCurCookieJar.Cookies(req.URL)

			//write to file
			resp, err := client.Do(req)
			if err != nil {
				errorChan <- err
				return
			}
			defer resp.Body.Close()
			f, err := os.OpenFile(part.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

			defer f.Close()
			if err != nil {
				errorChan <- err
				return
			}

			var writer io.Writer
			writer = io.Writer(f)

			//make copy interruptable by copy 100 bytes each loop
			current := int64(0)
			for {
				select {
				case <-interruptChan:
					bar.Abort(false)
					stateSaveChan <- Part{Url: d.url, Path: part.Path, RangeFrom: current + part.RangeFrom, RangeTo: part.RangeTo}
					return
				default:
					var written int64
					if DisplayProgressBar() {
						reader := bar.ProxyReader(resp.Body)
						written, err = io.CopyN(writer, reader, 100)
					} else {
						written, err = io.CopyN(writer, io.Reader(resp.Body), 100)
					}
					current += written
					if err != nil {
						if err != io.EOF && DisplayProgressBar() {
							errorChan <- err
						} else if err != io.EOF {
							errorChan <- err
						}
						if DisplayProgressBar() {
							bar.Completed()
						}
						fileChan <- part.Path
						return
					}
				}
			}
		}(d, int64(i), p)
	}

	ws.Wait()
	doneChan <- true
}
