package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"qiniupkg.com/text/tpl.v1/extractor"
)

const help = `
Usage:
  tplconv -o output-format -ext output-ext input-file1 [... input-fileN]
`

var (
	outputFormat = flag.String("o", "json", "output format")
	outputExt    = flag.String("ext", "", "output extension")
)

// tplconv -o output-format -ext output-ext input-file1 [... input-fileN]
//
func main() {

	root := os.Getenv("HOME") + "/.tplconv"

	flag.Parse()
	n := flag.NArg()
	if n == 0 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	for i := 0; i < n; i++ {
		err := processFile(root, flag.Arg(i))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}

func processFile(root, inputFile string) (err error) {

	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return
	}

	ext := filepath.Ext(inputFile)
	if ext == "" {
		return fmt.Errorf("unknown file format: %s", inputFile)
	}

	inputTempFile := filepath.Join(root, "input", ext[1:]+".tpl")
	tpl, err := ioutil.ReadFile(inputTempFile)
	if err != nil {
		return fmt.Errorf("unknown file format: %s", ext)
	}

	e, err := extractor.New(tpl, nil)
	if err != nil {
		err = fmt.Errorf("invalid tpl file %s:\n ==> %v", inputTempFile, err)
		return
	}

	doc, err := e.SafeParse(input, inputFile)
	if err != nil {
		return
	}

	w := os.Stdout
	if *outputExt != "" {
		outFile := inputFile + "." + *outputExt
		f, err2 := os.Create(outFile)
		if err2 != nil {
			return err2
		}
		defer f.Close()
		w = f
	}

	outext := filepath.Ext(*outputFormat)
	if outext == "" {
		print, ok := outInternals[*outputFormat]
		if !ok {
			return fmt.Errorf("unknown output format: %s", *outputFormat)
		}
		io.WriteString(w, print(doc))
		return nil
	}

	gen, ok := outTemplates[outext]
	if !ok {
		return fmt.Errorf("unknown output format: %s", outext)
	}

	outputTempFile := filepath.Join(root, "output", *outputFormat)
	out, err := gen(outputTempFile)
	if err != nil {
		return
	}

	return out.Execute(w, doc)
}

// -----------------------------------------------------------------------------

func printJson(v interface{}) string {

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

var funcMap = template.FuncMap{
	"json": printJson,
}

func makeHtmlTempl(outputTempFile string) (out iExecutor, err error) {

	return template.New(filepath.Base(outputTempFile)).Funcs(funcMap).ParseFiles(outputTempFile)
}

// -----------------------------------------------------------------------------

type iExecutor interface {
	Execute(w io.Writer, doc interface{}) error
}

var outTemplates = map[string]func(outputTempFile string) (iExecutor, error){
	".htm":  makeHtmlTempl,
	".html": makeHtmlTempl,
}

var outInternals = map[string]func(v interface{}) string{
	"json": printJson,
}

// -----------------------------------------------------------------------------
