package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	kacker "gitlab.inahga.org/aghani/kacker/pkg"
)

func main() {
	verbose := flag.Bool("verbose", false, "")
	insecure := flag.Bool("no-verify-ssl", false, "Do not verify SSL certificates when resolving URLs")
	packerFlags := flag.String("packer-flags", "", "Send these flags to Packer when executing")
	keepTemp := flag.Bool("keep-temp", false, "Do not delete temporary files")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Args()[0]
	absPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	conf := &kacker.Configuration{
		VerboseLogging: *verbose,
		NoVerifySSL:    *insecure,
		KeepTempFiles:  *keepTemp,
		RelativeDir:    path.Dir(absPath),
	}
	kacker.UseConfiguration(conf)

	if err := kacker.Run(flag.Args()[0], strings.Split(*packerFlags, " ")); err != nil {
		panic(err)
	}
}
