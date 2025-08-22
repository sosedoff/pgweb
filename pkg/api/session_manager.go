package api

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/metrics"
)

type Session struct {
	Client         *client.Client
	SessionExpiry  time.Time
	SessionRefresh func(*Session) (*Session, error)
}

type SessionManager struct {
	logger      *logrus.Logger
	sessions    map[string]*Session
	mu          sync.Mutex
	idleTimeout time.Duration
}

func NewSessionManager(logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		logger:   logger,
		sessions: map[string]*Session{},
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

	mapOfClients := map[string]*client.Client{}
	for k, v := range sessions {
		mapOfClients[k] = v.Client
	}

	return mapOfClients
}

func (m *SessionManager) GetSession(id string) *Session {
	m.mu.Lock()
	c := m.sessions[id]
	m.mu.Unlock()

	return c
}

func (m *SessionManager) Get(id string) *client.Client {
	m.mu.Lock()
	c := m.sessions[id]
	m.mu.Unlock()

	if c == nil {
		return nil
	}

	return c.Client
}

func (m *SessionManager) Add(id string, conn *client.Client) {
	m.AddSession(id, &Session{
		Client: conn,
	})
}

func (m *SessionManager) AddSession(id string, session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session.Client == nil && session.SessionRefresh != nil {
		var err error
		session, err = session.SessionRefresh(session)
		if err != nil {
			return err
		}
	}
	m.sessions[id] = session

	metrics.SetSessionsCount(len(m.sessions))
	return nil
}

func (m *SessionManager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if ok {
		session.Client.Close()
		delete(m.sessions, id)
	}

	metrics.SetSessionsCount(len(m.sessions))
	return ok
}

func (m *SessionManager) RefreshSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		// session not found
		return nil
	}

	if session.SessionRefresh == nil || session.SessionExpiry.IsZero() {
		// ClientFactory or SessionExpiry is not set so it is impossible to refresh
		// the session
		return nil
	}

	if session.SessionExpiry.After(time.Now()) {
		// session has not expired yet
		return nil
	}

	session, err := session.SessionRefresh(session)
	if err != nil {
		return err
	}

	m.sessions[id] = session

	return nil
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
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	ids := []string{}

	for id, session := range m.sessions {
		if now.Sub(session.Client.LastQueryTime()) > m.idleTimeout {
			ids = append(ids, id)
		}
	}

	return ids
}

func (m *SessionManager) RefreshSessions() error {
	m.mu.Lock()
	sessions := m.sessions
	m.mu.Unlock()

	for id := range sessions {
		if err := m.RefreshSession(id); err != nil {
			return err
		}
	}

	return nil
}

func (m *SessionManager) RunPeriodicRefresh() {
	m.logger.Info("session manager refresh enabled")

	for range time.Tick(time.Minute) {
		if err := m.RefreshSessions(); err != nil {
			// TODO: better error handling and logging
			m.logger.Error("error refreshing sessions:", err)
		}
	}
}
