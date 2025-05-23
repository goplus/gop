package defender

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func genGopDefenderExcludePSFile(psFile string) error {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	resultMap := map[string]string{
		"gop env GOCACHE": "",
		"gop env GOPROOT": "",
		"gop env GOBIN":   "",
		"gop env GOENV":   "",
		"gop env GOPATH":  ""}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer stdin.Close()
		for k := range resultMap {
			fmt.Fprintln(stdin, k)
		}
		wg.Done()
	}()
	wg.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(out)
	scanner := bufio.NewScanner(buf)
	var key = ""
	for scanner.Scan() {
		lineText := scanner.Text()
		index := strings.Index(lineText, "gop env")
		if index > 0 {
			key = lineText[index:]
		} else {
			if len(key) > 0 {
				resultMap[key] = lineText
			}
			key = ""
		}
	}
	f, err := os.Create(psFile)
	if err != nil {
		return err
	}
	defer f.Close()
	for k, v := range resultMap {
		if len(v) <= 0 || k == "gop env GOPROOT" {
			continue
		}
		if k == "gop env GOCACHE" {
			localDir := filepath.Dir(v)
			goplsDir := filepath.Join(localDir, "/gopls")
			if err := writeExcludeDirToFile(f, goplsDir); err != nil {
				return err
			}
		}
		if err := writeExcludeDirToFile(f, v); err != nil {
			return err
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		if err := writeExcludeDirToFile(f, filepath.Join(home, ".gopls")); err != nil {
			return err
		}
		if err := writeExcludeDirToFile(f, filepath.Join(home, ".gop")); err != nil {
			return err
		}
		var exts = [11]string{
			"gop",
			"go",
			"yap",
			"spx",
			"rdx",
			"gmx",
			"gox",
			"gsh",
			"mod",
			"sum",
		}
		for _, extStr := range exts {
			if err := writeExcludeExtToFile(f, extStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func genGopDefenderRemovePSFile(exFile string, file string) error {
	if _, err := os.Stat(exFile); err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer stdin.Close()
		cmdStr := fmt.Sprintf("powershell -Command \"(gc %s) -replace 'Add-MpPreference', 'Remove-MpPreference' | Out-File -encoding ASCII %s\"", exFile, file)
		fmt.Fprintln(stdin, cmdStr)
		wg.Done()
	}()
	wg.Wait()
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func writeExcludeDirToFile(f *os.File, dir string) error {
	_, err := f.WriteString(excludeDirString(dir))
	if err != nil {
		return err
	}
	_, newLineError := f.WriteString("\r\n")
	if newLineError != nil {
		return newLineError
	}
	return nil
}

func writeExcludeExtToFile(f *os.File, dir string) error {
	_, err := f.WriteString(excludeExtString(dir))
	if err != nil {
		return err
	}
	_, newLineError := f.WriteString("\r\n")
	if newLineError != nil {
		return newLineError
	}
	return nil
}

func runDefenderExcludePSFile(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer stdin.Close()
		startProcess := fmt.Sprintf("start-process -WindowStyle hidden powershell -verb runas -ArgumentList \"-file %s\"", file)
		fmt.Fprintln(stdin, startProcess)
		wg.Done()
	}()
	wg.Wait()
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
