package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	kacker "gitlab.inahga.org/aghani/kacker/internal"

	"gopkg.in/yaml.v3"
)

// Accept a customization file
// Resolve the kickstart file with the templating engine and store as temporary file
// Merge customization and packer YAML, then convert to json and store as temporary file
// Execute packer

// Config stores the various paths needed for operation
type Config struct {
	CustomizationPath string `yaml:"customizationPath"`
	FragmentsPath     string `yaml:"fragmentsPath"`
	KickstartPath     string `yaml:"kickstartPath"`
	PackerfilePath    string `yaml:"packerfilePath"`
}

func main() {
	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Could not read customization file")
	}

	var cust kacker.Customization
	err = yaml.Unmarshal(f, &cust)
	if err != nil {
		log.Fatalln("Invalid customization file")
	}

	fmt.Printf("%+v\n", cust)
}
