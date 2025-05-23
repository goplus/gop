package main

import (
	"fmt"
	"os"

	"github.com/goplus/gop/defender"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Start processing, please wait for a moment ...")
		if err := defender.Exclude(); err != nil {
			fmt.Println(err)
		}
		return
	}
	switch os.Args[1] {
	case "add":
		fmt.Println("Start processing, please wait for a moment ...")
		if err := defender.Exclude(); err != nil {
			fmt.Println(err)
		}
	case "remove":
		fmt.Println("Start processing, please wait for a moment ...")
		if err := defender.Remove(); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("expected 'add' or 'remove' subcommands")
		os.Exit(1)
	}
}
