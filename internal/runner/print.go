package runner

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// A PrintRunner only prints out what it WOULD run, rather than actually
// running it.
type PrintRunner struct {
	logger *log.Logger
}

func NewPrintRuner(l *log.Logger) *PrintRunner {
	return &PrintRunner{logger: l}
}

func (r PrintRunner) Run(cmd string, options Options) error {
	fmt.Fprintln(options.Stdout, cmd)
	return nil
}

func (r PrintRunner) RunArgs(cmd string, args []string, options Options) error {
	_, e := fmt.Fprintf(options.Stdout, "cmd: %s args: %s\n", cmd, args)
	return e
}
