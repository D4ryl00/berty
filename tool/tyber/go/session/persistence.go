package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (m *Manager) SaveSessionFile(sessionID string, path string) error {
	m.sessionsLock.RLock()
	session, ok := m.sessions.Get(sessionID)
	m.sessionsLock.RUnlock()

	if !ok {
		return errors.New("session ID not found")
	}

	return m.saveSessionFile(session.(*Session), path)
}

func (m *Manager) saveSessionFile(s *Session, path string) error {
	content, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')

	return ioutil.WriteFile(path, content, 0644)
}

func (m *Manager) restoreSessionFile(sessionID string, path string) (*Session, error) {
	s := &Session{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fmt.Println("content file", string(content))
	if err = json.Unmarshal(content, &s); err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Manager) deleteSessionFile(sessionID string) error {
	path := filepath.Join(m.sessionPath, fmt.Sprintf("%s.json", sessionID))
	return os.Remove(path)
}

type SessionsIndex struct {
	Index []string `json:"index"`
}

func (m *Manager) saveSessionsIndexFile() error {
	index := &SessionsIndex{}
	for pair := m.sessions.Oldest(); pair != nil; pair = pair.Next() {
		index.Index = append(index.Index, pair.Key.(string))
	}

	content, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')

	path := filepath.Join(m.sessionPath, "sessions.json")

	return ioutil.WriteFile(path, content, 0644)
}

func (m *Manager) restoreSessionsIndexFile() ([]string, error) {
	path := filepath.Join(m.sessionPath, "sessions.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	index := &SessionsIndex{}
	if err = json.Unmarshal(content, index); err != nil {
		return nil, err
	}

	return index.Index, nil
}
