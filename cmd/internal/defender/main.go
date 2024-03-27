package main

import (
	"flag"
	"log"
	"os"
	"strings"

	denfender "github.com/goplus/gop/defender"
)

func main() {
	excludeDirs := []string{}
	excludeString := flag.String("e", "", "exclude the `[path]` to exclude defender scan.")
	flag.Parse()
	flag.Usage = func() {
		println("exclude [-e] dir1,dir2,file1,file2...")
	}
	if len(*excludeString) > 0 {
		excludeDirList := strings.Split(*excludeString, ",")
		for _, excludeDir := range excludeDirList {
			_, err := os.Stat(excludeDir)
			if err != nil {
				continue
			}
			excludeDirs = append(excludeDirs, excludeDir)
		}
		file, err := denfender.GenDefenderExcludePSFile(excludeDirs)
		if err != nil {
			log.Fatalln(err)
		}
		denfender.RunDefenderExcludePSFile(file)
	} else {
		if len(flag.Args()) > 0 {
			flag.Usage()
			flag.PrintDefaults()
		} else {
			file, err := denfender.GenGopDefenderExcludePSFile()
			if err != nil {
				log.Fatalln(err)
			}
			denfender.RunDefenderExcludePSFile(file)
		}
	}

}
