package hget

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	cio "github.com/openbiox/ligo/io"
	mpb "github.com/vbauerster/mpb/v5"
)

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

func (s *State) Save(dest string) error {
	//make temp folder
	j, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest+".st", j, 0644)
}

func Read(task string, pbg *mpb.Progress, dest string) (*State, error) {
	file := dest + ".st"
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
