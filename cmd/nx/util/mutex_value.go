package util

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

type MutexValue struct {
	selected string
	switches []string
}

func (mutex *MutexValue) Selected() string {
	return mutex.selected
}

func (mutex *MutexValue) Switcher(value string) *mutexFlag {
	if slices.Contains(mutex.switches, value) {
		panic(fmt.Errorf("duplicate switcher defined: %q", value))
	}
	mutex.switches = append(mutex.switches, value)
	return &mutexFlag{
		mutex: mutex,
		value: value,
	}
}

func (mutex *MutexValue) Switch(fs *pflag.FlagSet, value string, usage string) {
	fs.Var(mutex.Switcher(value), value, usage)
	fs.Lookup(value).NoOptDefVal = "true"
}

type mutexFlag struct {
	mutex *MutexValue
	value string
}

func (s *mutexFlag) Type() string {
	return "bool"
}

func (*mutexFlag) String() string {
	return "false"
}

func (s *mutexFlag) Set(value string) error {
	if value == "true" {
		if s.mutex.selected != "" {
			switches := strings.Join(s.mutex.switches, ", ")
			return fmt.Errorf("only one of %s can be set", switches)
		}
		s.mutex.selected = s.value
	}
	return nil
}

func (*mutexFlag) IsBoolFlag() bool {
	return true
}
