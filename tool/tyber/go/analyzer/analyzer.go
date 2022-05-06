package analyzer

import (
	"sync"

	"berty.tech/berty/tool/tyber/go/logger"
	"berty.tech/berty/tool/tyber/go/parser"
)

type Analyzer struct {
	logger       *logger.Logger
	sessions     []*parser.Session
	sessionsLock sync.RWMutex
}

func New(l *logger.Logger) *Analyzer {
	return &Analyzer{
		logger: l.Named("analyzer"),
	}
}

func (a *Analyzer) AddSessions(s []*parser.Session) {
	a.sessionsLock.Lock()
	a.sessions = append(a.sessions, s...)
	a.sessionsLock.Unlock()
}
