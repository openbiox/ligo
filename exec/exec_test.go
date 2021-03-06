package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"

	cio "github.com/openbiox/ligo/io"
)

func TestRunExecCmdConsole(t *testing.T) {
	args := []string{"-c"}
	cmd := exec.Command("wget", args...)
	dir := os.TempDir()
	logPath := path.Join(dir, "TestRunExecCmdConsole.log")
	System(cmd, logPath, false)
	if hasFile, _ := cio.PathExists(logPath); hasFile {
		if !(cio.FileSize(logPath) > 0) {
			log.Fatalf("%s not created.", logPath)
			os.Exit(1)
		}
		err := os.Remove(logPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	cmd = exec.Command("wget", args...)
	System(cmd, "", false)
	if hasFile, _ := cio.PathExists(logPath); hasFile {
		log.Fatalf("%s should not be created.", logPath)
	}

	fmt.Println("quite mode")
	args = []string{"-c", "https://raw.githubusercontent.com/openbiox/bioctl/master/main.go"}
	cmd = exec.Command("wget", args...)
	System(cmd, "", true)
	if hasFile, _ := cio.PathExists(logPath); hasFile {
		log.Fatalf("%s should not be created.", logPath)
	}
	fmt.Println("quite mode and save log")
	cmd = exec.Command("wget", args...)
	System(cmd, logPath, true)
	if hasFile, _ := cio.PathExists(logPath); !hasFile {
		log.Fatalf("%s should be created.", logPath)
	}
	//err := os.Remove(logPath)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
