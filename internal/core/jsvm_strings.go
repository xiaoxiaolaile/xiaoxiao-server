package core

import "strings"

type Strings struct {
}

func (sender *Strings) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (sender *Strings) Replace(s string, old string, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func (sender *Strings) ReplaceAll(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func (sender *Strings) Split(s, sep string, n int) []string {
	return strings.SplitN(s, sep, n)
}
