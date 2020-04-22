package hget

import (
	"io/ioutil"
	"os"
	"path/filepath"

	mpb "github.com/vbauerster/mpb/v5"
)

func TaskPrint() error {
	downloading, err := ioutil.ReadDir(filepath.Join(os.Getenv("HOME"), dataFolder))
	if err != nil {
		return err
	}

	folders := make([]string, 0)
	for _, d := range downloading {
		if d.IsDir() {
			folders = append(folders, d.Name())
		}
	}

	return nil
}

func Resume(task string, pbg *mpb.Progress) (*State, error) {
	return Read(task, pbg)
}
