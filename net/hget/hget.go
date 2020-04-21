package hget

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	clog "github.com/openbiox/ligo/log"
	mpb "github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

var log = clog.Logger
var displayProgress = true
var pbg *mpb.Progress

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
)

func Pull(taskname string, dest string, pg bool, thread int, mode string, timeout int, proxy string, p *mpb.Progress) (err error) {
	if pbg == nil {
		pbg = p
	}
	client = NewHTTPClient(timeout, proxy)
	if !pg {
		displayProgress = false
	}
	conn := thread
	skiptls := true

	if mode == "tasks" {
		if err = TaskPrint(); err != nil {
			return err
		}
	} else if mode == "resume" {
		var task string
		if IsUrl(taskname) {
			task = TaskFromUrl(taskname)
		} else {
			task = taskname
		}
		state, err := Resume(task)
		if err != nil && err.Error() == "state not existed" {
			os.RemoveAll(FolderOf(taskname))
			err = Execute(taskname, nil, conn, skiptls, dest)
			return err
		}
		err = Execute(state.Url, state, conn, skiptls, dest)
	} else {
		err = Execute(taskname, nil, conn, skiptls, dest)
	}
	return err
}

func Execute(url string, state *State, conn int, skiptls bool, dest string) (err error) {
	//otherwise is hget <URL> command

	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

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
	go downloader.Do(doneChan, fileChan, errorChan, interruptChan, stateChan, bars)

	for {
		select {
		case <-signal_chan:
			//send par number of interrupt for each routine
			isInterrupted = true
			for i := range bars {
				bars[i].Abort(false)
			}
			for i := 0; i < conn; i++ {
				interruptChan <- true
			}
		case file := <-fileChan:
			files = append(files, file)
		case err := <-errorChan:
			filler := makeLogBar(err.Error())
			pbg.Add(0, filler).SetTotal(0, true)
			return err
		case part := <-stateChan:
			parts = append(parts, part)
		case <-doneChan:
			if isInterrupted {
				if downloader.resumable {
					for i := range bars {
						bars[i].Abort(false)
					}
					filler := makeLogBar(fmt.Sprintf("Interrupted, saving state ..."))
					pbg.Add(0, filler).SetTotal(0, true)
					s := &State{Url: url, Parts: parts}
					err = s.Save()
					if err != nil {
						filler := makeLogBar(err.Error())
						pbg.Add(0, filler).SetTotal(0, true)
					}
					return err
				} else {
					for i := range bars {
						bars[i].Abort(false)
					}
					filler := makeLogBar(fmt.Sprintf("Interrupted, but downloading url is not resumable, silently die"))
					pbg.Add(0, filler).SetTotal(0, true)
					os.RemoveAll(FolderOf(url))
					return nil
				}
			} else {
				err := JoinFile(files, dest)
				if err != nil {
					return err
				}
				os.RemoveAll(FolderOf(url))
				return nil
			}
		}
	}
}

func makeLogBar(msg string) mpb.BarFiller {
	return mpb.BarFillerFunc(func(w io.Writer, width int, st *decor.Statistics) {
		fmt.Fprintf(w, msg)
	})
}
