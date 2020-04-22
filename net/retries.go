package net

import (
	"errors"
	"net/http"
	"os/exec"
	"time"

	bexec "github.com/openbiox/ligo/exec"
)

// RetriesClient http.Client and requenst with retries
func RetriesClient(client *http.Client, req *http.Request, opt *Params) (resp *http.Response, err error) {
	for t := 0; t < opt.Retries+1; t++ {
		resp, err = client.Do(req)
		if err != nil {
			if opt.Retries != 0 {
				log.Warnf("Failed to retrieve on attempt %d... error: %v ... retrying after %d seconds.", t+1, err, opt.RetSleepTime)
				time.Sleep(time.Duration(opt.RetSleepTime) * time.Second)
			} else {
				log.Warnf("Failed to retrieve on attempt %d... error: %v", t+1, err)
			}
			continue
		} else if err2 := checkResp(resp); err2 != nil {
			return nil, err2
		} else {
			break
		}
	}
	return resp, err
}

// RetriesURL access URL with retries
func RetriesURL(url string, cmd *exec.Cmd, logPath string, opt *Params) (err error) {
	var t int
	success := false
	var cmdRun exec.Cmd
	for t = 0; t < opt.Retries+1; t++ {
		cmdRun = *cmd
		err := bexec.System(&cmdRun, logPath, opt.Quiet)
		if err == nil {
			success = true
			break
		} else {
			if opt.Retries != 0 {
				if err.Error() == "signal: interrupt" {
					break
				}
				log.Warnf("Failed to retrive on attempt %d... error: %v ... retrying after %d seconds.", t+1, err, opt.RetSleepTime)
				time.Sleep(time.Duration(opt.RetSleepTime) * time.Second)
			} else {
				log.Warnf("Failed to retrieve on attempt %d... error: %v", t+1, err)
			}
			continue
		}
	}
	if !success {
		return errors.New("Faild to access: " + url)
	}
	return nil
}

// RetriesTask task with retries
func RetriesTask(taskName string, cmd *exec.Cmd, logPath string, opt *Params) (err error) {
	var t int
	success := false
	var cmdRun exec.Cmd
	for t = 0; t < opt.Retries+1; t++ {
		cmdRun = *cmd
		err := bexec.System(&cmdRun, logPath, opt.Quiet)
		if err == nil {
			success = true
			break
		} else {
			if opt.Retries != 0 {
				log.Warnf("Failed on attempt %d... error: %v ... retrying after %d seconds.", t+1, err, opt.RetSleepTime)
			} else {
				log.Warnf("Failed on attempt %d... error: %v", t+1, err)
			}
			time.Sleep(time.Duration(opt.RetSleepTime) * time.Second)
		}
	}
	if !success {
		return errors.New("Faild to access: " + taskName)
	}
	return nil
}
