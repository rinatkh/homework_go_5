package main

import (
	"fmt"
	"os"

	"github.com/rinatkh/homework_go_5/internal/cliargs"
)

func main() {
	args := cliargs.UserArgs(os.Args)
	if err := cliargs.Run(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
