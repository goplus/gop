// +build windows

package readline

import "syscall"

func SuspendMe() {
}

func GetStdin() int {
	return int(syscall.Stdin)
}

func init() {
	isWindows = true
}

// get width of the terminal
func GetScreenWidth() int {
	info, _ := GetConsoleScreenBufferInfo()
	if info == nil {
		return -1
	}
	return int(info.dwSize.x)
}

func DefaultIsTerminal() bool {
	return true
}

func DefaultOnWidthChanged(func()) {

}
