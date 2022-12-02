package api

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sosedoff/pgweb/pkg/client"
)

type SessionManager struct {
	logger   *logrus.Logger
	sessions map[string]*client.Client
	mu       sync.Mutex
}

func NewSessionManager(logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		logger:   logger,
		sessions: map[string]*client.Client{},
		mu:       sync.Mutex{},
	}
}

func (m *SessionManager) IDs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := []string{}
	for k := range m.sessions {
		ids = append(ids, k)
	}

	return ids
}

func (m *SessionManager) Sessions() map[string]*client.Client {
	m.mu.Lock()
	sessions := m.sessions
	defer m.mu.Unlock()

	return sessions
}

func (m *SessionManager) Get(id string) *client.Client {
	m.mu.Lock()
	c := m.sessions[id]
	m.mu.Unlock()

	return c
}

func (m *SessionManager) Add(id string, conn *client.Client) {
	m.mu.Lock()
	m.sessions[id] = conn
	m.mu.Unlock()
}

func (m *SessionManager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.sessions[id]
	if ok {
		conn.Close()
		delete(m.sessions, id)
	}

	return ok
}

func (m *SessionManager) Len() int {
	m.mu.Lock()
	sz := len(m.sessions)
	m.mu.Unlock()

	return sz
}

func (m *SessionManager) Cleanup() int {
	removed := 0

	m.logger.Debug("starting idle sessions cleanup")
	defer func() {
		m.logger.Debug("removed idle sessions:", removed)
	}()

	for _, id := range m.staleSessions() {
		m.logger.WithField("id", id).Debug("closing stale session")
		if m.Remove(id) {
			removed++
		}
	}

	return removed
}

func (m *SessionManager) RunPeriodicCleanup() {
	for range time.Tick(time.Minute) {
		m.Cleanup()
	}
}

func (m *SessionManager) staleSessions() []string {
	m.mu.TryLock()
	defer m.mu.Unlock()

	ids := []string{}
	for id, conn := range m.sessions {
		if conn.IsIdle() {
			ids = append(ids, id)
		}
	}

	return ids
}
