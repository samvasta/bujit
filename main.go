package main

import (
	"os"

	"samvasta.com/bujit/cli"
)

func main() {
	args := os.Args[1:]

	if args[0] == "cli" {
		cli.StartInteractive()
	}
}
