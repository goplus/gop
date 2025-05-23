package defender

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func excludePSFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gop/exclude.ps1"), nil
}

func removePSFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gop/remove.ps1"), nil
}

func Exclude() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
	file, err := excludePSFile()
	if err != nil {
		return err
	}
	err = genGopDefenderExcludePSFile(file)
	if err != nil {
		return err
	}
	return runDefenderExcludePSFile(file)
}

func Remove() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not supported for %s", runtime.GOOS)
	}
	exFile, err := excludePSFile()
	if err != nil {
		return err
	}
	file, err := removePSFile()
	if err != nil {
		return err
	}
	err = genGopDefenderRemovePSFile(exFile, file)
	if err != nil {
		return err
	}
	return runDefenderExcludePSFile(file)
}
