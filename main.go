package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	kacker "gitlab.inahga.org/aghani/kacker/pkg"

	"gopkg.in/yaml.v3"
)

// Accept a customization file
// Resolve the kickstart file with the templating engine and store as temporary file
// Merge customization and packer YAML, then convert to json and store as temporary file
// Execute packer

func main() {
	conf := &kacker.Configuration{
		UseInsecureSSL: true,
	}
	kacker.UseConfiguration(conf)

	f, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("Could not read customization file")
	}

	var cust kacker.Customization
	err = yaml.Unmarshal(f, &cust)
	if err != nil {
		log.Fatalf("Invalid customization file: %s", err.Error())
	}

	file, err := cust.Kickstart.ResolveFile()
	if err != nil {
		panic(err)
	}
	fmt.Println(file)
}
