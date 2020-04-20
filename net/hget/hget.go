package hget

import (
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	clog "github.com/openbiox/ligo/log"
	mpb "github.com/vbauerster/mpb/v4"
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
		if err != nil {
			return err
		}
		Execute(state.Url, state, conn, skiptls, dest)
	} else {
		Execute(taskname, nil, conn, skiptls, dest)
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
		downloader = NewHttpDownloader(url, conn, skiptls, dest)
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
			for i := range bars {
				bars[i].Abort(false)
			}
			return err
		case part := <-stateChan:
			parts = append(parts, part)
		case <-doneChan:
			if isInterrupted {
				if downloader.resumable {
					for i := range bars {
						bars[i].Abort(false)
					}
					log.Infof("Interrupted, saving state ...")
					s := &State{Url: url, Parts: parts}
					err = s.Save()
					if err != nil {
						log.Warnln(err)
					}
					return
				} else {
					for i := range bars {
						bars[i].Abort(false)
					}
					log.Warnln("Interrupted, but downloading url is not resumable, silently die")
					os.RemoveAll(FolderOf(url))
					return
				}
			} else {
				JoinFile(files, dest)
				os.RemoveAll(FolderOf(url))
				return
			}
		}
	}
}
