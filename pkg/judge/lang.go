package judge

import (
	"os/exec"
	"path"
)

// Language represents a programming language's means of compiling/interpreting
// a program and running it
type Language struct {
	Name        string
	Executable  string
	Extension   string
	GetCommands func(string, string) []*exec.Cmd
}

func findLanguage(name string) *Language {
	for _, lang := range languages {
		if lang.Name == name {
			return lang
		}
	}
	return nil
}

var languages = []*Language{
	{
		Name:       "c",
		Executable: "/usr/bin/gcc",
		Extension:  "c",
		GetCommands: func(p string, flags string) []*exec.Cmd {
			bin := path.Join(path.Dir(p), "a.out")
			return []*exec.Cmd{
				exec.Command("/usr/bin/gcc", "-o", bin, p),
				exec.Command(bin),
			}
		},
	},
	{
		Name:       "python3",
		Executable: "/usr/bin/python3",
		Extension:  "py",
		GetCommands: func(p string, flags string) []*exec.Cmd {
			return []*exec.Cmd{
				exec.Command("/usr/bin/python3", p),
			}
		},
	},
}
