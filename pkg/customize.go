package kacker

import (
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

// Customization represents the customization file: modifiers for the kickstart
// file and the packer file.
type Customization struct {
	Kickstart Kickstart `yaml:"kickstart"`
	Packer    Packer    `yaml:"packer"`
}

func (c *Customization) Run(packerFlags []string) error {
	kickstart, err := c.Kickstart.ResolveTempFile()
	if err != nil {
		return err
	}
	os.Setenv("KICKSTART", kickstart)
	if !globalConf.KeepTempFiles {
		defer os.Remove(kickstart)
	}

	packer, err := c.Packer.ResolveTempFile()
	if err != nil {
		return err
	}
	if !globalConf.KeepTempFiles {
		defer os.Remove(packer)
	}

	var commands []*exec.Cmd
	commands = append(commands, exec.Command("packer",
		append(append([]string{"validate"}, packerFlags...), packer)...))
	commands = append(commands, exec.Command("packer",
		append(append([]string{"build"}, packerFlags...), packer)...))

	if err = runCommands(commands); err != nil {
		return err
	}
	return nil
}

func Run(filename string, packerFlags []string) error {
	var cust Customization
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(f, &cust); err != nil {
		return err
	}
	return cust.Run(packerFlags)
}
