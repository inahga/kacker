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
}
