package hget

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cio "github.com/openbiox/ligo/io"
	clog "github.com/openbiox/ligo/log"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var log = clog.Logger
var displayProgress = true

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
)

func Pull(url string, dest string, pg bool, thread int, mode string, timeout int, proxy string, pbg *mpb.Progress) (err error) {

	client = NewHTTPClient(timeout, proxy)
	if !pg {
		displayProgress = false
	}
	conn := thread
	skiptls := true

	if mode == "resume" {
		state, err := Resume(url, pbg, dest)
		if err != nil && err.Error() == "state not existed" {
			err = Execute(url, nil, conn, skiptls, dest, pbg)
			return err
		}
		err = Execute(state.Url, state, conn, skiptls, dest, pbg)
	} else {
		err = Execute(url, nil, conn, skiptls, dest, pbg)
	}
	return err
}

func Execute(url string, state *State, conn int, skiptls bool, dest string, pbg *mpb.Progress) (err error) {
	//otherwise is hget <URL> command

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

	//set up parallel

	var files = make([]string, 0)
	var parts = make([]Part, 0)
	var isInterrupted = false

	doneChan := make(chan bool, conn)
	fileChan := make(chan string, conn)
	errorChan := make(chan error, 1)
	stateChan := make(chan Part, 1)
	interruptChan := make(chan bool, conn)

	var downloader *HttpDownloader
	if state == nil {
		downloader, err = NewHttpDownloader(url, conn, skiptls, dest)
		if err != nil {
			return err
		}
	} else {
		downloader = &HttpDownloader{url: state.Url, file: dest, par: int64(len(state.Parts)), parts: state.Parts, resumable: true}
	}
	bars := make([]*mpb.Bar, 0)
	go downloader.Do(doneChan, fileChan, errorChan, interruptChan, stateChan, bars, pbg)

	for {
		select {
		case <-signalChan:
			//send par number of interrupt for each routine
			isInterrupted = true
			if DisplayProgressBar() {
				for i := range bars {
					bars[i].Abort(false)
				}
			}
			for i := 0; i < conn; i++ {
				interruptChan <- true
			}
		case file := <-fileChan:
			files = append(files, file)
		case err := <-errorChan:
			filler := makeLogBar(err.Error())
			pbg.Add(0, filler).SetTotal(0, true)
			s := &State{Url: url, Parts: parts}
			s.Save(dest)
			return err
		case part := <-stateChan:
			parts = append(parts, part)
		case <-doneChan:
			if isInterrupted {
				if downloader.resumable {
					if DisplayProgressBar() {
						for i := range bars {
							bars[i].Abort(false)
						}
					}
					filler := makeLogBar(fmt.Sprintf("Interrupted, saving state ..."))
					pbg.Add(0, filler).SetTotal(0, true)
					s := &State{Url: url, Parts: parts}
					err = s.Save(dest)
					if err != nil {
						filler := makeLogBar(err.Error())
						pbg.Add(0, filler).SetTotal(0, true)
					}
					time.Sleep(1 * time.Second)
					os.Exit(130)
					return err
				} else {
					for i := range bars {
						bars[i].Abort(false)
					}
					filler := makeLogBar(fmt.Sprintf("Interrupted, but downloading url is not resumable, silently die"))
					pbg.Add(0, filler).SetTotal(0, true)
					return err
				}
			} else {
				if hasStFile, _ := cio.PathExists(dest + ".st"); hasStFile {
					err = os.Remove(dest + ".st")
				}
				return nil
			}
		}
	}
}

func makeLogBar(msg string) mpb.BarFiller {
	return mpb.BarFillerFunc(func(w io.Writer, _ int, _ decor.Statistics) {
		fmt.Fprintf(w, msg)
	})
}
