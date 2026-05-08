// Command gograph is the entrypoint for the gograph CLI tool.
package main

import (
	"os"

	"gograph/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
