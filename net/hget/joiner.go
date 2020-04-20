package hget

import (
	"io"
	"os"
	"sort"
)

func JoinFile(files []string, out string) error {
	//sort with file name or we will join files with wrong order
	sort.Strings(files)
	outf, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0600)
	defer outf.Close()
	if err != nil {
		return err
	}
	for _, f := range files {
		copy(f, outf)
	}
	return nil
}

//this function split just to use defer
func copy(from string, to io.Writer) error {
	f, err := os.OpenFile(from, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	io.Copy(to, f)
	return nil
}
