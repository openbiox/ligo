package par

import "testing"

func TestTasks(t *testing.T) {
	a := ClisT{
		LogDir:  "/tmp/_log",
		Quiet:   "false",
		SaveLog: "false",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	b := ClisT{
		LogDir:  "/tmp/_log",
		Quiet:   "false",
		SaveLog: "true",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	c := ClisT{
		LogDir:  "/tmp/_log",
		Quiet:   "true",
		SaveLog: "true",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&a)
	Tasks(&b)
	Tasks(&c)
}
