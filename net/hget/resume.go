package hget

import (
	mpb "github.com/vbauerster/mpb/v7"
)

func Resume(task string, pbg *mpb.Progress, dest string) (*State, error) {
	return Read(task, pbg, dest)
}
