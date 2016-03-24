package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var version string = "develop"

func init() {
	version = strings.TrimRight(version, " ")
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "-version") {
		Copyright()
		os.Exit(0)
	}
}

func Copyright() {
	fmt.Printf("Q-language - http://qlang.io, version qlang-%s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	fmt.Println("Copyright (C) 2015 Qiniu.com - Shanghai Qiniu Information Technologies Co., Ltd.\n")
}

func Version() string {
	return version
}

