package defender

import (
	"log"
	"runtime"
)

func ExcludeWithDirs(excludeDirs []string) {
	if runtime.GOOS != "windows" {
		return
	}
	file, err := genDefenderExcludePSFile(excludeDirs)
	if err != nil {
		log.Fatalln(err)
	}
	runDefenderExcludePSFile(file)
}

func Exclude() {
	if runtime.GOOS != "windows" {
		return
	}
	file, err := genGopDefenderExcludePSFile()
	if err != nil {
		log.Fatalln(err)
	}
	runDefenderExcludePSFile(file)
}
