package main

import (
	"fmt"
	"runtime"
	"time"
)

func startRuntimeProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		fmt.Println("-----------------------")
		fmt.Println("Goroutines:", runtime.NumGoroutine())
		fmt.Println("Memory acquired:", m.Sys)
		fmt.Println("Memory used:", m.Alloc)

		time.Sleep(time.Minute)
	}
}
