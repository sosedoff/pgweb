package api

import (
	"log"
	"time"

	"github.com/sosedoff/pgweb/pkg/command"
)

func cleanupIdleSessions() {
	ids := []string{}

	// Figure out which sessions are idle
	for id, client := range DbSessions {
		if client.IsIdle() {
			ids = append(ids, id)
		}
	}

	// Close and delete idle sessions
	if len(ids) == 0 {
		return
	}
	log.Println("Closing", len(ids), "idle sessions")
	for _, id := range ids {
		// TODO: concurrent map edit will trigger panic
		if err := DbSessions[id].Close(); err == nil {
			delete(DbSessions, id)
		}
	}
}

func StartSessionCleanup() {
	ticker := time.NewTicker(time.Minute)

	for {
		<-ticker.C

		if command.Opts.Debug {
			log.Println("Triggering idle session deletion")
		}

		cleanupIdleSessions()
	}
}
