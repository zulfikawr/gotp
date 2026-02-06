package main

import (
	"fmt"
	"os"

	"github.com/zulfikawr/gotp/internal/cli"
	"github.com/zulfikawr/gotp/internal/cli/ui"
)

var version = "0.1.1"

func main() {
	rootCmd := cli.NewRootCmd()
	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", ui.DangerBright, err, ui.Reset)
		os.Exit(1)
	}
}
