package judge

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func installBogusLanguage() {
	languages = append(languages, &Language{
		Name:       "sleep",
		Executable: "/usr/bin/sleep",
		Extension:  ".sleep",
		GetCommands: func(string, string) []*exec.Cmd {
			return []*exec.Cmd{
				exec.Command("/usr/bin/sleep", "5"),
			}
		},
	})
}

func TestRun(t *testing.T) {
	installBogusLanguage()

	t.Run("Returns error when invalid language is specified", func(t *testing.T) {
		verdict, err := Run("fakelang", "")
		if verdict != nil || err == nil {
			t.Errorf("Did not receive expected error")
		}
	})

	t.Run("Executes successfully", func(t *testing.T) {
		verdict, err := Run("sleep", "")
		if verdict == nil || err != nil {
			t.Errorf("Received unexpected error: %s", err)
		}

		time.Sleep(time.Second)
		tmpDir, _ := ioutil.ReadDir("/tmp")
		var found bool
		for _, tmp := range tmpDir {
			if strings.Contains(tmp.Name(), ".sleep") {
				t.Log("Found temporary file")
				found = true
			}
		}
		if !found {
			t.Log("Didn't find temporary file")
		}

		for i := 0; i < 6; i++ {
			t.Log("Checking if ready...")
			if verdict.IsReady() {
				t.Log(verdict.Result)
				return
			}
			time.Sleep(time.Second)
		}
		t.Error("Verdict was not ready after 7 seconds")
	})
}
