// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

package readline

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// SuspendMe use to send suspend signal to myself, when we in the raw mode.
// For OSX it need to send to parent's pid
// For Linux it need to send to myself
func SuspendMe() {
	p, _ := os.FindProcess(os.Getppid())
	p.Signal(syscall.SIGTSTP)
	p, _ = os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTSTP)
}

// get width of the terminal
func getWidth(stdoutFd int) int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(stdoutFd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		_ = errno
		return -1
	}
	return int(ws.Col)
}

func GetScreenWidth() int {
	return getWidth(syscall.Stdout)
}

func DefaultIsTerminal() bool {
	return IsTerminal(syscall.Stdin) && IsTerminal(syscall.Stdout)
}

func GetStdin() int {
	return syscall.Stdin
}

// -----------------------------------------------------------------------------

var (
	widthChange         sync.Once
	widthChangeCallback func()
)

func DefaultOnWidthChanged(f func()) {
	widthChangeCallback = f
	widthChange.Do(func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)

		go func() {
			for {
				_, ok := <-ch
				if !ok {
					break
				}
				widthChangeCallback()
			}
		}()
	})
}
