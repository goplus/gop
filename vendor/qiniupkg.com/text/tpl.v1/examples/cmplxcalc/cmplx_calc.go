package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"qiniupkg.com/text/tpl.v1/cmplx"
	"qiniupkg.com/text/tpl.v1/interpreter"
)

var (
	calc   = cmplx.New()
	engine *interpreter.Engine
)

func eval(line string) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	line = strings.Trim(line, " \t\r\n")
	if line == "" {
		return
	}

	if err := engine.Eval(line); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	v, _ := calc.Ret()
	if imag(v) == 0 {
		fmt.Printf("> %v\n\n", real(v))
	} else {
		fmt.Printf("> %v\n\n", v)
	}
}

func main() {

	var err error
	if engine, err = interpreter.New(calc, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		eval(scanner.Text())
	}
}
