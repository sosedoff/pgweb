package api

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/metrics"
)

type SessionManager struct {
	logger      *logrus.Logger
	sessions    map[string]*client.Client
	mu          sync.Mutex
	idleTimeout time.Duration
}

func NewSessionManager(logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		logger:   logger,
		sessions: map[string]*client.Client{},
		mu:       sync.Mutex{},
	}
}

func (m *SessionManager) SetIdleTimeout(timeout time.Duration) {
	m.idleTimeout = timeout
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
	m.mu.Unlock()

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
	defer m.mu.Unlock()

	m.sessions[id] = conn
	metrics.SetSessionsCount(len(m.sessions))
}

func (m *SessionManager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.sessions[id]
	if ok {
		conn.Close()
		delete(m.sessions, id)
	}

	metrics.SetSessionsCount(len(m.sessions))
	return ok
}

func (m *SessionManager) Len() int {
	m.mu.Lock()
	sz := len(m.sessions)
	m.mu.Unlock()

	return sz
}

func (m *SessionManager) Cleanup() int {
	if m.idleTimeout == 0 {
		return 0
	}

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
	m.logger.WithField("timeout", m.idleTimeout).Info("session manager cleanup enabled")

	for range time.Tick(time.Minute) {
		m.Cleanup()
	}
}

func (m *SessionManager) staleSessions() []string {
	m.mu.TryLock()
	defer m.mu.Unlock()

	now := time.Now()
	ids := []string{}

	for id, conn := range m.sessions {
		if now.Sub(conn.LastQueryTime()) > m.idleTimeout {
			ids = append(ids, id)
		}
	}

	return ids
}
