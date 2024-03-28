package defender

import (
	"fmt"
	"log"
	"runtime"
)

func ExcludeWithDirs(excludeDirs []string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
	file, err := genDefenderExcludePSFile(excludeDirs)
	if err != nil {
		log.Fatalln(err)
	}
	_, resultError := runDefenderExcludePSFile(file)
	return resultError
}

func Exclude() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
	file, err := genGopDefenderExcludePSFile()
	if err != nil {
		log.Fatalln(err)
	}
	_, resultError := runDefenderExcludePSFile(file)
	return resultError
}
