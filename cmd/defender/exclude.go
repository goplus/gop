package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goplus/gop/defender"
)

func main() {
	excludeDirs := []string{}
	excludeString := flag.String("e", "", "exclude the `[path...]` to exclude defender scan.")
	flag.Parse()
	flag.Usage = func() {
		println("Usage of exclude.exe")
	}

	if len(*excludeString) > 0 {
		fmt.Println("Start processing, please wait for a moment...")
		excludeDirList := strings.Split(*excludeString, ",")
		for _, excludeDir := range excludeDirList {
			_, err := os.Stat(excludeDir)
			if err != nil {
				continue
			}
			excludeDirs = append(excludeDirs, excludeDir)
		}
		defender.ExcludeWithDirs(excludeDirs)
	} else {
		if len(flag.Args()) > 0 {
			flag.Usage()
			flag.PrintDefaults()
		} else {
			fmt.Println("Start processing, please wait for a moment...")
			defender.Exclude()
		}
	}

}
