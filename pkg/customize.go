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

func (c *Customization) Run() error {
	kickstart, err := c.Kickstart.ResolveTempFile()
	if err != nil {
		return err
	}
	os.Setenv("KICKSTART", kickstart)

	packer, err := c.Packer.ResolveTempFile()
	if err != nil {
		return err
	}

	var commands []*exec.Cmd
	commands = append(commands, exec.Command("packer", "validate", packer))
	commands = append(commands, exec.Command("packer", "build", packer))

	if err = runCommands(commands); err != nil {
		return err
	}
	return nil
}

func Run(filename string) error {
	var cust Customization
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(f, &cust); err != nil {
		return err
	}
	return cust.Run()
}
