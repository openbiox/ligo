package par

import "testing"

func TestTasks(t *testing.T) {
	a := ClisT{
		LogDir:  "/tmp/_log1",
		Quiet:   "false",
		SaveLog: "false",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&a)
	b := ClisT{
		LogDir:  "/tmp/_log2",
		Quiet:   "false",
		SaveLog: "true",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&b)
	c := ClisT{
		LogDir:  "/tmp/_log3",
		Quiet:   "true",
		SaveLog: "true",
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&c)
}
