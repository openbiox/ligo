package par

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	bexec "github.com/openbiox/ligo/exec"
	cio "github.com/openbiox/ligo/io"
	clog "github.com/openbiox/ligo/log"
	"github.com/openbiox/ligo/stringo"
	"github.com/vbauerster/mpb/v4"
	"github.com/vbauerster/mpb/v4/decor"
)

var wg sync.WaitGroup

// ClisT is the type of parameters of parTasks
type ClisT struct {
	Script      string
	Index       string
	Env         string
	ForceAddIdx string
	Thread      int
	TaskID      string
	LogDir      string
	SaveLog     bool
	// Verbose: verbose level (0: no output, 1: basic output)
	Verbose int
}

// Tasks parallel run tasks
func Tasks(parClis *ClisT) (err error) {
	var log = clog.Logger

	index2 := []int{}
	var logCon io.Writer
	var logDir = parClis.LogDir
	var logPrefix string
	var quiet = parClis.Verbose == 0
	var saveLog = parClis.SaveLog

	if parClis.SaveLog {
		if logDir == "" {
			logDir = filepath.Join(os.TempDir(), "_log")
		}
		logPrefix = fmt.Sprintf("%s/%s", logDir, parClis.TaskID)
		cio.CreateDir(logDir)
		logCon, _ = cio.Open(logPrefix + ".log")
	}
	clog.SetLogStream(log, quiet, saveLog, &logCon)

	if parClis.Index != "" {
		index := strings.Split(parClis.Index, ",")
		for i := range index {
			if strings.Contains(index[i], "-") {
				startEnd := strings.Split(index[i], "-")
				start, _ := strconv.ParseInt(startEnd[0], 10, 64)
				end, _ := strconv.ParseInt(startEnd[1], 10, 64)
				for j := start; j < end+1; j++ {
					index2 = append(index2, int(j))
				}
			} else {
				val, _ := strconv.ParseInt(index[i], 10, 64)
				index2 = append(index2, int(val))
			}
		}
	} else {
		index2 = append(index2, 1)
	}
	envObj := make(map[string]string)
	if parClis.Env != "" {
		envSlice := strings.Split(parClis.Env, ",")
		for k := range envSlice {
			envSlice2 := strings.Split(envSlice[k], ":")
			envObj[envSlice2[0]] = envSlice2[1]
		}
	}
	sort.Sort(sort.IntSlice(index2))

	sem := make(chan bool, parClis.Thread)
	p := NewMpb(quiet, saveLog, &logCon)
	wg.Add(len(index2))

	logSlice := []string{}
	for i := range index2 {
		if parClis.SaveLog {
			logSlice = append(logSlice, fmt.Sprintf("%s-%d.log", logPrefix, index2[i]))
		} else {
			logSlice = append(logSlice, "")
		}
	}
	log.Infof("Total %d tasks were submited (%s | %s).", len(index2), parClis.TaskID, parClis.Index)
	log.Infof("Timestamp: %s", time.Now())
	hostname, _ := os.Hostname()
	user, _ := user.Current()
	platform := runtime.GOOS
	log.Infof("Hostname: %s, Username: %s", hostname, user.Username)
	wd, _ := os.Getwd()
	log.Infof("Platform: %s, Working: %s", platform, wd)
	if parClis.SaveLog {
		log.Infof("Task log from #1 to #%d will be saved in %s-*.log", len(index2), logPrefix)
	}
	errorMsg := make(map[int]string)
	var lock sync.Mutex
	total := 100
	for i := 0; i < len(index2); i++ {
		var ind = index2[i]
		sem <- true
		name := fmt.Sprintf("Job: #%-3d", i+1)
		bar := p.AddBar(int64(total), mpb.BarID(i),
			mpb.BarStyle("╢=>-╟"),
			// override mpb.DefaultSpinnerStyle
			mpb.PrependDecorators(
				// simple name decorator
				decor.OnComplete(decor.Spinner(nil, decor.WCSyncSpace), "√"),
				decor.Name(" "+name+fmt.Sprintf("| index: #%-3d", ind)),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.Name(" | elapsed: "), decor.Elapsed(decor.ET_STYLE_HHMMSS),
			),
		)
		// simulating some work
		go func(i int, bar *mpb.Bar) {
			defer func() {
				<-sem
			}()
			var logPath = logSlice[i]
			defer wg.Done()
			go func(bar *mpb.Bar) {
				count := 1
				for {
					if bar.Completed() {
						break
					}
					bar.SetCurrent(int64(count))
					count++
					if count == 99 {
						count = 0
					} else if count < 10 {
						time.Sleep(time.Second / 3)
					} else if count < 20 {
						time.Sleep(time.Second)
					} else if count < 30 {
						time.Sleep(time.Second * 2)
					} else if count < 50 {
						time.Sleep(time.Second * 3)
					} else if count < 70 {
						time.Sleep(time.Second * 4)
					} else if count < 80 {
						time.Sleep(time.Second * 5)
					} else {
						time.Sleep(time.Second * 6)
					}
				}
			}(bar)
			var cmd *exec.Cmd
			script := stringo.StrReplaceAll(parClis.Script, "\n$", "")

			for k, v := range envObj {
				script = stringo.StrReplaceAll(script, fmt.Sprintf("{{%s}}", k), v)
			}
			indStr := fmt.Sprintf("%d", ind)
			if strings.Contains(script, "{{index}}") {
				script = stringo.StrReplaceAll(script, "{{index}}", indStr)
			} else if parClis.ForceAddIdx == "true" {
				script = script + " " + indStr
			}

			cmd = exec.Command("bash", "-c", script)
			err := bexec.System(cmd, logPath, true)
			if err != nil {
				//fmt.Println(err)
				lock.Lock()
				errorMsg[i+1] = fmt.Sprintf("Task #%d error: %s", i+1, err)
				lock.Unlock()
				bar.SetCurrent(int64(0))
				bar.Abort(false)
				time.Sleep(time.Second * 1)
			} else {
				bar.SetCurrent(int64(total))
			}
		}(i, bar)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	// wait for all bars to complete and flush
	p.Wait()
	for i := 0; i < len(index2); i++ {
		if errorMsg[i+1] != "" {
			log.Warn(errorMsg[i+1])
		}
	}
	return nil
}

// NewMpb create mpb.Progress and with log context
func NewMpb(quiet bool, saveLog bool, logCon *io.Writer) (p *mpb.Progress) {
	writers := []io.Writer{
		*logCon,
		os.Stderr}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if !quiet && saveLog {
		p = mpb.New(
			mpb.WithOutput(fileAndStdoutWriter),
			mpb.WithDebugOutput(fileAndStdoutWriter),
			mpb.WithWaitGroup(&wg),
			mpb.WithWidth(108),
		)
	} else if !quiet && !saveLog {
		p = mpb.New(
			mpb.WithWaitGroup(&wg),
			mpb.WithWidth(60),
			mpb.WithOutput(os.Stderr),
			mpb.WithDebugOutput(os.Stderr),
		)
	} else if quiet && saveLog {
		p = mpb.New(
			mpb.WithOutput(*logCon),
			mpb.WithDebugOutput(*logCon),
			mpb.WithWaitGroup(&wg),
			mpb.WithWidth(108),
		)
	} else if quiet && !saveLog {
		p = mpb.New(
			mpb.WithWaitGroup(&wg),
			mpb.WithWidth(60),
			mpb.WithOutput(ioutil.Discard),
			mpb.WithDebugOutput(ioutil.Discard),
		)
	}
	return p
}
