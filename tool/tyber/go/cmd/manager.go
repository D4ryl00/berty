package cmd

import (
	"context"
	"errors"
	"log"
	"time"

	"berty.tech/berty/tool/tyber/go/config"
	"berty.tech/berty/tool/tyber/go/logger"
	"berty.tech/berty/tool/tyber/go/parser"
	"berty.tech/berty/tool/tyber/go/session"
)

type Manager struct {
	ctx    context.Context
	cancel func()

	Config         *config.Config
	SessionManager *session.Manager
	Parser         *parser.Parser
	DataPath       string

	logger *logger.Logger
}

func New(ctx context.Context, cancel func()) *Manager {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())
	logger := logger.New(l, func(log *logger.Log) {}).Named("manager")

	return &Manager{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
	}
}

func (m *Manager) Init() error {
	m.Config = config.New(m.ctx, m.logger)
	m.SessionManager = session.New(m.ctx, m.logger)
	m.Parser = parser.New(m.ctx, m.logger, m.SessionManager)

	// init
	if err := m.Config.Init(m.DataPath); err != nil {
		m.cancel()
		return err
	}

	sessionPath, err := m.Config.GetSessionsPath()
	if err != nil {
		m.logger.Errorf("config getting session path failed: %v", err)
		m.cancel()
		return err
	}

	cerr := make(chan error)
	go func() {
		if err = m.SessionManager.Init(sessionPath); err != nil {
			m.logger.Errorf("parser init error: %v", err)
			m.cancel()
			cerr <- err
		}
		cerr <- nil
	}()

	select {
	case evt := <-m.SessionManager.EventChan:
		if _, ok := evt.([]session.CreateSessionEvent); !ok {
			m.cancel()
			return errors.New("parser init: wrong event received")
		}
	case <-time.After(2 * time.Second):
		m.cancel()
		return errors.New("parser init timeout")
	}

	if err = <-cerr; err != nil {
		m.cancel()
		return err
	}

	return nil
}

func (m *Manager) Cancel() {
	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
}
