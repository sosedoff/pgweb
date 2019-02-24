package main

import (
	"os"
	"github.com/sosedoff/pgweb/pkg/cli"
)

func main() {
	cli.InitOptions(os.Args)
	cli.Run()
}
