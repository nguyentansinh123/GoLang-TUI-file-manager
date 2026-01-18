package main

import (
	"fmt"
	"os"

	"github.com/tansinhnguyen123/my-tui-explorer/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	explorer := ui.NewExplorer()
	return explorer.Run()
}
