package cliargs

import (
	"errors"
	"io"
	"strings"
)

// ErrUsage marks invalid command-line arguments.
var ErrUsage = errors.New("usage: args <name> [repeat]")

// Options contains validated user arguments.
type Options struct {
	Name   string
	Repeat int
}

// UserArgs returns a new slice without os.Args[0].
func UserArgs(all []string) []string {
	// TODO: implement according to docs/task.md.
	return nil
}

// Parse validates the arguments passed after os.Args[0].
func Parse(args []string) (Options, error) {
	// TODO: implement according to docs/task.md.
	return Options{}, nil
}

// Run writes the command result to the supplied writer.
func Run(args []string, out io.Writer) error {
	// TODO: implement according to docs/task.md.
	return nil
}

// Example returns deterministic output used by cmd and the example test.
func Example() string {
	var out strings.Builder
	_ = Run([]string{"Maria", "2"}, &out)
	return out.String()
}
