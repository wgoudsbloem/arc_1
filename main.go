package main

import (
	"os"

	cli "arcessio/cli"
)

func main() {
	cli.Start(os.Stdin, os.Stdout)
}
