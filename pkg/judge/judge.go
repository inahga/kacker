package judge

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync/atomic"
)

// Verdict represents the execution of a code snippet
type Verdict struct {
	Lang string
	Code string

	// Result must only be checked if IsReady() returns true
	Result string
	ready  *atomic.Value
}

// IsReady checks whether the Verdict goroutine has finished and stored a
// result
func (v *Verdict) IsReady() bool {
	if v.ready.Load() == 0 {
		return false
	}
	return true
}

func runSync(c []*exec.Cmd) (string, error) {
	var ret string
	for _, command := range c {
		output, err := command.CombinedOutput()
		ret += string(output)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

// Run executes the code of the given language. It returns a pending verdict
// that finishes processing asynchronously. This will eventually be in charge
// of jailing
func Run(lang string, code string) (*Verdict, error) {
	language := findLanguage(lang)
	if language == nil {
		return nil, fmt.Errorf("Could not find language %s", lang)
	}

	ret := &Verdict{
		Lang:   lang,
		Code:   code,
		Result: "",
		ready:  &atomic.Value{},
	}
	ret.ready.Store(0)

	go func() {
		defer ret.ready.Store(1)

		f, err := ioutil.TempFile("/tmp", "defendant-*."+language.Extension)
		if err != nil {
			ret.Result += fmt.Sprintf("Failed to open temporary file: %s", err)
			return
		}
		defer os.Remove(f.Name())
		_, err = f.WriteString(code)
		if err != nil {
			ret.Result += fmt.Sprintf("Failed to write to temporary file: %s", err)
			return
		}
		err = f.Close()
		if err != nil {
			ret.Result += fmt.Sprintf("Failed to close temporary file: %s", err)
			return
		}

		output, err := runSync(language.GetCommands(f.Name(), ""))
		if err != nil {
			ret.Result += fmt.Sprintf("Command exited with error: %s\n", err)
		}
		ret.Result += output
	}()

	return ret, nil
}
