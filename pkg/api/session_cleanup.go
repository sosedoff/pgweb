package api

import (
	"log"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"
)

// StartSessionCleanup starts a goroutine to cleanup idle database sessions
func StartSessionCleanup() {
	for range time.Tick(time.Minute) {
		if command.Opts.Debug {
			log.Println("Triggering idle session deletion")
		}
		cleanupIdleSessions()
	}
}

func cleanupIdleSessions() {
	ids := []string{}

	// Figure out which sessions are idle
	for id, client := range DbSessions {
		if client.IsIdle() {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return
	}

	// Close and delete idle sessions
	log.Println("Closing", len(ids), "idle sessions")
	for _, id := range ids {
		// TODO: concurrent map edit will trigger panic
		if err := DbSessions[id].Close(); err == nil {
			delete(DbSessions, id)
		}
	}
}
