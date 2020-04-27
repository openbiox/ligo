package hget

import (
	"bytes"
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

	ret := new(HttpDownloader)
	ret.url = url
	ret.file = dest
	ret.par = int64(par)
	ret.len = len
	ret.ips = ipstr
	ret.skipTls = skipTls
	ret.parts = partCalculate(int64(par), len, url, dest)
	ret.resumable = resumable

	return ret, nil
}

func partCalculate(par int64, len int64, url string, dest string) []Part {
	ret := make([]Part, 0)
	for j := int64(0); j < par; j++ {
		from := (len / par) * j
		var to int64
		if j < par-1 {
			to = (len/par)*(j+1) - 1
		} else {
			to = len
		}
		path := dest // ~/.hget/download-file-name/part-name
		ret = append(ret, Part{Url: url, Path: path, RangeFrom: from, RangeTo: to})
	}
	return ret
}

func (d *HttpDownloader) Do(doneChan chan bool, fileChan chan string, errorChan chan error, interruptChan chan bool, stateSaveChan chan Part, bars []*mpb.Bar, pbg *mpb.Progress) {

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
			SciencedirectassetsRed(d.url, req, client)

			if d.resumable { //support range download just in case parallel factor is over 1
				req.Header.Add("Range", ranges)
			}

			gCurCookies = gCurCookieJar.Cookies(req.URL)

			//write to file
			resp, err := client.Do(req)
			if err != nil {
				errorChan <- err
				return
			}
			defer resp.Body.Close()
			var f *os.File
			f, err = os.OpenFile(part.Path, os.O_CREATE|os.O_WRONLY, 0600)

			defer f.Close()
			if err != nil {
				errorChan <- err
				return
			}

			var writer io.WriterAt
			writer = io.WriterAt(f)
			//make copy interruptable by copy 100 bytes each loop
			current := int64(0)
			for {
				from := current + part.RangeFrom
				select {
				case <-interruptChan:
					if DisplayProgressBar() {
						bar.Abort(false)
					}
					stateSaveChan <- Part{Url: d.url, Path: part.Path, RangeFrom: from, RangeTo: part.RangeTo}
					return
				default:
					var written int64
					var buf = &bytes.Buffer{}
					if DisplayProgressBar() && d.resumable {
						reader := bar.ProxyReader(resp.Body)
						written, err = io.CopyN(buf, reader, 100)
						writer.WriteAt(buf.Bytes(), from)
					} else if !DisplayProgressBar() && d.resumable {
						written, err = io.CopyN(buf, io.Reader(resp.Body), 100)
						writer.WriteAt(buf.Bytes(), from)
					} else if DisplayProgressBar() && !d.resumable {
						reader := bar.ProxyReader(resp.Body)
						written, err = io.CopyN(buf, reader, 100)
						writer.WriteAt(buf.Bytes(), current)
					} else if !DisplayProgressBar() && !d.resumable {
						written, err = io.CopyN(buf, io.Reader(resp.Body), 100)
						writer.WriteAt(buf.Bytes(), current)
					}
					current += written
					if err != nil {
						if err != io.EOF && DisplayProgressBar() {
							stateSaveChan <- Part{Url: d.url, Path: part.Path, RangeFrom: from, RangeTo: part.RangeTo}
							bar.Abort(false)
							errorChan <- err
						} else if err != io.EOF {
							stateSaveChan <- Part{Url: d.url, Path: part.Path, RangeFrom: from, RangeTo: part.RangeTo}
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
