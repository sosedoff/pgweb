package main

import (
	"github.com/sosedoff/pgweb/pkg/cli"
)

func main() {
	dummyAuxCloser := make(chan int);
	cli.Run(dummyAuxCloser)
}
