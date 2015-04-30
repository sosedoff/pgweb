package util

import (
	"log"
	"os"
	"runtime"
	"time"
)

const MEGABYTE = 1024 * 1024

func runProfiler() {
	logger := log.New(os.Stdout, "", 0)
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.Printf(
			"[DEBUG] Goroutines: %v, Mem used: %v (%v mb), Mem acquired: %v (%v mb)\n",
			runtime.NumGoroutine(),
			m.Alloc, m.Alloc/MEGABYTE,
			m.Sys, m.Sys/MEGABYTE,
		)

		time.Sleep(time.Second * 30)
	}
}

func StartProfiler() {
	go runProfiler()
}
