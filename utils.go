package main

import (
	"fmt"
	"runtime"
	"time"
)

const MEGABYTE = 1024 * 1024

func startRuntimeProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		fmt.Println("-----------------------")
		fmt.Println("Goroutines:", runtime.NumGoroutine())
		fmt.Println("Memory acquired:", m.Sys, "bytes,", m.Sys/MEGABYTE, "mb")
		fmt.Println("Memory used:", m.Alloc, "bytes,", m.Alloc/MEGABYTE, "mb")

		time.Sleep(time.Minute)
	}
}
