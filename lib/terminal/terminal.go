package terminal

import (
	"os"
	"strings"

	"github.com/peterh/liner"
)

type Terminal struct {
	*liner.State
	promptFirst string
	promptNext  string
	fnReadMore  func(expr string, line string) (string, bool)
}

func New(promptFirst, promptNext string, fnReadMore func(expr string, line string) (string, bool)) *Terminal {
	term := &Terminal{
		State:       liner.NewLiner(),
		promptFirst: promptFirst,
		promptNext:  promptNext,
		fnReadMore:  fnReadMore,
	}
	term.SetCtrlCAborts(true)
	return term
}

func (term *Terminal) LoadHistroy(historyFile string) error {
	f, err := os.Open(historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	term.ReadHistory(f)
	return nil
}

func (term *Terminal) SaveHistroy(historyFile string) error {
	f, err := os.Create(historyFile)
	if err != nil {
		return err
	}
	defer f.Close()
	term.WriteHistory(f)
	return nil
}

var (
	ErrPromptAborted = liner.ErrPromptAborted
)

func (term *Terminal) Scan() (string, error) {
	var all string
	var more bool
	prompt := term.promptFirst
	fnReadMore := term.fnReadMore
	for {
		line, err := term.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				return all, ErrPromptAborted
			}
			return all, err
		}
		if strings.TrimSpace(line) != "" {
			term.AppendHistory(line)
		}
		if fnReadMore == nil {
			return line, nil
		}
		all, more = fnReadMore(all, line)
		if !more {
			break
		}
		prompt = term.promptNext
	}
	return all, nil
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"new":       New,
	"supported": liner.TerminalSupported,
	"mode":      liner.TerminalMode,

	"New":       New,
	"Supported": liner.TerminalSupported,
	"Mode":      liner.TerminalMode,

	"ErrPromptAborted": ErrPromptAborted,
}
