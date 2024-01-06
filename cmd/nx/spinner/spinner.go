package spinner

import (
	"fmt"
	"os"
	"time"

	"github.com/theckman/yacspin"
)

var spinner *yacspin.Spinner

var isDevice bool

func init() {
	s, _ := os.Stderr.Stat()
	isDevice = (s.Mode() & os.ModeCharDevice) > 0

	var err error
	spinner, err = yacspin.New(yacspin.Config{
		TerminalMode: yacspin.ForceTTYMode,
		Writer:       os.Stderr,
		Frequency:    50 * time.Millisecond,
		CharSet:      yacspin.CharSets[14],
		Suffix:       " ",
	})

	if err != nil {
		panic(err)
	}
}

func Start() error {
	if isDevice {
		return spinner.Start()
	} else {
		return nil
	}
}

func Stop() error {
	if isDevice {
		return spinner.Stop()
	} else {
		return nil
	}
}

func Message(message string) {
	if isDevice {
		spinner.Message(message)
	}
}

func Printf(format string, a ...any) {
	err := spinner.Stop()
	fmt.Printf(format, a...)
	if err == nil {
		spinner.Start()
	}
}

func Do(message string, action func()) {
	Message(message)
	Start()
	defer Stop()
	action()
}

func DoErr(message string, action func() error) error {
	Message(message)
	Start()
	defer Stop()
	return action()
}
