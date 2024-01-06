package util

import (
	"regexp"
	"strings"
)

type Filter interface {
	// Checks whether the specified string matches the filter pattern.
	Match(string) bool
	// Checks whether the specified string should be filtered out.
	// Returns false for empty patterns.
	Filter(string) bool
}

type FilterSet []Filter

func (fs *FilterSet) Match(s string) bool {
	if len(*fs) == 0 {
		return false
	}
	for _, f := range *fs {
		if !f.Match(s) {
			return false
		}
	}
	return true
}

func (fs *FilterSet) Filter(s string) bool {
	for _, f := range *fs {
		if f.Filter(s) {
			return true
		}
	}
	return false
}

type Wildcard struct {
	Anchor  bool
	Pattern string
	rgx     *regexp.Regexp
}

func (w *Wildcard) String() string {
	return w.Pattern
}

func (w *Wildcard) Set(value string) (err error) {
	w.Pattern = strings.TrimSpace(value)
	if w.Pattern == "" {
		w.rgx = nil
	} else {
		p := strings.ReplaceAll(regexp.QuoteMeta(w.Pattern), `\*`, `.*`)
		if w.Anchor {
			p = "(?i)^" + p + "$"
		} else {
			p = "(?i)" + p
		}
		w.rgx, err = regexp.Compile(p)
	}
	return
}

func (w *Wildcard) Type() string {
	return "string"
}

// Returns true if the search pattern is non-empty and matches the specified string.
func (w *Wildcard) Match(s string) bool {
	if w.rgx == nil {
		return false
	}
	return w.rgx.MatchString(s)
}

// Returns true if the search pattern is non-empty and does not match the specified string.
func (w *Wildcard) Filter(s string) bool {
	if w.rgx == nil {
		return false
	}
	return !w.rgx.MatchString(s)
}
