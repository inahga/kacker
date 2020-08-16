package main

import (
	"fmt"
	"io/ioutil"
	"os"

	kacker "gitlab.inahga.org/aghani/kacker/pkg"
	"gopkg.in/yaml.v2"
)

// Accept a customization file
// Resolve the kickstart file with the templating engine and store as temporary file
// Merge customization and packer YAML, then convert to json and store as temporary file
// Execute packer

func main() {
	// conf := &kacker.Configuration{
	// 	VerboseLogging: true,
	// 	UseInsecureSSL: true,
	// }
	// kacker.UseConfiguration(conf)

	// f, err := ioutil.ReadFile(os.Args[1])
	// if err != nil {
	// 	log.Fatalln("Could not read customization file")
	// }

	// var cust kacker.Customization
	// err = yaml.Unmarshal(f, &cust)
	// if err != nil {
	// 	log.Fatalf("Invalid customization file: %s", err.Error())
	// }

	// file, err := cust.Kickstart.ResolveFile()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(file)

	// kacker.MergeProperties([]string{"test1.yml", "test2.yml"})
	// kacker.ReadYamlFile("test1.yml")
	var cust kacker.Customization
	f, _ := ioutil.ReadFile(os.Args[1])
	yaml.Unmarshal(f, &cust)
	baseNodes, err := kacker.ReadYamlFile(cust.Packer.From)
	if err != nil {
		panic(err)
	}
	fmt.Println(cust.Packer.From)
	childNodes, _ := kacker.ReadYamlField(cust.Packer.Merge)
	fmt.Printf("%v", baseNodes)

	commands := kacker.GetMergeCommands(append(baseNodes, childNodes...))
	for _, command := range commands {
		fmt.Printf("%v\n", command)
		fmt.Printf("%v\n", command.Value)
		fmt.Println("")
	}
	// err = kacker.ExecuteMergeCommands(os.Args[1], os.Stdout, commands)
	// if err != nil {
	// 	panic(err)
	// }
}
