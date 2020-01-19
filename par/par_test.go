package par

import "testing"

func TestTasks(t *testing.T) {
	a := ClisT{
		LogDir:  "/tmp/_log1",
		Verbose: 1,
		SaveLog: false,
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}} && sleep 200",
		Thread:  3,
	}
	Tasks(&a)
	b := ClisT{
		LogDir:  "/tmp/_log2",
		Verbose: 1,
		SaveLog: true,
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&b)
	c := ClisT{
		LogDir:  "/tmp/_log3",
		Verbose: 0,
		SaveLog: true,
		TaskID:  "test123",
		Index:   "1,2,3",
		Script:  "echo {{index}}",
		Thread:  3,
	}
	Tasks(&c)
}
