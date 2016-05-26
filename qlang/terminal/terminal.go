package terminal

import (
	"os"
	"strings"

	"github.com/peterh/liner"
)

type Terminal struct {
	state *liner.State
}

func New() *Terminal {
	term := &Terminal{liner.NewLiner()}
	term.state.SetCtrlCAborts(true)
	return term
}

func (term *Terminal) LoadHistroy(historyFile string) error {
	f, err := os.Open(historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	term.state.ReadHistory(f)
	return nil
}

func (term *Terminal) SaveHistroy(historyFile string) error {
	f, err := os.Create(historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	term.state.WriteHistory(f)
	return nil
}

var (
	ErrPromptAborted = liner.ErrPromptAborted
)

func (term *Terminal) Scan(prompt string, fnReadMore func(expr string, line string) (string, bool)) (string, error) {
	var all string
	var more bool
	for {
		line, err := term.state.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				return all, ErrPromptAborted
			}
			return all, err
		}
		if strings.TrimSpace(line) != "" {
			term.state.AppendHistory(line)
		}
		if fnReadMore == nil {
			return line, nil
		}
		all, more = fnReadMore(all, line)
		if !more {
			break
		}
		prompt = "..."
	}
	return all, nil
}
