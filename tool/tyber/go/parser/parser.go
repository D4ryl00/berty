package parser

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"berty.tech/berty/tool/tyber/go/logger"
	session "berty.tech/berty/tool/tyber/go/session_manager"
	"berty.tech/berty/v2/go/pkg/tyber"
	"github.com/pkg/errors"
)

type Parser struct {
	ctx            context.Context
	logger         *logger.Logger
	listener       *net.TCPListener
	cancelNetwork  context.CancelFunc
	sessionManager *session.Manager
	EventChan      chan interface{} // TODO: base event definition
}

func New(ctx context.Context, l *logger.Logger, sm *session.Manager) *Parser {
	return &Parser{
		ctx:            ctx,
		logger:         l.Named("parser"),
		sessionManager: sm,
		EventChan:      make(chan interface{}),
	}
}

func (p *Parser) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		p.logger.Errorf("opening file failed: %v", err)
		return err
	}

	p.startSession(path, FileType, file)

	return nil
}

func (p *Parser) NetworkListen(address, port string) error {
	if p.listener != nil {
		p.cancelNetwork()
		p.listener = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancelNetwork = cancel

	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return err
	}

	p.listener, err = net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}

	go func(l *net.TCPListener) {
		p.logger.Infof("started listening on %s", l.Addr().String())
		defer l.Close()

		for {
			select {
			case <-ctx.Done():
				p.logger.Infof("stopped listening on %s", l.Addr().String())
				return

			default:
				if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
					p.logger.Errorf("can't set deadline on tcp listener: %v", err)
					return
				}

				conn, err := l.Accept()
				if err != nil {
					if os.IsTimeout(err) {
						continue
					}
					p.logger.Errorf("TCP accept error: %v", err)
					return
				}

				p.startSession(conn.RemoteAddr().String(), NetworkType, conn)
			}
		}
	}(p.listener)

	return nil
}

func (p *Parser) OpenSession(sessionID string) error {
	s, ok := p.sessionManager.OpenSession(sessionID)
	if !ok {
		return errors.New(fmt.Sprintf("session %s not found", sessionID))
	}

	var events []CreateTraceEvent
	s.tracesLock.Lock()
	for _, t := range s.Traces {
		events = append(events, t.ToCreateTraceEvent())
	}
	p.EventChan <- events
	s.tracesLock.Unlock()

	return nil
}

func (p *Parser) startSession(srcName string, srcType SrcType, srcIO io.ReadCloser) {
	s, err := newSession(srcName, srcType, srcIO, p.EventChan, p.logger)
	if err != nil {
		p.logger.Errorf("starting session failed with logs from %s (%s): %v", srcType, srcName, err)
		return
	}
	p.logger.Infof("started session %s with logs from %s (%s)", s.ID, srcType, srcName)

	s.StatusType = tyber.Running
	p.sessionManager.AddSession(s)

	wait := make(chan interface{})

	go func() {
		if err := s.parseLogs(); err != nil {
			p.logger.Errorf("parsing session %s logs from %s (%s) failed: %v", s.ID, srcType, srcName, err)
		} else {
			p.logger.Infof("successfully parsed session %s logs from %s (%s)", s.ID, srcType, srcName)
		}

		if s.canceled {
			return
		}

		p.sessionManager.SaveSession(s)

		if srcType == FileType {
			wait <- struct{}{}
		}
	}()

	if srcType == FileType {
		<-wait
	}

	p.EventChan <- SessionToCreateEvent(s)
}
