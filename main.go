package main

import (
	"os"

	kacker "gitlab.inahga.org/aghani/kacker/pkg"
)

func main() {
	conf := &kacker.Configuration{
		VerboseLogging: false,
		UseInsecureSSL: true,
	}
	kacker.UseConfiguration(conf)
	if err := kacker.Run(os.Args[1]); err != nil {
		panic(err)
	}

	// f, err := ioutil.ReadFile(os.Args[1])
	// if err != nil {
	// 	log.Fatalln("Could not read customization file")
	// }

	// var cust kacker.Customization
	// err = yaml.Unmarshal(f, &cust)
	// if err != nil {
	// 	log.Fatalf("Invalid customization file: %s", err.Error())
	// }

	// file, err := cust.Kickstart.ResolveTempFile()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(file)

	// err = cust.Kickstart.Resolve(os.Stdout)

	// kacker.MergeProperties([]string{"test1.yml", "test2.yml"})
	// kacker.ReadYamlFile("test1.yml")

	// var cust kacker.Customization
	// f, _ := ioutil.ReadFile(os.Args[1])
	// yaml.Unmarshal(f, &cust)
	// err := cust.Packer.Resolve(os.Stdout)
	// if err != nil {
	// 	panic(err)
	// }

	// baseNodes, err := kacker.ReadYamlFile(cust.Packer.From)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(cust.Packer.From)
	// childNodes, _ := kacker.ReadYamlField(cust.Packer.Merge)
	// fmt.Printf("%v", baseNodes)

	// commands := kacker.GetMergeCommands(append(baseNodes, childNodes...))
	// for _, command := range commands {
	// 	fmt.Printf("%v\n", command)
	// 	fmt.Printf("%v\n", command.Value)
	// 	fmt.Println("")
	// }

	// err = kacker.UpdateDocument(cust.Packer.From, os.Stdout, commands)
	// if err != nil {
	// 	panic(err)
	// }
}
