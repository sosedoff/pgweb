package main

import (
	"os"
	"github.com/sosedoff/pgweb/pkg/cli"
)

func main() {
	cli.InitOptions(os.Args)
	dummyAuxCloser := make(chan int);
	cli.Run(dummyAuxCloser)
}
