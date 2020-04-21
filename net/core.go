package net

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	cio "github.com/openbiox/ligo/io"
	"github.com/openbiox/ligo/net/hget"
	mpb "github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

// HTTPGetURLs can use golang http.Get and external commandline tools including wget, curl, axel, git and rsync
// to query URL with progress bar
func HTTPGetURLs(urls []string, destDir []string, opt *Params) (destFns []string) {
	sem := make(chan bool, opt.Thread)
	newOpt := *opt
	if len(urls) > 1 && opt.Thread > 1 && opt.Engine != "simplego" && opt.Engine != "default" {
		newOpt.Quiet = true
	}
	wg := sync.WaitGroup{}
	destMap := make(map[string]string)
	wg.Add(len(urls))
	for j := range urls {
		url := urls[j]
		dest := destDir[j]
		go func(url, dest string) {
			filename := FormatURLfileName(url, newOpt.RemoteName, newOpt.Timeout, newOpt.Proxy)
			lock.Lock()
			destMap[url] = path.Join(dest, filename)
			lock.Unlock()
			wg.Done()
		}(url, dest)
	}
	wg.Wait()
	for k, v := range destMap {
		log.Infof("Trying %s => %s", k, v)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	for j := range urls {
		cio.CreateDir(destDir[j])
		destFn := destMap[urls[j]]
		if newOpt.Overwrite {
			err := os.RemoveAll(destFn)
			if err != nil {
				log.Warnf("Can not remove %s.", destFn)
			}
		}
		if hasDestFn, _ := cio.PathExists(destFn); !hasDestFn || newOpt.Ignore {
			url := urls[j]
			sem <- true
			go func(url string, destFn string, signalChan chan os.Signal) {
				defer func() {
					select {
					case <-signalChan:
						//send par number of interrupt for each routine
						time.Sleep(3 * time.Second)
						filler := makeLogBar(fmt.Sprintf("Existing ..."))
						newOpt.Pbar.Add(0, filler).SetTotal(0, true)
						os.Exit(130)
						return
					case <-sem:
						return
					}
				}()
				err := AsyncURL2(url, destFn, &newOpt)
				if err == nil {
					destFns = append(destFns, destFn)
				}
			}(url, destFn, signalChan)
		} else {
			destFns = append(destFns, destFn)
			log.Infof("%s existed.", destFn)
		}
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	opt.Pbar.Wait()
	return destFns
}

// NewHTTPClient create http.Client with timeout, proxy, and gCurCookieJar
func NewHTTPClient(timeout int, proxy string) *http.Client {
	transPort := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		urlproxy, _ := RandProxy(proxy)
		transPort.Proxy = http.ProxyURL(urlproxy)
	}
	return &http.Client{
		CheckRedirect: defaultCheckRedirect,
		Jar:           gCurCookieJar,
		Transport:     transPort,
	}
}

// HTTPGetURL can use golang http.Get to query URL with progress bar
func HTTPGetURL(url string, destFn string, opt *Params) error {

	client := NewHTTPClient(opt.Timeout, opt.Proxy)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	if err != nil {
		// handle error
		log.Warn(err)
		return err
	}
	gCurCookies = gCurCookieJar.Cookies(req.URL)

	var t int
	var success = false
	for t = 0; t < opt.Retries; t++ {
		err = downloadWorker(client, req, url, destFn, opt)
		if err == nil {
			success = true
			break
		} else {
			log.Warnf("Failed to retrive on attempt %d... error: %v ... retrying after %d seconds.", t+1, err, opt.RetSleepTime)
			time.Sleep(time.Duration(opt.RetSleepTime) * time.Second)
		}
	}
	if !success {
		return err
	}
	return nil
}

// ParseOutfnFromHeader get filename from response header
func ParseOutfnFromHeader(outfn string, resp *http.Response, useRemoteName bool) string {
	contentDis := resp.Header.Get("Content-Disposition")
	if outfn == "" && contentDis != "" && useRemoteName &&
		strings.Contains(contentDis, "filename") {
		_, params, err := mime.ParseMediaType(contentDis)
		if err != nil {
			log.Warn(err)
		} else {
			outfn = params["filename"]
		}
	}
	return outfn
}

// AsyncURL can access URL via using external commandline tools including
// wget, curl, axel, git and rsync
func AsyncURL(url string, destFn string, opt *Params) error {
	if opt.Mirror != "" {
		if !strings.HasSuffix(opt.Mirror, "/") {
			opt.Mirror = opt.Mirror + "/"
		}
		url = opt.Mirror + filepath.Base(url)
	}
	gitEngine := checkGitEngine(url)
	engine := opt.Engine
	if gitEngine == "git" {
		engine = "git"
	} else if gitEngine == "hg" {
		engine = "hg"
	}
	if engine == "wget" {
		return Wget(url, destFn, opt)
	} else if engine == "curl" {
		return Curl(url, destFn, opt)
	} else if engine == "axel" {
		return Axel(url, destFn, opt)
	} else if engine == "rsync" {
		return Rsync(url, destFn, opt)
	} else if engine == "git" {
		return Git(url, destFn, opt)
	} else if engine == "hg" {
		return Hg(url, destFn, opt)
	}
	log.Infof("Not found valid download engine name.")
	return nil
}

// AsyncURL2 can access URL via using golang http library (with mbp progress bar) and
// external commandline tools including wget, curl, axel, git and rsync
func AsyncURL2(url string, destFn string, opt *Params) error {
	engine := opt.Engine
	gitEngine := checkGitEngine(url)
	if gitEngine == "git" {
		engine = "git"
	} else if gitEngine == "hg" {
		engine = "hg"
	}
	if engine == "simplego" {
		if opt.Mirror != "" {
			if !strings.HasSuffix(opt.Mirror, "/") {
				opt.Mirror = opt.Mirror + "/"
			}
			url = opt.Mirror + filepath.Base(url)
		}
		return HTTPGetURL(url, destFn, opt)
	} else if engine == "default" {
		mode := ""
		stateDir := path.Join(os.Getenv("HOME"), ".config/bget/data", filepath.Base(destFn))
		hasDest, _ := cio.PathExists(destFn)
		hasState, _ := cio.PathExists(stateDir)
		if hasDest || hasState {
			mode = "resume"
		}
		return hget.Pull(url, destFn, !opt.Quiet, opt.ThreadQuery, mode, opt.Timeout, opt.Proxy, opt.Pbar)
	}
	return AsyncURL(url, destFn, opt)
}

// AsyncURL3 can access URL via using golang http library (with mbp progress bar) and
// external commandline tools including wget, curl, axel, git and rsync
func AsyncURL3(url string, destFn string, opt *Params) (err error) {
	engine := opt.Engine
	gitEngine := checkGitEngine(url)
	if gitEngine == "git" {
		engine = "git"
	} else if gitEngine == "hg" {
		engine = "hg"
	}
	if engine == "simplego" {
		if opt.Mirror != "" {
			if !strings.HasSuffix(opt.Mirror, "/") {
				opt.Mirror = opt.Mirror + "/"
			}
			url = opt.Mirror + filepath.Base(url)
		}
		err = HTTPGetURL(url, destFn, opt)
		opt.Pbar.Wait()
	} else {
		err = AsyncURL(url, destFn, opt)
	}
	return nil
}

func makeLogBar(msg string) mpb.BarFiller {
	limit := "%%.%ds"
	return mpb.BarFillerFunc(func(w io.Writer, width int, st *decor.Statistics) {
		fmt.Fprintf(w, fmt.Sprintf(limit, width), msg)
	})
}

func init() {
	gCurCookies = nil
	//var err error;
	gCurCookieJar, _ = cookiejar.New(nil)
}
