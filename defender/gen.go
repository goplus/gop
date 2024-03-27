package defender

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func genGopDefenderExcludePSFile() (string, error) {
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return "", err
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
		log.Fatal(err)
		return "", err
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
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	psFile := filepath.Join(home, "ex.ps1")
	f, err := os.Create(psFile)
	if err != nil {
		return "", err
	}
	defer f.Close()
	for k, v := range resultMap {
		if len(v) <= 0 || k == "gop env GOPROOT" {
			continue
		}
		if k == "gop env GOCACHE" {
			localDir := filepath.Dir(v)
			goplsDir := filepath.Join(localDir, "/gopls")
			writeExcludeDirToFile(f, goplsDir)
		}
		writeExcludeDirToFile(f, v)
	}
	return f.Name(), nil
}

func writeExcludeDirToFile(f *os.File, dir string) {
	_, err := f.WriteString(excludeCmdString(dir))
	if err != nil {
		log.Println(err)
	}
	_, newLineError := f.WriteString("\r\n")
	if newLineError != nil {
		log.Println(newLineError)
	}
}

func genDefenderExcludePSFile(excludeDirs []string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	psFile := filepath.Join(home, "ex.ps1")
	f, err := os.Create(psFile)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	defer f.Close()
	for _, v := range excludeDirs {
		writeExcludeDirToFile(f, v);
	}
	return f.Name(), nil
}

func runDefenderExcludePSFile(file string) (string, error) {
	_, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("powershell", "-nologo", "-noprofile")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer stdin.Close()
		startProcess := fmt.Sprintf("start-process powershell -verb runas -ArgumentList \"-file %s\"", file)
		fmt.Fprintln(stdin, startProcess)
		wg.Done()
	}()
	wg.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(out), nil
}
