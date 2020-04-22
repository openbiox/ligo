package hget

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	cio "github.com/openbiox/ligo/io"
	mpb "github.com/vbauerster/mpb/v5"
)

var dataFolder = ".config/bget/data"
var stateFileName = "state.json"

type State struct {
	Url   string
	Parts []Part
}

type Part struct {
	Url       string
	Path      string
	RangeFrom int64
	RangeTo   int64
}

func (s *State) Save() error {
	//make temp folder
	//only working in unix with env HOME
	folder := FolderOf(s.Url)
	if err := MkdirIfNotExist(folder); err != nil {
		return err
	}

	//move current downloading file to data folder
	for _, part := range s.Parts {
		os.Rename(part.Path, filepath.Join(folder, filepath.Base(part.Path)))
	}

	//save state file
	j, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(folder, stateFileName), j, 0644)
}

func Read(task string, pbg *mpb.Progress) (*State, error) {
	file := filepath.Join(os.Getenv("HOME"), dataFolder, task, stateFileName)
	if hasStateFile, _ := cio.PathExists(file); hasStateFile {
		filler := makeLogBar(fmt.Sprintf("Getting state data from %s.", file))
		pbg.Add(0, filler).SetTotal(0, true)
	} else {
		filler := makeLogBar(fmt.Sprintf("State file of %s not existed.", path.Base(path.Dir(file))))
		pbg.Add(0, filler).SetTotal(0, true)
		return nil, errors.New("state not existed")
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	s := new(State)
	err = json.Unmarshal(bytes, s)
	return s, err
}
