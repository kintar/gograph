// Command gograph is the entrypoint for the gograph CLI tool.
package main

import (
	"os"

	"github.com/ozgurcd/gograph/internal/cli"
)

// version is set at compile time via -ldflags "-X main.version=x.y.z".
// Falls back to "dev" when built without ldflags.
var version = "dev"

func main() {
	if version != "dev" {
		cli.Version = version
	}
	os.Exit(cli.Run(os.Args[1:]))
}
