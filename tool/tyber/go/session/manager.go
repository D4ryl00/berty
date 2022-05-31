package session

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"berty.tech/berty/tool/tyber/go/logger"
	"github.com/pkg/errors"
	orderedmap "github.com/wk8/go-ordered-map"
)

type Manager struct {
	ctx           context.Context
	logger        *logger.Logger
	sessionPath   string
	sessions      *orderedmap.OrderedMap
	sessionsLock  sync.RWMutex
	openedSession *Session
	initialized   bool
	initLock      sync.RWMutex
	EventChan     chan interface{}
}

func New(ctx context.Context, l *logger.Logger) *Manager {
	return &Manager{
		ctx:       ctx,
		logger:    l,
		sessions:  orderedmap.New(),
		EventChan: make(chan interface{}),
	}
}

func (m *Manager) Init(sessionPath string) error {
	m.sessionPath = sessionPath

	sessionsIndex, err := m.restoreSessionsIndexFile()
	if err != nil {
		return errors.Wrap(err, "listing persisted sessions failed")
	}

	start := time.Now()
	m.logger.Debug("partially restoring sessions started")

	var events []CreateSessionEvent
	for _, sessionID := range sessionsIndex {
		path := filepath.Join(m.sessionPath, fmt.Sprintf("%s.json", sessionID))
		s, err := m.restoreSessionFile(sessionID, path)
		if err != nil {
			m.logger.Errorf("restoring session %s failed: %v", sessionID, err)
			continue
		}

		m.logger.Debugf("session %s restored successfully", sessionID)

		m.sessions.Set(sessionID, s)
		events = append(events, SessionToCreateEvent(s))
	}

	elapsed := time.Since(start)
	m.logger.Debugf("restoring sessions took: %s", elapsed)

	select {
	case m.EventChan <- events:
	case <-m.ctx.Done():
		return m.ctx.Err()
	}

	m.initLock.Lock()
	m.initialized = true
	m.initLock.Unlock()

	m.logger.Infof("initialization successful with session path %s", m.sessionPath)

	return nil
}

func (m *Manager) isInitialized() bool {
	m.initLock.RLock()
	defer m.initLock.RUnlock()
	return m.initialized
}

func (m *Manager) AddSession(s *Session) error {
	if !m.isInitialized() {
		return errors.New("Manager not initialized")
	}

	m.logger.Debugf("request save session %s logs from %s (%s)", s.ID, s.SrcType, s.SrcName)

	m.sessionsLock.Lock()
	m.sessions.Set(s.ID, s)
	m.sessionsLock.Unlock()

	m.sessionsLock.Lock()
	path := filepath.Join(m.sessionPath, fmt.Sprintf("%s.json", s.ID))
	if err := m.saveSessionFile(s, path); err != nil {
		m.logger.Errorf("saving session %s logs from %s (%s) failed: %v", s.ID, s.SrcType, s.SrcName, err)
	} else {
		m.logger.Debugf("successfully saved session %s logs from %s (%s)", s.ID, s.SrcType, s.SrcName)
		if err = m.saveSessionsIndexFile(); err != nil {
			m.logger.Errorf("saving sessions index file failed: %v", err)
		}
	}
	m.sessionsLock.Unlock()

	return nil
}

func (m *Manager) OpenSession(sessionID string) error {
	if !m.isInitialized() {
		return errors.New("Manager not initialized")
	}

	m.sessionsLock.RLock()
	v, ok := m.sessions.Get(sessionID)
	m.sessionsLock.RUnlock()

	if !ok {
		return errors.New(fmt.Sprintf("session %s not found", sessionID))
	}

	s := v.(*Session)
	m.sessionsLock.Lock()
	if m.openedSession != nil && m.openedSession.ID != s.ID && m.openedSession.isRunning() {
		m.openedSession.SetOpenned(false)
	}
	m.openedSession = s
	m.sessionsLock.Unlock()

	s.SetOpenned(true)

	m.EventChan <- s.GenerateCreateTraceEvents()

	return nil
}

func (m *Manager) ListSessions() {
	if m.isInitialized() {
		var events []CreateSessionEvent
		m.sessionsLock.RLock()
		for pair := m.sessions.Oldest(); pair != nil; {
			s := pair.Value.(*Session)
			events = append(events, SessionToCreateEvent(s))
			pair = pair.Next()
		}
		m.sessionsLock.RUnlock()

		m.EventChan <- events
	}
}

func (m *Manager) DeleteSession(sessionID string) {
	m.logger.Debugf("delete session requested: %s", sessionID)

	if m.isInitialized() {
		m.sessionsLock.Lock()
		if m.openedSession != nil && m.openedSession.ID == sessionID {
			m.openedSession.SetOpenned(false)
			m.openedSession = nil
		}

		v, ok := m.sessions.Delete(sessionID)
		if ok {
			s := v.(*Session)

			if s.isRunning() {
				s.SrcCloser().Close()
			} else {
				if err := m.deleteSessionFile(sessionID); err != nil {
					m.logger.Errorf("deleting session %s file failed: %v", sessionID, err)
				}
				if err := m.saveSessionsIndexFile(); err != nil {
					m.logger.Errorf("saving sessions index file failed: %v", err)
				}
			}
		}
		m.EventChan <- DeleteSessionEvent{ID: sessionID}
		m.sessionsLock.Unlock()
	}
}

func (m *Manager) DeleteAllSessions() {
	if m.isInitialized() {
		m.sessionsLock.Lock()
		for pair := m.sessions.Oldest(); pair != nil; {
			s := pair.Value.(*Session)
			pair = pair.Next()

			m.logger.Debugf("delete session requested: %s", s.ID)

			if m.openedSession != nil && m.openedSession.ID == s.ID {
				m.openedSession.SetOpenned(false)
				m.openedSession = nil
			}

			if s.isRunning() {
				s.SrcCloser().Close()
			} else {
				if err := m.deleteSessionFile(s.ID); err != nil {
					m.logger.Errorf("deleting session %s file failed: %v", s.ID, err)
				}
			}
			m.sessions.Delete(s.ID)
		}
		if err := m.saveSessionsIndexFile(); err != nil {
			m.logger.Errorf("saving sessions index file failed: %v", err)
		}
		m.EventChan <- []CreateSessionEvent{}
		m.sessionsLock.Unlock()
	}
}
