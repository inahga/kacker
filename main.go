package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"gopkg.in/yaml.v3"
)

// Accept a customization file
// Resolve the kickstart file with the templating engine and store as temporary file
// Merge customization and packer YAML, then convert to json and store as temporary file
// Execute packer

// Config stores the various paths needed for operation
type Config struct {
	CustomizationPath string `yaml:customizationPath`
	FragmentsPath     string `yaml:fragmentsPath`
	KickstartPath     string `yaml:kickstartPath`
	PackerfilePath    string `yaml:packerfilePath`
}

// Variable represents a template variable that will be provided to the kickstart
// template.
type Variable struct {
	Name      string   `yaml:"name"`
	Value     string   `yaml:"value"`
	URL       string   `yaml:"url"`
	Fragments []string `yaml:"fragments"`
}

// KickstartCustomization collects the kickstart template file and any variable substitutions.
type KickstartCustomization struct {
	From      string     `yaml:"from"`
	Variables []Variable `yaml:"variables"`
}

// PackerCustomization collects the packer file and any YAML substitutions that
// need to be made.
type PackerCustomization struct {
	From  string      `yaml:"from"`
	Merge interface{} `yaml:"merge"`
}

// Customization represents the customization file: modifiers for the kickstart
// file and the packer file.
type Customization struct {
	Kickstart KickstartCustomization `yaml:"kickstart"`
	Packer    PackerCustomization    `yaml:"packer"`
}

func main() {
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Could not read customization file")
	}

	var cust Customization
	err = yaml.Unmarshal(f, &cust)
	if err != nil {
		log.Fatalln("Invalid customization file")
	}

	yqlib.Get()

	fmt.Printf("%+v\n", cust)
}
